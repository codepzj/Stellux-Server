package service

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/codepzj/Stellux-Server/internal/file/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/file/internal/repository"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/codepzj/Stellux-Server/internal/pkg/logger"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IFileService interface {
	UploadFile(ctx *gin.Context, file *multipart.FileHeader) error
	QueryFileList(ctx *gin.Context, page *apiwrap.Page) ([]*domain.File, int64, error)
	DeleteFiles(ctx *gin.Context, idList []string) error
}

var _ IFileService = (*FileService)(nil)

func NewFileService(repo repository.IFileRepository) *FileService {
	return &FileService{
		repo: repo,
	}
}

type FileService struct {
	repo repository.IFileRepository
}

func (s *FileService) UploadFile(ctx *gin.Context, file *multipart.FileHeader) error {
	err := os.MkdirAll("static/images", os.ModePerm)
	if err != nil {
		logger.Error("创建目录失败",
			logger.WithError(err),
		)
		return err
	}

	// 获取domain文件信息

	fileName := file.Filename
	// 生成新的文件名
	timestamp := time.Now().Unix()
	newFileName := strconv.FormatInt(timestamp, 10) + utils.RandString(10) + filepath.Ext(fileName)
	networkPath := "/images/" + newFileName
	filePath := "static/images/" + newFileName
	uploadFile := &domain.File{
		FileName: fileName,
		Url:      networkPath,
		Dst:      filePath,
	}

	// 保存文件
	os.MkdirAll(filepath.Dir(uploadFile.Dst), 0755)
	err = ctx.SaveUploadedFile(file, uploadFile.Dst)
	if err != nil {
		logger.Error("保存文件失败",
			logger.WithError(err),
			logger.WithString("filename", uploadFile.FileName),
		)
		return errors.Wrapf(err, "保存文件失败: %s", uploadFile.FileName)
	}

	// 5. 存入数据库
	err = s.repo.Create(ctx, uploadFile)
	if err != nil {
		logger.Error("存入数据库失败",
			logger.WithError(err),
			logger.WithString("filename", uploadFile.FileName),
		)
		return err
	}

	logger.Info("上传文件成功",
		logger.WithString("filename", fileName),
		logger.WithString("url", networkPath),
	)

	return nil
}

func (s *FileService) QueryFileList(ctx *gin.Context, page *apiwrap.Page) ([]*domain.File, int64, error) {
	files, total, err := s.repo.GetList(ctx, page)
	if err != nil {
		logger.Error("查询文件列表失败",
			logger.WithError(err),
		)
		return nil, 0, err
	}
	return files, total, nil
}

func (s *FileService) DeleteFiles(ctx *gin.Context, idList []string) error {
	var objIdList []bson.ObjectID
	for _, id := range idList {
		objId, err := bson.ObjectIDFromHex(id)
		if err != nil {
			logger.Error("转换ObjectID失败",
				logger.WithError(err),
				logger.WithString("id", id),
			)
			return err
		}
		objIdList = append(objIdList, objId)
	}

	files, err := s.repo.GetListByIDList(ctx, objIdList)
	if err != nil {
		logger.Error("查询文件列表失败",
			logger.WithError(err),
		)
		return err
	}

	for _, file := range files {
		err := os.Remove(file.Dst)
		if err != nil {
			logger.Warn("删除物理文件失败",
				logger.WithError(err),
				logger.WithString("path", file.Dst),
			)
		}
	}

	err = s.repo.DeleteMany(ctx, objIdList)
	if err != nil {
		logger.Error("删除文件记录失败",
			logger.WithError(err),
		)
		return err
	}

	logger.Info("批量删除文件成功",
		logger.WithInt("count", len(idList)),
	)

	return nil
}
