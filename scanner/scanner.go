package scanner

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var licenses []string

const (
	GPL        = "GENERAL PUBLIC LICENSE"
	LGPL       = "LESSER GENERAL PUBLIC LICENSE"
	MIT        = "MIT LICENSE"
	APACHE     = "APACHE LICENSE"
	COPYRIGHT  = "COPYRIGHT"
	BSD3       = "BSD 3-CLAUSE LICENSE"
	DIR        = "./tmp"
	ProjectTmp = "./project_tmp"
)

func main() {
	var param = ""
	//for _, arg := range os.Args{
	//	fmt.Println(arg)
	//}
	if len(os.Args) < 2 {
		fmt.Println("No file specify")
		os.Exit(0)
	} else {
		param = os.Args[1]
	}
	licenses = append(licenses, MIT, APACHE, COPYRIGHT, GPL, LGPL, BSD3)

	suffix := strings.ToUpper(filepath.Ext(param))
	//fmt.Println(suffix)

	if suffix == ".ZIP" || suffix == ".JAR" { //if file type is .zip, make tmp dir
		if _, err := os.ReadDir(ProjectTmp); err != nil {
			if err := os.MkdirAll(ProjectTmp, os.ModePerm); err != nil {
				log.Fatalln(err)
			}
		}

		result := make(map[string]interface{})
		external := make([]map[string]interface{}, 0)
		local := make([]map[string]interface{}, 0)

		zipDeCompress(param, ProjectTmp)

		//result = shallowScan(ProjectTmp)
		deepScan(ProjectTmp, &result)

		re1 := regexp.MustCompile(".*?\\.[jarJAR]*")
		for k, v := range result {
			if strings.HasPrefix(k, "tmp") {
				tk := re1.FindString(k)[len("tmp/project_tmp/"):]
				external = append(external, map[string]interface{}{tk: v})
			} else {
				local = append(local, map[string]interface{}{k[len("project_tmp/"):]: v})
			}
		}

		t := map[string]interface{}{"external": external, "local": local}

		marshal, err := json.Marshal(t)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Scan result: %v", string(marshal))

		defer os.RemoveAll(ProjectTmp)

	} else if suffix == "" || suffix == ".TXT" || suffix == ".LICENSE" {
		fmt.Println(param, scan(param))
	} else {
		fmt.Println("Not support file")
		os.Exit(0)
	}
}

func zipDeCompress(fileName string, destDir string) {
	zipReader, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		fPath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				log.Fatalln(err)
			}
			inFile, err := f.Open()
			if err != nil {
				log.Fatalln(err)
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatalln(err)
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func scan(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		lineText := strings.ToUpper(line)
		for _, license := range licenses {
			if strings.Contains(lineText, license) {
				if strings.EqualFold(license, APACHE) || strings.EqualFold(license, LGPL) ||
					strings.EqualFold(license, GPL) {
					scanner.Scan()
					version := scanner.Text()
					version = strings.TrimLeft(version, " ")
					version = strings.Split(version, ",")[0]
					//fmt.Println(license, version)
					return license + " " + version
				} else {
					//fmt.Println(license)
					return license
				}
			}
		}
	}
	//fmt.Println("other license")
	return "other license"
}

func deepScan(filePath string, result *map[string]interface{}) {
	if info, err := os.Stat(filePath); err == nil {
		if !info.IsDir() {
			name := strings.Split(strings.ToUpper(info.Name()), ".")
			if name[0] == "LICENSE" || name[0] == "COPYING" || name[len(name)-1] == "LICENSE" {
				(*result)[filePath] = scan(filePath)
			} else if strings.EqualFold(name[len(name)-1], "jar") { // judge equal with case-insensitivity
				tmp := filepath.Join(DIR, filePath)
				os.MkdirAll(tmp, os.ModePerm)
				zipDeCompress(filePath, tmp)
				deepScan(tmp, result)
				defer os.RemoveAll(tmp)
			}
			return
		} else {
			//fmt.Println(filePath)
			f, _ := os.Open(filePath)
			defer f.Close()
			names, _ := f.Readdirnames(0)
			for _, name := range names {
				newPath := path.Join(filePath, name)
				deepScan(newPath, result)
			}
		}
	}
}

func shallowScan(path string) map[string]interface{} {
	//If the dependent module only needs to scan the outermost layer, use this
	result := make(map[string]interface{})
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var tmpDir = path
	if len(files) == 1 && files[0].IsDir() {
		tmpDir = filepath.Join(path, files[0].Name())
		files, err = os.ReadDir(tmpDir)
	}
	for _, file := range files {
		name := strings.Split(strings.ToUpper(file.Name()), ".")
		if name[0] == "LICENSE" || name[0] == "COPYING" || name[len(name)-1] == "LICENSE" {
			//fileName, _ := filepath.Abs(file.Name())
			fileName := filepath.Join(tmpDir, file.Name())
			result[fileName] = scan(fileName)
		}
	}
	return result
}
