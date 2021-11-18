package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "HelloWorld")
}

func check(l1 string, l2 string) string {
	if l1 == "GPL-2.0-only" && l2 == "LGPL-2.0-only" {
		return "冲突"
	}
	switch l1 {
	case "GPL-2.0-only":
		switch l2 {
		case "LGPL-3.0-only":
			return "冲突"
		case "GPL-3.0-only":
			return "冲突"
		case "AGPL-3.0-only":
			return "冲突"
		case "MS-RL":
			return "冲突"
		case "ODbL-1.0":
			return "冲突"
		case "OSL-3.0":
			return "冲突"
		case "Vim":
			return "冲突"
		case "LPL-1.02":
			return "冲突"
		case "Apache-2.0":
			return "冲突"
		case "ECL-2.0":
			return "冲突"
		case "php-3.01":
			return "冲突"
		}
	}
	return "合法"
}

func licenseCheck(c *gin.Context) {
	l1 := c.Query("l1")
	l2 := c.Query("l2")
	c.IndentedJSON(http.StatusOK, check(l1, l2))
}

func main() {
	router := gin.Default()
	router.GET("/helloWorld", helloWorld)
	router.GET("/check", licenseCheck)
	router.Run("localhost:8080")
}
