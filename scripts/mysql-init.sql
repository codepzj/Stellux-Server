/*!40101 SET NAMES utf8 */;

CREATE DATABASE stellux;

USE stellux;

/* 创建文章表 */
CREATE TABLE post (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp COMMENT '逻辑删除标记',
    `title` varchar(255) NOT NULL COMMENT '标题',
    `content` text NOT NULL COMMENT '内容',
    `description` varchar(255) NOT NULL COMMENT '描述',
    `alias` varchar(255) NOT NULL COMMENT '别名',
    `category_id` bigint(32) NOT NULL DEFAULT '0' COMMENT '分类id',
    `is_publish` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否发布',
    `is_top` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否置顶',
    `thumbnail` varchar(255) NOT NULL COMMENT '缩略图',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_alias` (`alias`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';

/* 创建文章分类表 */
CREATE TABLE post_category (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp COMMENT '逻辑删除标记',
    `category_name` varchar(255) NOT NULL COMMENT '分类名称',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_category_name` (`category_name`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类表';

/* 创建文档表 */
CREATE TABLE document (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp COMMENT '逻辑删除标记',
    `title` varchar(255) NOT NULL COMMENT '标题',
    `description` varchar(255) NOT NULL COMMENT '描述',
    `alias` varchar(255) NOT NULL COMMENT '别名',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    `is_public` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否公开',
    `thumbnail` varchar(255) NOT NULL COMMENT '缩略图',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_alias` (`alias`),
    KEY `idx_is_public` (`is_public`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档表';

/* 创建文档内容表 */
CREATE TABLE document_content (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp COMMENT '逻辑删除标记',
    `document_id` bigint(32) NOT NULL COMMENT '文档id',
    `title` varchar(255) NOT NULL COMMENT '标题',
    `content` text NOT NULL COMMENT '内容',
    `description` varchar(255) NOT NULL COMMENT '描述',
    `alias` varchar(255) NOT NULL COMMENT '别名',
    `parent_id` bigint(32) NOT NULL COMMENT '父级id',
    `is_dir` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否是目录',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_alias` (`alias`),
    KEY `idx_document_id` (`document_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档内容表';

/* 创建本地文件表 */
CREATE TABLE file (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp COMMENT '逻辑删除标记',
    `file_name` varchar(255) NOT NULL COMMENT '文件名称',
    `url` varchar(255) NOT NULL COMMENT '文件url',
    `dst` varchar(255) NOT NULL COMMENT '文件路径',
    PRIMARY KEY (`id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='本地文件表';


/* 创建友情链接表 */
CREATE TABLE friend (
    `id` bigint(32) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp COMMENT '逻辑删除标记',
    `name` varchar(255) NOT NULL COMMENT '名称',
    `description` varchar(255) NOT NULL COMMENT '描述',
    `site_url` varchar(255) NOT NULL COMMENT '站点url',
    `avatar_url` varchar(255) NOT NULL COMMENT '头像url',
    `website_type` int(11) NOT NULL DEFAULT '0' COMMENT '网站类型',
    `is_active` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否激活',
    PRIMARY KEY (`id`),
    KEY `idx_is_active` (`is_active`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='友情链接表';

