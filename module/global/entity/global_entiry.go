package global_entity

type SendEmailFormat struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

type (
	EmailServiceInterface interface {
		SendEmail(sendEmailFormat SendEmailFormat) error
	}

	StorageServiceInterface interface {
		UploadFile(filePath string, bucketName string, fileName string) error
		DownlaodFile(bucketName string, key string, downloadPath string) error
		DeleteFile(bucketName string, fileName string) error
		GetUrl(bucketName string, fileName string) (string, error)
	}
)
