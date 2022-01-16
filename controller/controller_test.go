package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
	"io"
	"license-analyzer/scanner"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type TestCase struct {
	name           string
	filePath       string
	expectedResult scanner.TaskResult
}

var testCaseScannerTest1 = TestCase{
	name:     "ScannerTest1",
	filePath: "..\\testProject\\ScannerTest1.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest1\\License.txt",
				License: "GENERAL PUBLIC LICENSE Version 3",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name: "ScannerTest1",
				Dependencies: []scanner.JarPackage{
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest1\\lib\\DependedProject.jar",
							License: "GENERAL PUBLIC LICENSE Version 2",
						},
						Package: []string{"pri"},
					},
				},
			},
			Modules: []scanner.UnitDependency{
				{
					Name: "ScannerTest1",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest1\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2",
							},
							Package: []string{"pri"},
						},
					},
				},
			},
		},
		PomLicense:        nil,
		LocalConflicts:    nil,
		ExternalConflicts: nil,
		PomConflicts:      nil,
		RecommendLicenses: []string{"GPL-2.0-only", "GPL-2.0-or-later",
			"GPL-3.0-only", "AGPL-3.0-only", "Apache-2.0", "ECL-2.0"},
	},
}

var testCaseScannerTest2 = TestCase{
	name:     "ScannerTest2",
	filePath: "..\\testProject\\ScannerTest2.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest2\\LICENSE.txt",
				License: "GENERAL PUBLIC LICENSE Version 2",
			},
			{
				Path:    "ScannerTest2\\src\\work\\gpl3\\License.txt",
				License: "GENERAL PUBLIC LICENSE Version 3",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name: "ScannerTest2",
				Dependencies: []scanner.JarPackage{
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest2\\lib\\DependedProject.jar",
							License: "GENERAL PUBLIC LICENSE Version 2"},
						Package: []string{"pri"},
					},
				},
			},
			Modules: []scanner.UnitDependency{
				{
					Name: "ScannerTest2",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest2\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2"},
							Package: []string{"pri"},
						},
					},
				},
				{
					Name: "ScannerTest2\\src\\work\\gpl3",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest2\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2"},
							Package: []string{"pri"},
						},
					},
				},
			},
		},
		PomLicense:        nil,
		LocalConflicts:    nil,
		ExternalConflicts: nil,
		PomConflicts:      nil,
		RecommendLicenses: []string{"GPL-2.0-or-later",
			"GPL-3.0-only", "AGPL-3.0-only", "Apache-2.0", "ECL-2.0"}, // 已删去“GPL-2.0-only”
	},
}

var testCaseScannerTest3 = TestCase{
	name:     "ScannerTest3",
	filePath: "..\\testProject\\ScannerTest3.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest3\\LICENSE",
				License: "GENERAL PUBLIC LICENSE Version 2",
			},
			{
				Path:    "ScannerTest3\\Module1\\LICENSE",
				License: "GENERAL PUBLIC LICENSE Version 3",
			},
			{
				Path:    "ScannerTest3\\Module2\\LICENSE.txt",
				License: "APACHE LICENSE Version 2.0",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name: "ScannerTest3",
				Dependencies: []scanner.JarPackage{
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest3\\lib\\DependedProject.jar",
							License: "GENERAL PUBLIC LICENSE Version 2"},
						Package: []string{"pri"},
					},
				},
			},
			Modules: []scanner.UnitDependency{
				{
					Name: "ScannerTest3",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest3\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2"},
							Package: []string{"pri"},
						},
					},
				},
				{
					Name: "ScannerTest3\\Module1",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest3\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2"},
							Package: []string{"pri"},
						},
					},
				},
				{
					Name: "ScannerTest3\\Module2",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest3\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2"},
							Package: []string{"pri"},
						},
					},
				},
			},
		},
		PomLicense:        nil,
		LocalConflicts:    nil,
		ExternalConflicts: nil,
		PomConflicts:      nil,
		RecommendLicenses: []string{"GPL-2.0-or-later",
			"GPL-3.0-only", "AGPL-3.0-only", "Apache-2.0", "ECL-2.0"},
	},
}

var testCaseScannerTest4 = TestCase{
	name:     "ScannerTest4",
	filePath: "..\\testProject\\ScannerTest4.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest4\\LICENSE",
				License: "GENERAL PUBLIC LICENSE Version 2",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name: "ScannerTest4",
				Dependencies: []scanner.JarPackage{
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest4\\lib\\DependedProject.jar",
							License: "GENERAL PUBLIC LICENSE Version 2"},
						Package: []string{"pri"},
					},
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest4\\antlib\\DependedProject2.jar",
							License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1"},
						Package: []string{"dp2"},
					},
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest4\\antlib2\\DependedProject3.jar",
							License: "MICROSOFT RECIPROCAL LICENSE"},
						Package: []string{"dp3"},
					},
				},
			},
			Modules: []scanner.UnitDependency{
				{
					Name: "ScannerTest4",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest4\\lib\\DependedProject.jar",
								License: "GENERAL PUBLIC LICENSE Version 2"},
							Package: []string{"pri"},
						},
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest4\\antlib\\DependedProject2.jar",
								License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1"},
							Package: []string{"dp2"},
						},
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest4\\antlib2\\DependedProject3.jar",
								License: "MICROSOFT RECIPROCAL LICENSE"},
							Package: []string{"dp3"},
						},
					},
				},
			},
		},
		PomLicense:     nil,
		LocalConflicts: nil,
		ExternalConflicts: []scanner.ExternalConflict{
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest4\\LICENSE",
					License: "GENERAL PUBLIC LICENSE Version 2",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest4\\antlib\\DependedProject2.jar",
					License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1",
				},
				Result: scanner.ConflictResult{Unknown: false, Pass: false, Message: "冲突"},
			},
			{

				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest4\\LICENSE",
					License: "GENERAL PUBLIC LICENSE Version 2",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest4\\antlib2\\DependedProject3.jar",
					License: "MICROSOFT RECIPROCAL LICENSE",
				},
				Result: scanner.ConflictResult{Unknown: false, Pass: false, Message: "冲突"},
			},
		},
		PomConflicts:      nil,
		RecommendLicenses: []string{"Apache-2.0", "ECL-2.0"},
	},
}

var testCaseScannerTest5 = TestCase{
	name:     "ScannerTest5",
	filePath: "..\\testProject\\ScannerTest5.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest5\\License.txt",
				License: "GENERAL PUBLIC LICENSE Version 3",
			},
			{
				Path:    "ScannerTest5\\Module1\\LICENSE",
				License: "LESSER GENERAL PUBLIC LICENSE Version 2.1",
			},
			{
				Path:    "ScannerTest5\\Module2\\LICENSE",
				License: "LESSER GENERAL PUBLIC LICENSE Version 2.1",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name: "ScannerTest5",
				Dependencies: []scanner.JarPackage{
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest5\\Module1\\lib\\DependedProject3.jar",
							License: "MICROSOFT RECIPROCAL LICENSE"},
						Package: []string{"dp3"},
					},
					{
						PathLicense: scanner.PathLicense{
							Path:    "ScannerTest5\\Module2\\antlib\\DependedProject2.jar",
							License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1"},
						Package: []string{"dp2"},
					},
				},
			},
			Modules: []scanner.UnitDependency{
				{
					Name: "ScannerTest5",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest5\\Module1\\lib\\DependedProject3.jar",
								License: "MICROSOFT RECIPROCAL LICENSE"},
							Package: []string{"dp3"},
						},
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest5\\Module2\\antlib\\DependedProject2.jar",
								License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1"},
							Package: []string{"dp2"},
						},
					},
				},
				{
					Name: "ScannerTest5\\Module1",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest5\\Module1\\lib\\DependedProject3.jar",
								License: "MICROSOFT RECIPROCAL LICENSE"},
							Package: []string{"dp3"},
						},
					},
				},
				{
					Name: "ScannerTest5\\Module2",
					Dependencies: []scanner.JarPackage{
						{
							PathLicense: scanner.PathLicense{
								Path:    "ScannerTest5\\Module2\\antlib\\DependedProject2.jar",
								License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1"},
							Package: []string{"dp2"},
						},
					},
				},
			},
		},
		PomLicense:     nil,
		LocalConflicts: nil,
		ExternalConflicts: []scanner.ExternalConflict{
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\License.txt",
					License: "GENERAL PUBLIC LICENSE Version 3",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\Module1\\lib\\DependedProject3.jar",
					License: "MICROSOFT RECIPROCAL LICENSE",
				},
				Result: scanner.ConflictResult{Unknown: false, Pass: false, Message: "冲突"},
			},
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\License.txt",
					License: "GENERAL PUBLIC LICENSE Version 3",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\Module2\\antlib\\DependedProject2.jar",
					License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1",
				},
				Result: scanner.ConflictResult{Unknown: false, Pass: false, Message: "冲突"},
			},
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\Module1\\LICENSE",
					License: "LESSER GENERAL PUBLIC LICENSE Version 2.1",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\Module1\\lib\\DependedProject3.jar",
					License: "MICROSOFT RECIPROCAL LICENSE",
				},
				Result: scanner.ConflictResult{Unknown: false, Pass: false, Message: "冲突"},
			},
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\Module2\\LICENSE",
					License: "LESSER GENERAL PUBLIC LICENSE Version 2.1",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest5\\Module2\\antlib\\DependedProject2.jar",
					License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1",
				},
				Result: scanner.ConflictResult{Unknown: false, Pass: false, Message: "冲突"},
			},
		},
		PomConflicts:      nil,
		RecommendLicenses: []string{"Apache-2.0", "ECL-2.0"},
	},
}

var testCaseScannerTest6 = TestCase{
	name:     "ScannerTest6",
	filePath: "..\\testProject\\ScannerTest6.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{{
			Path: "ScannerTest6\\LICENSE", License: "GENERAL PUBLIC LICENSE Version 2",
		}},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{Name: "ScannerTest6", Dependencies: nil},
			Modules: []scanner.UnitDependency{{Name: "ScannerTest6", Dependencies: nil}},
		},
		PomLicense: []scanner.PomLicense{{
			XMLDependency: scanner.XMLDependency{
				GroupID: "com.google.code.gson", ArtifactID: "gson", Version: "2.8.9",
			},
			License: "Apache-2.0",
		}},
		LocalConflicts:    nil,
		ExternalConflicts: nil,
		PomConflicts:      nil,
		RecommendLicenses: []string{"LGPL-2.1-only", "LGPL-2.1-or-later", "LGPL-3.0-only", "GPL-2.0-or-later",
			"GPL-3.0-only", "AGPL-3.0-only", "Artistic-2.0", "CECILL-2.1", "EPL-1.0", "EPL-2.0", "EUPL-1.1",
			"EUPL-1.2", "MPL-2.0", "MS-RL", "OSL-3.0", "Vim", "CDDL-1.0", "CPAL-1.0", "CPL-1.0", "IPL-1.0",
			"LPL-1.02", "Nokia", "RPSL-1.0", "SISSL", "Sleepycat", "SPL-1.0", "php-3.01", "Apache-2.0", "ECL-2.0"},
	},
}
var testCaseScannerTest7 = TestCase{
	name:     "ScannerTest7",
	filePath: "..\\testProject\\ScannerTest7.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest7\\LICENSE.txt",
				License: "APACHE LICENSE Version 2.0",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name:         "ScannerTest7",
				Dependencies: nil,
			},
			Modules: []scanner.UnitDependency{
				{
					Name:         "ScannerTest7",
					Dependencies: nil,
				},
			},
		},
		PomLicense: []scanner.PomLicense{
			{
				XMLDependency: scanner.XMLDependency{
					GroupID:    "org.jboss.aop",
					ArtifactID: "jboss-aop",
					Version:    "2.2.2.GA",
				},
				License: "LGPL-2.1",
			},
		},
		LocalConflicts:    nil,
		ExternalConflicts: nil,
		PomConflicts: []scanner.PomConflict{
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest7\\LICENSE.txt",
					License: "APACHE LICENSE Version 2.0",
				},
				PomLicense: scanner.PomLicense{
					XMLDependency: scanner.XMLDependency{
						GroupID:    "org.jboss.aop",
						ArtifactID: "jboss-aop",
						Version:    "2.2.2.GA",
					},
					License: "LGPL-2.1",
				},
				Result: scanner.ConflictResult{Unknown: true, Pass: false, Message: "分析失败，许可证库不包含对应信息"},
			},
		},
		RecommendLicenses: nil,
	},
}
var testCaseScannerTest8 = TestCase{
	name:     "ScannerTest8",
	filePath: "..\\testProject\\ScannerTest8.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest8\\LICENSE",
				License: "LESSER GENERAL PUBLIC LICENSE Version 3",
			},
			{
				Path:    "ScannerTest8\\src\\main\\LICENSE",
				License: "GENERAL PUBLIC LICENSE Version 2",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name:         "ScannerTest8",
				Dependencies: nil,
			},
			Modules: []scanner.UnitDependency{
				{
					Name:         "ScannerTest8",
					Dependencies: nil,
				},
				{
					Name:         "ScannerTest8\\src\\main",
					Dependencies: nil,
				},
			},
		},
		PomLicense: []scanner.PomLicense{
			{
				XMLDependency: scanner.XMLDependency{
					GroupID:    "com.google.code.gson",
					ArtifactID: "gson",
					Version:    "2.8.9",
				},
				License: "Apache-2.0",
			},
		},
		LocalConflicts: []scanner.ExternalConflict{
			{
				MainLicense: scanner.PathLicense{
					Path:    "ScannerTest8\\LICENSE",
					License: "LESSER GENERAL PUBLIC LICENSE Version 3",
				},
				ExternalLicense: scanner.PathLicense{
					Path:    "ScannerTest8\\src\\main\\LICENSE",
					License: "GENERAL PUBLIC LICENSE Version 2",
				},
				Result: scanner.ConflictResult{
					Unknown: false,
					Pass:    false,
					Message: "冲突, 组合遵循GPL-2.0-or-later",
				},
			},
		},
		ExternalConflicts: nil,
		PomConflicts:      nil,
		RecommendLicenses: []string{
			"GPL-2.0-or-later",
			"GPL-3.0-only",
			"AGPL-3.0-only",
			"Apache-2.0",
			"ECL-2.0",
		},
	},
}

var testCaseScannerTest9 = TestCase{
	name:     "ScannerTest9",
	filePath: "..\\testProject\\ScannerTest9.zip",
	expectedResult: scanner.TaskResult{
		IsFinish:     true,
		ErrorMessage: "",
		Local: []scanner.PathLicense{
			{
				Path:    "ScannerTest9\\LICENSE",
				License: "GENERAL PUBLIC LICENSE Version 3",
			},
			{
				Path:    "ScannerTest9\\src\\main\\LICENSE",
				License: "LESSER GENERAL PUBLIC LICENSE Version 3",
			},
		},
		Dependency: scanner.AllModuleDependency{
			Project: scanner.UnitDependency{
				Name:         "ScannerTest9",
				Dependencies: nil,
			},
			Modules: []scanner.UnitDependency{
				{
					Name:         "ScannerTest9",
					Dependencies: nil,
				},
				{
					Name:         "ScannerTest9\\src\\main",
					Dependencies: nil,
				},
			},
		},
		PomLicense: []scanner.PomLicense{
			{
				XMLDependency: scanner.XMLDependency{
					GroupID:    "com.google.code.gson",
					ArtifactID: "gson",
					Version:    "2.8.9",
				},
				License: "Apache-2.0",
			},
		},
		LocalConflicts:    nil,
		ExternalConflicts: nil,
		PomConflicts:      nil,
		RecommendLicenses: []string{
			"LGPL-2.1-only",
			"LGPL-2.1-or-later",
			"LGPL-3.0-only",
			"GPL-3.0-only",
			"AGPL-3.0-only",
			"Artistic-2.0",
			"CECILL-2.1",
			"EPL-1.0",
			"EPL-2.0",
			"EUPL-1.1",
			"EUPL-1.2",
			"MPL-2.0",
			"MS-RL",
			"OSL-3.0",
			"Vim",
			"CDDL-1.0",
			"CPAL-1.0",
			"CPL-1.0",
			"IPL-1.0",
			"LPL-1.02",
			"Nokia",
			"RPSL-1.0",
			"SISSL",
			"Sleepycat",
			"SPL-1.0",
			"Apache-2.0",
			"ECL-2.0",
		},
	},
}

func TestSubmitScanTaskAndGetTaskResult(t *testing.T) {
	// 设置全局日志(输出接下来操作的错误信息)
	//logger.SetLoggerConfig()
	// 运行任务处理队列
	scanner.StartTaskSystem()
	// 初始化许可证分析
	scanner.InitializeAnalyzer()

	//将gin设置为测试模式
	gin.SetMode(gin.TestMode)

	//待测试的接口：获取验证码，获取验证邮件，创建用户
	router := gin.New()
	const urlTask = "/task"
	router.POST(urlTask, SubmitScanTask)
	router.GET(urlTask, GetTaskResult)

	testCases := []TestCase{
		testCaseScannerTest1,
		testCaseScannerTest2,
		testCaseScannerTest3,
		testCaseScannerTest4,
		testCaseScannerTest5,
		testCaseScannerTest6,
		testCaseScannerTest7,
		testCaseScannerTest8,
		testCaseScannerTest9,
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("TestCase%d_%s", i, testCase.name), func(t *testing.T) {
			convey.Convey("Test", t, func() {
				var response *httptest.ResponseRecorder
				var request *http.Request
				var err error
				//// 读取项目压缩包，构造并提交请求
				// 读取项目压缩包
				file, err := os.Open(testCase.filePath)
				convey.So(err, convey.ShouldBeNil)
				// 构造POST 请求multipart/form-data请求体
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				fileWriter, err := writer.CreateFormFile("file", filepath.Base(testCase.filePath))
				_, err = io.Copy(fileWriter, file)
				convey.So(err, convey.ShouldBeNil)

				err = writer.Close()
				convey.So(err, convey.ShouldBeNil)
				// 提交项目分析任务
				response = httptest.NewRecorder()
				request, _ = http.NewRequest("POST", urlTask, body)
				request.Header.Set("Content-Type", writer.FormDataContentType())
				router.ServeHTTP(response, request)

				fmt.Printf("submit analyze task\n")

				// !！测试过程中输出内容较长时，不论是否以换行\n结束，不论使用log.Printf fmt.Printf 还是 t.Logf
				// 均可能导致Goland（或GoLand？）出现"No tests were run"，无法正常测试

				//fmt.Printf("response %#v\n", response)
				//fmt.Printf("response.Body.String() %#v\n", response.Body.String())
				convey.So(response.Code, convey.ShouldEqual, http.StatusOK)
				convey.So(response.Body.String(), convey.ShouldEqual, strconv.Itoa(i+1))

				//// 获取项目分析结果
				for {
					response = httptest.NewRecorder()
					request, _ = http.NewRequest("GET",
						fmt.Sprintf("%s?id=%d", urlTask, i+1), nil)
					router.ServeHTTP(response, request)

					fmt.Printf("get task analyze result\n")
					//fmt.Printf("response %#v\n", response)
					//fmt.Printf("response.Body.String() %#v\n", response.Body.String())
					convey.So(response.Code, convey.ShouldEqual, http.StatusOK)

					fmt.Printf("json.Unmarshal\n")
					var result scanner.TaskResult
					convey.So(json.Unmarshal(response.Body.Bytes(), &result), convey.ShouldBeNil)
					if result.IsFinish {
						actual := fmt.Sprintf("%+v", result)
						expected := fmt.Sprintf("%+v", testCase.expectedResult)
						convey.So(actual, convey.ShouldEqual, expected)
						break
					}
					time.Sleep(time.Second)
				}
			})
		})
	}
}
