package script

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yumimobi/trace/util/json"

	"github.com/yumimobi/trace/config"
)

type Message struct {
	Msg string `json:"msg"`
	IP  string `json:"ip"`
	Err string `json:"err"`
}

func Command(m map[string]string) string {
	msg := make(chan *Message, 100)
	msgs := make([]*Message, 0)

	cmd, tmp := getCmd(m)
	execGrepCmd(cmd, tmp, msg)

	for {
		select {
		case message, ok := <-msg:
			if ok != true {
				data, _ := json.Marshal(msgs)
				return string(data)
			}
			msgs = append(msgs, message)
		}
	}
	return ""
}

func getCmd(m map[string]string) (string, string) {

	id := m["ID"]
	sspId := m["SspID"]
	slotId := m["SlotID"]
	appId := m["AppID"]
	adType := m["AdType"]
	sspAppIdKey := m["SspAppIdKey"]
	sspAppPlaceIdKey := m["SspAppPlaceIdKey"]
	sspAppSecretKey := m["SspAppSecretKey"]
	timestamp := m["Timestamp"]

	if 12-len(timestamp) > 0 {
		timestamp = timestamp + strings.Repeat("[0-9]", 12-len(timestamp))
	}
	file := config.Conf.Client.Target.Dir + "api.log." + timestamp

	grep := ""
	if sspId != "" {
		grep = grep + `| grep -a "\"sspId\":` + sspId + `"`
	}
	if slotId != "" {
		grep = grep + `| grep -a "` + slotId + `"`
	}
	if appId != "" {
		grep = grep + `| grep -a "` + appId + `"`
	}
	if adType != "" {
		grep = grep + `| grep -a "\"adType\":` + adType + `"`
	}
	if sspAppIdKey != "" {
		grep = grep + `| grep -a "` + sspAppIdKey + `"`
	}
	if sspAppPlaceIdKey != "" {
		grep = grep + `| grep -a "` + sspAppPlaceIdKey + `"`
	}
	if sspAppSecretKey != "" {
		grep = grep + `| grep -a "` + sspAppSecretKey + `"`
	}

	prefixCmd := ""
	switch m["Type"] {
	case "cat":
		prefixCmd = "cat "
	case "tail":
		prefixCmd = "tail -f "
	}
	grep = prefixCmd + file + grep + " >" + config.Conf.Client.Target.Dir + id + ".tmp"

	return grep, config.Conf.Client.Target.Dir + id + ".tmp"
}

func execGrepCmd(cmd string, tmp string, msg chan *Message) {
	fmt.Println("-----cmd=", cmd)

	// message := &Message{}
	// ips, err := util.GetLocalIP()
	// if err != nil {
	// 	message.Err = err.Error()
	// 	msg <- message
	// }

	// message.IP = strings.Join(ips, ",")
	// IsCmd := exec.Command("bash", "-c", cmd)
	// err = IsCmd.Run()
	// if err != nil {
	// 	log.Println("exec bash shell is failed, err: ", err)
	// 	message.Err = err.Error()
	// 	msg <- message
	// 	return
	// }

	// f, err := os.Open(tmp)
	// if err != nil {
	// 	message.Err = err.Error()
	// 	msg <- message
	// 	return
	// }

	// ReadTmpFile(f, msg, message)
	// close(msg)

	removeTmpFile(tmp)
	return
}

func ReadTmpFile(f *os.File, msg chan *Message, message *Message) {
	r := bufio.NewReader(f)
	for {
		str, err := r.ReadString('\n')
		// 走到此处，如果是EOF，则在读取到EOF之前会把有效的数据读入到str中
		if err == io.EOF {
			//此时str中是有数据的需要处理
			message.Msg = str
			msg <- message
			break
		}
		if err != nil {
			message.Err = err.Error()
			msg <- message
			//应该返回错误,先判断是否为nil，不为在判断是否等于EOF,等于break出for，然后继续执行；其他错误直接return
			// fmt.Println("read string is failed, err: ", err)
		}
		//str中会包含'\n'
		message.Msg = str
		msg <- message
		// fmt.Println("len(str)=", len(str), "str=", str)
	}
	return
}

func removeTmpFile(file string) {
	os.Remove(file)
}
