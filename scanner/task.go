package scanner

import (
	"github.com/sirupsen/logrus"
	"license-analyzer/util"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Task struct {
	ID       int
	FullPath string
}

type PomLicense struct {
	XMLDependency
	License string
}
type TaskResult struct {
	IsFinish          bool
	ErrorMessage      string
	Local             []PathLicense
	External          []PathLicense
	AllExternalModule []JarPackages
	AllLocalModule    []ModuleLicense
	Dependency        AllModuleDependency
	PomLicense        []PomLicense
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

	runOneTaskInQueue := func() {
		select {
		case task := <-taskQueue:
			// 删除任务文件夹，节省磁盘空间
			defer os.RemoveAll(filepath.Dir(task.FullPath))

			logrus.Debugf("startTaskQueue %+v", task)

			// 执行分析
			external := make([]PathLicense, 0)
			local := make([]PathLicense, 0)

			// 创建解压到的临时文件夹
			if err := os.MkdirAll(ProjectTmp, os.ModePerm); err != nil {
				logrus.Fatalln(err)
			}
			// 删除解压到的文件夹
			defer os.RemoveAll(ProjectTmp)

			// 解压
			zipDecompress(task.FullPath, ProjectTmp)

			// 扫描 项目本身/jar包，获得许可证路径/包 和 许可证类型
			deepScan(ProjectTmp, "", true, &local, &external)
			// 删除jar包解压到的文件夹
			defer os.RemoveAll(DIR)

			logrus.Debugf("Scan result external: %+v", external)
			logrus.Debugf("Scan result local: %+v", local)

			// 获取jar包中的包名
			allExternalModule := findAllExternalModule(external)
			logrus.Debugf("all external modules: %+v", allExternalModule)

			// 获取项目本身包名和许可证类型
			allLocalModule := findAllLocalModule(local)
			logrus.Debugf("all local modules: %+v", allLocalModule)

			// 得到项目本身每个包对jar包的依赖
			dependency := dependencyAnalyze(findAllExternalModule(external), local)
			logrus.Debugf("dependency result: %+v", dependency)

			// 得到pom.xml定义的外部依赖
			packageFileName := filepath.Base(task.FullPath)
			packageType := path.Ext(packageFileName)
			packageName := strings.TrimSuffix(packageFileName, packageType)
			pomFullPath := filepath.Join(ProjectTmp, packageName, "pom.xml")
			logrus.Debugf("pom full path: %+v", pomFullPath)
			p := PomScan(pomFullPath)
			logrus.Debugf("pom result: %+v", p)

			// 抓取pom.xml外部依赖的许可证
			errorMessage := ""
			var pomLicenses []PomLicense
			for _, d := range p.Dependencies.Dependencies {
				result, err := util.FetchMvnPackageLicense(d.GroupID, d.ArtifactID, d.Version)
				if err != nil {
					errorMessage += err.Error()
					continue
				}
				pomLicenses = append(pomLicenses, PomLicense{XMLDependency: d, License: result})
			}
			// 汇总结果
			taskResultMap[task.ID] = TaskResult{
				IsFinish:          true,
				ErrorMessage:      "",
				Local:             local,
				External:          external,
				AllExternalModule: allExternalModule,
				AllLocalModule:    allLocalModule,
				Dependency:        dependency,
				PomLicense:        pomLicenses,
			}
		}
	}

	// 后台运行
	go func() {
		// 不断读取已提交的任务
		for {
			runOneTaskInQueue()
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
