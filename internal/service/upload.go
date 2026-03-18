package service

type UploadService interface {
	GenerateSign(fileName, contentType string, fileSize int64) (uploadURL, objectKey string, expiredAt string, err error)
}
