/*
Package controller
编写路由到的控制器函数
*/
package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"license-analyzer/scanner"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func HelloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "HelloWorld")
}

func LicenseCheck(c *gin.Context) {
	l1 := c.Query("l1")
	l2 := c.Query("l2")
	c.IndentedJSON(http.StatusOK, scanner.Check(l1, l2))
}

func ScannerLicenseFile(c *gin.Context) {
	// 用户通过表单项设置服务端文件保存路径
	dst := c.PostForm("dst")
	logrus.Debug("dst", dst)
	// 获取表单中的文件
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Error("ScannerLicenseFile get form err: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	logrus.Debug("file.Filename", file.Filename, "file.Size", file.Size)
	// 服务端文件完整路径
	fileFullPath := path.Join(dst, filepath.Base(file.Filename))
	logrus.Debug("fileFullPath", fileFullPath)
	// 创建文件保存路径 对于windows系统，/开头为从磁盘根目录，否则从项目根目录
	if err = os.MkdirAll(dst, os.ModePerm); err != nil {
		logrus.Error("ScannerLicenseFile create file directory err: ", dst, " ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 保存到磁盘
	if err = c.SaveUploadedFile(file, fileFullPath); err != nil {
		logrus.Error("ScannerLicenseFile upload file err: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 扫描协议文件
	result := scanner.ScanFile(fileFullPath)
	// 返回结果
	logrus.Debug(result)
	c.String(http.StatusOK, result)
}

// ScannerPackage 扫描jar或zip项目包
func ScannerPackage(c *gin.Context) {
	// 用户通过表单项设置服务端文件保存路径
	dst := c.PostForm("dst")
	// 获取表单中的文件
	file, err := c.FormFile("file")
	if err != nil {
		logrus.Error("ScannerPackage get form err: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	// 服务端文件完整路径
	fileFullPath := path.Join(dst, filepath.Base(file.Filename))
	logrus.Debug("fileFullPath ", fileFullPath)
	// 创建文件保存路径 对于windows系统，/开头为从磁盘根目录，否则从项目根目录
	if dst != "" {
		if err = os.MkdirAll(dst, os.ModePerm); err != nil {
			logrus.Error("ScannerPackage create file directory err: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}
	}
	// 保存到磁盘
	if err = c.SaveUploadedFile(file, fileFullPath); err != nil {
		logrus.Error("ScannerPackage upload file err: ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	// 检查文件类型
	suffix := strings.ToUpper(filepath.Ext(file.Filename))
	if suffix != ".ZIP" && suffix != ".JAR" {
		logrus.Error("ScannerPackage upload file type err: ", err.Error())
		c.Status(http.StatusBadRequest)
	}
	// 扫描协议文件
	result := scanner.ScanPackage(fileFullPath)
	// 返回结果
	marshal, err := json.Marshal(result)
	if err != nil {
		logrus.Error("ScannerPackage json.Marshal err: ", err.Error())
		c.Status(http.StatusInternalServerError)
	}
	c.String(http.StatusOK, string(marshal))
}
