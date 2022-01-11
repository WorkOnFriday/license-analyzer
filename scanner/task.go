package scanner

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Task struct {
	ID       int
	FullPath string
}
type TaskResult struct {
	IsFinish     bool
	ErrorMessage string
}

var taskCounter chan int

const channelBufferSize = 100

var taskQueue chan Task

var taskResultMap = make(map[int]TaskResult)

func startTaskCounter() {
	taskCounter = make(chan int)
	go func() {
		id := 0
		for {
			taskCounter <- id
			id++
		}
	}()
}

func startTaskQueue() {
	taskQueue = make(chan Task, channelBufferSize)
	go func() {
		select {
		case task := <-taskQueue:
			logrus.Debugf("startTaskQueue %+v", task)
			taskResultMap[task.ID] = TaskResult{IsFinish: true, ErrorMessage: "test"}

			// 删除任务文件夹，节省磁盘空间
			os.RemoveAll(filepath.Dir(task.FullPath))
		}
	}()
}

func StartTaskSystem() {
	startTaskCounter()
	startTaskQueue()
}

func GetTaskID() int {
	return <-taskCounter
}

func SubmitTask(taskID int, fileFullPath string) {
	taskQueue <- Task{ID: taskID, FullPath: fileFullPath}
}

func GetTaskResult(taskID int) (result TaskResult) {
	result = taskResultMap[taskID]
	return
}
