package scanner

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCase(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		convey.Convey("Test1", t, func() {
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

				t := map[string]interface{}{"external": external, "local": local}

				marshal, err := json.Marshal(t)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("Scan result: %v\n", string(marshal))
				fmt.Println("all external modules: ", findAllExternalModule(external))
				fmt.Println("all local modules: ", findAllLocalModule(local))
				dependency, _ := json.Marshal(dependencyAnalyze(findAllExternalModule(external), local))
				fmt.Println("dependency result: ", string(dependency))

				p := PomScan("./pom.xml")
				fmt.Println(p)

			} else if suffix == "" || suffix == ".TXT" || suffix == ".LICENSE" {
				fmt.Println(param, ScanFile(param))
			} else {
				fmt.Println("Not support file")
				os.Exit(0)
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
