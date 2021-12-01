# 许可证分析器 license-analyzer

## 简介 Abstract

一个基于Gin (go web框架) 的网络服务，
提供分析项目许可证冲突的功能

a gin web service to analyze project license


## 项目结构
调用关系:

*表示可能非必须

main.go: *conf router

*conf: 读取系统配置文件

router: 设置路由=>控制器 controller 

logger: 全局日志配置

> controller: *session scan
>
> > *session: 会话 *mysql 或 *redis 
> > 
> > scanner: 扫描 *mysql *redis
> > > *mysql: 对mysql的增删改查，考虑使用gorm
> > >
> > > *redis: 对redis的增删改查，考虑使用go-redis

util 其它杂项，如: 生成验证码 发邮件等

testProject 用于测试功能的java项目

### scanner输出格式调整：
Scan result as follows:
```json
{
  "external":
    [{"ScannerTest1/lib/DependedProject.jar":"GENERAL PUBLIC LICENSE Version 2"}],
  "local":
    [{"ScannerTest1/License.txt":"GENERAL PUBLIC LICENSE Version 3"}]
}
```

### dependency analyze result：
```json
{
  "ScannerTest4": [
    "ScannerTest4/lib/DependedProject.jar",
    "ScannerTest4/antlib/DependedProject2.jar",
    "ScannerTest4/antlib2/DependedProject3.jar",
    "ScannerTest4/lib/mockito-core-2.28.2.jar",
    "ScannerTest4/lib/jmh-core-1.33.jar"
  ]
}
```