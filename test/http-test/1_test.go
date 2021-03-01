package main

import (
	"encoding/json"
	"fmt"
	"games.agamestudio.com/http/etcd"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	etcdUrl := []string{"192.168.2.47:2379"}
	dialTimeout := time.Second * 3
	userName := ""
	passWord := ""
	key := "/www/test/"
	fmt.Println(key)
	e := new(etcd.EtcdClient)
	err := e.Open(etcdUrl, userName, passWord, dialTimeout)
	if err != nil {
		fmt.Println(err)
	}
	sa := A{"test001", 1}
	ss, _ := json.Marshal(sa)
	a, err := e.PutValue(key, string(ss))
	fmt.Println("put", a, err)
	res, err := e.GetValueWithPrefix(key)
	fmt.Println("get", res, err)
	for i := int64(0); i < res.Count; i++ {
		var cfg A
		str := res.Kvs[i].Value
		fmt.Println(str)
		fmt.Println(string(str))
		err = json.Unmarshal(res.Kvs[i].Value, &cfg)
		fmt.Println("===", cfg)
	}
	n := e.WatchWithPrefix(key)
	fmt.Println(n, len(n))
}

type A struct {
	A string
	B int
}
