package logger

import (
	"os"
	"sync"
	"log"
	"time"
	"path/filepath"
)

type LEVEL int 					//日志等级
type COLOR int 					//显示颜色
type STYLE int 					//显示样式

type LogFile struct {
	sync.RWMutex             	//线程锁
	logdir_path     string      //日志存放目录
	logfilename 	string      //日志基础名字
	timestamp    	time.Time   //日志创建时的时间戳
	lastlog_path    string		//上一次保存日志的路径
	logfile_path  	string      //当前日志路径
	fp      		*os.File    //当前日志文件实例
	logger       	*log.Logger //当前日志操作实例
	split_mode      int
}

const (
	SPLIT_BY_FILESIZE = iota  	// 按文件大小拆分
	SPLIT_BY_DATE				// 按日期分割
	SPLIT_BY_DEFAULT			// 默认分割方式：即按文件大小又按日期分割
)

const (
	CLR_BLACK   = COLOR(30) 	// 黑色
	CLR_RED     = COLOR(31) 	// 红色
	CLR_GREEN   = COLOR(32) 	// 绿色
	CLR_YELLOW  = COLOR(33) 	// 黄色
	CLR_BLUE    = COLOR(34) 	// 蓝色
	CLR_PURPLE  = COLOR(35) 	// 紫红色
	CLR_CYAN    = COLOR(36) 	// 青蓝色
	CLR_WHITE   = COLOR(37) 	// 白色
	CLR_DEFAULT = COLOR(39) 	// 默认
)

const (
	STYLE_DEFAULT   = STYLE(0) 	//终端默认设置
	STYLE_HIGHLIGHT = STYLE(1) 	//高亮显示
	SYTLE_UNDERLINE = STYLE(4) 	//使用下划线
	SYTLE_BLINK     = STYLE(5) 	//闪烁
	STYLE_INVERSE   = STYLE(7) 	//反白显示
	STYLE_INVISIBLE = STYLE(8) 	//不可见
)

const (
	ALL   LEVEL = iota 			//所有日志
	DEBUG              			//调试
	INFO               			//信息
	WARN               			//警告
	ERROR              			//错误
	FATAL              			//崩溃
)

const (
	LogFlags             = log.Ldate | log.Lmicroseconds | log.Lshortfile //日志输出flag
	LogConsoleFlag       = 0                                              //console输出flag
	LogDumpExceptionFlag = 0
)

func init()  {
	// 输出执行程序名
	running_name := filepath.Base(os.Args[0])
	Initialize("/tmp", running_name, SPLIT_BY_FILESIZE)
}
