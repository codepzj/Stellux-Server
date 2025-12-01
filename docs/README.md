# Stellux-Server 文档中心

## 📚 文档导航

### 日志体系
- 📖 [日志规范指南](./LOGGING_GUIDE.md) - 完整的日志规范和最佳实践
- 🚀 [快速实施指南](./QUICK_LOGGING_GUIDE.md) - 5分钟快速为Service添加日志
- 📊 [实施总结](./LOGGING_IMPLEMENTATION_SUMMARY.md) - 日志体系实施情况和进度

### 快速开始

#### 查看日志规范
```bash
# 阅读完整规范
cat docs/LOGGING_GUIDE.md

# 查看快速指南
cat docs/QUICK_LOGGING_GUIDE.md
```

#### 为Service添加日志
1. 阅读 [快速实施指南](./QUICK_LOGGING_GUIDE.md)
2. 参考示例: `internal/user/internal/service/user.go`
3. 按模板添加日志
4. 运行 `go build` 验证

### 示例代码
- ✅ User Service: `internal/user/internal/service/user.go`
- ✅ Post Service: `internal/post/internal/service/post.go`

---

**更新时间**: 2025-06-01
