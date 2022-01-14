package scanner

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

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

var licenses = []string{MIT, APACHE, COPYRIGHT, LGPL, GPL, BSD3, EUPL, MsRl, MsPl}

type PathLicense struct {
	Path    string
	License string
}

type AllPathLicense struct {
	Local    []PathLicense
	External []PathLicense
}

func zipDecompress(filePath string, destDir string) {
	// 打开压缩文件
	zipReader, err := zip.OpenReader(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = zipReader.Close(); err != nil {
			panic(err)
		}
	}()
	// 扫描所有文件和文件夹
	for _, f := range zipReader.File {
		// 顺序似乎总是先文件夹，后其中的文件，但我不确定
		fPath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			// 不创建空文件夹
			continue
		}
		// 创建对应的外部文件所在目录
		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			panic(err)
		}
		// 将压缩包中文件复制到外部指定位置
		func() {
			inFile, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer func() {
				if err = inFile.Close(); err != nil {
					panic(err)
				}
			}()

			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				panic(err)
			}
			defer func() {
				if err := outFile.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				panic(err)
			}
			// 函数结束时关闭Close()
		}()
	}
}

// ScanPackage 解压并扫描整个zip或jar包
func ScanPackage(zipFilePath string) AllPathLicense {
	// 目录已存在返回nil
	if err := os.MkdirAll(ProjectTmp, os.ModePerm); err != nil {
		logrus.Panic(err)
		panic(err)
	}
	defer os.RemoveAll(ProjectTmp)

	zipDecompress(zipFilePath, ProjectTmp)

	local := make([]PathLicense, 0)
	external := make([]PathLicense, 0)
	deepScan(ProjectTmp, "", true, &local, &external)
	defer os.RemoveAll(DIR)

	logrus.Debug("local ", local)
	logrus.Debug("external ", external)

	return AllPathLicense{Local: local, External: external}
}

// ScanFile 扫描证书文本文件内容
func ScanFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		lineUpper := strings.ToUpper(line)
		for _, license := range licenses {
			if !strings.Contains(lineUpper, license) {
				// 此行不含此证书
				continue
			}
			if !(strings.EqualFold(license, APACHE) || strings.EqualFold(license, LGPL) ||
				strings.EqualFold(license, GPL) || strings.EqualFold(license, EUPL)) {
				// 证书不涉及版本号
				return license
			}
			// 版本号形如"  Version 3, 29 June 2007"
			scanner.Scan()
			version := scanner.Text()
			for strings.TrimSpace(version) == "" {
				scanner.Scan()
				version = scanner.Text()
			}
			version = strings.TrimLeft(version, " ")
			version = strings.Split(version, ",")[0]
			return license + " " + version
		}
	}
	return "other license"
}

// 递归搜索所有目录的证书文件
// local: real path and license
// external: jar path and license
func deepScan(realFilePath, resultPath string, isLocal bool, local, external *[]PathLicense) {
	info, err := os.Stat(realFilePath)
	if err != nil {
		panic(err)
	}
	// 跳过out目录
	if strings.EqualFold(info.Name(), "out") {
		return
	}
	// 对于目录，进入目录搜索
	if info.IsDir() {
		f, _ := os.Open(realFilePath)
		defer func() {
			if err = f.Close(); err != nil {
				panic(err)
			}
		}()
		// 获取所有文件夹名字
		names, _ := f.Readdirnames(0)
		for _, name := range names {
			newRealPath := filepath.Join(realFilePath, name)
			newResultPath := filepath.Join(resultPath, name)
			if !isLocal {
				newResultPath = resultPath
			}
			deepScan(newRealPath, newResultPath, isLocal, local, external)
		}
		return
	}
	// 文件
	name := strings.Split(strings.ToUpper(info.Name()), ".")
	// 文本证书文件
	if name[0] == "LICENSE" || name[0] == "COPYING" || name[len(name)-1] == "LICENSE" {
		license := ScanFile(realFilePath)
		if isLocal {
			*local = append(*local, PathLicense{Path: resultPath, License: license})
		} else {
			*external = append(*external, PathLicense{Path: resultPath, License: license})
		}
		return
	}
	// jar包
	if strings.EqualFold(name[len(name)-1], "jar") { // judge equal with case-insensitivity
		newRealPath := filepath.Join(DIR, realFilePath)
		newResultPath := resultPath
		os.MkdirAll(newRealPath, os.ModePerm)
		zipDecompress(realFilePath, newRealPath)
		deepScan(newRealPath, newResultPath, false, local, external)
	}
}

type UnitDependency struct {
	Name         string
	Dependencies []JarPackage
}
type AllModuleDependency struct {
	Project UnitDependency
	Modules []UnitDependency
}

func dependencyAnalyze(externalModules []JarPackage, local []PathLicense) (all AllModuleDependency) {
	// 字符串去重
	uniqueString := func(strings []string) (result []string) {
		existMap := make(map[string]bool)
		for _, str := range strings {
			_, ok := existMap[str]
			if !ok {
				result = append(result, str)
				existMap[str] = true
			}
		}
		return
	}
	// jar包信息按路径去重
	uniqueJarPackage := func(jars []JarPackage) (result []JarPackage) {
		existMap := make(map[string]bool)
		for _, jar := range jars {
			_, ok := existMap[jar.Path]
			if !ok {
				result = append(result, jar)
				existMap[jar.Path] = true
			}
		}
		return
	}

	projectName := ""
	for _, nowPathLicense := range local {
		logrus.Debug("path: ", nowPathLicense.Path, " license: ", nowPathLicense.License)

		modulePath := filepath.Dir(nowPathLicense.Path)
		logrus.Debug("modulePath: ", modulePath)
		// 保存项目名
		splitToken := "\\"
		if runtime.GOOS == "linux" {
			splitToken = "/"
		}
		if !strings.Contains(modulePath, splitToken) {
			projectName = modulePath
		}
		// 扫描此模块的依赖
		var importArr []string
		javaScan(filepath.Join(ProjectTmp, modulePath), &importArr)
		logrus.Debug("import arr: ", importArr)

		for i := range importArr {
			importArr[i] = removeModuleSuffix(importArr[i])
		}
		importArr = uniqueString(importArr)
		logrus.Debug("import arr: ", importArr)
		// 此模块非java依赖对应的外部依赖
		var md UnitDependency
		md.Name = modulePath
		for _, module := range importArr {
			if strings.HasPrefix(module, "java.") {
				continue
			}
			// 找每个import依赖的jar包
			func() {
				for _, jarPackage := range externalModules {
					for _, pkg := range jarPackage.Package {
						if moduleEqualInDepthLevel(3, strings.ReplaceAll(pkg, "/", "."), module) {
							md.Dependencies = append(md.Dependencies, jarPackage)
							return // 已找到
						}
					}
				}
			}()
		}
		// 对被依赖的jar包去重
		md.Dependencies = uniqueJarPackage(md.Dependencies)
		all.Modules = append(all.Modules, md)
	}
	// 合并形成项目对jar包的依赖信息
	var projectDependencies []JarPackage
	for _, md := range all.Modules {
		projectDependencies = append(projectDependencies, md.Dependencies...)
	}
	projectDependencies = uniqueJarPackage(projectDependencies)
	all.Project = UnitDependency{Name: projectName, Dependencies: projectDependencies}
	return
}

/**
scan all java file to find import modules
@param: filePath -> module path
@param: result -> an array to store result
*/
func javaScan(filePath string, result *[]string) {
	info, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	// 对于目录
	if info.IsDir() {
		f, _ := os.Open(filePath)
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		names, _ := f.Readdirnames(0)
		for _, name := range names {
			newPath := filepath.Join(filePath, name)
			javaScan(newPath, result)
		}
		return
	}
	// 对于java文件
	if strings.HasSuffix(strings.ToLower(info.Name()), ".java") {
		//fmt.Println(info.Name(), "-----------")
		file, err := os.Open(filepath.Join(filePath))
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.Trim(scanner.Text(), " ")
			if !strings.HasPrefix(line, "import ") {
				continue
			}
			*result = append(*result, strings.Trim(line[6:], " "))
		}
	}
	return
}

type JarPackage struct {
	PathLicense
	Package []string
}

/**
find all external modules from array of module path
return package name in jar
*/
func findAllExternalModule(external []PathLicense) (result []JarPackage) {
	for _, jar := range external {
		var jarPackages JarPackage
		jarPackages.PathLicense = jar
		filePath := filepath.Join(DIR, ProjectTmp, jar.Path)
		/* 应该不存在此情况
		if !strings.HasSuffix(jarPath, ".jar") {
			filePath = filepath.Dir(filePath)
		}*/
		files, err := os.ReadDir(filePath)
		if err != nil {
			log.Fatalln(err)
		}
		for _, file := range files {
			if !file.IsDir() || strings.EqualFold(file.Name(), "META-INF") {
				continue
			}
			fileFullPath := filepath.Join(filePath, file.Name())
			f, err := os.ReadDir(fileFullPath)
			if err != nil {
				log.Fatalln(err)
			}
			for {
				if !(len(f) == 1 && f[0].IsDir()) {
					// 不只有一个文件 或 只有一个非目录文件（即搜到了包中）
					break
				}
				fileFullPath = filepath.Join(fileFullPath, f[0].Name())
				f, err = os.ReadDir(fileFullPath)
				if err != nil {
					log.Fatalln(err)
				}
			}
			splitToken := ".jar\\"
			if runtime.GOOS == "linux" {
				splitToken = ".jar/"
			}
			jarPackages.Package = append(jarPackages.Package, strings.Split(fileFullPath, splitToken)[1])
			result = append(result, jarPackages)
		}
	}
	return result
}

type ModuleLicense struct {
	Module  string
	License string
}

func findAllLocalModule(local []PathLicense) (result []ModuleLicense) {
	for _, nowPathLicense := range local {
		path := nowPathLicense.Path
		license := nowPathLicense.License
		// 项目模块的完整路径
		module := filepath.Dir(path)
		result = append(result, ModuleLicense{Module: module, License: license})
	}
	return
}

/**
judge s1, s2 module equal or not (only in first "depth" level)
@param: depth -> the max depth
*/
func moduleEqualInDepthLevel(depth int, s1, s2 string) bool {
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

// removeModuleSuffix if only one level, do nothing
func removeModuleSuffix(module string) string {
	re := regexp.MustCompile("\\..[^.]+;")
	suffix := re.FindStringIndex(module)
	if len(suffix) == 0 {
		return module
	}
	return module[:suffix[0]]
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="...">
  <modules>
    <module>zstd-proxy</module>
  </modules>
  <dependencies>
    <dependency>
      <groupId>org.neo4j</groupId>
      <artifactId>neo4j-kernel</artifactId>
      <version>${project.version}</version>
    </dependency>
    <dependency>
      <groupId>org.junit.jupiter</groupId>
      <artifactId>junit-jupiter</artifactId>
    </dependency>
  </dependencies>
</project>
*/

type XMLModules struct {
	Modules []string `xml:"module"`
}

type XMLDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type XMLDependencies struct {
	Dependencies []XMLDependency `xml:"dependency"`
}

type XMLProject struct {
	//XMLName      xml.Name        `xml:"project"`
	Modules      XMLModules      `xml:"modules"`
	Dependencies XMLDependencies `xml:"dependencies"`
}

func PomScan(fileName string) (project XMLProject) {
	// 支持文件为空
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return XMLProject{}
	}
	// 打开文件
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		panic(err)
	}
	// 解析XML
	if err = xml.Unmarshal(b, &project); err != nil {
		panic(err)
	}
	// 设置默认值（暂时没找到更好的做法）
	for i := range project.Dependencies.Dependencies {
		if project.Dependencies.Dependencies[i].Version == "" {
			project.Dependencies.Dependencies[i].Version = "latest"
		}
	}
	return
}
