package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
	"io"
	"license-analyzer/logger"
	"license-analyzer/scanner"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSubmitScanTaskAndGetTaskResult(t *testing.T) {
	// 设置全局日志(输出接下来操作的错误信息)
	logger.SetLoggerConfig()
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

	t.Run("TestSubmitScanTaskAndGetTaskResult", func(t *testing.T) {
		convey.Convey("Test", t, func() {
			var response *httptest.ResponseRecorder
			var request *http.Request
			var err error
			//// 读取项目压缩包，构造并提交请求
			// 读取项目压缩包
			filePath := "..\\testProject\\ScannerTest6.zip"
			file, err := os.Open(filePath)
			convey.So(err, convey.ShouldBeNil)
			// 构造POST 请求multipart/form-data请求体
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			fileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
			_, err = io.Copy(fileWriter, file)
			convey.So(err, convey.ShouldBeNil)
			err = writer.Close()
			convey.So(err, convey.ShouldBeNil)
			// 提交项目分析任务
			response = httptest.NewRecorder()
			request, _ = http.NewRequest("POST", urlTask, body)
			request.Header.Set("Content-Type", writer.FormDataContentType())
			router.ServeHTTP(response, request)
			log.Printf("submit analyze task\n")
			log.Printf("response %#v\n", response)
			log.Printf("response.Body.String() %#v\n", response.Body.String())
			convey.So(response.Code, convey.ShouldEqual, http.StatusOK)
			convey.So(response.Body.String(), convey.ShouldEqual, "1")

			//// 获取项目分析结果
			var expectedResult = scanner.TaskResult{
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
			}
			for {
				response = httptest.NewRecorder()
				request, _ = http.NewRequest("GET", fmt.Sprintf("%s?id=%d", urlTask, 1), nil)
				router.ServeHTTP(response, request)
				log.Printf("get task analyze result\n")
				log.Printf("response %#v\n", response)
				log.Printf("response.Body.String() %#v\n", response.Body.String())
				convey.So(response.Code, convey.ShouldEqual, http.StatusOK)
				var result scanner.TaskResult
				convey.So(json.Unmarshal(response.Body.Bytes(), &result), convey.ShouldBeNil)
				if result.IsFinish {
					convey.So(fmt.Sprintf("%+v", result), convey.ShouldEqual,
						fmt.Sprintf("%+v", expectedResult))
					break
				}
				time.Sleep(time.Second)
			}
		})
	})
}
