/*
Package logger
全局日志配置
*/
package logger

import (
	"github.com/sirupsen/logrus"
)

func SetLoggerConfig() {
	// 设置日志格式为JSON格式 默认text格式
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	//logrus.SetOutput(os.Stdout)
	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debugln("set logger config")
}
