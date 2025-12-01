# MongoDB 官方驱动重构总结

## 项目概述
将项目从 `go-mongox` 迁移到 MongoDB 官方驱动 (`go.mongodb.org/mongo-driver/v2`)。

## 已完成工作 ✅

### 1. 基础设施更新
- ✅ 更新 `go.mod`，移除 `go-mongox` 依赖
- ✅ 重构 `internal/infra/mongodb.go` 使用 `*mongo.Database`

### 2. Wire 依赖注入更新
已更新所有模块的 wire.go 文件，将 `*mongox.Database` 改为 `*mongo.Database`:
- ✅ `internal/user/wire.go`
- ✅ `internal/comment/wire.go`
- ✅ `internal/document/wire.go`
- ✅ `internal/document_content/wire.go`
- ✅ `internal/file/wire.go`
- ✅ `internal/friend/wire.go`
- ✅ `internal/label/wire.go`
- ✅ `internal/post/wire.go`

### 3. 代码生成模板更新
- ✅ `cmd/gen/templates/dao.tmpl`
- ✅ `cmd/gen/templates/wire.tmpl`

### 4. 部分 DAO 层重构
- ✅ `internal/user/internal/repository/dao/user.go` - 已完成基本重构
- ✅ `internal/comment/internal/repository/dao/comment.go` - 已完成基本重构

## 待完成工作 🔧

### 1. 剩余 DAO 文件需要重构
以下DAO文件需要按照相同模式进行重构：

#### 需要修改的关键点：
1. **导入包**：
   ```go
   // 移除
   "github.com/chenmingyong0423/go-mongox/v2"
   "github.com/chenmingyong0423/go-mongox/v2/builder/query"
   "github.com/chenmingyong0423/go-mongox/v2/builder/update"
   "github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"
   
   // 添加
   "go.mongodb.org/mongo-driver/v2/mongo"
   "go.mongodb.org/mongo-driver/v2/mongo/options"
   "time"
   ```

2. **结构体定义**：
   ```go
   // 旧方式
   type Entity struct {
       mongox.Model `bson:",inline"`
       Field string `bson:"field"`
   }
   
   // 新方式
   type Entity struct {
       ID        bson.ObjectID  `bson:"_id,omitempty"`
       CreatedAt time.Time      `bson:"created_at"`
       UpdatedAt time.Time      `bson:"updated_at"`
       DeletedAt *time.Time     `bson:"deleted_at,omitempty"`
       Field     string         `bson:"field"`
   }
   ```

3. **DAO 构造函数**：
   ```go
   // 旧方式
   func NewDao(db *mongox.Database) *Dao {
       return &Dao{coll: mongox.NewCollection[Entity](db, "collection")}
   }
   type Dao struct {
       coll *mongox.Collection[Entity]
   }
   
   // 新方式
   func NewDao(db *mongo.Database) *Dao {
       return &Dao{coll: db.Collection("collection")}
   }
   type Dao struct {
       coll *mongo.Collection
   }
   ```

4. **CRUD 操作**：
   ```go
   // 插入操作
   entity.ID = bson.NewObjectID()
   entity.CreatedAt = time.Now()
   entity.UpdatedAt = time.Now()
   result, err := d.coll.InsertOne(ctx, entity)
   
   // 查询操作
   var entity Entity
   err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
   
   // 更新操作
   update := bson.M{
       "$set": bson.M{
           "field": value,
           "updated_at": time.Now(),
       },
   }
   result, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
   
   // 删除操作
   result, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
   
   // 列表查询
   opts := options.Find().SetSort(bson.M{"created_at": -1})
   cursor, err := d.coll.Find(ctx, filter, opts)
   defer cursor.Close(ctx)
   var results []*Entity
   err = cursor.All(ctx, &results)
   ```

#### 待重构文件列表：
- ❌ `internal/file/internal/repository/dao/file.go`
- ❌ `internal/friend/internal/repository/dao/friend.go`
- ❌ `internal/label/internal/repository/dao/label.go`
- ❌ `internal/document/internal/repository/dao/document.go`
- ❌ `internal/document_content/internal/repository/dao/document_content.go`
- ❌ `internal/post/internal/repository/dao/post.go` (这是最复杂的，包含聚合查询)

### 2. Repository 层更新
部分 repository 文件需要更新类型转换逻辑：
- ❌ `internal/user/internal/repository/user.go` - 修复 FindOptions 使用
- ❌ 其他所有 repository 文件中涉及 DAO 结构体字段访问的地方

### 3. 聚合查询重构
`post` 模块使用了复杂的聚合管道，需要特别注意：
- 使用 `mongo.Pipeline` 代替 `aggregation.NewStageBuilder()`
- 手动构建 pipeline 阶段

### 4. 测试和验证
- ❌ 运行 `go build` 确保编译通过
- ❌ 更新单元测试（如果有）
- ❌ 运行应用确保功能正常

## 已知问题 ⚠️

### 1. FindOptions 类型问题
`internal/user/internal/repository/user.go:71` 有类型不匹配：
- `options.Find()` 返回 `*FindOptionsBuilder`
- DAO 接口需要 `*FindOptions`

**解决方案**：将 DAO 接口改为接受 `*FindOptions` 而不是 `*FindOptionsBuilder`

### 2. DeletedAt 字段类型
`*time.Time` 和 `time.Time` 之间的转换需要注意处理 nil 值

## 重构策略建议 📋

### 按模块优先级：
1. **简单模块**（先做）：
   - file
   - friend
   
2. **中等复杂度**：
   - label
   - document
   - document_content
   
3. **复杂模块**（最后）：
   - post（包含复杂聚合查询）

### 每个模块的步骤：
1. 更新 DAO 文件的 imports
2. 修改结构体定义添加时间戳字段
3. 更新构造函数
4. 重构所有 CRUD 方法
5. 更新 repository 层的类型转换
6. 测试编译

## 参考资源 📚
- [MongoDB Go Driver v2 官方文档](https://www.mongodb.com/zh-cn/docs/drivers/go/current/)
- [CRUD Operations](https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/)
- [Aggregation](https://www.mongodb.com/docs/drivers/go/current/fundamentals/aggregation/)

## 下一步行动 🎯
1. 按照上述策略逐个模块完成 DAO 重构
2. 修复所有编译错误
3. 运行 `go mod tidy` 清理依赖
4. 测试应用功能
5. 删除此文档和 REFACTOR_PROGRESS.md

## 注意事项 ⚡
- 所有时间戳字段在插入时需要手动设置
- 删除操作不再自动设置 deleted_at，需要手动实现软删除
- 聚合查询需要手动构建 pipeline，不能再使用 builder
- Find 操作需要手动处理 cursor 和 decode
