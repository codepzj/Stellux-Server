# 快速修复指南

## 当前状态
✅ 已完成模块：user, comment, file, friend
❌ 待完成模块：label, document, document_content, post

## 剩余工作
需要重构4个模块的DAO文件，每个文件的修改模式相同：

### 修改步骤（每个DAO文件）：

1. **imports** - 替换
```go
"github.com/chenmingyong0423/go-mongox/v2"
"github.com/chenmingyong0423/go-mongox/v2/builder/query"
"github.com/chenmingyong0423/go-mongox/v2/builder/update"
"github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"
```
为：
```go
"time"
"go.mongodb.org/mongo-driver/v2/mongo"
"go.mongodb.org/mongo-driver/v2/mongo/options"
"go.mongodb.org/mongo-driver/v2/bson"
```

2. **结构体** - 添加时间戳字段
```go
// 旧
type Entity struct {
    mongox.Model `bson:",inline"`
    ...
}

// 新
type Entity struct {
    ID        bson.ObjectID  `bson:"_id,omitempty"`
    CreatedAt time.Time      `bson:"created_at"`
    UpdatedAt time.Time      `bson:"updated_at"`
    DeletedAt *time.Time     `bson:"deleted_at,omitempty"`
    ...
}
```

3. **构造函数**
```go
// 旧
func NewDao(db *mongox.Database) *Dao {
    return &Dao{coll: mongox.NewCollection[Entity](db, "table")}
}

// 新
func NewDao(db *mongo.Database) *Dao {
    return &Dao{coll: db.Collection("table")}
}
```

4. **集合类型**
```go
// 旧
type Dao struct {
    coll *mongox.Collection[Entity]
}

// 新
type Dao struct {
    coll *mongo.Collection
}
```

5. **CRUD 操作** - 使用原生MongoDB API

完成后运行：
```bash
go run github.com/google/wire/cmd/wire@latest gen ./internal/label ./internal/document ./internal/document_content ./internal/post
go build
```
