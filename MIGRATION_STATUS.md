# MongoDB 官方驱动迁移状态

## ✅ 已完成 (4/8 模块)

### 基础设施
- ✅ go.mod - 已更新移除 go-mongox
- ✅ internal/infra/mongodb.go - 返回 *mongo.Database
- ✅ 所有 wire.go 文件 - 已更新为 *mongo.Database
- ✅ 代码生成模板 - dao.tmpl, wire.tmpl

### 重构完成的模块
1. ✅ **user** - DAO + Repository + Wire 生成
2. ✅ **comment** - DAO + Repository + Wire 生成  
3. ✅ **file** - DAO + Wire 生成
4. ✅ **friend** - DAO + Wire 生成

## ⏳ 待完成 (4/8 模块)

需要重构以下4个模块的 DAO 文件：

### 5. ❌ label  
文件：`internal/label/internal/repository/dao/label.go`
- 包含聚合查询（Lookup, AddFields）
- 需要手动构建 pipeline

### 6. ❌ document
文件：`internal/document/internal/repository/dao/document.go`
- 标准 CRUD 操作
- 较简单

### 7. ❌ document_content  
文件：`internal/document_content/internal/repository/dao/document_content.go`
- 标准 CRUD 操作
- 较简单

### 8. ❌ post (最复杂)
文件：`internal/post/internal/repository/dao/post.go`
- 大量聚合管道
- Lookup、Unwind、Match、Sort 等
- Repository 层也需要更新

## 🔧 每个模块的重构步骤

### 1. 更新 imports
```go
// 删除
"github.com/chenmingyong0423/go-mongox/v2"
"github.com/chenmingyong0423/go-mongox/v2/builder/query"  
"github.com/chenmingyong0423/go-mongox/v2/builder/update"
"github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"

// 添加
"time"
"go.mongodb.org/mongo-driver/v2/mongo"
"go.mongodb.org/mongo-driver/v2/mongo/options"  
```

### 2. 更新结构体
```go
type Entity struct {
    ID        bson.ObjectID  `bson:"_id,omitempty"`
    CreatedAt time.Time      `bson:"created_at"`
    UpdatedAt time.Time      `bson:"updated_at"`
    DeletedAt *time.Time     `bson:"deleted_at,omitempty"`
    // 其他字段...
}
```

### 3. 更新构造函数和集合
```go
func NewDao(db *mongo.Database) *Dao {
    return &Dao{coll: db.Collection("table_name")}
}

type Dao struct {
    coll *mongo.Collection
}
```

### 4. 重写 CRUD 方法
参考已完成的 user/comment/file/friend 模块

### 5. 处理聚合查询（label, post）
手动构建 `mongo.Pipeline`:
```go
pipeline := mongo.Pipeline{
    {{Key: "$lookup", Value: bson.D{
        {Key: "from", Value: "collection"},
        {Key: "localField", Value: "field"},
        {Key: "foreignField", Value: "_id"},
        {Key: "as", Value: "result"},
    }}},
    {{Key: "$unwind", Value: "$result"}},
}
```

### 6. 生成 wire 代码
```bash
go run github.com/google/wire/cmd/wire@latest gen ./internal/label ./internal/document ./internal/document_content ./internal/post
```

### 7. 测试编译
```bash
go build
```

## 📊 进度统计
- 总模块数：8
- 已完成：4 (50%)
- 待完成：4 (50%)

## ⚡ 估计剩余工作量
- label: ~30分钟（有聚合）
- document: ~15分钟
- document_content: ~20分钟  
- post: ~60分钟（最复杂）
- **总计：约2小时**

## 🎯 建议执行顺序
1. document（最简单）
2. document_content
3. label（中等复杂度）
4. post（最复杂，最后处理）

## 📝 参考文档
已完成的示例文件可作为参考：
- `internal/user/internal/repository/dao/user.go`
- `internal/comment/internal/repository/dao/comment.go`
- `internal/file/internal/repository/dao/file.go`
- `internal/friend/internal/repository/dao/friend.go`
