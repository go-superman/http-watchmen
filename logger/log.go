package logger

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var (
	innerLog = logs.NewLogger(10000)
)

func init() {
	//debug: 7
	//fmt.Printf("logmod:%v\n", logs.LevelDebug)
	innerLog.SetLevel(logs.LevelDebug)
	innerLog.EnableFuncCallDepth(true)
	innerLog.SetLogFuncCallDepth(3)
}

func Warn(v ...interface{}) {
	innerLog.Warn(fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	innerLog.Warn(format, v...)
}

func Error(v ...interface{}) {
	innerLog.Error(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	innerLog.Error(format, v...)
}

func Info(v ...interface{}) {
	innerLog.Info(fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	innerLog.Info(format, v...)
}

func Debug(v ...interface{}) {
	innerLog.Debug(fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	innerLog.Debug(format, v...)
}
