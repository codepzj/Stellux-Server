#!/bin/bash

mongosh -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase admin <<EOF

// 创建目标数据库用户
db = db.getSiblingDB('$MONGO_INITDB_DATABASE');
db.createUser({
    user: '$MONGO_USERNAME',
    pwd: '$MONGO_PASSWORD',
    roles: [{ role: 'readWrite', db: '$MONGO_INITDB_DATABASE' }]
});

db.auth('$MONGO_USERNAME', '$MONGO_PASSWORD');

let AdminId = ObjectId("67c453eda04b00c407b431fd");
let UserId = ObjectId("67c453eda04b00c407b431fe");
let TestId = ObjectId("67c453eda04b00c407b431ff");

// 管理员所有权限
db.casbin_rule.insertMany([{
    "_id": ObjectId(),
    "ptype": "p",
    "v0": "admin",
    "v1": "*",
    "v2": "GET"
},{
    "_id": ObjectId(),
    "ptype": "p",
    "v0": "admin",
    "v1": "*",
    "v2": "POST"
},{
    "_id": ObjectId(),
    "ptype": "p",
    "v0": "admin",
    "v1": "*",
    "v2": "PATCH"
},{
    "_id": ObjectId(),
    "ptype": "p",
    "v0": "admin",
    "v1": "*",
    "v2": "PUT"
},{
    "_id": ObjectId(),
    "ptype": "p",
    "v0": "admin",
    "v1": "*",
    "v2": "DELETE"
}]);

// 为用户授权
db.casbin_rule.insertMany([{
    "_id": ObjectId(),
    "ptype": "g",
    "v0": "67c453eda04b00c407b431fd",
    "v1": "admin"
}, {
    "_id": ObjectId(),
    "ptype": "g",
    "v0": "67c453eda04b00c407b431fe",
    "v1": "user"
}, {
    "_id": ObjectId(),
    "ptype": "g",
    "v0": "67c453eda04b00c407b431ff",
    "v1": "test"
}]);

// 初始化用户
db.user.insertMany([{
    "_id": AdminId,
    "username": "admin",
    "password": "\$2a\$10\$SLcnDmaJc1nLtUOsZS4yquXyVeu5E6qJHNTVeKSzTk4JO4Xq/FPSy",
    "nickname": "codezj",
    "role_id": 0,
    "created_at": new Date(),
    "updated_at": new Date(),
    "avatar": "https://github.githubassets.com/assets/pull-shark-default-498c279a747d.png",
    "email": "admin@example.com",
}, {
    "_id": UserId,
    "username": "alice",
    "password": "\$2a\$10\$SLcnDmaJc1nLtUOsZS4yquXyVeu5E6qJHNTVeKSzTk4JO4Xq/FPSy",
    "nickname": "小李",
    "role_id": 1,
    "created_at": new Date(),
    "updated_at": new Date(),
    "avatar": "https://github.githubassets.com/assets/quickdraw-default-39c6aec8ff89.png",
    "email": "alice@example.com",
}, {
    "_id": TestId,
    "username": "test",
    "password": "\$2a\$10\$SLcnDmaJc1nLtUOsZS4yquXyVeu5E6qJHNTVeKSzTk4JO4Xq/FPSy",
    "nickname": "小张",
    "role_id": 2,
    "created_at": new Date(),
    "updated_at": new Date(),
    "avatar": "https://github.githubassets.com/assets/yolo-default-be0bbff04951.png",
    "email": "test@example.com",
}]);

// 初始化分类
db.label.insertMany([
  { _id: ObjectId("67c453eda04b00c407b43197"), type: "category", name: "默认分类" },
  { _id: ObjectId("67c453eda04b00c407b43198"), type: "category", name: "前端开发" },
  { _id: ObjectId("67c453eda04b00c407b43199"), type: "category", name: "后端开发" }
]);

// 初始化标签
db.label.insertMany([
  { _id: ObjectId("67c453eda04b00c407b43200"), type: "tag", name: "JavaScript" },
  { _id: ObjectId("67c453eda04b00c407b43201"), type: "tag", name: "Node.js" },
  { _id: ObjectId("67c453eda04b00c407b43202"), type: "tag", name: "MongoDB" },
  { _id: ObjectId("67c453eda04b00c407b43203"), type: "tag", name: "Go" },
  { _id: ObjectId("67c453eda04b00c407b43204"), type: "tag", name: "部署运维" }
]);

// 初始化文章
db.post.insertMany([
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "Stellux知识库系统简介",
    content: "欢迎使用Stellux知识库系统，这是一款支持标签分类、全文搜索和可视化编辑的轻量级系统。",
    description: "系统简介",
    category_id: ObjectId("67c453eda04b00c407b43197"),
    tags_id: [ObjectId("67c453eda04b00c407b43200"), ObjectId("67c453eda04b00c407b43201")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: true
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "前端开发规范汇总",
    content: "本文介绍了一套适用于团队协作的前端编码规范，包括命名规范、注释规范等。",
    description: "团队前端规范文档",
    category_id: ObjectId("67c453eda04b00c407b43198"),
    tags_id: [ObjectId("67c453eda04b00c407b43200")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "Node.js 入门教程",
    content: "从零开始学习 Node.js，包括模块系统、异步 IO、Express 框架等。",
    description: "Node.js 学习笔记",
    category_id: ObjectId("67c453eda04b00c407b43199"),
    tags_id: [ObjectId("67c453eda04b00c407b43201")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "MongoDB 常用操作指令",
    content: "总结了 MongoDB 的基础操作，如增删改查、索引、聚合等。",
    description: "MongoDB 使用手册",
    category_id: ObjectId("67c453eda04b00c407b43199"),
    tags_id: [ObjectId("67c453eda04b00c407b43202")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "Go 后端服务架构设计",
    content: "介绍如何使用 Go 编写一个可扩展的微服务，包括 RESTful API、路由中间件、日志追踪等内容。",
    description: "Go 微服务架构方案",
    category_id: ObjectId("67c453eda04b00c407b43199"),
    tags_id: [ObjectId("67c453eda04b00c407b43203")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "使用 PM2 管理 Node 应用",
    content: "介绍如何使用 PM2 工具部署并守护 Node.js 应用。",
    description: "Node 应用部署方案",
    category_id: ObjectId("67c453eda04b00c407b43199"),
    tags_id: [ObjectId("67c453eda04b00c407b43204"), ObjectId("67c453eda04b00c407b43201")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "Git 工作流与分支策略",
    content: "讲解 Git Flow、GitHub Flow 等主流工作流的使用方式及其适用场景。",
    description: "Git 多人协作最佳实践",
    category_id: ObjectId("67c453eda04b00c407b43197"),
    tags_id: [ObjectId("67c453eda04b00c407b43204")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "前端部署流程总结",
    content: "介绍了前端从打包构建到 Nginx 部署上线的整个流程。",
    description: "完整部署流程",
    category_id: ObjectId("67c453eda04b00c407b43198"),
    tags_id: [ObjectId("67c453eda04b00c407b43200"), ObjectId("67c453eda04b00c407b43204")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "Go 单元测试与覆盖率",
    content: "如何编写 Go 的测试代码，使用 `go test` 和 `-cover` 进行覆盖率分析。",
    description: "Go 测试实践",
    category_id: ObjectId("67c453eda04b00c407b43199"),
    tags_id: [ObjectId("67c453eda04b00c407b43203")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  },
  {
    _id: ObjectId(),
    created_at: new Date(),
    updated_at: new Date(),
    author: "codepzj",
    title: "MongoDB 聚合管道实战",
    content: "使用 MongoDB 聚合管道处理复杂数据查询，实战案例包含 `$lookup`、`$group` 等。",
    description: "复杂数据聚合处理",
    category_id: ObjectId("67c453eda04b00c407b43199"),
    tags_id: [ObjectId("67c453eda04b00c407b43202")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: false
  }
]);

// 初始化文档
db.document.insert( {
    _id: ObjectId("68382e116d357d131691114c"),
    created_at: ISODate("2025-05-29T09:51:13.468Z"),
    updated_at: ISODate("2025-05-29T10:23:55.88Z"),
    title: "测试文档",
    content: "",
    alias: "test",
    description: "一篇测试文档~",
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    document_type: "root",
    is_public: true
} );

db.document.insert( {
    _id: ObjectId("68382e1a6d357d131691114d"),
    created_at: ISODate("2025-05-29T09:51:22.679Z"),
    updated_at: ISODate("2025-05-29T10:24:31.407Z"),
    title: "mytest",
    content: "这是我的第一篇文档😀",
    alias: "",
    description: "",
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    document_type: "leaf",
    is_public: false,
    parent_id: ObjectId("68382e116d357d131691114c"),
    document_id: ObjectId("68382e116d357d131691114c")
} );

// 初始化设置
db.setting.insert( {
    _id: ObjectId("687312fa9484dab0aed07741"),
    key: "site_config",
    value: {
        site_title: "codepzj",
        site_subtitle: "codepzj",
        site_favicon: "/favicon.ico",
        site_author: "codepzj",
        site_animate_text: "浩瀚星河",
        site_avatar: "https://cdn.codepzj.cn/image/20250529174726187.jpeg",
        site_description: "88",
        site_copyright: "345",
        site_url: "http://localhost:9003",
        site_keywords: "6666",
        site_icp: "43543",
        site_icplink: "https://baicu.com",
        github_username: "codepzj"
    }
} );


EOF