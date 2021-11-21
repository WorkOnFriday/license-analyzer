在lib文件夹中放入DependedProject.jar依赖，该jar包提供pri包，使用GPL-2.0-only许可证。

ScannerTest2项目整体使用GPL-2.0-only许可证，位于根目录License.txt，其子包gpl3使用GPL-3.0-only许可证，位于src/work/gpl3/License.txt

使用依赖的配置在ScannerTest2.iml中（由idea生成，作用于idea开发的项目）

对此项目应能检测出其子包gpl3所使用的GPL-3.0-only许可证与其依赖包使用的GPL-2.0-only冲突。