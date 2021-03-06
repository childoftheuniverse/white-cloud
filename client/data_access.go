package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/childoftheuniverse/red-cloud"
	"github.com/childoftheuniverse/red-cloud/common"
	"github.com/golang/protobuf/proto"
	etcd "go.etcd.io/etcd/clientv3"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
ErrNoDataNodes indicates that no data nodes could be found covering the
requested range.
*/
var ErrNoDataNodes = errors.New(
	"No data nodes could be found for the requested range")

type krClientConn struct {
	KeyRange   *common.KeyRange
	ClientConn *grpc.ClientConn
}

/*
DataAccessClient is a client library allowing users to have convenient
access to data stored in red-cloud without having to implement data node
discovery and other related algorithms.
*/
type DataAccessClient struct {
	instance   string
	etcdClient *etcd.Client
	tlsConfig  *tls.Config
	lock       sync.RWMutex

	/*
		Cache for table -> key range to data node mappings.
		This will be updated only in case a mapping error is reported by a
		client RPC.
	*/
	dataNodeRangeCache map[string][]*krClientConn

	// Cache for host:port -> client connection.
	clientConnCache map[string]*grpc.ClientConn
}

/*
NewDataAccessClient creates a new DataAccessClient.
*/
func NewDataAccessClient(
	instance string,
	etcdClient *etcd.Client,
	tlsConfig *tls.Config) *DataAccessClient {
	var rv = &DataAccessClient{
		instance:           instance,
		etcdClient:         etcdClient,
		tlsConfig:          tlsConfig,
		dataNodeRangeCache: make(map[string][]*krClientConn),
		clientConnCache:    make(map[string]*grpc.ClientConn),
	}
	return rv
}

/*
getRangeClients finds the gRPC client connections for the data nodes holding
parts of the specified key range.
*/
func (d *DataAccessClient) getRangeClients(
	ctx context.Context, table string, keyRange *common.KeyRange,
	forceFetch bool) ([]*grpc.ClientConn, error) {
	var span = trace.FromContext(ctx)
	var dialOpts []grpc.DialOption
	var resp *etcd.GetResponse
	var rangeClients []*krClientConn
	var md *redcloud.ServerTableMetadata
	var td *redcloud.ServerTabletMetadata
	var rv = make([]*grpc.ClientConn, 0)
	var ok bool
	var err error

	if d.tlsConfig == nil {
		dialOpts = append(dialOpts, grpc.WithInsecure())
		span.AddAttributes(trace.BoolAttribute("insecure", true))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(
			credentials.NewTLS(d.tlsConfig)))
	}

	span.AddAttributes(trace.BoolAttribute("force-fetch", forceFetch))
	if !forceFetch {
		d.lock.RLock()

		/*
			Attempt to find the covered table in the cache.
			TODO: check that the cover computed from cache actually covers the range
			completely.
		*/
		if rangeClients, ok = d.dataNodeRangeCache[table]; ok {
			var krcli *krClientConn
			for _, krcli = range rangeClients {
				if keyRange.ContainsRange(krcli.KeyRange) {
					rv = append(rv, krcli.ClientConn)
				}
			}

			d.lock.RUnlock()

			span.AddAttributes(
				trace.BoolAttribute("from-cache", true),
				trace.Int64Attribute("num-ranges", int64(len(rv))))
			return rv, nil
		}

		d.lock.RUnlock()
	}

	span.AddAttributes(
		trace.BoolAttribute("from-cache", false))

	// Fetch the table definition from etcd.
	md = new(redcloud.ServerTableMetadata)
	if resp, err = d.etcdClient.Get(
		ctx, common.EtcdTableConfigPath(d.instance, table)); err != nil {
		span.AddAttributes(trace.StringAttribute("error", err.Error()))
		span.Annotate(nil, "Error communicating with etcd")
		return []*grpc.ClientConn{}, err
	}

	if len(resp.Kvs) == 0 {
		span.Annotate(nil, "No data nodes assigned to table")
		return []*grpc.ClientConn{}, ErrNoDataNodes
	}

	if err = proto.Unmarshal(resp.Kvs[0].Value, md); err != nil {
		span.AddAttributes(trace.StringAttribute("error", err.Error()))
		span.Annotate(nil, "Table metadata corruption detected")
		return []*grpc.ClientConn{}, err
	}

	// Time to take the write lock and fill the cache.
	d.lock.Lock()
	defer d.lock.Unlock()

	// Start from an empty cache of the table.
	d.dataNodeRangeCache[table] = make([]*krClientConn, 0)

	for _, td = range md.Tablet {
		var kr = common.NewKeyRange(td.StartKey, td.EndKey)
		var hostPort = net.JoinHostPort(
			td.Host, strconv.FormatInt(int64(td.Port), 10))
		var client *grpc.ClientConn

		if client, ok = d.clientConnCache[hostPort]; !ok {
			if client, err = grpc.Dial(hostPort, dialOpts...); err != nil {
				span.AddAttributes(
					trace.StringAttribute("failing-host", hostPort),
					trace.StringAttribute("error", err.Error()))
				span.Annotate(nil, "Data node connection failed")
				return []*grpc.ClientConn{}, err
			}
		}

		d.dataNodeRangeCache[table] = append(
			d.dataNodeRangeCache[table], &krClientConn{
				KeyRange:   kr,
				ClientConn: client,
			})

		if keyRange.ContainsRange(kr) {
			rv = append(rv, client)
		}
	}

	return rv, nil
}

/*
Get requests the latest version of a single key of data from the specified
column path (table, row, column family, column).
*/
func (d *DataAccessClient) Get(
	parentCtx context.Context, req *redcloud.GetRequest,
	opts ...grpc.CallOption) (*redcloud.Column, error) {
	// The key range is just 1 key wide.
	var kr = common.NewKeyRange(req.Key, req.Key)
	var conns []*grpc.ClientConn
	var conn *grpc.ClientConn
	var ctx context.Context
	var span *trace.Span
	var err error

	ctx, span = trace.StartSpan(parentCtx, "red-cloud.DataAccessClient/Get")
	defer span.End()

	if conns, err = d.getRangeClients(ctx, req.Table, kr, false); err != nil {
		return nil, err
	} else if len(conns) > 1 {
		log.Printf("Error: multiple data nodes registered for key? %v", req)
	}

	for {
		for _, conn = range conns {
			var dnsc = redcloud.NewDataNodeServiceClient(conn)
			var col *redcloud.Column

			if col, err = dnsc.Get(
				ctx, req, opts...); err == common.ErrTabletNotLoaded {
				span.Annotate(nil, "Tablet not loaded")
				// Refresh data nodes covering the key and retry.
				if conns, err = d.getRangeClients(
					ctx, req.Table, kr, true); err != nil {
					return nil, err
				}
			} else if err != nil {
				span.AddAttributes(
					trace.StringAttribute("error", err.Error()))
				span.Annotate(nil, "Error communicating with data node")
				return nil, err
			} else {
				return col, nil
			}
		}

		if ctx.Err() != nil {
			span.AddAttributes(
				trace.StringAttribute("error", ctx.Err().Error()))
			span.Annotate(nil, "Context expired")
			return nil, ctx.Err()
		}
	}
}

/*
GetRange requests all versions of a key range of data from the specified
column paths (table, row, column family, columns).
*/
func (d *DataAccessClient) GetRange(
	parentCtx context.Context, req *redcloud.GetRangeRequest,
	resp chan *redcloud.ColumnSet, opts ...grpc.CallOption) error {
	var kr = common.NewKeyRange(req.StartKey, req.EndKey)
	var conns []*grpc.ClientConn
	var conn *grpc.ClientConn
	var ctx context.Context
	var span *trace.Span
	var err error

	ctx, span = trace.StartSpan(parentCtx, "red-cloud.DataAccessClient/GetRange")
	defer span.End()

	if conns, err = d.getRangeClients(ctx, req.Table, kr, false); err != nil {
		return err
	}

	for _, conn = range conns {
		var dnsc = redcloud.NewDataNodeServiceClient(conn)
		var rstream redcloud.DataNodeService_GetRangeClient

		if rstream, err = dnsc.GetRange(ctx, req, opts...); err != nil {
			span.AddAttributes(trace.StringAttribute("error", err.Error()))
			span.Annotate(nil, "Data node communication error")
			return err
		}

		for {
			var colset *redcloud.ColumnSet
			if colset, err = rstream.Recv(); err != nil {
				span.AddAttributes(
					trace.StringAttribute("error", err.Error()))
				span.Annotate(nil, "Data node result stream error")
				return err
			} else if colset == nil {
				break
			}
			resp <- colset
		}
	}

	span.Annotate(nil, "Result stream complete")
	return nil
}

/*
Insert requests to place a new version of a column into the database.
The destination of the column must be specified as a
(table, row, column family, column) tuple.
*/
func (d *DataAccessClient) Insert(
	parentCtx context.Context, req *redcloud.InsertRequest,
	opts ...grpc.CallOption) error {
	// The key range is just 1 key wide.
	var kr = common.NewKeyRange(req.Key, req.Key)
	var conns []*grpc.ClientConn
	var conn *grpc.ClientConn
	var ctx context.Context
	var span *trace.Span
	var err error

	ctx, span = trace.StartSpan(parentCtx, "red-cloud.DataAccessClient/Insert")
	defer span.End()

	if conns, err = d.getRangeClients(ctx, req.Table, kr, false); err != nil {
		return err
	}

	if len(conns) > 1 {
		span.AddAttributes(
			trace.Int64Attribute("num-range-covers", int64(len(conns))))
		span.Annotate(
			nil, "More than one range cover registered for a single key")
		return fmt.Errorf("Error: multiple data nodes registered for key? %v",
			req)
	}

	for {
		for _, conn = range conns {
			var dnsc = redcloud.NewDataNodeServiceClient(conn)

			if _, err = dnsc.Insert(
				ctx, req, opts...); err == common.ErrTabletNotLoaded {
				span.Annotate(nil, "Tablet not loaded")
				// Refresh data nodes covering the key and retry.
				if conns, err = d.getRangeClients(
					ctx, req.Table, kr, true); err != nil {
					return err
				}
			} else if err != nil {
				span.AddAttributes(
					trace.StringAttribute("error", err.Error()))
				span.Annotate(nil, "Data node communication error")
				return err
			} else {
				return nil
			}
		}

		// Check TTL / RPC cancelled.
		if ctx.Err() != nil {
			span.AddAttributes(
				trace.StringAttribute("error", ctx.Err().Error()))
			span.Annotate(nil, "Context expired")
			return ctx.Err()
		}
	}
}
