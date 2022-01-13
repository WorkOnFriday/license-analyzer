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

	file, err := os.Open("conf/config.json")
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
