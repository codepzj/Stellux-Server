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

db.getCollection("user").insert( {
    _id: ObjectId("67c453eda04b00c407b431fd"),
    username: "admin",
    password: "$2a$10$SLcnDmaJc1nLtUOsZS4yquXyVeu5E6qJHNTVeKSzTk4JO4Xq/FPSy",
    nickname: "codezj",
    role_id: Int32("0"),
    created_at: ISODate("2025-07-13T02:19:26.865Z"),
    updated_at: ISODate("2025-07-13T02:19:26.865Z"),
    avatar: "https://github.githubassets.com/assets/pull-shark-default-498c279a747d.png",
    email: "admin@example.com"
} );

db.getCollection("post").insert( {
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
} );

db.getCollection("label").insert( {
    _id: ObjectId("688361837e8253c25c4c20c7"),
    type: "category",
    name: "后端开发"
} );
db.getCollection("label").insert( {
    _id: ObjectId("6883618b7e8253c25c4c20c8"),
    type: "tag",
    name: "golang"
} );

db.getCollection("document").insert( {
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
} );

db.getCollection("document_content").insert( {
    _id: ObjectId("688363165fc67ae7db9e8894"),
    created_at: ISODate("2025-07-25T10:57:26.365Z"),
    updated_at: ISODate("2025-07-26T04:17:50.304Z"),
    document_id: ObjectId("688362fe5fc67ae7db9e8893"),
    title: "测试文档1",
    content: "## 1 系统概述\n\n### 1.1 背景\n\n随着中国企业数字化转型加速，客服需求激增，传统第三方客服平台存在隐私泄露风险和高昂的 PaaS 平台接入成本。企业内部知识系统（如 ERP、CRM）通常孤立，难以与客服系统高效集成，导致服务效率低下。设计一个智能客服系统，旨在通过自建系统保障数据隐私，降低长期运营成本，并打通企业内部知识系统，提供高效、智能的客服体验，满足中国企业对数据安全和成本控制的迫切需求。\n\n### 1.2 系统功能\n\n本系统是一个集 AI 机器人、人工客服和视频会议功能于一体的智能客服平台，主要功能包括：\n\n- **AI 机器人客服**：基于自然语言处理（NLP）提供 24/7 自动问答，处理常见问题。\n- **人工客服转接**：支持 AI 转人工，1 对 1 文本或语音交互。\n- **高级用户视频会议**：为高级用户提供 1 对 1 WebRTC 视频对接。\n- **企业知识库集成**：打通企业内部 CRM、ERP 系统，实现知识查询和个性化服务。\n- **多端支持**：支持 Web 和移动端，适配企业员工和客户使用场景。\n\n**与现有系统的不同**：\n\n- 隐私保护：数据完全托管于企业内部，符合中国《数据安全法》和《个人信息保护法》。\n- 成本优化：避免第三方 PaaS 平台高额订阅费用。\n- 灵活性：支持企业知识库无缝集成，定制化程度高。\n\n## 2 总体设计\n\n以下是系统的总体功能架构图（使用 Mermaid 绘制）：\n\n```mermaid\ngraph TD\n    A[用户端] -->|Web/移动端| B[前端界面: Vue]\n    B -->|API 调用| C[后端服务: Go]\n    C -->|NLP 处理| D[AI 客服模块]\n    C -->|实时通信| E[WebRTC 视频/语音模块]\n    C -->|数据查询| F[数据库: MySQL]\n    C -->|知识库集成| G[企业内部系统: CRM/ERP]\n    D -->|复杂问题| H[人工客服模块]\n    E -->|高级用户| I[视频会议]\n    F -->|用户信息/记录| G\n```\n\n## 3 功能设计\n\n### 3.1 AI 客服模块\n\n**功能描述**：\n\n- 提供 24/7 自动问答，处理常见问题（如订单查询、产品信息）。\n- 支持多轮对话，理解用户意图。\n- 复杂问题自动转接人工客服。\n  **技术描述**：\n- 使用 Go 开发后端服务，集成开源 NLP 模型（如 Rasa 或 BERT 模型的轻量化版本）。\n- 通过 RESTful API 与前端交互。\n- 数据库查询使用 MySQL，存储常见问题和对话记录。\n\n### 3.2 人工客服模块\n\n**功能描述**：\n\n- 支持 AI 转人工，1 对 1 文本或语音交互。\n- 人工客服可查看用户历史对话记录和企业知识库信息。\n  **技术描述**：\n- Go 实现 WebSocket 实时通信，处理文本和语音消息。\n- 集成企业内部 CRM 系统，查询用户信息。\n- 使用 Redis 缓存对话记录，提升响应速度。\n\n### 3.3 视频会议模块\n\n**功能描述**：\n\n- 为高级用户提供 1 对 1 视频会议功能。\n- 支持实时音视频交互，屏幕共享。\n- 使用 WebRTC 实现点对点音视频通信，Go 提供信令服务器。\n- Vue 前端集成 WebRTC 客户端，支持多端适配。\n- 确保视频流加密，符合数据隐私要求。\n\n### 3.4 知识库集成模块\n\n**功能描述**：\n\n- 打通企业内部 CRM、ERP 系统，提供个性化客服响应。\n- 支持客服查询订单、库存、员工信息等。\n  **技术描述**：\n- Go 开发 API 网关，集成企业内部系统（如 SAP、用友）。\n- 使用 gRPC 或 RESTful API 实现高效数据交互。\n- MySQL 存储临时数据，Redis 缓存频繁查询结果。\n\n### 3.5 用户管理模块\n\n**功能描述**：\n\n- 用户注册、登录、权限管理。\n- 支持普通用户和高级用户分级服务。\n- Go 开发 JWT 认证，保障用户登录安全。\n- Vue 前端实现响应式登录注册页面。\n- MySQL 存储用户数据，Redis 缓存会话信息。\n\n## 4 数据库设计\n\n### 4.1 数据库技术\n\n**技术选择**：MySQL + Redis\n\n- **MySQL**：作为主数据库，存储用户信息、对话记录、知识库索引。MySQL 成熟稳定，适合结构化数据，社区支持广泛。\n- **Redis**：作为缓存，存储会话数据和频繁查询的知识库结果，提升响应速度，降低数据库压力。\n  **原因**：\n- MySQL 满足企业级数据存储需求，支持复杂查询，易于与 Go 集成。\n- Redis 提供高性能缓存，适合实时客服场景。\n- 两者结合平衡了性能和开发复杂度，适合毕业设计周期。\n\n### 4.2 表格设计\n\n以下是初步设计的表格及其关系：\n\n- **users**：存储用户信息（用户 ID、用户名、密码哈希、角色、注册时间）。\n- **conversations**：存储对话记录（会话 ID、用户 ID、客服类型、时间戳、内容）。\n- **knowledge_base**：存储企业知识库索引（知识 ID、分类、内容、关联系统）。\n- **video_sessions**：存储视频会议记录（会话 ID、用户 ID、客服 ID、开始时间、结束时间）。\n  **关系**：\n- users 一对多 conversations（一个用户有多个会话）。\n- conversations 一对多 knowledge_base（一个会话可能引用多个知识点）。\n- users 一对多 video_sessions（高级用户可发起多个视频会话）。\n\n## 5 用户分类\n\n### 5.1 管理员\n\n**功能描述**：\n\n- 管理整个数据库的增删改查（用户、会话、知识库）。\n- 配置 AI 模型参数，分配人工客服。\n- 查看系统运行状态和统计报表。\n- Vue 实现管理员仪表板，Go 提供管理 API。\n- MySQL 存储管理日志，Redis 缓存报表数据。\n\n### 5.2 普通用户\n\n**功能描述**：\n\n- 注册、登录，查看个人信息和历史对话。\n- 使用 AI 客服，必要时转接人工客服。\n- Vue 响应式前端，Go JWT 认证。\n- MySQL 存储用户信息，Redis 缓存会话。\n\n### 5.3 高级用户\n\n**功能描述**：\n\n- 除普通用户功能外，可发起 1 对 1 视频会议。\n- 优先接入人工客服，获取个性化服务。\n- WebRTC 实现视频功能，Go 提供信令服务。\n- Vue 前端集成 WebRTC 客户端。\n\n## 6 计划表\n\n以下是系统设计的甘特图（使用 Mermaid 绘制），从 2025 年 9 月到 2026 年 3 月，约 7 个月：\n\n```mermaid\ngantt\n    title 智能客服系统开发计划\n    dateFormat  YYYY-MM-DD\n    section 需求分析\n    需求收集与分析         :a1, 2025-09-01, 14d\n    概要设计文档         :a2, 2025-09-15, 10d\n    section 技术选型\n    技术调研与选型       :b1, 2025-09-25, 10d\n    section 后端开发\n    数据库设计与实现     :c1, 2025-10-05, 20d\n    API 开发与测试       :c2, 2025-10-25, 30d\n    WebRTC 信令服务器    :c3, 2025-11-24, 20d\n    section 前端开发\n    Vue 界面开发         :d1, 2025-10-25, 30d\n    WebRTC 客户端集成    :d2, 2025-11-24, 15d\n    section 系统集成\n    知识库集成           :e1, 2025-12-09, 20d\n    AI 模型部署          :e2, 2025-12-29, 15d\n    section 测试与优化\n    单元测试与集成测试   :f1, 2026-01-13, 30d\n    性能优化             :f2, 2026-02-12, 15d\n    section 文档与答辩\n    完善文档与准备答辩   :g1, 2026-02-27, 20d\n```",
    description: "测试文档描述1",
    alias: "test01",
    parent_id: ObjectId("688362fe5fc67ae7db9e8893"),
    is_dir: false,
    sort: Int32("1"),
    is_deleted: false
} );
db.getCollection("document_content").insert( {
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
} );
db.getCollection("document_content").insert( {
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
} );
db.getCollection("document_content").insert( {
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
} );

EOF