package grpc

import (
	context "context"
	fmt "fmt"
	"io"
	"os"

	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/models"
	grpc "google.golang.org/grpc"
)

var Clients []TraceClient

func NewGRPCClien() {
	conf := config.Conf

	if len(conf.Server.GRPC.Address) != len(conf.Server.GRPC.Port) {
		log.Entry.Error().Msg("trace server grpc addr is not match port.")
		os.Exit(0)
	}

	Clients = make([]TraceClient, len(conf.Server.GRPC.Address))

	for index, _ := range conf.Server.GRPC.Address {
		conn, err := grpc.Dial(conf.Server.GRPC.Address[index]+":"+conf.Server.GRPC.Port[index], grpc.WithInsecure())
		if err != nil {
			log.Entry.Error().Err(err).Msg("grpc dial is failed.")
		}
		// defer conn.Close()
		Clients[index] = NewTraceClient(conn)
	}
}

func sendMsg2All(req *Request, client TraceClient, resp chan *Response) {
	response, err := client.TransportLog(context.Background(), req)
	if err != nil {
		r := &Response{Code: models.StatusGRPCResponseFiled}
		resp <- r
	}

	resp <- response
}

func SendStreamMsg2All(req *Request, client TraceClient, resp chan *Response) {
	stream, err := client.ListTransportLog(context.Background(), req)
	if err != nil {
		fmt.Println("----send grpc stream is failed, err:", err)
		r := &Response{Code: models.StatusGRPCResponseFiled}
		resp <- r
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("SendStreamMsg2All  is EOF")
			break
		}
		if err != nil {
			fmt.Println("SendStreamMsg2All  recv is failed,err: ", err)
			r := &Response{Code: models.StatusGRPCResponseFiled}
			resp <- r
		}

		resp <- feature
		fmt.Println("--SendStreamMsg2All---", feature)
	}
}
