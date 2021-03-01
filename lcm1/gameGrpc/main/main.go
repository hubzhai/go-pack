package main

import (
	"lcm1/gameGrpc/gameshz"
	"lcm1/gameGrpc/grpc"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	lRPC := new(grpcLcm.LRPC)
	rpc.Register(lRPC)
	GSRPC := &gamegrpc.GSRPC{
		Data:    make(map[int32]*gamegrpc.GamePool),
		ShzCard: make(map[int32][]int64),
	}
	rpc.Register(GSRPC)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:8000")

	if e != nil {
		log.Println("listen error:", e)
	}
	http.Serve(l, nil)
}
