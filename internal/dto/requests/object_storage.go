package requests

import "mime/multipart"

type FileStorage struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}
