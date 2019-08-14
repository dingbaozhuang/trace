package script

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/yumimobi/trace/config"
	"github.com/yumimobi/trace/util"
	"github.com/yumimobi/trace/util/json"
)

type Message struct {
	Msg string `json:"msg"`
	IP  string `json:"ip"`
	Err string `json:"err"`
}

type FileInfo struct {
	Name    string
	ModTime time.Time
}

type fileInfo []*FileInfo

func (a fileInfo) Len() int           { return len(a) }
func (a fileInfo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfo) Less(i, j int) bool { return a[i].ModTime.Before(a[j].ModTime) }

func Command(m map[string]string) string {
	msg := make(chan Message, 100)
	msgs := make([]Message, 0)

	cmd, tmp := GetCmd(m)
	if cmd == "" {
		return "Required parameter is missing."
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer func() {
		cancle()
	}()

	go ExecGrepCmd(ctx, cmd, tmp, msg)

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

// func StreamCommand(m map[string]string, stream chan string) string {
// 	msg := make(chan Message, 100)

// 	cmd, tmp := getCmd(m)
// 	if cmd == "" {
// 		return "Required parameter is missing."
// 	}

// 	go ExecGrepCmd(cmd, tmp, msg)

// 	for {
// 		select {
// 		case message, ok := <-msg:
// 			if ok != true {
// 				break
// 			}
// 			data, _ := json.Marshal(message)
// 			stream <- string(data)
// 		}
// 	}
// 	return ""
// }

func GetCmd(m map[string]string) (string, string) {

	id := m["UUID"]
	sspId := m["SspID"]
	slotId := m["SlotID"]
	appId := m["AppID"]
	adType := m["AdType"]
	sspAppIdKey := m["SspAppIdKey"]
	sspAppPlaceIdKey := m["SspAppPlaceIdKey"]
	sspAppSecretKey := m["SspAppSecretKey"]
	timestamp := m["Timestamp"]
	sId := m["SID"]

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
	if sId != "" {
		grep = grep + `| grep -a "` + sId + `"`
	}

	if grep == "" {
		return "", ""
	}

	prefixCmd := ""
	switch m["Type"] {
	case "cat":
		prefixCmd = "cat "

	case "tail":
		prefixCmd = "tail -f "
		f, err := getLatestFileName(config.Conf.Client.Target.Dir)
		if err != nil {
			return "", ""
		}
		file = f
	}
	grep = prefixCmd + file + grep + " >" + config.Conf.Client.Target.Dir + id + ".tmp"

	return grep, config.Conf.Client.Target.Dir + id + ".tmp"
}

func ExecGrepCmd(ctx context.Context, cmd string, tmp string, msg chan Message) {
	message := Message{}
	defer close(msg)

	fmt.Println("~~~~~~~~~~", cmd)

	ips, err := util.GetLocalIP()
	if err != nil {
		message.Err = err.Error()
		msg <- message
		return
	}

	message.IP = strings.Join(ips, ",")

	fmt.Println("----------cmd=", cmd)
	IsCmd := exec.Command("bash", "-c", cmd)

	method := ""
	if strings.HasPrefix(cmd, "cat") {
		method = "cat"
		err = Cat(IsCmd, tmp)

	} else if strings.HasPrefix(cmd, "tail") {
		method = "tail"
		err = Tail(IsCmd, tmp)
	}

	if err != nil {
		fmt.Println("exec bash shell is failed, err: ", err)
		message.Err = err.Error()
		msg <- message
		return
	}

	f, err := os.Open(tmp)
	if err != nil {
		fmt.Println("------err:", err)
		message.Err = err.Error()
		msg <- message
		return
	}
	defer f.Close()

	fmt.Println("===read tmp file")
	switch method {
	case "cat":
		ReadCatTmpFile(f, msg, message)
	case "tail":
		ReadTailTmpFile(ctx, f, msg, message)
	}

	return
}

// func getLatestFileName() string {

// 	t := time.Now()
// 	year, month, date := t.Date()
// 	hour := t.Hour()
// 	minute := t.Minute()
// 	name := config.Conf.Client.Target.Dir + "api.log." + strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", date) + fmt.Sprintf("%02d", hour) + fmt.Sprintf("%02d", minute)

// 	return name
// }

func getLatestFileName(dir string) (string, error) {
	fileInformation, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("read dir is failed, err:", err)
		return "", nil
	}

	var fileInfos fileInfo
	for _, file := range fileInformation {
		fi := &FileInfo{
			Name:    file.Name(),
			ModTime: file.ModTime(),
		}

		// 针对 server_api 日志名称格式
		if strings.Contains(fi.Name, "api.log.") {
			fileInfos = append(fileInfos, fi)
		}
	}
	sort.Sort(fileInfos)

	fmt.Println("------name:", fileInfos[len(fileInfos)-1].Name, "modtime:", fileInfos[len(fileInfos)-1].ModTime)

	return config.Conf.Client.Target.Dir + fileInfos[len(fileInfos)-1].Name, nil
}

func removeTmpFile(file string) {

	// 注册延时处理函数
	util.AddDelayFunc("remove", func(file interface{}) {
		os.Remove(file.(string))
	})
	util.Tw.AddTimer(600*time.Second, "remove", util.GenerateDelayParameter("remove", file))
}
