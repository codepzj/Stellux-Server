package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Config 网站配置
type Config struct {
	Id        bson.ObjectID `bson:"_id,omitempty"` // 配置ID
	CreatedAt time.Time     `bson:"created_at"`    // 创建时间
	UpdatedAt time.Time     `bson:"updated_at"`    // 更新时间
	Type      string        `bson:"type"`          // 配置类型: home, about
	Content   Content       `bson:"content"`       // 网站内容配置
}

// Content 网站内容配置
type Content struct {
	// 通用配置
	Title       string `bson:"title"`       // 网站标题
	Description string `bson:"description"` // 网站描述

	// 主页配置
	Avatar           string   `bson:"avatar,omitempty"`             // 头像URL
	Name             string   `bson:"name,omitempty"`               // 姓名
	Bio              string   `bson:"bio,omitempty"`                // 个人简介
	Github           string   `bson:"github,omitempty"`             // GitHub地址
	Blog             string   `bson:"blog,omitempty"`               // 博客地址
	Location         string   `bson:"location,omitempty"`           // 位置
	TechStacks       []string `bson:"tech_stacks,omitempty"`        // 技术栈
	Repositories     []Repo   `bson:"repositories,omitempty"`       // 开源项目
	Quote            string   `bson:"quote,omitempty"`              // 名言
	ShowRecentPosts  bool     `bson:"show_recent_posts,omitempty"`  // 是否显示最新文章
	RecentPostsCount int      `bson:"recent_posts_count,omitempty"` // 最新文章数量

	// About页面配置
	Skills     []Skill    `bson:"skills,omitempty"`      // 技能配置
	Timeline   []Timeline `bson:"timeline,omitempty"`    // 时间线
	Interests  []string   `bson:"interests,omitempty"`   // 兴趣爱好
	FocusItems []string   `bson:"focus_items,omitempty"` // 当前专注事项
}

// Repo 开源项目
type Repo struct {
	Name string `bson:"name"` // 项目名称
	URL  string `bson:"url"`  // 项目地址
	Desc string `bson:"desc"` // 项目描述
}

// Skill 技能
type Skill struct {
	Category string   `bson:"category"` // 分类
	Items    []string `bson:"items"`    // 技能列表
}

// Timeline 时间线项目
type Timeline struct {
	Year  string `bson:"year"`  // 年份
	Title string `bson:"title"` // 标题
	Desc  string `bson:"desc"`  // 描述
}
