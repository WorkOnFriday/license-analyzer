/*
Package controller
编写路由到的控制器函数
*/
package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"license-analyzer/scanner"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const uploadPath = "upload"

func HelloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "HelloWorld")
}

func LicenseCheck(c *gin.Context) {
	l1 := c.Query("l1")
	l2 := c.Query("l2")
	c.IndentedJSON(http.StatusOK, scanner.Check(l1, l2))
}

// SubmitScanTask 提交扫描jar或zip项目包的任务
func SubmitScanTask(c *gin.Context) {
	// 获取项目包文件
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Error("SubmitScanTask get form err: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	// 获取任务号
	taskID := scanner.GetTaskID()
	// 服务端文件保存到任务号路径
	taskIDString := fmt.Sprint(taskID)
	// 服务端文件完整路径
	taskPath := path.Join(uploadPath, taskIDString)
	fileFullPath := path.Join(taskPath, filepath.Base(file.Filename))
	logrus.Debug("fileFullPath ", fileFullPath)
	// 创建文件保存路径 对于windows系统，/开头为从磁盘根目录，否则从项目根目录
	if err = os.MkdirAll(taskPath, os.ModePerm); err != nil {
		logrus.Error("SubmitScanTask create file directory err: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 保存到磁盘
	if err = c.SaveUploadedFile(file, fileFullPath); err != nil {
		logrus.Error("SubmitScanTask upload file err: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 检查文件类型
	suffix := strings.ToUpper(filepath.Ext(file.Filename))
	if suffix != ".ZIP" && suffix != ".JAR" {
		logrus.Error("SubmitScanTask upload file type err: ", err.Error())
		c.Status(http.StatusBadRequest)
	}
	// 提交任务
	scanner.SubmitTask(taskID, fileFullPath)
	// 未发现易发现错误，返回任务号
	c.String(http.StatusOK, taskIDString)
}

// GetTaskResult 获取任务执行结果
func GetTaskResult(c *gin.Context) {
	// 获取任务号
	taskIDString := c.Query("id")
	// 转换任务号
	taskID, err := strconv.Atoi(taskIDString)
	if err != nil {
		logrus.Error("GetTaskResult strconv error: ", err.Error())
		c.Status(http.StatusBadRequest)
	}
	// 获取任务结果
	taskResult := scanner.GetTaskResult(taskID)
	c.JSON(http.StatusOK, taskResult)
}
