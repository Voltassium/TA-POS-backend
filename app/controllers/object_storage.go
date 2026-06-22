package controllers

import (
	"backend-ta/pkg/config"
	"backend-ta/pkg/errors"
	"backend-ta/pkg/http/server/http_response"
	"backend-ta/pkg/storage"
	utils "backend-ta/pkg/utils"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type StorageCtl struct {
	minioClient storage.Storage
}

func NewStorageController(client storage.Storage) StorageCtl {
	return StorageCtl{
		minioClient: client,
	}
}

func (ctl *StorageCtl) UploadFile(c *gin.Context) {
	// Limit file size (in MB → bytes)
	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		config.LoadConfig().ObjectStorage.MaxFileSize*1024*1024,
	)

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		http_response.SendError(c, errors.StorageErrorToAppError("Failed to read uploaded file"))
		return
	}

	// Validate file format
	if ok, expectedFormat, actualFormat := utils.IsDocumentFile(fileHeader.Filename); !ok {
		http_response.SendError(c, errors.StorageErrorToAppError(
			fmt.Sprintf("Cannot upload file with format: %s, expected: %s", actualFormat, expectedFormat),
		))
		return
	}

	// Generate a safe/unique filename
	filename := utils.GenerateKeyFile(fileHeader.Filename)

	// Save file to local storage (e.g., ./uploads)
	savePath := path.Join("uploads", filename)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		http_response.SendError(c, errors.StorageErrorToAppError("Failed to save file"))
		return
	}

	// Return success response
	http_response.SendSuccess(c, http.StatusOK, "Success upload file", map[string]interface{}{
		"filename": filename,
		"path":     savePath,
	})
}
