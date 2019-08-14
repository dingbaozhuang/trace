package service

import (
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/yumimobi/trace/config"
)

type Client struct {
	Server    *http.Server
	RPCClient *rpc.Client
	Status    chan int
}

func NewClient() *Client {
	// c := config.Conf
	c := &Client{
		Status: make(chan int, 1),
	}
	return c
}

func (c *Client) StartHTTP() error {
	conf := config.Conf

	gin.SetMode(gin.DebugMode)
	r := gin.New()

	pprof.Register(r)

	c.Router(r)

	c.Server = &http.Server{
		Addr:           conf.Client.HTTP.Address + ":" + conf.Client.HTTP.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go c.Server.ListenAndServe()

	return nil
}

var timer *time.Timer

func (c *Client) Timer() {
	timer = time.NewTimer(time.Second * 10 * 60)

	for {
		select {
		case <-timer.C:
			os.Exit(0)
		}
	}
}

func ResetTimer() {
	timer.Reset(time.Second * 10 * 60)
}
