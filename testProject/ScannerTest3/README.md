对于大型的java项目，可能会出现具有多个Module的情况，这些Module相对独立，但是可以共享主项目的依赖。除此之外，他们也可以有自己独立的Module级别依赖。

在lib文件夹中放入DependedProject.jar依赖，该jar包提供pri包，使用GPL-2.0-only许可证。主项目（主Module），Module1, Module2都会使用该依赖。

主项目（主Module）使用GPL-2.0-only许可证，位于根目录下LICENSE。Module1使用GPL-3.0-only许可证，位于/Module1/LICENSE。Module2使用Apache-2.0许可证，位于Module2/LICENSE.txt。

主项目，Module1，Module2使用依赖的配置分别在在ScannerTest3.iml，Module1/Module1.iml，Module2/Module2.iml中（iml由idea生成，作用于idea开发的项目）

对此项目应能检测出Module1所使用的GPL-3.0-only许可证与其依赖包使用的GPL-2.0-only冲突，Module2所使用的Apache-2.0许可证与其依赖包使用的GPL-2.0-only冲突，检测结果也可以包为单位，即Module1.gpl3与其依赖包冲突。