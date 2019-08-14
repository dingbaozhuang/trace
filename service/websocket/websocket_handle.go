package websocket

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/yumimobi/trace/util/json"

	"github.com/gorilla/websocket"
	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/log"
	"github.com/yumimobi/trace/service/grpc"
)

var upgrader = websocket.Upgrader{
	// 跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func WebSocketInit() {
	conf := config.Conf

	http.HandleFunc("/trace", WebSocketHandler)

	err := http.ListenAndServe(conf.Server.WebSocket.Address+":"+conf.Server.WebSocket.Port, nil)
	if err != nil {
		log.Entry.Error().Err(err).Msg("websocket listen and serve is failed")
	}
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Entry.Error().Err(err).Msg("websocket connect is failed")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}

		c.Close()
	}()

	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Entry.Error().Err(err).Int("msg type", mt).Msg("websocket read is failed")
		return
	}

	log.Entry.Debug().Str("req", string(message)).Msg("websocket request msg")

	// request := &grpc.Request{}
	// err = json.Unmarshal(message, request)
	// if err != nil {
	// 	return
	// }

	SendGRPCMsg(mt, c, message)

	return
}

func SendGRPCMsg(mt int, c *websocket.Conn, req []byte) {
	// 因为是一对多，只能用close channel
	stop := make(chan struct{})
	resps := make(chan *grpc.Response, 1000)

	request := &grpc.Request{}
	err := json.Unmarshal(req, request)
	if err != nil {
		return
	}

	defer func() {
		close(stop)
	}()

	for i, _ := range grpc.Clients {
		go grpc.SendStreamMsg2All(request, grpc.Clients[i], resps, stop)
	}

	var data []byte
	for {

		select {
		case msg, ok := <-resps:
			if !ok {
				fmt.Println("recv grpc msg is failed")
				return
			}

			data, err = json.Marshal(msg)
			if err != nil {
				fmt.Println("marshal recv grpc msg is failed,err:", err)
				return
			}

		case <-time.Tick(time.Second * 5):
			pingFunc := c.PingHandler()
			err = pingFunc("ping")
			if err != nil {
				// client 端关闭
				fmt.Println("websocket ping client is failed")
				return
			}

			fmt.Println("ping")
			continue
		}

		data = trimResponse(data)
		err = c.WriteMessage(mt, data)
		if err != nil {
			log.Entry.Error().Err(err).Int("msg type", mt).Msg("websocket write is failed")
			fmt.Println("websocket write is failed")
			return
		}
	}

	return
}

// 去掉转义字符和日志中的无效数据
func trimResponse(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte(`\x03`), []byte(`\n`))
	data = bytes.ReplaceAll(data, []byte(`\`), []byte(``))
	return data
}
