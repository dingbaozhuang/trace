package service

import (
	"fmt"

	"github.com/yumimobi/trace/models"
)

func (p *Postback) PostbackMessage(req *models.RPCRequest, resp *models.RPCResponse) error {

	fmt.Println("------req msg=", req.IPAddr, req.Msg)
	resp.Code = 2000
	return nil
}
