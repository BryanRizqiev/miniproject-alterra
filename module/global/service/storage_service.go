package global_service

import (
	global_entity "miniproject-alterra/module/global/entity"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type StorageService struct {
	uploader   *s3manager.Uploader
	downlaoder *s3manager.Downloader
	client     *s3.S3
}

func NewStorageService(uploader *s3manager.Uploader, downlaoder *s3manager.Downloader, client *s3.S3) global_entity.StorageServiceInterface {

	return &StorageService{
		uploader:   uploader,
		downlaoder: downlaoder,
		client:     client,
	}

}

// DeleteFile implements global_entity.StorageServiceInterface.
func (*StorageService) DeleteFile(bucketName string, fileName string) error {
	panic("unimplemented")
}

// DownlaodFile implements global_entity.StorageServiceInterface.
func (*StorageService) DownlaodFile(bucketName string, key string, downloadPath string) error {
	panic("unimplemented")
}

// GetUrl implements global_entity.StorageServiceInterface.
func (*StorageService) GetUrl(bucketName string, fileName string) (string, error) {
	panic("unimplemented")
}

// UploadFile implements global_entity.StorageServiceInterface.
func (*StorageService) UploadFile(filePath string, bucketName string, fileName string) error {
	panic("unimplemented")
}
