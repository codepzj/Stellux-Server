package web

import (
	"github.com/codepzj/Stellux-Server/internal/file/internal/service"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
)

func NewFileHandler(serv service.IFileService) *FileHandler {
	return &FileHandler{
		serv: serv,
	}
}

type FileHandler struct {
	serv service.IFileService
}

func (h *FileHandler) RegisterGinRoutes(engine *gin.Engine) {
	engine.Static("/images", "./static/images")
	fileGroup := engine.Group("/file")
	{
		fileGroup.GET("/list", apiwrap.WrapWithQuery(h.QueryFileList))
	}
	fileAdminGroup := engine.Group("/admin-api/file")
	{
		fileAdminGroup.POST("/upload", apiwrap.Wrap(h.UploadFile))
		fileAdminGroup.DELETE("/delete", apiwrap.WrapWithJson(h.DeleteFiles))
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) (int, string, any) {
	file, err := c.FormFile("file")
	if err != nil {
		return 400, err.Error(), nil
	}
	if file == nil {
		return 400, "未找到上传的文件", nil
	}
	err = h.serv.UploadFile(c, file)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "文件上传成功", nil
}

func (h *FileHandler) QueryFileList(c *gin.Context, page *apiwrap.Page) (int, string, any) {
	files, count, err := h.serv.QueryFileList(c, page)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "文件列表查询成功", apiwrap.ToPageVO(page.PageNo, page.PageSize, count, h.FileDomainToVOList(files))
}

func (h *FileHandler) DeleteFiles(c *gin.Context, deleteFilesRequest *DeleteFilesRequest) (int, string, any) {
	err := h.serv.DeleteFiles(c, deleteFilesRequest.IDList)
	if err != nil {
		return 500, err.Error(), nil
	}
	return 200, "文件删除成功", nil
}
