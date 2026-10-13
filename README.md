# bk-apigateway-sdks
[English](README_EN.md)

## 简介
该仓库包含了蓝鲸 API 网关相关的 SDK 及工具集，简化开发者接入和使用蓝鲸 API 网关的开发过程。

## 总览
- [manager](manager)：网关管理工具集
- [bkapi](bkapi)：网关 API 调用客户端，处理认证、响应、状态码、日志等公共逻辑
- [gin_contrib](gin_contrib)：基于gin的网关编程框架

## 从 v1 升级

v2 包含不兼容变更，升级前请阅读 [v2 迁移指南](MIGRATION_V2.md)，覆盖模块路径、客户端调用、manager、权限配置和 Gin 集成。

## 蓝鲸社区

- [BK-CI](https://github.com/Tencent/bk-ci)：蓝鲸持续集成平台是一个开源的持续集成和持续交付系统，可以轻松将你的研发流程呈现到你面前。
- [BK-BCS](https://github.com/Tencent/bk-bcs)：蓝鲸容器管理平台是以容器技术为基础，为微服务业务提供编排管理的基础服务平台。
- [BK-PaaS](https://github.com/Tencent/bk-PaaS)：蓝鲸PaaS平台是一个开放式的开发平台，让开发者可以方便快捷地创建、开发、部署和管理SaaS应用。
- [BK-SOPS](https://github.com/Tencent/bk-sops)：标准运维（SOPS）是通过可视化的图形界面进行任务流程编排和执行的系统，是蓝鲸体系中一款轻量级的调度编排类SaaS产品。
- [BK-CMDB](https://github.com/Tencent/bk-cmdb)：蓝鲸配置平台是一个面向资产及应用的企业级配置管理平台。

## 贡献

如果你有好的意见或建议，欢迎给我们提 Issues 或 Pull Requests，为蓝鲸开源社区贡献力量。

## 协议

基于 MIT 协议，详细请参考[LICENSE](LICENSE)
