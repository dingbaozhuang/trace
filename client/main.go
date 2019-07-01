package main

import (
	"fmt"

	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/service"
	"github.com/yumimobi/trace/service/grpc"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Println("init config is failed, err:", err)
		return
	}

	log.Init(&config.Conf.Client.Log)

	c := service.NewClient()
	c.StartHTTP()
	// go c.StartRPC()
	// go c.Timer()

	go grpc.GRPCServerInit()

	service.Shutdown(c.Server, c.Status)
}
