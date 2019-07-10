package grpc

import (
	"context"
	fmt "fmt"
	"net"
	"os"
	"strings"

	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/models"
	"github.com/yumimobi/trace/service/script"
	"github.com/yumimobi/trace/util"

	"github.com/yumimobi/trace/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var clients []TraceClient

func GRPCServerInit() {
	conf := config.Conf
	lis, err := net.Listen("tcp", conf.Client.GRPC.Address+":"+conf.Client.GRPC.Port)
	if err != nil {
		log.Entry.Error().Err(err).Msg("server grpc listen is failed.")
	}

	s := grpc.NewServer()
	RegisterTraceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Entry.Error().Err(err).Msg("failed to serve.")
	}
}

func NewGRPCClien() {
	conf := config.Conf

	if len(conf.Server.GRPC.Address) != len(conf.Server.GRPC.Port) {
		log.Entry.Error().Msg("trace server grpc addr is not match port.")
		os.Exit(0)
	}

	clients = make([]TraceClient, len(conf.Server.GRPC.Address))

	for index, _ := range conf.Server.GRPC.Address {
		conn, err := grpc.Dial(conf.Server.GRPC.Address[index]+":"+conf.Server.GRPC.Port[index], grpc.WithInsecure())
		if err != nil {
			log.Entry.Error().Err(err).Msg("grpc dial is failed.")
		}
		// defer conn.Close()
		clients[index] = NewTraceClient(conn)
	}
}

// 简单rpc
func (s *server) TransportLog(ctx context.Context, req *Request) (*Response, error) {
	ips, _ := util.GetLocalIP()

	if req == nil {
		resp := &Response{
			Code: models.StatusGRPCRequestIsNil,
			ID:   req.ID,
			Data: "",
			IP:   strings.Join(ips, ","),
		}
		return resp, nil
	}

	m := requestConvert2Map(req)
	msg := script.Command(m)

	resp := &Response{
		Code: 0,
		ID:   req.ID,
		Data: msg,
		IP:   strings.Join(ips, ","),
	}
	return resp, nil
}

// // Server-side streaming RPC
// func (s *server) ListTransportLog(req *Request, stream Trace_ListTransportLogServer) error {
// 	ips, _ := util.GetLocalIP()
// 	if req == nil {
// 		resp := &Response{
// 			Code: models.StatusGRPCRequestIsNil,
// 			ID:   req.ID,
// 			Data: "",
// 			IP:   strings.Join(ips, ","),
// 		}

// 		if err := stream.Send(resp); err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// }

func SendMsg(req *Request) ([]*Response, error) {

	fmt.Println("xxxxxxgrpcxxxxx******", req.String())
	conf := config.Conf

	serverNum := len(conf.Server.GRPC.Address)
	resps := make(chan *Response, serverNum)
	responses := make([]*Response, 0)
	for i, _ := range clients {
		go sendMsg2All(req, clients[i], resps)
	}

	count := 0
	for {
		if count >= serverNum {
			break
		}

		select {
		case msg, _ := <-resps:
			responses = append(responses, msg)
		}
		count++
	}

	return responses, nil
}

func requestConvert2Map(req *Request) map[string]string {
	m := make(map[string]string)
	if req.SID != "" {
		m["SID"] = req.SID
	} else {
		m["ID"] = req.ID
		m["SspID"] = req.SspID
		m["SlotID"] = req.SlotID
		m["AppID"] = req.AppID
		m["AdType"] = fmt.Sprint(req.AdType)
		m["SspAppIdKey"] = req.SspAppIdKey
		m["SspAppPlaceIdKey"] = req.SspAppPlaceIdKey
		m["SspAppSecretKey"] = req.SspAppSecretKey
		m["Timestamp"] = req.Timestamp
	}

	m["Type"] = req.Type
	m["UUID"] = req.Uuid

	return m
}

func sendMsg2All(req *Request, client TraceClient, resp chan *Response) {
	response, err := client.TransportLog(context.Background(), req)
	if err != nil {
		r := &Response{Code: models.StatusGRPCResponseFiled}
		resp <- r
	}

	resp <- response
}
