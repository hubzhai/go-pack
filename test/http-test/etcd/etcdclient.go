package etcd

import (
	"context"
	"github.com/idealeak/goserver.v3/core/logger"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type EtcdClient struct {
	cli *clientv3.Client
}

func (this *EtcdClient) Open(etcdUrl []string, userName, passWord string, dialTimeout time.Duration) error {
	var err error
	this.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdUrl,
		DialTimeout: dialTimeout,
		Username:    userName,
		Password:    passWord,
	})

	if err != nil {
		logger.Logger.Warnf("EtcdClient.open(%v) err:%v", etcdUrl, err)
	}
	return err
}

func (this *EtcdClient) Close() error {
	logger.Logger.Warnf("EtcdClient.close(%v) err:%v")
	if this.cli != nil {
		return this.cli.Close()
	}
	return nil
}

//添加键值对
func (this *EtcdClient) PutValue(key, value string) (*clientv3.PutResponse, error) {
	resp, err := this.cli.Put(context.TODO(), key, value)
	if err != nil {
		logger.Logger.Warnf("EtcdClient.PutValue(%v,%v) error:%v", key, value, err)
	}
	return resp, err
}

//查询
func (this *EtcdClient) GetValue(key string) (*clientv3.GetResponse, error) {
	resp, err := this.cli.Get(context.TODO(), key)
	if err != nil {
		logger.Logger.Warnf("EtcdClient.GetValue(%v) error:%v", key, err)
	}
	return resp, err
}

// 返回删除了几条数据
func (this *EtcdClient) DelValue(key string) (*clientv3.DeleteResponse, error) {
	res, err := this.cli.Delete(context.TODO(), key)
	if err != nil {
		logger.Logger.Warnf("EtcdClient.DelValue(%v) error:%v", key, err)
	}
	return res, err
}

//按照前缀删除
func (this *EtcdClient) DelValueWithPrefix(prefix string) (*clientv3.DeleteResponse, error) {
	res, err := this.cli.Delete(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Logger.Warnf("EtcdClient.DelValueWithPrefix(%v) error:%v", prefix, err)
	}
	return res, err
}

//获取前缀
func (this *EtcdClient) GetValueWithPrefix(prefix string) (*clientv3.GetResponse, error) {
	resp, err := this.cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Logger.Warnf("EtcdClient.GetValueWIthPrefix(%v) error:%v", prefix, err)
	}
	return resp, err
}

func (this *EtcdClient) WatchWithPrefix(prefix string) clientv3.WatchChan {
	if this.cli != nil {
		return this.cli.Watch(clientv3.WithRequireLeader(context.Background()), prefix, clientv3.WithPrefix())
	}
	return nil
}
