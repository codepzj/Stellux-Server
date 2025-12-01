# Stellux-Server 日志规范指南

## 📋 目录

- [日志级别](#日志级别)
- [日志场景](#日志场景)
- [Service层日志规范](#service层日志规范)
- [日志字段规范](#日志字段规范)
- [最佳实践](#最佳实践)
- [示例代码](#示例代码)

---

## 🎯 日志级别

### Debug
- **用途**: 详细的调试信息
- **场景**: 开发阶段追踪程序执行流程
- **示例**: 打印函数参数、中间变量值

### Info
- **用途**: 重要的业务流程信息
- **场景**: 
  - 方法调用开始/结束
  - 业务操作成功
  - 关键数据查询
- **示例**: "用户登录成功", "创建文章成功"

### Warn
- **用途**: 警告信息，不影响核心功能
- **场景**:
  - 业务逻辑异常但可恢复
  - 数据不一致但不致命
  - 外部依赖超时重试
- **示例**: "用户不存在", "缓存未命中"

### Error
- **用途**: 错误信息，影响功能正常运行
- **场景**:
  - 数据库操作失败
  - 外部服务调用失败
  - 数据验证失败
- **示例**: "数据库查询失败", "密码加密失败"

### Fatal
- **用途**: 致命错误，程序无法继续运行
- **场景**: 系统初始化失败、配置加载失败
- **示例**: "数据库连接失败"

---

## 📍 日志场景

### Service层必须记录的日志

#### 1. 方法入口 (Info级别)
```go
logger.Info("开始创建用户",
    logger.WithString("username", user.Username),
    logger.WithString("method", "AdminCreate"),
)
```

#### 2. 方法成功 (Info级别)
```go
logger.Info("创建用户成功",
    logger.WithString("username", user.Username),
    logger.WithString("userId", id),
)
```

#### 3. 业务验证失败 (Warn级别)
```go
logger.Warn("用户已存在",
    logger.WithString("username", user.Username),
)
```

#### 4. 系统错误 (Error级别)
```go
logger.Error("查询用户失败",
    logger.WithError(err),
    logger.WithString("username", user.Username),
)
```

---

## 🎨 Service层日志规范

### 标准模板

```go
package service

import (
    "context"
    "github.com/codepzj/Stellux-Server/internal/pkg/logger"
    // ... other imports
)

type UserService struct {
    repo repository.IUserRepository
}

func (s *UserService) AdminCreate(ctx context.Context, user *domain.User) error {
    // 1. 记录方法入口
    logger.Info("开始创建用户",
        logger.WithString("method", "AdminCreate"),
        logger.WithString("username", user.Username),
    )
    
    // 2. 业务逻辑
    u, err := s.repo.GetByUsername(ctx, user.Username)
    if err != nil && err != mongo.ErrNoDocuments {
        // 3. 记录系统错误
        logger.Error("查询用户失败",
            logger.WithError(err),
            logger.WithString("username", user.Username),
        )
        return err
    }
    
    if u != nil {
        // 4. 记录业务验证失败
        logger.Warn("用户已存在",
            logger.WithString("username", user.Username),
        )
        return errors.New("用户已存在")
    }
    
    // 加密密码
    user.Password, err = utils.GenerateHashPassword(user.Password)
    if err != nil {
        logger.Error("密码加密失败",
            logger.WithError(err),
            logger.WithString("username", user.Username),
        )
        return err
    }
    
    // 创建用户
    id, err := s.repo.Create(ctx, user)
    if err != nil {
        logger.Error("创建用户失败",
            logger.WithError(err),
            logger.WithString("username", user.Username),
        )
        return err
    }
    
    // 5. 记录操作成功
    logger.Info("创建用户成功",
        logger.WithString("username", user.Username),
        logger.WithString("userId", id),
    )
    
    return nil
}
```

---

## 🏷️ 日志字段规范

### 必选字段
- `method`: 方法名
- `error`: 错误信息（使用`logger.WithError(err)`）

### 推荐字段
- `username`: 用户名
- `userId`: 用户ID
- `postId`: 文章ID
- `documentId`: 文档ID
- `operation`: 操作类型（create/update/delete）
- `duration`: 操作耗时（毫秒）

### 字段命名规范
- 使用驼峰命名: `userId`, `postTitle`
- 保持一致性: 整个项目统一使用相同的字段名
- 避免敏感信息: 不要记录密码、token等

---

## ✅ 最佳实践

### 1. 关键操作必须记录
- ✅ 用户登录/注册
- ✅ 数据创建/更新/删除
- ✅ 权限验证
- ✅ 外部API调用

### 2. 错误必须记录完整信息
```go
// ❌ 错误示例
logger.Error("操作失败")

// ✅ 正确示例
logger.Error("创建文章失败",
    logger.WithError(err),
    logger.WithString("method", "AdminCreatePost"),
    logger.WithString("postTitle", post.Title),
)
```

### 3. 避免过度日志
```go
// ❌ 不要在循环中大量打印
for _, item := range items {
    logger.Info("处理item", logger.WithAny("item", item)) // 会产生大量日志
}

// ✅ 记录汇总信息
logger.Info("批量处理完成",
    logger.WithInt("totalCount", len(items)),
    logger.WithInt("successCount", successCount),
)
```

### 4. 使用结构化日志
```go
// ❌ 字符串拼接
logger.Info("用户" + username + "创建成功")

// ✅ 结构化字段
logger.Info("用户创建成功",
    logger.WithString("username", username),
)
```

### 5. 日志级别使用原则
- **90%** 使用 Info 和 Error
- **5%** 使用 Warn
- **5%** 使用 Debug
- **0.1%** 使用 Fatal

---

## 📝 示例代码

### CRUD操作日志模板

#### Create
```go
func (s *Service) Create(ctx context.Context, entity *domain.Entity) error {
    logger.Info("开始创建实体",
        logger.WithString("method", "Create"),
        logger.WithString("name", entity.Name),
    )
    
    err := s.repo.Create(ctx, entity)
    if err != nil {
        logger.Error("创建实体失败",
            logger.WithError(err),
            logger.WithString("name", entity.Name),
        )
        return err
    }
    
    logger.Info("创建实体成功",
        logger.WithString("entityId", entity.ID.Hex()),
    )
    return nil
}
```

#### Update
```go
func (s *Service) Update(ctx context.Context, entity *domain.Entity) error {
    logger.Info("开始更新实体",
        logger.WithString("method", "Update"),
        logger.WithString("entityId", entity.ID.Hex()),
    )
    
    err := s.repo.Update(ctx, entity)
    if err != nil {
        logger.Error("更新实体失败",
            logger.WithError(err),
            logger.WithString("entityId", entity.ID.Hex()),
        )
        return err
    }
    
    logger.Info("更新实体成功",
        logger.WithString("entityId", entity.ID.Hex()),
    )
    return nil
}
```

#### Delete
```go
func (s *Service) Delete(ctx context.Context, id string) error {
    logger.Info("开始删除实体",
        logger.WithString("method", "Delete"),
        logger.WithString("entityId", id),
    )
    
    err := s.repo.Delete(ctx, id)
    if err != nil {
        logger.Error("删除实体失败",
            logger.WithError(err),
            logger.WithString("entityId", id),
        )
        return err
    }
    
    logger.Info("删除实体成功",
        logger.WithString("entityId", id),
    )
    return nil
}
```

#### Query
```go
func (s *Service) GetByID(ctx context.Context, id string) (*domain.Entity, error) {
    logger.Info("查询实体",
        logger.WithString("method", "GetByID"),
        logger.WithString("entityId", id),
    )
    
    entity, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            logger.Warn("实体不存在",
                logger.WithString("entityId", id),
            )
        } else {
            logger.Error("查询实体失败",
                logger.WithError(err),
                logger.WithString("entityId", id),
            )
        }
        return nil, err
    }
    
    logger.Info("查询实体成功",
        logger.WithString("entityId", id),
    )
    return entity, nil
}
```

---

## 🚀 快速开始

### 步骤1: 导入logger包
```go
import "github.com/codepzj/Stellux-Server/internal/pkg/logger"
```

### 步骤2: 在service方法中添加日志
参考上述模板，在关键位置添加日志

### 步骤3: 运行测试
```bash
go test ./internal/.../service/...
```

### 步骤4: 查看日志输出
日志会输出到控制台和配置的日志文件中

---

## 📊 日志查看

### 开发环境
- 控制台彩色输出
- 日志文件: `logs/app.log`
- 错误日志: `logs/error.log`

### 生产环境
- JSON格式输出
- 日志轮转: 每10MB切割
- 保留30天

---

## 🎓 参考资源

- [Zap Logger文档](https://pkg.go.dev/go.uber.org/zap)
- [Uber日志最佳实践](https://github.com/uber-go/guide/blob/master/style.md#logging)
- [12-Factor应用日志](https://12factor.net/logs)

---

**更新时间**: 2025-06-01  
**维护者**: Stellux Team
