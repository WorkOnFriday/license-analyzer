/*
Package conf
支持手动配置，解析项目的配置文本文件
比如json格式:
import "encoding/json"
或
import jsonIter "github.com/json-iterator/go"
ini形式:
"gopkg.in/gcfg.v1"
*/
package conf

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type MysqlConfigure struct {
	IsUsed       bool
	Username     string
	Password     string
	DatabaseName string
}

type Configure struct {
	MySQL MysqlConfigure
}

var Config Configure

func ReadConfFile() {
	log.Println("read conf file")

	file, err := os.Open(GetAppPath() + "/conf/config.json")
	if err != nil {
		log.Fatalln("Configure file error ", err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		log.Fatalln("decode error", err.Error())
	}
	log.Println("read conf file success", Config)
}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
