## 小微书Webook

### 简介

​	小微书是一个尝试模仿小红书的应用，包含有帖子，点赞，评论，收藏，聊天，榜单等基本功能，对事件的存储，分发，治理，容错，安全，可以理解成一个大型的论坛项目，虽然项目创新点不足，但在应用上可学习相关内容，项目采用微服务架构，并且对高并发，高可用，高性能进行全方位思考和实践，此项目只供个人学习相关，前端部分相对简陋，以下设计针对后端进行实现，根据提交内容可以查看对应的进度，为每个开发阶段制定测试方案，包括单元测试、集成测试和负载测试，确保系统稳定可靠，当前仍在建设中。

### 开发环境

IDE🧑‍💻： [GoLand](https://www.jetbrains.com/go/)

go 1.22.0

OS🪟🐧：[Ubuntu 22.04.3 LTS (WSL2)](https://ubuntu.com/desktop/wsl)

### 业务分析

#### 用户服务

安全:
1.密码不采用明文存储，用md5加密后存入数据库

2.访问资源采用jwt，并且这个jwt其实可以做成一个middleware给gin使用，并且可以结合长短token一起使用

3.还需要请求头携带user-agent进行，更深的可以上升到浏览器指纹这块让前端去操作

4.注册和登录的逻辑采用insert for update，其实就是那句如果遇到id冲突就成更新
注册操作需要校验是否之前注册过，这里不要用先去查一次库的方式，而是直接插入由于手机号是独立的，去根据mysql的错误重复数据1042进行判断create是否成立的
这里面还有有个findById的操作，思路是这里采用查数据库后写进缓存，触发降级就不走数据库的快慢路径操作先走缓存，缓存找不到就从数据库找，找完重新set进去，如果redis出现错误，我就不让他流量直接打到数据库，因此直接返回错误

5.微信登录，这个确实进行访问后微信给你一个

url，然后扫后等待微信回调一个方法等待接受，然后进行注册或者登录操作，主要是里面的unionid和openid

#### 文章服务

#### 关注服务

#### 交互服务

#### 支付服务

#### 热榜服务

#### 搜索服务

#### 短信服务

#### 标签服务

#### 打赏服务

#### 权限服务

#### 聊天服务

#### 评论服务

#### 权限服务

#### 自定义三方扩展pkg

**logger**：这个主要是采用zapLogger拓展，在不适用其他日志现成框架的前提下，单纯用接口定义出对应常用的方法，并且采用Field的字段约束，并且用这个给gorm和gin进行拓展打印日志，在其他业务代码中也是采用这个

**架构**：
![image-20241119154948019](https://github.com/Hwoss-Pe/Webook/blob/main/image-20241119154948019.png)

### 技术栈

- [Node.js](https://nodejs.org/en)
- Docker
  - [镜像源](https://yeasy.gitbook.io/docker_practice/install/mirror)（还是挂代理方便）
  - [mysql](https://hub.docker.com/_/mysql) - An open-source relational database management system (RDBMS)
  - [redis](https://hub.docker.com/r/bitnami/redis) - An open-source in-memory storage
  - [etcd](https://hub.docker.com/r/bitnami/etcd) - A distributed key-value store designed to securely store data across a cluster
  - [mongo](https://hub.docker.com/_/mongo) - MongoDB document databases provide high availability and easy scalability
  - [kafka](https://hub.docker.com/r/bitnami/kafka) - Apache Kafka is a distributed streaming platform used for building real-time applications
  - [prometheus](https://hub.docker.com/r/bitnami/prometheus) - The Prometheus monitoring system and time series database
  - *grafana - The open observability platform*
  - *zipkin - A distributed tracing system*
- kubernates
  - [Kubernetes cluster architecture](https://kubernetes.io/docs/concepts/architecture/)
  - [kubectl](https://kubernetes.io/docs/tasks/tools/) - The Kubernetes command-line tool
  - [HELM](https://helm.sh/) - The package manager for Kubernetes
  - [ingress-nignx](https://github.com/kubernetes/ingress-nginx) - Ingress-NGINX Controller for Kubernetes
- [wrk](https://github.com/wg/wrk) - Modern HTTP benchmarking tool
- [protobuf](https://github.com/protocolbuffers/protobuf) - Protocol Buffers - Google's data interchange format
- [buf](https://github.com/bufbuild/buf) - The best way of working with Protocol Buffers

### 编程能力

**分包架构**

这体现在对于各种服务之间的管理和应用，并且服务内遵循dao- repository-service-handler ，并且使用mock和wire对里面进行测试和依赖注入，并且采用DDD和TDD的一些实际运用，各个服务之间都采用rpc调用

**面向失败编程**

面向失败编程（Failure-oriented Programming，FOP）是一种编程范式，强调在编写逻辑时尽可能考虑边界条件和失败情况，以增强程序的稳定性。 简单来说，就是时刻考虑系统可能会崩溃。无论是系统本身、依赖的服务还是依赖的数据库，都可能会崩溃。 面向失败编程不仅仅是对输入进行校验，它还包括：

- 错误处理：需要严密处理各种可能的错误情况
- 容错设计：长期培养的能力是针对业务和系统特征设计容错策略。这通常是较难掌握的，而其余部分可以通过规范来达成。 在面向失败编程中，需要长期培养的能力是针对业务和系统特征设计容错策略。其他方面较容易掌握，或者公司可以通过规范来达成。 在项目中，讨论了许多容错方案，包括：
- 重试机制：需要考虑重试的间隔和次数，以及最后可能需要人工介入。
- 监控与告警：在追求高可用时，还要考虑自动修复的程度
- 限流：用于保护系统本身。
- 下游服务治理：如果下游服务可能崩溃，需使用一些治理技巧：
  - 轮询：可以是每次都轮询，也可以针对某个下游节点失败后的限流。
  - 客户端限流：限制客户端的请求速率以保护系统资源。
  - 同步转异步：在转为异步后，必须保证请求会被处理而不会遗漏
- 考虑安全性：例如，防止 token 泄露以增强系统的安全性。 在设计容错方案时，尽可能在平时收集别人使用的容错方案，以了解各种处理方式。根据自己实际处理的业务设计合适的容错方案。简单地生搬硬套别人的方案，效果可能不佳。

**灵活的缓存方案**

在整个单体应用中，已经充分接触了缓存方案。相比传统的缓存方案，项目中的缓存方案更具“趣味性”。在实践中，除非逼不得已，通常不会使用看起来非常特殊的缓存方案。 使用过和讨论过的缓存方案包括：

- 只使用 Redis：更新缓存的常见方案是更新数据库后删除缓存。
- 本地缓存与 Redis 缓存结合使用。大多数系统完成这些步骤即可，
  - 查找顺序：本地缓存 - Redis - 数据库
  - 更新顺序：数据库 - 本地缓存 - Redis
- 根据业务特征动态设置缓存的过期时间。例如，如果能判定某个用户是大 V，则他的数据过期时间应设得更长。
- 淘汰对象：根据业务特征来淘汰缓存对象。
- 缓存崩溃：需要考虑缓存崩溃的问题。在实践中，缓存崩溃可能导致数据库也一起崩溃。 在上述缓存方案的基础上，需要能够举一反三，根据业务特征设计针对性的解决方案。在整个职业生涯中，如果能有效使用缓存，就能解决 90% 的性能问题。剩下的 10% 则需要依靠各种技巧和优化手段。

**注意并发问题**

无论是代码中的 Go 实例，还是外部数据库，在实现任何功能时操作对象或 Redis 缓存数据时，都必须考虑并发问题。具体来说，需要关注是否有多个 goroutine 在同一时刻读写对象，这些 goroutine 可能在不同的实例（机器）上，也可能在同一实例（机器）上。 在项目中，使用了多种方法来解决并发问题：

- SELECT FOR UPDATE：用于确保读取的数据在操作期间不会被修改，简单且有效。
- 分布式锁：用于保证同一时刻只有一个 goroutine 可以执行特定操作。
- Lua 脚本：在 Redis 中使用 Lua 脚本来确保在执行多个操作时没有其他 goroutine 修改 Redis 数据。
- 乐观锁：使用数据库 version 加 CAS（Compare and Swap）机制来保证在修改数据时，数据未被其他操作修改过。
- Go 对象锁：使用 `sync.Mutex` 和 `sync.RWMutex` 来管理对 Go 对象的并发访问，在某些情况下，还可以使用原子操作（`atomic` 包）来处理简单的并发问题。 在实践中，只能通过长期训练来培养并发意识。在项目开始时，就应有意识地培养自己对并发问题的关注和敏感度。

**依赖注入**

首先要整体上领悟依赖注入和面向接口编程的优势，这些优点在项目中体现得非常明显：

- 依赖注入完全达成了控制反转的目标。不再关心如何创建依赖对象。例如，在 cache 模块中，虽然使用了 Redis 客户端，但 cache 实现并不关心具体的实现或客户端的相关参数。
- 依赖注入提高了代码的可测试性。可以在单元测试中注入由 `gomock` 生成的实例。在集成测试阶段，为了节省公司资源，第三方依赖通常被替换为内存实现或 mock 实现。
- 依赖注入叠加面向接口编程后，装饰器模式效果更佳。在 sms 模块中，有各种装饰器的实现，这些实现都是基于面向接口编程和依赖注入的。这使得装饰器可以自由组合，提升了系统的灵活性和扩展性。
