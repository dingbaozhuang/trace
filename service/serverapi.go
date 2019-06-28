package service

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/service/grpc"
	"github.com/yumimobi/trace/util"
	"github.com/yumimobi/trace/util/json"
)

type Message struct {
	Msg string `json:"msg"`
	IP  string `json:"ip"`
	Err string `json:"err"`
}

func serverAPIHTTPHandler(c *gin.Context) {
	cmd := getCmd(c)
	addrs := getAddrs(c)

	msg := make(chan *Message, 100)
	for _, ip := range addrs {
		ssh := `ssh ` + ip + ` '` + cmd + `' `
		go execCmd(ssh, ip, msg)
	}

}

// func ServerAPIRPCHandler(cmd string) ([]byte, error) {
// 	// req, err := parseParameter(cmd)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// resp, err := grpc.SendMsg(req)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// data, err := json.Marshal(resp)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	data,err:=requestMsg(cmd)

// 	return data, nil
// }

func RequestMsg(cmd string) ([]byte, error) {
	req, err := parseParameter(cmd)
	if err != nil {
		return nil, err
	}

	// 此处应是send to client
	resp, err := grpc.SendMsgToServer(req)
	// resp, err := grpc.SendMsg(req)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func ServerAPIRPCHandler() {

}

func parseParameter(cmd string) (*grpc.Request, error) {
	req := &grpc.Request{}
	err := json.Unmarshal([]byte(cmd), req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func getAddrs(c *gin.Context) []string {
	role := c.Param("role")
	addrs := []string{}

	switch role {
	case "test":
		addrs = config.Conf.Server.RemoteAddress.TestAddr
	default:
		addrs = config.Conf.Server.RemoteAddress.ProductionAddr
	}

	return addrs
}

func execCmd(cmd string, ip string, msg chan *Message) {
	fmt.Println("-----cmd=", cmd)

	message := &Message{}
	ips, err := util.GetLocalIP()
	if err != nil {
		message.Err = err.Error()
		msg <- message
	}

	message.IP = strings.Join(ips, ",")
	IsCmd := exec.Command("bash", "-c", cmd)
	IsOut, err := IsCmd.Output()
	if err != nil {
		fmt.Println("exec bash shell is failed, err: ", err)
		message.Err = err.Error()
		msg <- message
		return
	}

	fmt.Println("------msg=", string(IsOut))
	message.Msg = string(IsOut)
	msg <- message
	return
}

func getCmd(c *gin.Context) string {
	date := c.Param("date")

	grep := getPrimaryParameter(c)
	parameter := c.Request.URL.Query()
	for key, value := range parameter {
		if key == "sspId" {
			continue
		}

		if grep != "" {
			grep += ` | `
		}
		grep += `grep -a ` + `"\"` + key + `\":\"` + getValue(value) + `\""`
	}

	if 12-len(date) > 0 {
		date = date + strings.Repeat("[0-9]", 12-len(date))
	}
	file := `/data/logs/server_api/api.log.` + date
	grep = `cat ` + file + ` | ` + grep
	return grep
}

func getPrimaryParameter(c *gin.Context) string {
	grep := ""
	sspId := c.Query("sspId")
	if sspId != "" {
		grep = `grep -a ` + `"\"sspId\":` + sspId + `"`
	}

	return grep
}

func getValue(strs []string) string {
	if len(strs) > 0 {
		return strs[0]
	}
	return ""
}
