package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"context"
	"github.com/childoftheuniverse/etcd-discovery"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var numChangeReports = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "red_cloud",
	Subsystem: "caretaker_etcd_node_discovery",
	Name:      "num_change_reports",
	Help:      "Number of etcd changes reported",
}, []string{"change_type"})
var numChangeErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "red_cloud",
	Subsystem: "caretaker_etcd_node_discovery",
	Name:      "num_change_report_errors",
	Help:      "Number of errors watching etcd for changes",
}, []string{"error_type"})
var numEtcdHBCommErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "red_cloud",
	Subsystem: "caretaker_etcd_node_discovery",
	Name:      "num_etcd_hb_communication_errors",
	Help:      "Number of etcd communication errors checking for uncovered nodes",
})
var numHBProtoParsingErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "red_cloud",
	Subsystem: "caretaker_etcd_node_discovery",
	Name:      "num_hb_proto_parsing_errors",
	Help:      "Number of errors parsing discovery protocol buffers for heartbeat",
})
var numHBPeersDiscovered = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "red_cloud",
	Subsystem: "caretaker_etcd_node_discovery",
	Name:      "num_hb_peers_discovered",
	Help:      "Number of peers discovered via heartbeat rather than discovery",
})

func init() {
	prometheus.MustRegister(numChangeReports)
	prometheus.MustRegister(numChangeErrors)
	prometheus.MustRegister(numEtcdHBCommErrors)
	prometheus.MustRegister(numHBProtoParsingErrors)
	prometheus.MustRegister(numHBPeersDiscovered)
}

/*
EtcdExportedNodeDiscoveryStrategy uses etcd to discover new and disappearing
data nodes. To use this, data nodes must use etcd exported ports.
*/
type EtcdExportedNodeDiscoveryStrategy struct {
	path       string
	dnr        *DataNodeRegistry
	etcdClient *etcd.Client
}

/*
NewEtcdExportedNodeDiscoveryStrategy creates a new
EtcdExportedNodeDiscoveryStrategy discovering nodes using the given etcd
client and feeding them into the specified DataNodeRegistry.
*/
func NewEtcdExportedNodeDiscoveryStrategy(dnr *DataNodeRegistry,
	etcdClient *etcd.Client, instance string) *EtcdExportedNodeDiscoveryStrategy {
	var path = fmt.Sprintf("/ns/service/red-cloud/%s/data-node/", instance)
	var strategy = &EtcdExportedNodeDiscoveryStrategy{
		path:       path,
		dnr:        dnr,
		etcdClient: etcdClient,
	}
	discovery.MonitorExportedService(etcdClient, etcdClient, path, strategy)
	go strategy.monitorMissingNodes()
	return strategy
}

/*
monitorMissingNodes
*/
func (s *EtcdExportedNodeDiscoveryStrategy) monitorMissingNodes() {
	for {
		var resp *etcd.GetResponse
		var record *mvccpb.KeyValue
		var ctx context.Context
		var cancel context.CancelFunc
		var err error

		time.Sleep(5 * time.Minute)

		ctx, cancel = context.WithTimeout(context.Background(),
			20*time.Second)
		defer cancel()

		if resp, err = s.etcdClient.Get(
			ctx, s.path, etcd.WithPrefix()); err != nil {
			log.Print("Error communicating with etcd on ", s.path, ": ", err)
			numEtcdHBCommErrors.Inc()
			continue
		}

		for _, record = range resp.Kvs {
			var hostPort discovery.ExportedServiceRecord
			var peerAddr string

			if err = proto.Unmarshal(record.Value, &hostPort); err != nil {
				log.Print("Error parsing etcd discovery record for host ",
					string(record.Key), ": ", err)
				numHBProtoParsingErrors.Inc()
				break
			}

			peerAddr = net.JoinHostPort(
				hostPort.Address, strconv.Itoa(int(hostPort.Port)))

			if s.dnr.GetNodeByAddress(ctx, peerAddr) == nil {
				numHBPeersDiscovered.Inc()
				log.Print("Adding node ", string(record.Key), " (", peerAddr, ")")
				if err = s.dnr.Add(ctx, hostPort.Protocol, peerAddr); err != nil {
					log.Print("Error adding node ", peerAddr, ": ", err)
				}
			}
		}
	}
}

/*
ReportChange indicates that a new data node has appeared or an old one has been
updated or disappeared.
*/
func (s *EtcdExportedNodeDiscoveryStrategy) ReportChange(
	notification *discovery.ExportedServiceUpdateNotification) {
	var err error

	if notification.Update == discovery.ExportedServiceUpdateNotification_NEW {
		var dest = net.JoinHostPort(
			notification.UpdatedData.Address,
			strconv.Itoa(int(notification.UpdatedData.Port)))
		var ctx context.Context
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(),
			20*time.Second)
		defer cancel()
		log.Print("Adding node ", notification.Path, " (", dest, ")")
		err = s.dnr.Add(ctx, notification.UpdatedData.Protocol, dest)
		if err != nil {
			log.Print("Error adding node ", dest, ": ", err)
		}
		numChangeReports.With(
			prometheus.Labels{"change_type": "NEW"}).Inc()
	} else if notification.Update == discovery.ExportedServiceUpdateNotification_DELETED {
		var dest = net.JoinHostPort(
			notification.UpdatedData.Address,
			strconv.Itoa(int(notification.UpdatedData.Port)))
		var ctx context.Context
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(),
			20*time.Second)
		defer cancel()

		log.Print("Removing node ", notification.Path, " (", dest, ")")
		if err = s.dnr.Remove(ctx, dest); err != nil {
			log.Print("Error removing node ", dest, ": ", err)
		}
		numChangeReports.With(
			prometheus.Labels{"change_type": "DELETED"}).Inc()
	} else {
		numChangeReports.With(
			prometheus.Labels{"change_type": "UNKNOWN"}).Inc()
	}
}

/*
ReportError logs non-fatal errors to the console.
*/
func (s *EtcdExportedNodeDiscoveryStrategy) ReportError(err error) {
	numChangeErrors.With(prometheus.Labels{"error_type": "error"}).Inc()
	log.Print("Error discovering node changes from etcd: ", err)
}

/*
ReportFatal logs fatal errors to the console and ends the program with an error
exit code.
*/
func (s *EtcdExportedNodeDiscoveryStrategy) ReportFatal(err error) {
	numChangeErrors.With(prometheus.Labels{"error_type": "fatal"}).Inc()
	log.Fatal("Fatal error from etcd discovery: ", err)
}

/*
ReportCancelled reports that the watcher object has ceased to be active.
*/
func (s *EtcdExportedNodeDiscoveryStrategy) ReportCancelled() {
	numChangeErrors.With(prometheus.Labels{"error_type": "cancelled"}).Inc()
	/* TODO: probably do something here? */
}
