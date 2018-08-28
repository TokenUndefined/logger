package logger

import (
	"os"
	"time"
	"fmt"
)

func PathExits(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 当前日期
func Today() string {
	return time.Now().Format("060102")
}

//获取单个文件的大小
func GetFileSize(path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size() //获取size
	//fmt.Println(path+" 的大小为", fileSize, "byte")
	return fileSize
}