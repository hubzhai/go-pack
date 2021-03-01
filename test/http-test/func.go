package main

import (
	"games.agamestudio.com/http/etcd"
	"github.com/idealeak/goserver.v3/core"
	"github.com/idealeak/goserver.v3/core/basic"
	"github.com/idealeak/goserver.v3/core/logger"
	"github.com/idealeak/goserver.v3/core/module"
	"go.etcd.io/etcd/clientv3"
	"time"
)

///////////////////////////////config//////////////
var Config = Configuration{}

type Configuration struct {
	EtcdUrl  []string
	EtcdUser string
	EtcdPwd  string
}

func (this *Configuration) Name() string {
	return "Config"
}

func (this *Configuration) Init() error {
	return nil
}

func (this *Configuration) Close() error {
	return nil
}

///////////////////////////config end//////////////////
//////////////////////////etcd/////////////////////////////////
type EtcdMgr struct {
	*etcd.EtcdClient
	closed bool
}

func (this *EtcdMgr) goWatch(prefix string, f func(res clientv3.WatchResponse) error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Errorf("etcd watch WithPrefix(%v) panic:%v", prefix, err)
			}
			logger.Logger.Warnf("etcd watch WithPrefix(%v) quit!!!", prefix)
		}()
		var times int64
		for !this.closed {
			times++
			logger.Logger.Warnf("etcd watch WithPrefix(%v) start[%v]!!!", prefix, times)
			rch := this.WatchWithPrefix(prefix)
			for wresp := range rch {
				//if !model.GameParamData.UseEtcd {
				//	continue
				//}
				if wresp.Canceled {
					logger.Logger.Warnf("etcd watch WithPrefix(%v) be closed, reason:%v", prefix, wresp.Err())
					continue
				}
				obj := core.CoreObject()
				if obj != nil {
					func(res clientv3.WatchResponse) {
						obj.SendCommand(basic.CommandWrapper(func(*basic.Object) error {
							return f(res)
						}), true)
					}(wresp)
				}
			}
		}
	}()
}
func (this *EtcdMgr) Init() {
	logger.Logger.Infof("EtcdClient开始连接url:%v;etcduser:%v;etcdpwd:%v", Config.EtcdUrl, Config.EtcdUser, Config.EtcdPwd)
	err := this.Open(Config.EtcdUrl, Config.EtcdUser, Config.EtcdPwd, time.Second*5)
	if err != nil {
		logger.Logger.Tracef("EtcdMgr.Open err:%v", err)
	} else {
		this.check()
	}
	//if _, ok := common.DelayInvake(func() {
	//	this.check()
	//}, nil, 5*time.Second, 1); ok {
	//}
}
func (this *EtcdMgr) ModuleName() string {
	return "EtcdMgr"
}
func (this *EtcdMgr) Update() {

}
func (this *EtcdMgr) Shutdown() {
	this.closed = true
	this.Close()
}

//////////////////////////etcd  end/////////////////////////////////
//
func init() {
	//config
	core.RegistePackage(&Config)
	module.RegisteModule(EtcdMgrSington, time.Second, 1)
}
