package script

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/yumimobi/trace/util"
)

func Cat(IsCmd *exec.Cmd, fileName string) error {
	return IsCmd.Run()
}

func ReadCatTmpFile(f *os.File, msg chan Message, message Message) {
	// 延迟删除临时文件
	removeTmpFile(f.Name())
	fmt.Println("=+++++++++")

	r := bufio.NewReader(f)
	for {
		str, err := r.ReadString('\n')
		str = util.GreedyMatchJSONString(str)

		// 走到此处，如果是EOF，则在读取到EOF之前会把有效的数据读入到str中
		if err == io.EOF {
			//此时str中是有数据的需要处理
			message.Msg = str
			msg <- message
			fmt.Println("---- io.EOF")
			break
		}
		if err != nil {
			message.Err = err.Error()
			msg <- message
			//应该返回错误,先判断是否为nil，不为在判断是否等于EOF,等于break出for，然后继续执行；其他错误直接return
			// fmt.Println("read string is failed, err: ", err)
			fmt.Println("---- io.EOF is nil")
			break
		}
		//str中会包含'\n'
		message.Msg = str
		msg <- message
		fmt.Println("---- message <-")
	}
	return
}
