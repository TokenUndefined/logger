package logger

import (
	"time"
	"os"
	"path/filepath"
	"fmt"
	"strings"
	"strconv"
	"log"
	"runtime"
)

var logfile *LogFile
var backupCount int = 9						// 最大保存文件数
var logMaxSize int64 = 100 * 1024 * 1024  	// 100M
var logConsolePrefix string           		//终端控制台显示前缀
var logConsole bool = true 					// 终端控制台显示控制，默认为true
var logLevel LEVEL = ALL
var today string = Today()

func SetConsole(isConsole bool) {
	logConsole = isConsole
}

func SetLogMaxSize(filesize int64) {
	logMaxSize = filesize
}

func SetLevel(_level LEVEL) {
	logLevel = _level
}

func SetLogBackupCount(count int) {
	backupCount = count
}

func SetConsolePrefix (prefix string) {
	logConsolePrefix = prefix
}

func (f *LogFile)fileSizeCheck() (bool, string) {
	fileSize := GetFileSize(f.logfile_path)
	if fileSize >= logMaxSize {
		var fn string
		maxfile_path := fmt.Sprintf("%s.%d", f.logfile_path, backupCount)
		if f.logfile_path==f.lastlog_path || f.lastlog_path==maxfile_path {
			fn = fmt.Sprintf("%s.%d", f.logfile_path, 1)
		}else{
			s := strings.Split(f.lastlog_path, ".")
			num,err := strconv.Atoi(s[len(s)-1])
			if err != nil{
				panic(err)
			}
			fn = fmt.Sprintf("%s.%d", f.logfile_path, num+1)
		}
		return true, fn
	}
	return false, ""
}

func (f *LogFile) getLogfile() (logfile_path string) {
	if f.split_mode==SPLIT_BY_FILESIZE {
		logfile_path = filepath.Join(f.logdir_path, f.logfilename+".log")
	} else{
		logfile_path = filepath.Join(f.logdir_path, f.logfilename + "_" + today + ".log")
	}

	if PathExits(logfile_path)==false || f.split_mode==SPLIT_BY_DATE {
		return logfile_path
	}

	if len(f.logfile_path)==0 {
		f.logfile_path = logfile_path
		f.lastlog_path = logfile_path
	}

	flag, fn := f.fileSizeCheck()
	if flag {
		err1 := os.Rename(logfile_path, fn)
		if err1 !=nil {
			panic(err1)
		}else {
			f.lastlog_path = fn
		}
	}
	return logfile_path
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func fileCheck() {
	defer catchError()
	// 检查是否跨天
	if time.Now().YearDay() != logfile.timestamp.YearDay() {
		today = Today()
	}
	// 检查文件是否存在
	fn := logfile.getLogfile()
	if PathExits(fn)==false {
		logfile.Lock()
		defer logfile.Unlock()

		if logfile.fp != nil {
			logfile.fp.Close()
		}

		logfile.logfile_path = fn
		logfile.fp, _ = os.OpenFile(fn, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		logfile.logger = log.New(logfile.fp, "", LogFlags)
		logfile.timestamp = time.Now()
	}
}

func fileMonitor() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			fileCheck()
		}
	}
}

func Initialize(path string, fileName string, split_mode int) {
	// 判断日志存储目录是否存在
	if PathExits(path) == false {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// 结构体初始化
	logfile = &LogFile{
		logdir_path: path,
		logfilename: fileName,
		timestamp: time.Now(),
		split_mode: split_mode,
	}
	logfile.Lock()
	defer logfile.Unlock()
	// 创建文件
	fn := logfile.getLogfile()
	var err error
	logfile.fp, err = os.OpenFile(fn, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	// 初始化日志
	logfile.logger = log.New(logfile.fp, "", LogFlags)
	log.SetFlags(LogConsoleFlag)
	logfile.logfile_path = fn
	logfile.lastlog_path = fn
	//启动文件监控模块
	go fileMonitor()
}

func SprintColor(str string, s STYLE, fc, bc COLOR) string {
	_fc := int(fc)      	//前景色
	_bc := int(bc) + 10 	//背景色
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, int(s), _bc, _fc, str, 0x1B)
}

func console(ll LEVEL, args string) {

	if logConsole {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		now := time.Now()
		context := ""
		if len(logConsolePrefix) > 0 {
			context = fmt.Sprintf("[%04d/%02d/%02d %02d:%02d:%02d.%03d] %s #%s:%d ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), time.Duration(now.Nanosecond())/(time.Microsecond), logConsolePrefix, file, line)
		} else {
			context = fmt.Sprintf("[%04d/%02d/%02d %02d:%02d:%02d.%03d] %s:%d ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), time.Duration(now.Nanosecond())/(time.Microsecond), file, line)
		}
		prefix := "[default]"
		switch ll {
		case DEBUG:
			prefix = fmt.Sprintf("%s", SprintColor("[Debug]", STYLE_INVERSE, CLR_GREEN, CLR_DEFAULT))
			context = fmt.Sprintf("%s %s", SprintColor(context, STYLE_INVERSE, CLR_GREEN, CLR_DEFAULT), args)
		case INFO:
			prefix = fmt.Sprintf("%s", SprintColor("[Info]", STYLE_DEFAULT, CLR_DEFAULT, CLR_BLUE))
			context = fmt.Sprintf("%s %s", SprintColor(context, STYLE_DEFAULT, CLR_DEFAULT, CLR_BLUE), args)
		case WARN:
			prefix = fmt.Sprintf("%s", SprintColor("[Warn]", STYLE_DEFAULT, CLR_YELLOW, CLR_DEFAULT))
			context = fmt.Sprintf("%s %s", SprintColor(context, STYLE_DEFAULT, CLR_YELLOW, CLR_DEFAULT), args)
		case ERROR:
			prefix = fmt.Sprintf("%s", SprintColor("[Error]", STYLE_HIGHLIGHT, CLR_RED, CLR_DEFAULT))
			context = fmt.Sprintf("%s %s", SprintColor(context, STYLE_HIGHLIGHT, CLR_RED, CLR_DEFAULT), args)
		case FATAL:
			prefix = fmt.Sprintf("%s", SprintColor("[Fatal]", STYLE_HIGHLIGHT, CLR_PURPLE, CLR_DEFAULT))
			context = fmt.Sprintf("%s %s", SprintColor(context, STYLE_HIGHLIGHT, CLR_PURPLE, CLR_DEFAULT), args)
		default:
			context = fmt.Sprintf("%s %s", SprintColor(context, STYLE_DEFAULT, CLR_DEFAULT, CLR_DEFAULT), args)
		}
		log.SetPrefix(prefix)
		log.Printf(context)
	}
}

func Debug(arg interface{}) {
	defer catchError()
	if logLevel <= DEBUG {
		context := fmt.Sprintf("%s", fmt.Sprintln(arg))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Debug]")
			logfile.logger.Output(2, context)
		}
		console(DEBUG, context)
	}
}

func Debugf(format string, args ...interface{}) {
	defer catchError()
	if logLevel <= DEBUG {
		context := fmt.Sprintf("%s", fmt.Sprintf(format, args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Debug]")
			logfile.logger.Output(2, context)
		}
		console(DEBUG, context)
	}
}

func Debugln(args ...interface{}) {
	defer catchError()
	if logLevel <= DEBUG {
		context := fmt.Sprintf("%s", fmt.Sprintln(args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Debug]")
			logfile.logger.Output(2, context)
		}
		console(DEBUG, context)
	}
}

func Info(arg interface{}) {
	defer catchError()
	if logLevel <= INFO {
		context := fmt.Sprintf("%s", fmt.Sprintln(arg))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Info]")
			logfile.logger.Output(2, context)
		}
		console(INFO, context)
	}
}

func Infof(format string, args ...interface{}) {
	defer catchError()
	if logLevel <= INFO {
		context := fmt.Sprintf("%s", fmt.Sprintf(format, args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Info]")
			logfile.logger.Output(2, context)
		}
		console(INFO, context)
	}
}

func Infoln(args ...interface{}) {
	defer catchError()
	if logLevel <= INFO {
		context := fmt.Sprintf("%s", fmt.Sprintln(args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Info]")
			logfile.logger.Output(2, context)
		}
		console(INFO, context)
	}
}

func Warn(arg interface{}) {
	defer catchError()
	if logLevel <= WARN {
		context := fmt.Sprintf("%s", fmt.Sprintln(arg))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Warn]")
			logfile.logger.Output(2, context)
		}
		console(WARN, context)
	}
}

func Warnf(format string, args ...interface{}) {
	defer catchError()
	if logLevel <= WARN {
		context := fmt.Sprintf("%s", fmt.Sprintf(format, args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Warn]")
			logfile.logger.Output(2, context)
		}
		console(WARN, context)
	}
}

func Warnln(args ...interface{}) {
	defer catchError()
	if logLevel <= WARN {
		context := fmt.Sprintf("%s", fmt.Sprintln(args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Warn]")
			logfile.logger.Output(2, context)
		}
		console(WARN, context)
	}
}

func Error(arg interface{}) {
	defer catchError()
	if logLevel <= ERROR {
		context := fmt.Sprintf("%s", fmt.Sprintln(arg))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Error]")
			logfile.logger.Output(2, context)
		}
		console(ERROR, context)
	}
}

func Errorf(format string, args ...interface{}) {
	defer catchError()
	if logLevel <= ERROR {
		context := fmt.Sprintf("%s", fmt.Sprintf(format, args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Error]")
			logfile.logger.Output(2, context)
		}
		console(ERROR, context)
	}
}

func Errorln(args ...interface{}) {
	defer catchError()
	if logLevel <= ERROR {
		context := fmt.Sprintf("%s", fmt.Sprintln(args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Error]")
			logfile.logger.Output(2, context)
		}
		console(ERROR, context)
	}
}

func Fatal(arg interface{}) {
	defer catchError()
	if logLevel <= FATAL {
		context := fmt.Sprintf("%s", fmt.Sprintln(arg))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Fatal]")
			logfile.logger.Output(2, context)
		}
		console(FATAL, context)
	}
}

func Fatalf(format string, args ...interface{}) {
	defer catchError()
	if logLevel <= FATAL {
		context := fmt.Sprintf("%s", fmt.Sprintf(format, args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Fatal]")
			logfile.logger.Output(2, context)
		}
		console(FATAL, context)
	}
}

func Fatalln(args ...interface{}) {
	defer catchError()
	if logLevel <= FATAL {
		context := fmt.Sprintf("%s", fmt.Sprintln(args...))
		if logfile != nil {
			logfile.RLock()
			defer logfile.RUnlock()
			logfile.logger.SetPrefix("[Fatal]")
			logfile.logger.Output(2, context)
		}
		console(FATAL, context)
	}
}
