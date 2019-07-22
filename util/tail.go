package util

import "github.com/hpcloud/tail"

func Tail(fileName string) (chan *tail.Line, error) {
	tailfs, err := tail.TailFile(fileName, tail.Config{
		ReOpen:    true,                                 // 文件被移除或被打包，需要重新打开
		Follow:    true,                                 // 实时跟踪
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 如果程序出现异常，保存上次读取的位置，避免重新读取。
		MustExist: false,                                // 如果文件不存在，是否推出程序，false是不退出
		Poll:      true,
	})
	if err != nil {
		return nil, err
	}

	return tailfs.Lines, nil
}
