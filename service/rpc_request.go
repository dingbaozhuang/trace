package service

import (
	"fmt"

	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/models"
)

func (c *Client) RPCSend() {
	req := &models.RPCRequest{
		IPAddr: "192.168.0.1",
		Msg:    "log~~~~~~content",
	}
	resp := &models.RPCResponse{}
	err := c.RPCClient.Call("Postback.PostbackMessage", req, resp)
	if err != nil {
		log.Entry.Error().Err(err).Msg("client request rpc is failed.")
	}
	fmt.Println("-------resp=", *resp)
}
