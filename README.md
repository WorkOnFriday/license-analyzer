# 许可证分析器 license-analyzer

## 简介 Abstract

一个基于Gin (go web框架) 的网络服务， 提供分析项目许可证冲突的功能

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

### 提交任务任务 接口和返回值：

POST 127.0.0.1:8080/task

表单： file <要分析的项目压缩包>

返回Body为项目号

### 获取任务结果 接口和返回值：

GET 127.0.0.1:8080/task?id=<任务号>

例如：

GET 127.0.0.1:8080/task?id=0

```json
{
  "IsFinish": true,
  "ErrorMessage": "",
  "Local": [
    {
      "Path": "ScannerTest4\\LICENSE",
      "License": "GENERAL PUBLIC LICENSE Version 2"
    }
  ],
  "External": [
    {
      "Path": "ScannerTest4\\antlib\\DependedProject2.jar",
      "License": "EUROPEAN UNION PUBLIC LICENCE V. 1.1"
    },
    {
      "Path": "ScannerTest4\\antlib2\\DependedProject3.jar",
      "License": "MICROSOFT RECIPROCAL LICENSE"
    },
    {
      "Path": "ScannerTest4\\lib\\DependedProject.jar",
      "License": "GENERAL PUBLIC LICENSE Version 2"
    }
  ],
  "AllExternalModule": [
    {
      "JarPath": "ScannerTest4\\antlib\\DependedProject2.jar",
      "Package": [
        "dp2"
      ]
    },
    {
      "JarPath": "ScannerTest4\\antlib2\\DependedProject3.jar",
      "Package": [
        "dp3"
      ]
    },
    {
      "JarPath": "ScannerTest4\\lib\\DependedProject.jar",
      "Package": [
        "pri"
      ]
    }
  ],
  "AllLocalModule": [
    {
      "Module": "ScannerTest4",
      "License": "GENERAL PUBLIC LICENSE Version 2"
    }
  ],
  "Dependency": {
    "Project": {
      "Module": "ScannerTest4",
      "Dependencies": [
        "ScannerTest4\\lib\\DependedProject.jar",
        "ScannerTest4\\antlib\\DependedProject2.jar",
        "ScannerTest4\\antlib2\\DependedProject3.jar"
      ]
    },
    "Modules": [
      {
        "Module": "ScannerTest4",
        "Dependencies": [
          "ScannerTest4\\lib\\DependedProject.jar",
          "ScannerTest4\\antlib\\DependedProject2.jar",
          "ScannerTest4\\antlib2\\DependedProject3.jar"
        ]
      }
    ]
  },
  "PomLicense": null
}
```

```json
{
  "IsFinish": true,
  "ErrorMessage": "",
  "Local": [
    {
      "Path": "ScannerTest6\\LICENSE",
      "License": "GENERAL PUBLIC LICENSE Version 2"
    }
  ],
  "External": [],
  "AllExternalModule": null,
  "AllLocalModule": [
    {
      "Module": "ScannerTest6",
      "License": "GENERAL PUBLIC LICENSE Version 2"
    }
  ],
  "Dependency": {
    "Project": {
      "Module": "ScannerTest6",
      "Dependencies": []
    },
    "Modules": [
      {
        "Module": "ScannerTest6",
        "Dependencies": []
      }
    ]
  },
  "PomLicense": [
    {
      "GroupID": "com.google.code.gson",
      "ArtifactID": "gson",
      "Version": "2.8.9",
      "License": "Apache 2.0"
    }
  ]
}
```