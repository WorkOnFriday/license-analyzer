package scanner

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanPackage(t *testing.T) {

	t.Run("TestScanPackage_ScannerTest4", func(t *testing.T) {

		convey.Convey("Test4", t, func() {
			var param = "D:\\gopath\\src\\license-analyzer\\testProject\\ScannerTest4.zip"

			suffix := strings.ToUpper(filepath.Ext(param))

			if suffix == ".ZIP" || suffix == ".JAR" { //if file type is .zip, make tmp dir
				if _, err := os.ReadDir(ProjectTmp); err != nil {
					if err := os.MkdirAll(ProjectTmp, os.ModePerm); err != nil {
						log.Fatalln(err)
					}
					defer os.RemoveAll(ProjectTmp)
				}

				external := make([]PathLicense, 0)
				local := make([]PathLicense, 0)

				zipDecompress(param, ProjectTmp)

				//result = shallowScan(ProjectTmp)
				deepScan(ProjectTmp, "", true, &local, &external)
				defer os.RemoveAll(DIR)

				expectExternal := []PathLicense{
					{Path: "ScannerTest4\\antlib\\DependedProject2.jar", License: "EUROPEAN UNION PUBLIC LICENCE V. 1.1"},
					{Path: "ScannerTest4\\antlib2\\DependedProject3.jar", License: "MICROSOFT RECIPROCAL LICENSE"},
					{Path: "ScannerTest4\\lib\\DependedProject.jar", License: "GENERAL PUBLIC LICENSE Version 2"},
				}
				fmt.Printf("Scan result external: %+v\n", external)
				convey.So(fmt.Sprintf("%+v", external), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectExternal))

				expectLocal := []PathLicense{
					{Path: "ScannerTest4\\LICENSE", License: "GENERAL PUBLIC LICENSE Version 2"},
				}
				fmt.Printf("Scan result local: %+v\n", local)
				convey.So(fmt.Sprintf("%+v", local), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectLocal))

				allExternalModule := findAllExternalModule(external)
				expectAllExternalModule := []JarPackages{
					{JarPath: "ScannerTest4\\antlib\\DependedProject2.jar", Package: []string{"dp2"}},
					{JarPath: "ScannerTest4\\antlib2\\DependedProject3.jar", Package: []string{"dp3"}},
					{JarPath: "ScannerTest4\\lib\\DependedProject.jar", Package: []string{"pri"}},
				}
				fmt.Printf("all external modules: %+v\n", allExternalModule)
				convey.So(fmt.Sprintf("%+v", allExternalModule), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectAllExternalModule))

				allLocalModule := findAllLocalModule(local)
				expectAllLocalModule := []ModuleLicense{
					{Module: "ScannerTest4", License: "GENERAL PUBLIC LICENSE Version 2"},
				}
				fmt.Printf("all local modules: %+v\n", allLocalModule)
				convey.So(fmt.Sprintf("%+v", allLocalModule), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectAllLocalModule))

				dependency := dependencyAnalyze(findAllExternalModule(external), local)
				expectDependency := AllModuleDependency{Project: ModuleDependency{
					Module: "ScannerTest4", Dependencies: []string{
						"ScannerTest4\\lib\\DependedProject.jar",
						"ScannerTest4\\antlib\\DependedProject2.jar",
						"ScannerTest4\\antlib2\\DependedProject3.jar",
					}},
					Modules: []ModuleDependency{
						{Module: "ScannerTest4", Dependencies: []string{
							"ScannerTest4\\lib\\DependedProject.jar",
							"ScannerTest4\\antlib\\DependedProject2.jar",
							"ScannerTest4\\antlib2\\DependedProject3.jar",
						}},
					}}
				fmt.Printf("dependency result: %+v\n", dependency)
				convey.So(fmt.Sprintf("%+v", dependency), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectDependency))

				p := PomScan(filepath.Join(ProjectTmp, "ScannerTest4\\pom.xml"))
				expectP := XMLProject{}
				fmt.Printf("pom result: %+v\n", p)
				convey.So(fmt.Sprintf("%+v", p), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectP))
			}
		})
	})

	t.Run("TestScanPackage_ScannerTest6", func(t *testing.T) {
		convey.Convey("Test6", t, func() {
			var param = "D:\\gopath\\src\\license-analyzer\\testProject\\ScannerTest6.zip"

			suffix := strings.ToUpper(filepath.Ext(param))

			if suffix == ".ZIP" || suffix == ".JAR" { //if file type is .zip, make tmp dir
				if _, err := os.ReadDir(ProjectTmp); err != nil {
					if err := os.MkdirAll(ProjectTmp, os.ModePerm); err != nil {
						log.Fatalln(err)
					}
					defer os.RemoveAll(ProjectTmp)
				}

				external := make([]PathLicense, 0)
				local := make([]PathLicense, 0)

				zipDecompress(param, ProjectTmp)

				//result = shallowScan(ProjectTmp)
				deepScan(ProjectTmp, "", true, &local, &external)
				defer os.RemoveAll(DIR)

				var expectExternal []PathLicense
				fmt.Printf("Scan result external: %+v\n", external)
				convey.So(fmt.Sprintf("%+v", external), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectExternal))

				expectLocal := []PathLicense{
					{Path: "ScannerTest6\\LICENSE", License: "GENERAL PUBLIC LICENSE Version 2"},
				}
				fmt.Printf("Scan result local: %+v\n", local)
				convey.So(fmt.Sprintf("%+v", local), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectLocal))

				allExternalModule := findAllExternalModule(external)
				var expectAllExternalModule []JarPackages
				fmt.Printf("all external modules: %+v\n", allExternalModule)
				convey.So(fmt.Sprintf("%+v", allExternalModule), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectAllExternalModule))

				allLocalModule := findAllLocalModule(local)
				expectAllLocalModule := []ModuleLicense{
					{Module: "ScannerTest6", License: "GENERAL PUBLIC LICENSE Version 2"},
				}
				fmt.Printf("all local modules: %+v\n", allLocalModule)
				convey.So(fmt.Sprintf("%+v", allLocalModule), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectAllLocalModule))

				dependency := dependencyAnalyze(findAllExternalModule(external), local)
				expectDependency := AllModuleDependency{
					Project: ModuleDependency{Module: "ScannerTest6", Dependencies: []string{}},
					Modules: []ModuleDependency{
						{Module: "ScannerTest6", Dependencies: []string{}},
					}}
				fmt.Printf("dependency result: %+v\n", dependency)
				convey.So(fmt.Sprintf("%+v", dependency), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectDependency))

				p := PomScan(filepath.Join(ProjectTmp, "ScannerTest6\\pom.xml"))
				expectP := XMLProject{Dependencies: XMLDependencies{Dependencies: []XMLDependency{
					{GroupID: "com.google.code.gson",
						ArtifactID: "gson",
						Version:    "2.8.9"},
				}}}
				fmt.Printf("pom result: %+v\n", p)
				convey.So(fmt.Sprintf("%+v", p), convey.ShouldEqual,
					fmt.Sprintf("%+v", expectP))
			}
		})

	})
}

func TestPomScan(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
			project := PomScan("D:\\gopath\\src\\license-analyzer\\testFile\\pomDependency.xml")
			fmt.Printf("%+v\n", project)
			project2 := PomScan("D:\\gopath\\src\\license-analyzer\\testFile\\pomModule.xml")
			fmt.Printf("%+v\n", project2)
		})
	})
}

func TestScanFile(t *testing.T) {
	t.Run("TestScanFile", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
			filePath := "D:\\gopath\\src\\license-analyzer\\testFile\\EUPL.txt"
			license := ScanFile(filePath)
			fmt.Printf("license %+v\n", license)
			convey.So(license, convey.ShouldEqual, "EUROPEAN UNION PUBLIC LICENCE V. 1.1")
		})
	})
}
