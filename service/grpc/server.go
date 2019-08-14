package grpc

import (
	context "context"
	"encoding/json"
	fmt "fmt"
	"net"
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

// Server-side streaming RPC
func (s *server) ListTransportLog(req *Request, stream Trace_ListTransportLogServer) error {
	fmt.Println("-----server ListTransportLog,", *req)
	ips, _ := util.GetLocalIP()
	if req == nil {
		resp := &Response{
			Code: models.StatusGRPCRequestIsNil,
			ID:   req.ID,
			Data: "",
			IP:   strings.Join(ips, ","),
		}

		if err := stream.Send(resp); err != nil {
			fmt.Println("stream send msg is faild,err:", err)
			return err
		}
		return nil
	}

	m := requestConvert2Map(req)
	StreamCommand(m, stream)

	return nil
}

func StreamCommand(m map[string]string, stream Trace_ListTransportLogServer) string {
	msg := make(chan script.Message, 100)

	cmd, tmp := script.GetCmd(m)
	if cmd == "" {
		return "Required parameter is missing."
	}

	ctx, cancle := context.WithCancel(context.Background())
	go script.ExecGrepCmd(ctx, cmd, tmp, msg)

	resp := &Response{}
	for {
		select {
		case message, ok := <-msg:
			if ok != true {
				break
			}
			data, _ := json.Marshal(message)
			resp.Data = string(data)
			if err := stream.Send(resp); err != nil {
				fmt.Println("stream send tail msg is failed, err:", err)
				return ""
			}
		case <-stream.Context().Done():
			cancle()
			fmt.Println("client close context", stream.Context().Err())
			return ""
		}
	}
	return ""
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
