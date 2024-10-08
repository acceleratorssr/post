# post
本项目是一个基于微服务架构的系统，主要功能包括帖子、点赞、用户管理、基础搜索（Elasticsearch）、评论和任务调度。项目采用 BFF层统一提供对外 HTTP 接口，确保各项业务逻辑的灵活性与可扩展性。

## 文件结构
- api: grpc 定义文件;
- pkg: 框架扩展的一些常用方法;
- bff: 接口层，负责将 HTTP 请求转换为 gRPC 请求，将 gRPC 响应转换为 HTTP 响应，并实现业务逻辑;
- xxx: 业务模块，如：article、interactive、search、sso、user;
  - domain: 存放核心数据类型定义及相关逻辑;
  - events: 存放消息队列相关业务逻辑;
  - grpc: 存放 grpc 服务端的实现及部分业务逻辑;
  - ioc: 存放 ioc 配置文件及相关逻辑;
  - job: 定时任务相关逻辑;
  - repository: 存放数据库访问抽象逻辑;
    - cache: 基于redis的缓存操作相关逻辑;
    - dao: 基于mysql的数据操作逻辑
  - service: 存放主要业务逻辑实现;
- scripts: 编译、运行脚本，方便调试起见，没把服务独立化部署;

## 技术选型
RPC框架：gRPC(IDL语言: proto3)

HTTP框架：gin

ORM框架：GORM

数据库：MySQL

中间件：Kafka、Redis、Elasticsearch

服务发现：etcd

可观测性：OpenTelemetry, Zipkin, prometheus、Grafana

## 构建方式
通过 `Makefile` 运行脚本，如：`make build_py`;