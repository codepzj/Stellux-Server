package service

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/codepzj/stellux/server/internal/file/internal/domain"
	"github.com/codepzj/stellux/server/internal/file/internal/repository"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/codepzj/stellux/server/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "保存文件失败: %s", uploadFile.FileName)
	}

	// 5. 存入数据库
	err = s.repo.Create(ctx, uploadFile)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileService) QueryFileList(ctx *gin.Context, page *apiwrap.Page) ([]*domain.File, int64, error) {
	return s.repo.GetList(ctx, page)
}

func (s *FileService) DeleteFiles(ctx *gin.Context, idList []string) error {
	files, err := s.repo.GetListByIDList(ctx, apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(idList)))
	if err != nil {
		return err
	}
	for _, file := range files {
		_ = os.Remove(file.Dst)
	}

	return s.repo.DeleteMany(ctx, apiwrap.ToObjectIDList(apiwrap.ConvertBsonIDList(idList)))
}
