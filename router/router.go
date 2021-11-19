/*
Package router
设置路由，将请求转发到控制器函数
*/
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"license-analyzer/controller"
)

func SetRouterAndRun() {
	router := gin.Default()
	router.GET("/helloWorld", controller.HelloWorld)
	router.GET("/check", controller.LicenseCheck)
	logrus.Debugln("set router")

	if err := router.Run("localhost:8080"); err != nil {
		logrus.Fatal(err.Error())
	}
}
