package main

import (
	"fmt"

	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/service"
	"github.com/yumimobi/trace/service/grpc"
	"github.com/yumimobi/trace/util/json"
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
	go c.StartRPC()
	go c.Timer()

	grpc.NewGRPCClien()
	//-----------
	req := &grpc.Request{
		ID:               "1234567890",
		SspID:            "29",
		SlotID:           "qzntzwv",
		AppID:            "c989d0lc",
		AdType:           "2",
		SspAppIdKey:      "3-savsavd",
		SspAppPlaceIdKey: "3-place",
		Timestamp:        "2019061809",
		Type:             "grep",
	}
	resp, err := grpc.SendMsgToServer(req)
	if err != nil {
		fmt.Println("----grpc send msg is failed, err:", err)
	}
	data, _ := json.Marshal(resp)
	fmt.Println("grpc resp data is :", string(data))
	//-----------
	service.Shutdown(c.Server, c.Status)
}
