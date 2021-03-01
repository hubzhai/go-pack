package main

import (
	"lcm1/gameGrpc/gameshz"
	"lcm1/gameGrpc/rpc_proxy"
	"log"
	"os"
	"os/signal"
	"time"
)

var endtime int64

func main() {
	//rpc := utils.NewRpcProxy("tcp", "192.168.4.166:7777")
	rpc := utils.NewRpcProxy("tcp", "127.0.0.1:8000")

	if err := rpc.Start(); err != nil {
		log.Println("Start", err)
		return
	}
	log.Println("Start ", time.Now().UnixNano())
	for i := 0; i != 1; i++ {
		go testRPC(i, rpc)
	}

	test := make(chan os.Signal, 1)
	signal.Notify(test, os.Interrupt, os.Kill)
	sig := <-test
	log.Printf("webTest closing down (signal: %v)", sig)
}

func testRPC(i int, rpc *utils.RpcProxy) {
	gargs := &gamegrpc.UserData{
		GamefreeID: 27,
		Bet:        100,
		ChangeCoin: 1200,
	}

	greplyint := new(gamegrpc.WinData)
	if err := rpc.Call("GSRPC.CalcGive", gargs, greplyint); err != nil {
		log.Println("err:", err)
	}

	s := "localhost:27017"
	bi := new(bool)
	if err := rpc.Call("GSRPC.InitDB", &s, bi); err != nil {
		log.Println("InitDB err:", err)
	}

	rate := int32(2)
	card := int64(0)
	if err := rpc.Call("GSRPC.RetCard", &rate, &card); err != nil {
		log.Println("RetCard err:", err)
	}
	log.Println("RetCard :", card)
	// log.Println("End  ", time.Now().UnixNano())
}
