package main

import (
	"license-analyzer/conf"
	"license-analyzer/logger"
	"license-analyzer/mysql"
	"license-analyzer/router"
	"license-analyzer/scanner"
)

func main() {
	// 读取和解析配置文件(可以包括对日志的配置)
	conf.ReadConfFile()
	// 设置全局日志(输出接下来操作的错误信息)
	logger.SetLoggerConfig()
	// 设置数据库
	mysql.InitializeMySQL()
	// 运行任务处理队列
	scanner.StartTaskSystem()
	// 初始化许可证分析
	scanner.InitializeAnalyzer()
	// 运行服务
	router.SetRouterAndRun()
}
