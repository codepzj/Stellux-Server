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
db.config.insert( {
    _id: ObjectId("688364000000000000000001"),
    created_at: ISODate("2025-11-25T10:59:00.000Z"),
    updated_at: ISODate("2025-12-06T08:35:23.147Z"),
    type: "home",
    content: {
        SEOKeywords: null,
        SEODescription: "",
        RobotsMeta: "",
        CanonicalURL: "",
        OGTitle: "",
        Title: "主页",
        Location: "中国",
        ShowRecentPosts: true,
        RecentPostsCount: 4,
        SEOAuthor: "",
        OGDescription: "",
        OGImage: "",
        TwitterCard: "",
        Description: "欢迎来到我的个人网站",
        Bio: "一个热爱技术的开发者",
        Github: "https://github.com/codepzj",
        Blog: "https://www.golangblog.com",
        TechStacks: [
            "Go",
            "Gorm",
            "Kratos",
            "K8s",
            "Dataworks",
            "Maxcompute",
            "Kafka",
            "Redis",
            "Mysql",
            "MongoDB"
        ],
        Repositories: [
            {
                Name: "Stellux-Server",
                URL: "https://github.com/codepzj/Stellux-Server",
                Desc: "Stellux博客后端服务"
            },
            {
                Name: "hexo-graph",
                URL: "https://github.com/codepzj/hexo-graph",
                Desc: "hexo可视化插件"
            }
        ],
        Quote: "生活就像海洋,只有意志坚强的人才能到达彼岸",
        Skills: null,
        Avatar: "https://cdn.codepzj.cn/image/20250529174726187.jpeg",
        Name: "浩瀚星河",
        Motto: "低级的欲望通过放纵就可获得,高级的欲望通过自律方可获得,顶级的欲望通过煎熬才可获得。所谓自由,不是随心所欲,而是自我主宰。",
        Timeline: null,
        Interests: null,
        FocusItems: null,
        SEOTitle: ""
    }
} );
db.config.insert( {
    _id: ObjectId("688364000000000000000002"),
    created_at: ISODate("2025-11-25T10:59:00.000Z"),
    updated_at: ISODate("2025-12-06T08:28:18.482Z"),
    type: "about",
    content: {
        OGImage: "",
        Location: "",
        TechStacks: null,
        Skills: [
            {
                Category: "编程语言",
                Items: [
                    "Go",
                    "JavaScript",
                    "TypeScript",
                    "Python"
                ]
            },
            {
                Category: "前端技术",
                Items: [
                    "Vue.js",
                    "React",
                    "HTML",
                    "CSS"
                ]
            },
            {
                Category: "后端技术",
                Items: [
                    "Gin",
                    "MongoDB",
                    "Redis",
                    "Docker"
                ]
            }
        ],
        Timeline: [
            {
                Year: "2022",
                Title: "开始编程之旅",
                Desc: "学习编程基础知识"
            },
            {
                Year: "2023",
                Title: "全栈开发",
                Desc: "掌握前后端开发技能"
            },
            {
                Year: "2024",
                Title: "开源贡献",
                Desc: "开始参与开源项目"
            },
            {
                Year: "2025",
                Title: "技术探索",
                Desc: "深入学习 Go 微服务架构，实践云原生技术"
            },
            {
                Year: "如今",
                Title: "技术分享",
                Desc: "已经成为社畜🧐"
            }
        ],
        RobotsMeta: "",
        CanonicalURL: "",
        OGDescription: "",
        TwitterCard: "",
        Title: "关于我",
        Name: "",
        Github: "",
        Blog: "",
        Motto: "",
        Interests: [
            "阅读",
            "运动",
            "音乐",
            "旅行"
        ],
        SEODescription: "",
        SEOAuthor: "",
        Description: "了解更多关于我的信息",
        Bio: "",
        SEOKeywords: null,
        OGTitle: "",
        ShowRecentPosts: false,
        RecentPostsCount: 0,
        Avatar: "",
        Repositories: null,
        Quote: "",
        FocusItems: [
            "努力提升golang编程水平💪",
            "努力提升业务能力💪"
        ],
        SEOTitle: ""
    }
} );
db.config.insert( {
    _id: ObjectId("688364000000000000000003"),
    created_at: ISODate("2025-11-25T10:59:00.000Z"),
    updated_at: ISODate("2025-12-06T08:38:10.097Z"),
    type: "seo",
    content: {
        CanonicalURL: "https://www.golangblog.com/",
        OGImage: "https://cdn.codepzj.cn/image/20251206162201655.png",
        TwitterCard: "summary_large_image",
        Title: "网站配置",
        Description: "网站基础配置信息",
        Avatar: "",
        Blog: "",
        Location: "",
        Repositories: null,
        Quote: "",
        Motto: "",
        Name: "",
        TechStacks: null,
        Skills: null,
        FocusItems: null,
        SEOTitle: "浩瀚星河 - 个人技术博客",
        SEOKeywords: [
            "Go",
            "GoZero",
            "Kratos",
            "Echo",
            "Redis",
            "Mysql",
            "Pgsql",
            "Mongodb",
            "K8S",
            "微服务"
        ],
        SEODescription: "浩瀚星河的个人技术博客,记录Golang学习与开发实践。分享Go语言、微服务架构、前后端开发等技术经验。",
        SEOAuthor: "浩瀚星河",
        Bio: "",
        Github: "",
        ShowRecentPosts: false,
        RecentPostsCount: 0,
        Interests: null,
        RobotsMeta: "index,follow",
        OGTitle: "浩瀚星河 - 个人技术博客",
        OGDescription: "浩瀚星河的个人技术博客,记录Golang学习与开发实践。分享Go语言、微服务架构、前后端开发等技术经验。",
        Timeline: null
    }
} );

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