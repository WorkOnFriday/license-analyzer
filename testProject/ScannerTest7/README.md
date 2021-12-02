项目采用maven结构，其有两个子Module，两个Module也都是maven结构。

Module1依赖了slf4j-api，其依赖信息可通过Module1/pom.xml的<dependencies>标签查看。
slf4j-api使用MIT许可证，但其下载出的依赖包中似乎完全没有相关信息。可以考虑通过网页查询的方式获得该许可证信息。也可以考虑直接以相关命令得到的结果代替。

Module2依赖了junit-jupiter-api，其依赖信息可通过Module2/pom.xml的<dependencies>标签查看。
junit-jupiter-api使用EPL-2.0许可证，其存储在下载下来的junit-jupiter-api-5.8.2.jar/META-INF/LICENSE.md中

主项目依赖jboss-aop，其依赖信息可通过根目录下的pom.xml的<dependencies>标签查看。
jboss-aop使用LGPL-2.1许可证，其下载出的依赖包中似乎完全没有相关信息。可以考虑通过网页查询的方式获得该许可证信息。也可以考虑直接以相关命令得到的结果代替。

主项目使用Apache-2.0许可证，位于根目录下License.txt，不会产生依赖冲突。

