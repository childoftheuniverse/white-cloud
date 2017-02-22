package client

import (
	discovery "github.com/childoftheuniverse/etcd-discovery"
	"github.com/childoftheuniverse/red-cloud/common"
	etcd "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
GetMasterConnection finds the master for the red-cloud instance associated
with the giveh path and returns it.
*/
func GetMasterConnection(
	ctx context.Context, etcdClient *etcd.Client, path string) (
	*grpc.ClientConn, error) {
	var instance string
	var err error

	if instance, _, err = SplitTablePath(path); err != nil {
		return nil, err
	}

	// TODO: use actual credentials.
	return discovery.NewGrpcClient(
		ctx, etcdClient, common.EtcdMasterPrefix(instance), grpc.WithInsecure())
}
