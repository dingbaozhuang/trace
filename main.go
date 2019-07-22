package main

import (
	"fmt"

	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/service"
	"github.com/yumimobi/trace/service/grpc"
	"github.com/yumimobi/trace/service/websocket"
)

/*
	服务端重启，客户端自动重连
*/
func main() {
	err := config.Init()
	if err != nil {
		fmt.Println("init config is failed, err:", err)
		return
	}

	log.Init(&config.Conf.Server.Log)

	s := service.NewServer()
	s.StartHTTP()
	go websocket.WebSocketInit()
	// go s.StartRPC()
	grpc.NewGRPCClien()

	service.Shutdown(s.Server, s.Status)
}
