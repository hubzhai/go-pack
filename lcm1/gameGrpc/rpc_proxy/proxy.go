package utils

import (
	"errors"
	"net/rpc"
	"sync"
	"time"
)

type RpcProxy struct {
	net     string
	address string
	client  *rpc.Client  //rpc客户端
	mutex   sync.RWMutex //重连控制锁
	iskeep  bool
}

var _instance *RpcProxy

func (self *RpcProxy) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if !self.iskeep {
		return errors.New("rpc proxy not started:" + serviceMethod)
	}

	self.mutex.RLock()
	err := self.client.Call(serviceMethod, args, reply)
	if err != nil && err.Error() == "connection is shut down" {
		self.mutex.RUnlock()
		self.mutex.Lock()
		self.client.Close()
		ch := make(chan bool)
		go func(ch chan bool) {
			for self.iskeep {
				client, err := rpc.DialHTTP(self.net, self.address)
				if err == nil {
					self.client = client
					self.mutex.Unlock()
					ch <- true
					close(ch)
					break
				}
			}
		}(ch)

		//连接成功之后，重新执行。同时支持超时
		select {
		case <-time.After(time.Second * 3):

			return err
		case <-ch:
			self.mutex.RLock()
			err = self.client.Call(serviceMethod, args, reply)
			self.mutex.RUnlock()
			break
		}
	} else {
		self.mutex.RUnlock()
	}

	return err
}

func (self *RpcProxy) Start() error {
	client, err := rpc.DialHTTP(self.net, self.address)
	if err == nil {
		self.client = client
		self.iskeep = true
	} else {
		return err
	}

	return nil
}

func NewRpcProxy(network, address string) *RpcProxy {
	if _instance == nil {
		_instance = &RpcProxy{
			net:     network,
			address: address,
			iskeep:  false,
		}
	}

	return _instance
}

func GetRpc() *RpcProxy {
	if _instance != nil {
		return _instance
	}
	return nil
}
