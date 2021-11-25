此项目包含两个Module，主项目和两个Module都各自使用Ant来进行管理

Module1使用LGPL-2.1-only许可证，位于Module1/LICENSE中。其依赖于DependedProject3.jar，位于Module1/lib中，使用MS-RL许可证。
Module2使用LGPL-2.1-or-later许可证，位于Module2/LICENSE中。其依赖于DependedProject2.jar，位于Module2/antlib中，使用EUPL-1.1许可证。
主项目使用GPL-3许可证，位于根目录下License.txt

需要特别注意的是，LGPL-2.1-only和LGPL-2.1-or-later许可证文本信息几乎完全一样，其主要区别在于所在Module的源代码文件中会有文件头。在文件头中正文第一段的最后，LGPL-2.1-only会是version 2.1，而LGPL-2.1-or-later则会是either version 2.1 of the License, or (at your option) any later version.
另外LGPL-2.1-only文件头的一开始会注明著作权所有人，LGPL-2.1-or-later不仅会注明著作权所有人，还会对此Library的名字和简要说明。

此项目仍然是IDEA编写的项目，其依赖关系仍然会在各个iml文件中记录。
此项目的主项目的各个Ant命令会调用两个Module项目的Ant命令。即调用主项目的ant run会分别调用Moudle1的ant run和Module2的ant run.（这只是我实现的逻辑，不代表其他类似结构的Ant项目也具有这样的逻辑）

对此项目应能检测出Module1所使用的LGPL-2.1-only许可证与其依赖包使用的MS-RL冲突。
对此项目应能检测出Module2所使用的LGPL-2.1-or-later许可证与其依赖包使用的EUPL-1.1冲突。
对此项目应能检测出主项目所使用的GPL-3许可证与两个Module所依赖包使用的MS-RL，EUPL-1.1冲突。