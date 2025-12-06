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

// 批量插入用户数据
db.user.insertMany([
    {
        _id: ObjectId("67c453eda04b00c407b431fd"),
        username: "admin",
        password: "\$2a\$10\$EPbYKMyDA5RN9AaEEL7RqePI4BotBGCDvZ/ny/YHasEoU4vhU5n4e",
        nickname: "codepzj",
        role_id: Int32("0"),
        created_at: ISODate("2025-07-13T02:19:26.865Z"),
        updated_at: ISODate("2025-07-13T02:19:26.865Z"),
        avatar: "https://github.githubassets.com/assets/pull-shark-default-498c279a747d.png",
        email: "admin@example.com"
    }
]);

// 批量插入文章数据
db.post.insertMany([
    {
        _id: ObjectId("688361947e8253c25c4c20c9"),
        created_at: ISODate("2025-07-25T10:48:22.564Z"),
        updated_at: ISODate("2025-07-25T10:51:00.806Z"),
        title: "一篇测试文章",
        content: "这是一篇测试文章，赶紧开始使用吧",
        description: "测试文章描述",
        author: "codepzj",
        alias: "test-article",
        category_id: ObjectId("688361837e8253c25c4c20c7"),
        tags_id: [
            ObjectId("6883618b7e8253c25c4c20c8")
        ],
        is_publish: true,
        is_top: true,
        thumbnail: ""
    }
]);

// 批量插入标签数据
db.label.insertMany([
    {
        _id: ObjectId("688361837e8253c25c4c20c7"),
        type: "category",
        name: "后端开发"
    },
    {
        _id: ObjectId("6883618b7e8253c25c4c20c8"),
        type: "tag",
        name: "golang"
    }
]);

// 批量插入文档数据
db.document.insertMany([
    {
        _id: ObjectId("688362fe5fc67ae7db9e8893"),
        created_at: ISODate("2025-07-25T10:57:02.133Z"),
        updated_at: ISODate("2025-07-26T06:33:52.000Z"),
        title: "测试文档",
        description: "测试文档描述",
        thumbnail: "http://localhost:9002/api/images/1753500952NVRtKGpkIN.jpg",
        alias: "test-docs",
        sort: Int32("1"),
        is_public: true,
        is_deleted: false
    }
]);

// 批量插入网站配置数据
db.config.insertMany([
    {
        _id: ObjectId("688364000000000000000001"),
        created_at: ISODate("2025-07-25T10:59:00.000Z"),
        updated_at: ISODate("2025-07-25T10:59:00.000Z"),
        type: "home",
        content: {
            title: "主页",
            description: "欢迎来到我的个人网站",
            avatar: "https://github.githubassets.com/assets/pull-shark-default-498c279a747d.png",
            name: "codepzj",
            bio: "一个热爱技术的开发者",
            github: "https://github.com/codepzj",
            blog: "",
            location: "中国",
            tech_stacks: ["Go", "Vue.js", "React", "TypeScript"],
            quote: "代码改变世界",
            show_recent_posts: true,
            recent_posts_count: 5
        }
    },
    {
        _id: ObjectId("688364000000000000000002"),
        created_at: ISODate("2025-07-25T10:59:00.000Z"),
        updated_at: ISODate("2025-07-25T10:59:00.000Z"),
        type: "about",
        content: {
            title: "关于我",
            description: "了解更多关于我的信息",
            skills: [
                {
                    category: "编程语言",
                    items: ["Go", "JavaScript", "TypeScript", "Python"]
                },
                {
                    category: "前端技术",
                    items: ["Vue.js", "React", "HTML", "CSS"]
                },
                {
                    category: "后端技术",
                    items: ["Gin", "MongoDB", "Redis", "Docker"]
                }
            ],
            timeline: [
                {
                    year: "2020",
                    title: "开始编程之旅",
                    desc: "学习编程基础知识"
                },
                {
                    year: "2022",
                    title: "全栈开发",
                    desc: "掌握前后端开发技能"
                },
                {
                    year: "2024",
                    title: "开源贡献",
                    desc: "开始参与开源项目"
                }
            ],
            interests: ["阅读", "运动", "音乐", "旅行"],
            focus_items: ["学习新技术", "提升代码质量", "参与开源项目"]
        }
    }
]);

// 批量插入文档内容数据
db.document_content.insertMany([
    {
        _id: ObjectId("688363165fc67ae7db9e8894"),
        created_at: ISODate("2025-07-25T10:57:26.365Z"),
        updated_at: ISODate("2025-07-26T04:17:50.304Z"),
        document_id: ObjectId("688362fe5fc67ae7db9e8893"),
        title: "测试文档1",
        content: "## 1 系统概述\n\n### 1.1 背景\n\n随着中国企业数字化转型加速...",
        description: "测试文档描述1",
        alias: "test01",
        parent_id: ObjectId("688362fe5fc67ae7db9e8893"),
        is_dir: false,
        sort: Int32("1"),
        is_deleted: false
    },
    {
        _id: ObjectId("688363405fc67ae7db9e8895"),
        created_at: ISODate("2025-07-25T10:58:08.077Z"),
        updated_at: ISODate("2025-07-25T16:35:08.787Z"),
        document_id: ObjectId("688362fe5fc67ae7db9e8893"),
        title: "测试目录2",
        content: "",
        description: "",
        alias: "testdir03",
        parent_id: ObjectId("688362fe5fc67ae7db9e8893"),
        is_dir: true,
        sort: Int32("1"),
        is_deleted: false
    },
    {
        _id: ObjectId("6883b361cb7e6a4f53e28035"),
        created_at: ISODate("2025-07-25T16:40:01.79Z"),
        updated_at: ISODate("2025-07-26T04:59:21.289Z"),
        document_id: ObjectId("688362fe5fc67ae7db9e8893"),
        title: "测试文档3",
        content: "## 测试文档3",
        description: "测试文档3描述",
        alias: "test03",
        parent_id: ObjectId("688362fe5fc67ae7db9e8893"),
        is_dir: false,
        sort: Int32("1"),
        is_deleted: false
    },
    {
        _id: ObjectId("6883b45ccb7e6a4f53e28036"),
        created_at: ISODate("2025-07-25T16:44:12.791Z"),
        updated_at: ISODate("2025-07-25T16:44:27.739Z"),
        document_id: ObjectId("688362fe5fc67ae7db9e8893"),
        title: "测试文档2-1",
        content: "## 测试文档2-1",
        description: "",
        alias: "test02",
        parent_id: ObjectId("688363405fc67ae7db9e8895"),
        is_dir: false,
        sort: Int32("1"),
        is_deleted: false
    }
]);
EOF