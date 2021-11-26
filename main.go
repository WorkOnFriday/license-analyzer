package main

import (
	"github.com/sirupsen/logrus"
	"license-analyzer/conf"
	"license-analyzer/logger"
	"license-analyzer/mysql"
	"license-analyzer/redis"
	"license-analyzer/router"
	"license-analyzer/util"
)

func main() {
	// 读取和解析配置文件(可以包括对日志的配置)
	conf.ReadConfFile()
	// 设置全局日志(输出接下来操作的错误信息)
	logger.SetLoggerConfig()
	// 设置数据库
	mysql.SetMySQL()
	redis.SetRedis()
	logrus.Infoln("Fetch result: ",
		util.FetchHTMLItemInnerText("https://www.baidu.com/","span.title-content-title"))
	// 运行服务
	router.SetRouterAndRun()
}
