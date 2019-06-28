package service

import (
	"net"
	"net/rpc"

	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/models"
)

type Postback struct {
}

func (s *Server) StartRPC() {
	postback := new(Postback)
	err := rpc.Register(postback)
	if err != nil {
		log.Entry.Error().Err(err).Msg("server rpc register is failed.")
	}

	conf := config.Conf
	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.Server.RPC.Address+":"+conf.Server.RPC.Port)
	if err != nil {
		log.Entry.Error().Err(err).Msg("server rpc resolve tcp addr is failed.")
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Entry.Error().Err(err).Msg("server rpc listen tcp addr is failed.")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Entry.Error().Err(err).Msg("server rpc accept tcp addr is failed.")
			continue
		}
		rpc.ServeConn(conn)
	}
}

func (c *Client) InitRPC() error {
	conf := config.Conf
	client, err := rpc.Dial("tcp", conf.Client.RPC.Address+":"+conf.Client.RPC.Port)
	if err != nil {
		log.Entry.Error().Err(err).Msg("client rpc accept tcp addr is failed.")
		return err
	}

	c.RPCClient = client
	return nil
}

func (c *Client) StartRPC() {
	err := c.InitRPC()
	if err != nil {
		// 发送channel不让client服务启动
		c.Status <- models.StatusClientRPCStartFailed
		return
	}
	c.RPCSend()
}
