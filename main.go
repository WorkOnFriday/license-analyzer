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
	result, err := util.FetchDynamicHTMLItemInnerText("https://mvnrepository.com/artifact/org.springframework/spring-core",
		"#maincontent > table > tbody > tr:nth-child(1) > td > span")
	logrus.Infoln("Fetch result: ", result, err)
	// 运行服务
	router.SetRouterAndRun()
}
