package scanner

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"license-analyzer/conf"
	"license-analyzer/mysql"
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
type ExternalConflict struct {
	MainLicense     PathLicense
	ExternalLicense PathLicense
	Result          ConflictResult
}
type PomConflict struct {
	MainLicense PathLicense
	PomLicense  PomLicense
	Result      ConflictResult
}
type TaskResult struct {
	IsFinish          bool
	ErrorMessage      string
	Local             []PathLicense
	Dependency        AllModuleDependency
	PomLicense        []PomLicense
	LocalConflicts    []ExternalConflict
	ExternalConflicts []ExternalConflict
	PomConflicts      []PomConflict
	RecommendLicenses []string
}

var taskCounter chan int

const channelBufferSize = 100

var taskQueue chan Task

var taskResultMap = make(map[int]TaskResult)

func readTaskResultFromMySQL() (id int) {
	dbTaskResults := mysql.ReadAllTaskResult()
	taskResults := make([]TaskResult, len(dbTaskResults))
	for i, result := range dbTaskResults {
		err := json.Unmarshal(result.Result, &taskResults[i])
		if err != nil {
			logrus.Error("ReadAllTaskResult error: ", err.Error())
		}
		taskResultMap[result.ID] = taskResults[i]
		if result.ID > id {
			id = result.ID
		}
	}
	logrus.Debugf("readTaskResultFromMySQL taskResults %+v", taskResults)
	return
}

func startTaskCounter(id int) {
	taskCounter = make(chan int)
	go func() {
		for {
			id++
			taskCounter <- id
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
				result = strings.ReplaceAll(result, " ", "-")
				pomLicenses = append(pomLicenses, PomLicense{XMLDependency: d, License: result})
			}
			logrus.Debugf("pomLicenses %+v", pomLicenses)

			// 记录推荐许可证依据的其它许可证
			var otherLicenses []string

			// 获取整个项目的许可证
			// 检查项目深层许可证是否与浅层冲突
			var mainLicenseExist bool
			var mainLicense PathLicense
			var localConflicts []ExternalConflict
			for _, mainPathLicense := range local {
				if filepath.Dir(mainPathLicense.Path) == dependency.Project.Name {
					// 是整个项目的许可证
					mainLicenseExist = true
					mainLicense = mainPathLicense
					logrus.Debugf("mainLicense %+v", mainLicense)
				} else {
					// 不是整个项目的许可证
					otherLicenses = append(otherLicenses, LicenseLongNameToShort(mainPathLicense.License))
				}

				for _, pathLicense := range local {
					if !strings.HasPrefix(filepath.Dir(pathLicense.Path), filepath.Dir(mainPathLicense.Path)) {
						continue
					}
					result := CheckLicenseConflictByShortName(LicenseLongNameToShort(mainPathLicense.License),
						LicenseLongNameToShort(pathLicense.License))
					if !result.Pass {
						localConflicts = append(localConflicts,
							ExternalConflict{mainPathLicense, pathLicense, result})
					}
				}
			}
			logrus.Debugf("localConflicts: %+v", localConflicts)

			// 检查项目模块和依赖的jar包之间的冲突
			var externalConflicts []ExternalConflict
			for _, module := range dependency.Modules {
				for _, jarPackage := range module.Dependencies {
					// 每个模块依赖的每个jar包的许可证
					logrus.Debug(jarPackage.License)
					otherLicenses = append(otherLicenses, LicenseLongNameToShort(jarPackage.License))

					for _, license := range local {
						if module.Name != filepath.Dir(license.Path) {
							continue
						}
						// 模块在项目中许可证的范围内
						logrus.Debugf("- %s", license.License)
						// 检测和记录许可证冲突
						result := CheckLicenseConflictByShortName(LicenseLongNameToShort(license.License),
							LicenseLongNameToShort(jarPackage.License))
						if !result.Pass {
							externalConflicts = append(externalConflicts,
								ExternalConflict{MainLicense: license, ExternalLicense: jarPackage.PathLicense,
									Result: result})
						}
					}
				}
			}
			logrus.Debugf("externalConflicts: %+v", externalConflicts)

			// 若项目许可证存在，检查项目和pom.xml描述的依赖之间的冲突
			var recommendLicenses []string
			var pomConflicts []PomConflict
			if mainLicenseExist {
				for _, pomLicense := range pomLicenses {
					// pom.xml引入的依赖的许可证 应与 项目许可证不冲突
					logrus.Debugf("- %s", pomLicense.License)
					// 检测和记录许可证冲突
					result := CheckLicenseConflictByShortName(LicenseLongNameToShort(mainLicense.License),
						pomLicense.License)
					if !result.Pass {
						pomConflicts = append(pomConflicts,
							PomConflict{MainLicense: mainLicense, PomLicense: pomLicense,
								Result: result})
					}
				}
			}
			logrus.Debugf("pomConflicts: %+v", pomConflicts)

			// 不论有没有项目许可证，是否冲突，总是进行推荐
			for _, pomLicense := range pomLicenses {
				otherLicenses = append(otherLicenses, pomLicense.License)
			}
			logrus.Debugf("otherLicenses: %+v", otherLicenses)
			recommendLicenses = RecommendByLibraryLicenseShortName(otherLicenses)
			logrus.Debugf("recommendLicenses: %+v", recommendLicenses)

			// 汇总结果
			taskResult := TaskResult{
				IsFinish:          true,
				ErrorMessage:      "",
				Local:             local,
				Dependency:        dependency,
				PomLicense:        pomLicenses,
				LocalConflicts:    localConflicts,
				ExternalConflicts: externalConflicts,
				PomConflicts:      pomConflicts,
				RecommendLicenses: recommendLicenses,
			}
			taskResultMap[task.ID] = taskResult
			if conf.Config.MySQL.IsUsed {
				result, err := json.Marshal(&taskResult)
				if err != nil {
					logrus.Error("startTaskQueue error: ", err.Error())
					return
				}
				mysql.CreateTaskResult(task.ID, result)
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
	var id int
	if conf.Config.MySQL.IsUsed {
		// 读MySQL数据
		id = readTaskResultFromMySQL()
	}
	// 继续计数
	startTaskCounter(id)
	// 开始任务队列
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
