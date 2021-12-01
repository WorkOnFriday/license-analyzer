package scanner

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var licenses []string

const (
	GPL       = "GENERAL PUBLIC LICENSE"
	LGPL      = "LESSER GENERAL PUBLIC LICENSE"
	MIT       = "MIT LICENSE"
	APACHE    = "APACHE LICENSE"
	COPYRIGHT = "COPYRIGHT"
	BSD3      = "BSD 3-CLAUSE LICENSE"
	EUPL      = "EUROPEAN UNION PUBLIC LICENCE"
	MsPl      = "MICROSOFT PUBLIC LICENSE"
	MsRl      = "MICROSOFT RECIPROCAL LICENSE"

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
	licenses = append(licenses, MIT, APACHE, COPYRIGHT, LGPL, GPL, BSD3, EUPL, MsRl, MsPl)

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
		//fmt.Println(result, "----------------------------------")
		re1 := regexp.MustCompile(".*\\.[jarJAR]{3}")
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
		fmt.Printf("Scan result: %v\n", string(marshal))
		fmt.Println("all external modules: ", findAllExternalModule(external))
		fmt.Println("all local modules: ", findAllLocalModule(local))
		dependency, _ := json.Marshal(dependencyAnalyze(findAllExternalModule(external), local))
		fmt.Println("dependency result: ", string(dependency))

		defer os.RemoveAll(ProjectTmp)
		defer os.RemoveAll(DIR)
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
					strings.EqualFold(license, GPL) || strings.EqualFold(license, EUPL) {
					scanner.Scan()
					version := scanner.Text()
					for strings.TrimSpace(version) == "" {
						scanner.Scan()
						version = scanner.Text()
					}
					version = strings.TrimLeft(version, " ")
					version = strings.Split(version, ",")[0]
					//fmt.Println(license, "-------------", version)
					return license + " " + version
				} else {
					//fmt.Println(license)
					return license
				}
			}
		}
	}
	return "other license"
}

func deepScan(filePath string, result *map[string]interface{}) {
	if info, err := os.Stat(filePath); err == nil {
		if !info.IsDir() || strings.EqualFold(info.Name(), "out") {
			name := strings.Split(strings.ToUpper(info.Name()), ".")
			if name[0] == "LICENSE" || name[0] == "COPYING" || name[len(name)-1] == "LICENSE" {
				(*result)[filePath] = scan(filePath)
			} else if strings.EqualFold(name[len(name)-1], "jar") { // judge equal with case-insensitivity
				tmp := filepath.Join(DIR, filePath)
				os.MkdirAll(tmp, os.ModePerm)
				zipDeCompress(filePath, tmp)
				deepScan(tmp, result)
				//defer os.RemoveAll(tmp)
			}
			return
		} else {
			//fmt.Println(filePath)
			f, _ := os.Open(filePath)
			defer f.Close()
			names, _ := f.Readdirnames(0)
			for _, name := range names {
				newPath := filepath.Join(filePath, name)
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

func dependencyAnalyze(externalModules map[string][]string, local []map[string]interface{}) map[string][]string {
	dependency := make(map[string][]string)
	mainModule := ""
	for _, tmpMap := range local {
		for licensePath := range tmpMap {
			var tmpArr []string
			var tmpArr2 []string
			modulePath := filepath.Dir(licensePath)
			if !strings.Contains(modulePath, "/") && len(local) > 1 {
				mainModule = modulePath
				continue
			}
			javaScan(filepath.Join(ProjectTmp, modulePath), &tmpArr)
			//fmt.Println(tmpArr, "------------")
			for _, s := range tmpArr {
				tmpArr2 = append(tmpArr2, removeModuleSuffix(s))
			}
			for _, module := range removeRepeatElement(tmpArr2) {
				if !strings.HasPrefix(module, "java.") {
					for externalPath, external := range externalModules {
						for _, v := range external {
							if moduleEqual(3, strings.ReplaceAll(v, "/", "."), module) {
								dependency[modulePath] = append(dependency[modulePath], externalPath)
							}
						}
					}
				}
			}
		}
	}
	// finally, remove repeat module
	for k, v := range dependency {
		dependency[k] = removeRepeatElement(v)
	}

	if mainModule != "" {
		var mainModuleArr []string
		for _, v := range dependency {
			mainModuleArr = append(mainModuleArr, v...)
		}
		mainModuleArr = removeRepeatElement(mainModuleArr)
		dependency[mainModule] = mainModuleArr
	}
	return dependency
}

/**
scan all java file to find import modules
@param: filePath -> module path
@param: result -> an array to store result
*/
func javaScan(filePath string, result *[]string) {
	if info, err := os.Stat(filePath); err == nil {
		if !info.IsDir() {
			if strings.HasSuffix(strings.ToLower(info.Name()), ".java") {
				//fmt.Println(info.Name(), "-----------")
				file, err := os.Open(filepath.Join(filePath))
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
					line := strings.Trim(scanner.Text(), " ")
					if !strings.HasPrefix(line, "import ") {
						continue
					}
					for {
						if !strings.HasPrefix(line, "import ") {
							break
						}
						*result = append(*result, strings.Trim(line[6:], " "))
						scanner.Scan()
						line = strings.Trim(scanner.Text(), " ")
					}
				}
			}
			return
		} else {
			f, _ := os.Open(filePath)
			defer f.Close()
			names, _ := f.Readdirnames(0)
			for _, name := range names {
				newPath := filepath.Join(filePath, name)
				javaScan(newPath, result)
			}
		}
	}
}

/**
find all external modules from array of module path

*/
func findAllExternalModule(tmp []map[string]interface{}) map[string][]string {
	result := make(map[string][]string)
	arr := make([]string, 0)
	for _, tmpMap := range tmp {
		for key := range tmpMap {
			filePath := filepath.Join(DIR, ProjectTmp, key)
			if !strings.HasSuffix(key, ".jar") {
				filePath = filepath.Dir(filePath)
			}
			files, err := os.ReadDir(filePath)
			if err != nil {
				log.Fatalln(err)
			}
			for _, file := range files {
				if !file.IsDir() || strings.EqualFold(file.Name(), "META-INF") {
					continue
				}
				tmpDir := filepath.Join(filePath, file.Name())
				f, err := os.ReadDir(tmpDir)
				if err != nil {
					log.Fatalln(err)
				}
				for {
					if len(f) != 1 || (len(f) == 1 && !f[0].IsDir()) {
						result[key] = append(arr, strings.Split(tmpDir, ".jar/")[1])
						break
					}
					tmpDir = filepath.Join(tmpDir, f[0].Name())
					f, err = os.ReadDir(tmpDir)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}
	return result
}

func findAllLocalModule(tmp []map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, temMap := range tmp {
		for key, value := range temMap {
			module := filepath.Base(filepath.Dir(key))
			if strings.EqualFold(module, ".") {
				module = filepath.Dir(key)
			}
			result[module] = value
		}
	}
	return result
}

/**
remove repeat element from []string
*/
func removeRepeatElement(list []string) []string {
	flag := make(map[string]bool)
	tmp := make([]string, 0)
	index := 0
	for _, v := range list {
		_, ok := flag[v]
		if !ok {
			tmp = append(tmp, v)
			flag[v] = true
		}
		index++
	}
	return tmp
}

/**
judge s1, s2 module equal or not
@param: depth -> the max depth
*/
func moduleEqual(depth int, s1, s2 string) bool {
	t1 := strings.Split(s1, ".")
	t2 := strings.Split(s2, ".")
	for i := 0; i < depth; i++ {
		if len(t1) == i || len(t2) == i {
			break
		}
		if !strings.EqualFold(t1[i], t2[i]) {
			return false
		}
	}
	return true
}

func removeModuleSuffix(module string) string {
	re := regexp.MustCompile("\\..[^.]+;")
	suffix := re.FindStringIndex(module)
	if len(suffix) == 0 {
		return module
	}
	return module[:suffix[0]]
}
