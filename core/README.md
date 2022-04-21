# 蓝鲸云 API 客户端(Golang)
本项目基于 gentleman，提供了一种基于配置构建 SDK 的方案，作为蓝鲸 API 网关 SDK 底层实现。

## 项目结构：
- bkapi：可公开使用的各类场景插件
- define：核心模型定义
- demo：示例代码
- internal：内部实现模块（不公开）

## 命名约定：
在 bkapi 包中，作以下命名约定：

| 模式                  | 类型 | 含义                                    | 参考文件  |
| --------------------- | ---- | --------------------------------------- | --------- |
| bkapi.Opt*            | 函数 | 可同时作用于 Client 和 Operation 的选项 | option.go |
| bkapi.*BodyProvider   | 结构 | 提供请求体序列化封装的接口              | json.go   |
| bkapi.*ResultProvider | 结构 | 提供响应体反序列化封装的接口            | json.go   |

