package websocket

import (
	"bytes"
	"fmt"
	"net/http"

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
	}
	defer c.Close()

	// for {
	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Entry.Error().Err(err).Int("msg type", mt).Msg("websocket read is failed")
		return
		// break
	}

	log.Entry.Debug().Str("req", string(message)).Msg("websocket request msg")

	resp, err := convertMsgFormat(message)
	if err != nil {
		log.Entry.Error().Err(err).Str("req", string(message)).Msg("websocket convert msg is failed")
		return
		// break
	}

	log.Entry.Debug().Str("resp", string(resp)).Msg("websocket response msg")

	err = c.WriteMessage(mt, resp)
	if err != nil {
		log.Entry.Error().Err(err).Int("msg type", mt).Msg("websocket write is failed")
		return
		// break
	}
	// }
	return
}

func convertMsgFormat(req []byte) ([]byte, error) {
	request := &grpc.Request{}

	fmt.Println("-------", string(req))
	err := json.Unmarshal(req, request)
	if err != nil {
		return nil, err
	}

	response, err := grpc.SendMsg(request)
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	// fmt.Println("--------", string(resp))

	resp = bytes.ReplaceAll(resp, []byte(`\x03`), []byte(`\n`))
	resp = bytes.ReplaceAll(resp, []byte(`\`), []byte(``))
	return resp, nil
}
