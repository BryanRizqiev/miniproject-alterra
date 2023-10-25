package global_service

import (
	"io"
	global_entity "miniproject-alterra/module/global/entity"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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

func (this *StorageService) UploadFile(bucketName string, fileName string, body io.Reader) error {
	_, err := this.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   body,
	})

	return err
}

func (this *StorageService) DeleteFile(bucketName string, fileName string) error {
	_, err := this.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	return err
}

func (this *StorageService) GetUrl(bucketName string, fileName string) (string, error) {

	req, _ := this.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	urlStr, err := req.Presign(10 * time.Minute)
	if err != nil {
		return "", err
	}

	return urlStr, nil

}
