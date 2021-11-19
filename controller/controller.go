/*
Package controller
编写路由到的控制器函数
*/
package controller

import (
	"github.com/gin-gonic/gin"
	"license-analyzer/scanner"
	"net/http"
)

func HelloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "HelloWorld")
}

func LicenseCheck(c *gin.Context) {
	l1 := c.Query("l1")
	l2 := c.Query("l2")
	c.IndentedJSON(http.StatusOK, scanner.Check(l1, l2))
}
