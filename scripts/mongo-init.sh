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

// 初始化用户
db.user.insertMany([{
    "_id": ObjectId(),
    "username": "admin",
    "password": "\$2a\$10\$SLcnDmaJc1nLtUOsZS4yquXyVeu5E6qJHNTVeKSzTk4JO4Xq/FPSy",
    "nickname": "admin",
    "role_id": 0,
    "created_at": new Date(),
    "updated_at": new Date(),
    "avatar": "https://github.githubassets.com/assets/pull-shark-default-498c279a747d.png",
    "email": "admin@example.com",
}]);

// 初始化分类
db.label.insertMany([
  { _id: ObjectId("67c453eda04b00c407b43198"), type: "category", name: "前端开发" },
  { _id: ObjectId("67c453eda04b00c407b43199"), type: "category", name: "后端开发" }
]);

// 初始化标签
db.label.insertMany([
  { _id: ObjectId("67c453eda04b00c407b43203"), type: "tag", name: "Go" },
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
    alias: "stellux-knowledge-base-system-introduction",
    category_id: ObjectId("67c453eda04b00c407b43198"),
    tags_id: [ObjectId("67c453eda04b00c407b43203")],
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_publish: true,
    is_top: true
  }
]);

document_id = ObjectId();
// 初始化文档
db.document.insert( {
    _id: document_id,
    created_at: ISODate("2025-05-29T09:51:13.468Z"),
    updated_at: ISODate("2025-05-29T10:23:55.88Z"),
    title: "测试文档",
    alias: "test",
    description: "一篇测试文档~",
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    document_type: "root",
    is_public: true
} );

db.document_content.insert( {
    _id: document_id,
    created_at: ISODate("2025-05-29T09:51:22.679Z"),
    updated_at: ISODate("2025-05-29T10:24:31.407Z"),
    title: "mytest",
    content: "这是我的第一篇文档😀",
    alias: "",
    description: "",
    thumbnail: "https://cdn.codepzj.cn/image/202503041841864.png",
    is_public: false,
    parent_id: document_id,
    document_id: document_id
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