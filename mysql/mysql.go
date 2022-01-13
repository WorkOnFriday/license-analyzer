/*
Package mysql
连接和操作mysql
*/
package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"license-analyzer/conf"
)

var db *sqlx.DB

type TaskResult struct {
	ID     int
	Result []byte
}

func InitializeMySQL() {
	if !conf.Config.MySQL.IsUsed {
		return
	}

	logrus.Debugln("set mysql")

	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Config.MySQL.Username, conf.Config.MySQL.Password, conf.Config.MySQL.DatabaseName)
	//使用 Connect 连接,会验证是否连接成功,
	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		logrus.Fatalf("mysql connent error: %s", err.Error())
	}
	logrus.Debugln("set mysql success")
	// 创建数据表
	createTaskResultTableIfNotExist := "CREATE TABLE IF NOT EXISTS `task_result` (`id` INT NOT NULL AUTO_INCREMENT, `result` BLOB, PRIMARY KEY (`id`) ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4; "
	_, err = db.Exec(createTaskResultTableIfNotExist)
}

func ReadAllTaskResult() (taskResults []TaskResult) {
	selectAllTaskResultOrderByID := "SELECT * FROM `task_result` ORDER BY `id`;"
	db.Select(&taskResults, selectAllTaskResultOrderByID)
	return
}

func CreateTaskResult(id int, result []byte) {
	insertTaskResult := "INSERT INTO `task_result`(id, result) VALUES(?, ?);"
	_, err := db.Exec(insertTaskResult, id, result)
	if err != nil {
		logrus.Error("CreateTaskResult error: ", err.Error())
		return
	}
}
