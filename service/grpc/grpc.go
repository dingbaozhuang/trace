package grpc

import (
	"context"
	"errors"
	fmt "fmt"
	"net"

	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/service/script"

	"github.com/yumimobi/trace/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var c TraceClient

func GRPCServerInit() {
	conf := config.Conf
	lis, err := net.Listen("tcp", conf.Server.GRPC.Address+":"+conf.Server.GRPC.Port)
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
	conn, err := grpc.Dial(conf.Client.GRPC.Address+":"+conf.Client.GRPC.Port, grpc.WithInsecure())
	if err != nil {
		log.Entry.Error().Err(err).Msg("grpc dial is failed.")
	}
	// defer conn.Close()
	c = NewTraceClient(conn)
}

func (s *server) TransportLog(ctx context.Context, req *Request) (*Response, error) {
	if req == nil {
		return nil, errors.New("req is nil.")
	}

	m := requestConvert2Map(req)
	msg := script.Command(m)

	fmt.Println("----grpc--req msg=", req.IPAddr, req.Msg)
	resp := &Response{
		Code: 0,
		ID:   req.ID,
		Data: msg,
	}
	return resp, nil
}

func SendMsg(req *Request) (*Response, error) {

	fmt.Println("xxxxxxgrpcxxxxx******", req.String())
	resp, err := c.TransportLog(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func requestConvert2Map(req *Request) map[string]string {
	m := make(map[string]string)
	m["ID"] = req.ID
	m["SspID"] = req.SspID
	m["SlotID"] = req.SlotID
	m["AppID"] = req.AppID
	m["AdType"] = fmt.Sprint(req.AdType)
	m["SspAppIdKey"] = req.SspAppIdKey
	m["SspAppPlaceIdKey"] = req.SspAppPlaceIdKey
	m["SspAppSecretKey"] = req.SspAppSecretKey
	m["Timestamp"] = req.Timestamp
	m["Type"] = req.Type

	return m
}
