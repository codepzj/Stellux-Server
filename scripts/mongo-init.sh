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

// 文章表创建文本索引
db.post.createIndex({
    "title": "text",
    "description": "text",
    "content": "text",
});

// 文档表创建文本索引
db.document.createIndex({
    "title": "text",
    "description": "text",
    "content": "text",
});

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
db.label.insertMany([{
    "_id": ObjectId("67c453eda04b00c407b43197"),
    "label_type": "category",
    "name": "默认分类",
}]);

// 初始化标签
db.label.insertMany([{
    "_id": ObjectId("67c453eda04b00c407b43199"),
    "label_type": "tag",
    "name": "默认标签1",
}, {
    "_id": ObjectId("67c453eda04b00c407b43200"),
    "label_type": "tag",
    "name": "默认标签2",
}]);

// 初始化文章
db.post.insertMany([{
    "_id": ObjectId("67c453eda04b00c407b43202"),
    "created_at": new Date(),
    "updated_at": new Date(),
    "author": "codepzj",
    "title": "stellux知识库系统",
    "content": "如果你看到这篇文章,说明你已经成功安装了stellux,接下来你可以开始你的知识库之旅了😀",
    "description": "懂得都懂",
    "category_id": ObjectId("67c453eda04b00c407b43197"),
    "tags_id": [ObjectId("67c453eda04b00c407b43199"), ObjectId("67c453eda04b00c407b43201")],
    "thumbnail": "https://image.codepzj.cn/image/20250526184556014.png",
    "is_publish": true,
    "is_top": true,
}]);

// 初始化文档
db.document.insert( {
    _id: ObjectId("68382e116d357d131691114c"),
    created_at: ISODate("2025-05-29T09:51:13.468Z"),
    updated_at: ISODate("2025-05-29T10:23:55.88Z"),
    title: "测试文档",
    content: "",
    alias: "test",
    description: "一篇测试文档~",
    thumbnail: "https://image.codepzj.cn/image/20250529183152057.png",
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
    thumbnail: "",
    document_type: "leaf",
    is_public: false,
    parent_id: ObjectId("68382e116d357d131691114c"),
    document_id: ObjectId("68382e116d357d131691114c")
} );

// 初始化设置
db.setting.insert( {
    _id: ObjectId("6838191c8bc94e687f0a156d"),
    key: "basic_setting",
    value: {
        site_title: "stellux知识库",
        site_subtitle: "记录成长,点亮星途",
        site_favicon: "/favicon.ico"
    }
} );
db.setting.insert( {
    _id: ObjectId("68381dac8bc94e687f0a156e"),
    key: "seo_setting",
    value: {
        site_keywords: "知识库,golang,vue,nextjs,stellux,mongodb",
        twitter_card: "summary_large_image",
        site_description: "stellux是我的个人知识库,主要用来记录一些golang的零碎知识点以及文档,并且会记录一些生活,工作,学习中的点点滴滴,该博客会持续更新,欢迎关注,使用开源项目stellux构建",
        robots: "index,follow",
        og_image: "https://image.codepzj.cn/image/20250526184556014.png",
        site_author: "浩瀚星河",
        site_url: "https://gowiki.site",
        og_type: "website",
        twitter_site: "codepzj"
    }
} );
db.setting.insert( {
    _id: ObjectId("683822058bc94e687f0a156f"),
    key: "blog_setting",
    value: {
        blog_title: "浩瀚星河",
        blog_subtitle: "代码,日常,人生",
        blog_avatar: "https://image.codepzj.cn/image/20250529174726187.jpeg"
    }
} );

db.setting.insert({
    _id: ObjectId("683b4faac5451e983d5649d0"),
    key: "about_setting",
    value: {
        github_username: "codepzj",
        author: "浩瀚星河",
        avatar_url: "https://image.codepzj.cn/image/20250529174726187.jpeg",
        left_tags: [
            "🧠 技术探索者",
            "🛠️ 创意实践者",
            "🌐 架构与开发者"
        ],
        right_tags: [
            "兴趣点燃灵感火花 ✨",
            "开源协作推动者 🧐",
            "热情永不熄灭 🔥"
        ],
        know_me: "https://github.com/codepzj"
    }
} );


EOF