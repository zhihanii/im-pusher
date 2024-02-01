package rpc

import clientv3 "go.etcd.io/etcd/client/v3"

func Init(etcdClient *clientv3.Client) {
	initLogicClient(etcdClient)
}
