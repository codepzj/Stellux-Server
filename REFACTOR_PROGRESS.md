# MongoDB重构进度

## 已完成
- [x] 更新 go.mod 移除 go-mongox 依赖
- [x] 重构 internal/infra/mongodb.go 使用官方驱动  
- [x] internal/user/internal/repository/dao/user.go
- [x] internal/comment/internal/repository/dao/comment.go

## 进行中
- [ ] 重构所有 DAO 文件

## 待处理DAO文件
- [ ] internal/file/internal/repository/dao/file.go
- [ ] internal/friend/internal/repository/dao/friend.go  
- [ ] internal/label/internal/repository/dao/label.go
- [ ] internal/document/internal/repository/dao/document.go
- [ ] internal/document_content/internal/repository/dao/document_content.go
- [ ] internal/post/internal/repository/dao/post.go

## 待处理
- [ ] 更新所有 wire 文件使用 *mongo.Database
- [ ] 更新代码生成模板
- [ ] 更新 repository 层文件
