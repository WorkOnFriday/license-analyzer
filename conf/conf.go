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
	"log"
)

func ReadConfFile() {
	log.Println("read conf file")
}
