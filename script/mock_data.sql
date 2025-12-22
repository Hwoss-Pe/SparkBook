-- =============================================
-- 微书项目虚拟数据创建脚本
-- 包含用户、文章、互动数据，确保数据关联正确
-- =============================================

-- 使用用户数据库
USE webook_user;

-- 清空现有数据
TRUNCATE TABLE users;

-- 创建虚拟用户数据
INSERT INTO users (id, email, password, phone, birthday, nickname, about_me, wechat_open_id, wechat_union_id, ctime, utime) VALUES
(1, 'zhangsan@example.com', '$2a$10$encrypted_password_hash1', '13800138001', 631123200000, '张三', '热爱编程的后端工程师', NULL, NULL, 1672531200000, 1672531200000),
(2, 'lisi@example.com', '$2a$10$encrypted_password_hash2', '13800138002', 662659200000, '李四', '前端开发专家，Vue.js爱好者', NULL, NULL, 1672531200000, 1672531200000),
(3, 'wangwu@example.com', '$2a$10$encrypted_password_hash3', '13800138003', 694195200000, '王五', '全栈开发者，技术博主', NULL, NULL, 1672531200000, 1672531200000),
(4, 'zhaoliu@example.com', '$2a$10$encrypted_password_hash4', '13800138004', 725731200000, '赵六', 'Java架构师，微服务专家', NULL, NULL, 1672531200000, 1672531200000),
(5, 'sunqi@example.com', '$2a$10$encrypted_password_hash5', '13800138005', 757267200000, '孙七', 'DevOps工程师，云原生实践者', NULL, NULL, 1672531200000, 1672531200000),
(101, 'testuser@example.com', '$2a$10$encrypted_password_hash101', '13800138101', 788803200000, '测试用户', '用于测试的用户账号', NULL, NULL, 1672531200000, 1672531200000);

-- =============================================
-- 使用文章数据库
USE webook_article;

-- 清空现有数据
TRUNCATE TABLE articles;
TRUNCATE TABLE published_articles;

-- 创建虚拟文章数据（制作库）
INSERT INTO articles (id, title, content, cover_image, author_id, status, ctime, utime) VALUES
(111, '2025年Java后端面试高频考点总结', '### 一、JVM核心考点\n1. 垃圾回收算法（G1、CMS）\n2. 内存模型与线程安全\n### 二、Spring核心\n1. IoC容器初始化流程\n2. AOP实现原理\n### 三、数据库优化\n1. 索引设计原则\n2. SQL性能调优\n### 四、分布式系统\n1. 微服务架构设计\n2. 分布式事务处理', 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=800&h=400', 1, 2, 1704067200000, 1704067200000),
(112, 'Vue3 + TypeScript 项目最佳实践', '### 项目搭建\n使用Vite构建工具，配置TypeScript支持\n### 组合式API\n1. setup语法糖的使用\n2. 响应式数据管理\n### 状态管理\n1. Pinia的使用\n2. 模块化状态设计\n### 组件设计\n1. 可复用组件开发\n2. 类型安全的Props定义', 'https://images.unsplash.com/photo-1627398242454-45a1465c2479?w=800&h=400', 2, 2, 1704153600000, 1704153600000),
(113, '微服务架构下的数据一致性解决方案', '### 分布式事务问题\n在微服务架构中，数据一致性是一个重要挑战\n### 解决方案\n1. 两阶段提交（2PC）\n2. 三阶段提交（3PC）\n3. TCC模式\n4. Saga模式\n5. 事件驱动架构\n### 最佳实践\n1. 选择合适的一致性级别\n2. 设计幂等性接口\n3. 实现补偿机制', 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800&h=400', 4, 2, 1704240000000, 1704240000000),
(114, 'Docker容器化部署实战指南', '### 容器化基础\n1. Dockerfile编写规范\n2. 镜像优化技巧\n### 编排部署\n1. Docker Compose使用\n2. Kubernetes部署\n### 监控运维\n1. 日志收集\n2. 性能监控\n3. 健康检查', 'https://images.unsplash.com/photo-1605745341112-85968b19335b?w=800&h=400', 5, 2, 1704326400000, 1704326400000),
(115, 'Redis缓存设计模式与最佳实践', '### 缓存策略\n1. Cache-Aside模式\n2. Write-Through模式\n3. Write-Behind模式\n### 数据结构应用\n1. String的使用场景\n2. Hash的优势\n3. List和Set的应用\n### 高可用设计\n1. 主从复制\n2. 哨兵模式\n3. 集群部署', 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&h=400', 3, 2, 1704412800000, 1704412800000),
(116, 'Go语言并发编程深度解析', '### Goroutine原理\n1. 协程调度机制\n2. GMP模型详解\n### Channel通信\n1. 有缓冲vs无缓冲\n2. select语句使用\n### 并发安全\n1. Mutex互斥锁\n2. RWMutex读写锁\n3. sync.Once使用', 'https://images.unsplash.com/photo-1516259762381-22954d7d3ad2?w=800&h=400', 1, 2, 1704499200000, 1704499200000),
(117, 'Kubernetes生产环境最佳实践', '### 集群规划\n1. 节点规划与资源分配\n2. 网络架构设计\n### 应用部署\n1. Deployment策略\n2. Service网络配置\n3. Ingress路由规则\n### 运维监控\n1. Prometheus监控\n2. 日志聚合\n3. 故障排查', 'https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=800&h=400', 5, 2, 1704585600000, 1704585600000),
(118, 'MySQL性能优化实战经验分享', '### 索引优化\n1. B+树索引原理\n2. 复合索引设计\n3. 覆盖索引应用\n### 查询优化\n1. 执行计划分析\n2. SQL重写技巧\n### 架构优化\n1. 读写分离\n2. 分库分表\n3. 连接池配置', 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800&h=400', 4, 2, 1704672000000, 1704672000000),
(119, '前端性能优化全面指南', '### 加载优化\n1. 代码分割与懒加载\n2. 资源压缩与合并\n3. CDN加速\n### 渲染优化\n1. 虚拟滚动\n2. 防抖与节流\n3. 内存泄漏防范\n### 用户体验\n1. 骨架屏设计\n2. 错误边界处理\n3. 离线缓存策略', 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800&h=400', 2, 2, 1704758400000, 1704758400000),
(120, 'Spring Boot 3.0 新特性详解', '### 核心更新\n1. 原生镜像支持\n2. Java 17+ 要求\n3. Jakarta EE迁移\n### 性能提升\n1. 启动速度优化\n2. 内存占用减少\n### 新功能\n1. 可观测性增强\n2. 配置属性改进\n3. 安全性升级', 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&h=400', 1, 2, 1704844800000, 1704844800000);

-- 创建已发布文章数据（线上库）
INSERT INTO published_articles (id, title, content, cover_image, author_id, status, ctime, utime) VALUES
(111, '2025年Java后端面试高频考点总结', '### 一、JVM核心考点\n1. 垃圾回收算法（G1、CMS）\n2. 内存模型与线程安全\n### 二、Spring核心\n1. IoC容器初始化流程\n2. AOP实现原理\n### 三、数据库优化\n1. 索引设计原则\n2. SQL性能调优\n### 四、分布式系统\n1. 微服务架构设计\n2. 分布式事务处理', 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=800&h=400', 1, 2, 1704067200000, 1704067200000),
(112, 'Vue3 + TypeScript 项目最佳实践', '### 项目搭建\n使用Vite构建工具，配置TypeScript支持\n### 组合式API\n1. setup语法糖的使用\n2. 响应式数据管理\n### 状态管理\n1. Pinia的使用\n2. 模块化状态设计\n### 组件设计\n1. 可复用组件开发\n2. 类型安全的Props定义', 'https://images.unsplash.com/photo-1627398242454-45a1465c2479?w=800&h=400', 2, 2, 1704153600000, 1704153600000),
(113, '微服务架构下的数据一致性解决方案', '### 分布式事务问题\n在微服务架构中，数据一致性是一个重要挑战\n### 解决方案\n1. 两阶段提交（2PC）\n2. 三阶段提交（3PC）\n3. TCC模式\n4. Saga模式\n5. 事件驱动架构\n### 最佳实践\n1. 选择合适的一致性级别\n2. 设计幂等性接口\n3. 实现补偿机制', 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800&h=400', 4, 2, 1704240000000, 1704240000000),
(114, 'Docker容器化部署实战指南', '### 容器化基础\n1. Dockerfile编写规范\n2. 镜像优化技巧\n### 编排部署\n1. Docker Compose使用\n2. Kubernetes部署\n### 监控运维\n1. 日志收集\n2. 性能监控\n3. 健康检查', 'https://images.unsplash.com/photo-1605745341112-85968b19335b?w=800&h=400', 5, 2, 1704326400000, 1704326400000),
(115, 'Redis缓存设计模式与最佳实践', '### 缓存策略\n1. Cache-Aside模式\n2. Write-Through模式\n3. Write-Behind模式\n### 数据结构应用\n1. String的使用场景\n2. Hash的优势\n3. List和Set的应用\n### 高可用设计\n1. 主从复制\n2. 哨兵模式\n3. 集群部署', 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=800&h=400', 3, 2, 1704412800000, 1704412800000),
(116, 'Go语言并发编程深度解析', '### Goroutine原理\n1. 协程调度机制\n2. GMP模型详解\n### Channel通信\n1. 有缓冲vs无缓冲\n2. select语句使用\n### 并发安全\n1. Mutex互斥锁\n2. RWMutex读写锁\n3. sync.Once使用', 'https://images.unsplash.com/photo-1516259762381-22954d7d3ad2?w=800&h=400', 1, 2, 1704499200000, 1704499200000),
(117, 'Kubernetes生产环境最佳实践', '### 集群规划\n1. 节点规划与资源分配\n2. 网络架构设计\n### 应用部署\n1. Deployment策略\n2. Service网络配置\n3. Ingress路由规则\n### 运维监控\n1. Prometheus监控\n2. 日志聚合\n3. 故障排查', 'https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=800&h=400', 5, 2, 1704585600000, 1704585600000),
(118, 'MySQL性能优化实战经验分享', '### 索引优化\n1. B+树索引原理\n2. 复合索引设计\n3. 覆盖索引应用\n### 查询优化\n1. 执行计划分析\n2. SQL重写技巧\n### 架构优化\n1. 读写分离\n2. 分库分表\n3. 连接池配置', 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800&h=400', 4, 2, 1704672000000, 1704672000000),
(119, '前端性能优化全面指南', '### 加载优化\n1. 代码分割与懒加载\n2. 资源压缩与合并\n3. CDN加速\n### 渲染优化\n1. 虚拟滚动\n2. 防抖与节流\n3. 内存泄漏防范\n### 用户体验\n1. 骨架屏设计\n2. 错误边界处理\n3. 离线缓存策略', 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800&h=400', 2, 2, 1704758400000, 1704758400000),
(120, 'Spring Boot 3.0 新特性详解', '### 核心更新\n1. 原生镜像支持\n2. Java 17+ 要求\n3. Jakarta EE迁移\n### 性能提升\n1. 启动速度优化\n2. 内存占用减少\n### 新功能\n1. 可观测性增强\n2. 配置属性改进\n3. 安全性升级', 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&h=400', 1, 2, 1704844800000, 1704844800000);

-- =============================================
-- 使用互动数据库 (目标数据库 - 新架构)
USE webook_intr;

-- 同时也要在源数据库webook中创建数据（支持双写模式）
-- 如果Interactive服务当前使用PatternSrcOnly模式，需要在webook数据库中也有数据

-- 清空现有数据
TRUNCATE TABLE interactives;
TRUNCATE TABLE user_like_bizs;
TRUNCATE TABLE user_collection_bizs;
# TRUNCATE TABLE collections;

-- 创建收藏夹数据
# INSERT INTO collections (id, name, uid, ctime, utime) VALUES
# (1, '技术文章收藏', 101, 1704067200000, 1704067200000),
# (2, '面试准备', 101, 1704067200000, 1704067200000),
# (3, '架构设计', 2, 1704067200000, 1704067200000),
# (4, '前端技术', 2, 1704067200000, 1704067200000),
# (5, '后端开发', 3, 1704067200000, 1704067200000);

-- 创建文章互动数据
INSERT INTO interactives (id, biz_id, biz, read_cnt, collect_cnt, like_cnt, ctime, utime) VALUES
(1, 111, 'article', 156, 12, 23, 1704067200000, 1704931200000),
(2, 112, 'article', 89, 8, 15, 1704153600000, 1704931200000),
(3, 113, 'article', 234, 18, 45, 1704240000000, 1704931200000),
(4, 114, 'article', 167, 14, 28, 1704326400000, 1704931200000),
(5, 115, 'article', 198, 16, 35, 1704412800000, 1704931200000),
(6, 116, 'article', 145, 11, 22, 1704499200000, 1704931200000),
(7, 117, 'article', 276, 22, 52, 1704585600000, 1704931200000),
(8, 118, 'article', 189, 15, 31, 1704672000000, 1704931200000),
(9, 119, 'article', 134, 9, 19, 1704758400000, 1704931200000),
(10, 120, 'article', 201, 17, 38, 1704844800000, 1704931200000);

-- 创建用户点赞数据
INSERT INTO user_like_bizs (id, biz_id, biz, uid, status, ctime, utime) VALUES
(1, 111, 'article', 101, 1, 1704067200000, 1704067200000),
(2, 113, 'article', 101, 1, 1704240000000, 1704240000000),
(3, 115, 'article', 101, 1, 1704412800000, 1704412800000),
(4, 117, 'article', 101, 1, 1704585600000, 1704585600000),
(5, 120, 'article', 101, 1, 1704844800000, 1704844800000),
(6, 112, 'article', 2, 1, 1704153600000, 1704153600000),
(7, 114, 'article', 2, 1, 1704326400000, 1704326400000),
(8, 116, 'article', 2, 1, 1704499200000, 1704499200000),
(9, 118, 'article', 3, 1, 1704672000000, 1704672000000),
(10, 119, 'article', 3, 1, 1704758400000, 1704758400000);

-- 创建用户收藏数据（注意：每个用户对同一篇文章只能收藏一次，不同收藏夹）
INSERT INTO user_collection_bizs (id, cid, biz_id, biz, uid, ctime, utime) VALUES
(1, 1, 111, 'article', 101, 1704067200000, 1704067200000),
(2, 1, 113, 'article', 101, 1704240000000, 1704240000000),
(3, 1, 115, 'article', 101, 1704412800000, 1704412800000),
(4, 2, 117, 'article', 101, 1704585600000, 1704585600000),
(5, 2, 120, 'article', 101, 1704844800000, 1704844800000),
(6, 4, 112, 'article', 2, 1704153600000, 1704153600000),
(7, 4, 119, 'article', 2, 1704758400000, 1704758400000),
(8, 5, 114, 'article', 3, 1704326400000, 1704326400000),
(9, 5, 116, 'article', 3, 1704499200000, 1704499200000),
(10, 5, 118, 'article', 3, 1704672000000, 1704672000000);

-- =============================================
-- 数据验证查询（可选执行）
-- =============================================

-- 验证用户数据
-- SELECT id, nickname, email FROM webook_user.users;

-- 验证文章数据
-- SELECT id, title, author_id, status FROM webook_article.articles;

-- 验证互动数据
-- SELECT biz_id, read_cnt, like_cnt, collect_cnt FROM webook_intr.interactives WHERE biz = 'article';

-- 验证数据关联性
-- SELECT 
--     a.id as article_id,
--     a.title,
--     u.nickname as author_name,
--     i.read_cnt,
--     i.like_cnt,
--     i.collect_cnt
-- FROM webook_article.published_articles a
-- LEFT JOIN webook_user.users u ON a.author_id = u.id
-- LEFT JOIN webook_intr.interactives i ON a.id = i.biz_id AND i.biz = 'article'
-- ORDER BY a.id;

-- =============================================
-- 脚本执行完成
-- 数据创建完毕，各表数据已关联
-- =============================================
