package script

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hpcloud/tail"
)

func Tail(IsCmd *exec.Cmd, fileName string) error {

	err := IsCmd.Start()
	if err != nil {
		return err
	}

	waitCreateFile(fileName)
	go killProcess(IsCmd)
	return nil
}

func waitCreateFile(fileName string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("new watcher is failed, err:", err)
	}
	defer watcher.Close()

	// 获取临时文件文件夹
	path, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		fmt.Println("get file path is failed, err:", err)
	}

	// 监听外层文件夹
	err = watcher.Add(path)
	if err != nil {
		fmt.Println("watch add a.txt is failed, err:", err)
	}

	creatOK := false
	for {
		select {
		case ev := <-watcher.Events:
			path, _ = filepath.Abs(fileName)
			if ev.Op&fsnotify.Create == fsnotify.Create && path == ev.Name {
				fmt.Println("---create is ok")
				creatOK = true
				break
			}

		case <-time.After(time.Second * 200):
			fmt.Println("~~time out")
			creatOK = true
			break
		}

		if creatOK {
			break
		}
	}
}

func killProcess(cmd *exec.Cmd) {

	time.Sleep(time.Second * 500)

	err := cmd.Process.Kill()
	if err != nil {
		fmt.Println("kill is failed, err:", err)
	}

	// 不wait会造成僵尸进程
	_, err = cmd.Process.Wait()
	if err != nil {
		fmt.Println("cmd process wait is failed,err:", err)
	}

	fmt.Println("kill tail. ")
}

func ReadTailTmpFile(f *os.File, msg chan Message, message Message) {
	// 延迟删除临时文件
	removeTmpFile(f.Name())

	tailfs, err := tail.TailFile(f.Name(), tail.Config{
		ReOpen:    true,                                 // 文件被移除或被打包，需要重新打开
		Follow:    true,                                 // 实时跟踪
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 如果程序出现异常，保存上次读取的位置，避免重新读取。
		MustExist: false,                                // 如果文件不存在，是否推出程序，false是不退出
		Poll:      true,
	})

	if err != nil {
		fmt.Println("tailf failed, err:", err)
		return
	}

	for true {
		line, ok := <-tailfs.Lines
		// ok 是判断管道是否被关闭，如果关闭就是文件被重置了，需要重新读取新的管道

		if !ok {
			fmt.Println("tailf fail close reopen, fileName:", f.Name())
			continue
		}

		message.Msg = line.Text
		msg <- message

		fmt.Println("text:", line.Text)
	}

}
