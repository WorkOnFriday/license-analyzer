项目使用Maven项目结构

在pom.xml中定义了<dependency>,其中含有对gson的依赖。
在maven项目中不会含有gson的具体jar包，项目执行mvn compile命令会从maven远端仓库下载gson的jar包并放置在本地的一个固定目录下。
在本地电脑上找到该目录并且找到其中的gson-2.8.9.jar，其中含有META-INF文件夹。可在以下两个位置找到其许可证信息。
1. META-INF/MANIFEST.MF的Bundle-License子段
2. META-INF/maven/com.google.code.gson/gson/pom.xml中的<licenses>字段
也可以尝试采用从网页获取的方式得到其License信息。
该gson依赖使用Apache-2.0许可证

主项目使用GPL-2.0-only许可证，位于根目录LICENSE中。

应能扫描出主项目依赖的gson包使用了Apache-2.0许可证，并且和主项目使用的GPL-2.0-only产生了冲突。

