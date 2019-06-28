package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yumimobi/trace/log"
)

var upgrader = websocket.Upgrader{
	// 跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func WebSocketInit() {
	http.HandleFunc("/trace", WebSocketHandler)
	err := http.ListenAndServe(":8000", nil)
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

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Entry.Error().Err(err).Int("msg type", mt).Msg("websocket read is failed")
			break
		}
		fmt.Println("websocket read msg is:", string(message))

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Entry.Error().Err(err).Int("msg type", mt).Msg("websocket write is failed")
			break
		}
	}
}
