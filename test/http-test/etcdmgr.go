package main

import (
	"encoding/json"
	"fmt"
	"games.agamestudio.com/http/etcd"
	"github.com/idealeak/goserver.v3/core/logger"
	"go.etcd.io/etcd/clientv3"
)

var EtcdMgrSington = &EtcdMgr{EtcdClient: &etcd.EtcdClient{}}

func (this *EtcdMgr) check() {
	this.InitTest()
}

func (this *EtcdMgr) InitTest() {
	logger.Logger.Info("ETCD 拉取数据:", etcd.ETCDKEY_PLATFORM_PREFIX)
	res, err := this.GetValueWithPrefix(etcd.ETCDKEY_PLATFORM_PREFIX)
	if err == nil {
		for i := int64(0); i < res.Count; i++ {
			type A struct {
				A string
				B int
			}
			var a A
			err := json.Unmarshal(res.Kvs[i].Value, &a)
			if err == nil {
				fmt.Println("=get==", a)
			} else {
				logger.Logger.Errorf("etcd desc WithPrefix(%v) panic:%v", etcd.ETCDKEY_PLATFORM_PREFIX, err)
			}
		}
	} else {
		logger.Logger.Errorf("etcd get WithPrefix(%v) panic:%v", etcd.ETCDKEY_PLATFORM_PREFIX, err)
	}

	// 监控数据变动
	this.goWatch(etcd.ETCDKEY_PLATFORM_PREFIX, func(res clientv3.WatchResponse) error {
		for _, ev := range res.Events {
			switch ev.Type {
			case clientv3.EventTypeDelete:
			case clientv3.EventTypePut:
				type A struct {
					A string
					B int
				}
				var a A
				err := json.Unmarshal(ev.Kv.Value, &a)
				if err == nil {
					fmt.Println("==update=", a)
				} else {
					logger.Logger.Errorf("etcd desc WithPrefix(%v) panic:%v", etcd.ETCDKEY_PLATFORM_PREFIX, err)
				}
			}
		}
		return nil
	})
}
