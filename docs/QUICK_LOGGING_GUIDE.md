# 快速添加Service层日志指南

> 为其他Service模块快速添加日志的实操指南

## 📋 待完成模块

- [x] user service ✅
- [x] post service ✅
- [ ] document service
- [ ] document_content service
- [ ] comment service
- [ ] file service
- [ ] friend service
- [ ] label service

---

## 🚀 5分钟快速上手

### 步骤1: 导入logger包 (30秒)

在service文件顶部添加logger导入：

```go
import (
    "context"
    "errors"
    
    "github.com/codepzj/Stellux-Server/internal/pkg/logger"  // 👈 添加这行
    // ... 其他imports
)
```

### 步骤2: 为每个方法添加日志 (每个方法1-2分钟)

#### 模板A: Create/Update/Delete 方法

```go
func (s *Service) Create(ctx context.Context, entity *domain.Entity) error {
    // 1️⃣ 入口日志
    logger.Info("开始创建XXX",
        logger.WithString("method", "Create"),
        logger.WithString("name", entity.Name),  // 关键字段
    )
    
    // 业务逻辑...
    err := s.repo.Create(ctx, entity)
    
    // 2️⃣ 错误日志
    if err != nil {
        logger.Error("创建XXX失败",
            logger.WithError(err),
            logger.WithString("name", entity.Name),
        )
        return err
    }
    
    // 3️⃣ 成功日志
    logger.Info("创建XXX成功",
        logger.WithString("entityId", entity.ID.Hex()),
        logger.WithString("name", entity.Name),
    )
    
    return nil
}
```

#### 模板B: Query 方法

```go
func (s *Service) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
    // 1️⃣ 入口日志
    logger.Info("查询XXX",
        logger.WithString("method", "GetByID"),
        logger.WithString("entityId", id),
    )
    
    entity, err := s.repo.GetByID(ctx, id)
    
    // 2️⃣ 区分错误类型
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            logger.Warn("XXX不存在",  // 👈 业务异常用Warn
                logger.WithString("entityId", id),
            )
        } else {
            logger.Error("查询XXX失败",  // 👈 系统错误用Error
                logger.WithError(err),
                logger.WithString("entityId", id),
            )
        }
        return nil, err
    }
    
    // 3️⃣ 成功日志
    logger.Info("查询XXX成功",
        logger.WithString("entityId", id),
    )
    
    return entity, nil
}
```

#### 模板C: List 方法

```go
func (s *Service) GetList(ctx context.Context, page *domain.Page) ([]*domain.Entity, int64, error) {
    // 1️⃣ 入口日志
    logger.Info("查询XXX列表",
        logger.WithString("method", "GetList"),
        logger.WithInt("pageNo", int(page.PageNo)),
        logger.WithInt("pageSize", int(page.PageSize)),
    )
    
    entities, total, err := s.repo.GetList(ctx, page)
    
    // 2️⃣ 错误日志
    if err != nil {
        logger.Error("查询列表失败",
            logger.WithError(err),
        )
        return nil, 0, err
    }
    
    // 3️⃣ 成功日志（带统计）
    logger.Info("查询列表成功",
        logger.WithInt("count", len(entities)),
        logger.WithInt("total", int(total)),
    )
    
    return entities, total, nil
}
```

---

## 📝 日志字段速查表

### 常用方法名映射

| 操作 | method 字段值 | 示例 |
|------|--------------|------|
| 创建 | `AdminCreate`, `Create` | `logger.WithString("method", "AdminCreate")` |
| 更新 | `AdminUpdate`, `Update` | `logger.WithString("method", "AdminUpdate")` |
| 删除 | `AdminDelete`, `Delete` | `logger.WithString("method", "AdminDelete")` |
| 查询 | `GetByID`, `GetList` | `logger.WithString("method", "GetByID")` |
| 搜索 | `Search`, `GetByKeyword` | `logger.WithString("method", "Search")` |

### 常用字段类型

```go
// 字符串字段
logger.WithString("userId", user.ID.Hex())
logger.WithString("username", user.Username)
logger.WithString("title", post.Title)

// 整数字段
logger.WithInt("count", len(items))
logger.WithInt("pageNo", int(page.PageNo))

// 错误字段
logger.WithError(err)

// 任意类型字段
logger.WithAny("isPublish", isPublish)
logger.WithAny("status", status)
```

---

## ✅ 检查清单

为service添加日志后，检查以下项：

### 必须项 ✓
- [ ] 每个public方法都有入口日志
- [ ] 所有错误都记录了日志
- [ ] 所有成功操作都记录了日志
- [ ] 使用了正确的日志级别（Info/Warn/Error）

### 推荐项 ⭐
- [ ] 入口日志包含method字段
- [ ] 错误日志使用`logger.WithError(err)`
- [ ] 关键操作记录了业务实体ID
- [ ] 列表查询记录了count/total

### 避免项 ❌
- [ ] 不要在循环中大量打印日志
- [ ] 不要记录敏感信息（密码、token）
- [ ] 不要使用字符串拼接
- [ ] 不要遗漏错误处理的日志

---

## 🎯 针对特定模块的提示

### Document Service
```go
// 关键字段
logger.WithString("documentId", doc.ID.Hex())
logger.WithString("title", doc.Title)
logger.WithString("alias", doc.Alias)
```

### Document Content Service
```go
// 关键字段
logger.WithString("contentId", content.ID.Hex())
logger.WithString("documentId", content.DocumentId.Hex())
logger.WithString("parentId", content.ParentId.Hex())
```

### Comment Service
```go
// 关键字段
logger.WithString("commentId", comment.ID.Hex())
logger.WithString("postId", comment.PostId.Hex())
logger.WithString("userId", comment.UserId.Hex())
```

### File Service
```go
// 关键字段
logger.WithString("fileId", file.ID.Hex())
logger.WithString("filename", file.Filename)
logger.WithInt("fileSize", int(file.Size))
```

### Friend Service
```go
// 关键字段
logger.WithString("friendId", friend.ID.Hex())
logger.WithString("name", friend.Name)
logger.WithString("url", friend.Url)
```

### Label Service
```go
// 关键字段
logger.WithString("labelId", label.ID.Hex())
logger.WithString("name", label.Name)
logger.WithString("type", label.Type)
```

---

## 🔧 常见问题

### Q: 是否需要为private方法添加日志？
**A**: 不需要。只为public方法（接口方法）添加日志。private辅助方法不需要。

### Q: 日志太多会影响性能吗？
**A**: 适度的日志（按本指南添加）不会有明显影响。避免在循环中打印大量日志。

### Q: 如何处理批量操作的日志？
**A**: 只记录批量操作的汇总信息，不要为每个item打印日志。
```go
logger.Info("批量删除成功",
    logger.WithInt("count", len(ids)),  // ✅ 只记录数量
)
```

### Q: bson.ObjectID如何转为字符串？
**A**: 使用`.Hex()`方法
```go
logger.WithString("userId", user.ID.Hex())  // ✅
// 不要直接使用: logger.WithString("userId", user.ID) // ❌
```

---

## 📊 完成后验证

### 1. 编译检查
```bash
cd /Users/pzj/Desktop/stellux/Stellux-Server
go build
```

### 2. 查看日志输出
启动服务后调用API，观察日志格式是否正确：
```bash
# 开发环境日志示例
2025-06-01 12:00:00 INFO  开始创建用户  method=AdminCreate username=admin
2025-06-01 12:00:00 INFO  创建用户成功  userId=507f1f77bcf86cd799439011 username=admin
```

### 3. 错误场景测试
故意触发错误，检查错误日志是否完整：
```bash
# 错误日志示例
2025-06-01 12:00:00 ERROR 创建用户失败  error="用户已存在" username=admin
```

---

## 📚 参考资源

- [完整日志规范](./LOGGING_GUIDE.md)
- [User Service示例](../internal/user/internal/service/user.go)
- [Post Service示例](../internal/post/internal/service/post.go)

---

## 🎓 下一步

1. ✅ 为当前模块的service添加日志
2. ✅ 编译并测试
3. ✅ 提交代码
4. 🔄 重复以上步骤，完成其他模块

---

**预计时间**: 每个service文件 10-15分钟  
**维护者**: Stellux Team  
**更新时间**: 2025-06-01
