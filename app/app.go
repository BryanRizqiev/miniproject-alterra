package app

import (
	"miniproject-alterra/app/config"
	"miniproject-alterra/app/lib"
	event_controller "miniproject-alterra/module/events/controller"
	mysql_event_repository "miniproject-alterra/module/events/repository/mysql"
	event_service "miniproject-alterra/module/events/service"
	evd_controller "miniproject-alterra/module/evidence/controller"
	mysql_evd_repo "miniproject-alterra/module/evidence/repository/mysql"
	evd_svc "miniproject-alterra/module/evidence/service"
	global_service "miniproject-alterra/module/global/service"
	user_controller "miniproject-alterra/module/user/controller"
	mysql_user_repository "miniproject-alterra/module/user/repository/mysql"
	user_service "miniproject-alterra/module/user/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func newAWSSession(config *config.AppConfig) (*session.Session, error) {

	endpoint := config.ENDPOINT
	sess, err := session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   aws.String(config.AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			config.AWS_ACCESS_KEY_ID,
			config.AWS_SECRET_ACCESS_KEY,
			"",
		),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil

}

func Bootstrap(db *gorm.DB, e *echo.Echo, config *config.AppConfig) {

	sess, err := newAWSSession(config)
	if err != nil {
		panic("Failed to create AWS session")
	}
	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	emailService := global_service.NewEmailService(config)
	storageService := global_service.NewStorageService(uploader, downloader, s3Client)

	userRepository := mysql_user_repository.NewUserRepository(db)
	userService := user_service.NewUserService(userRepository, emailService, config)
	userController := user_controller.NewUserController(userService, emailService, storageService)

	evtRepo := mysql_event_repository.NewEventRepository(db)
	evtSvc := event_service.NewEventService(evtRepo, storageService)
	evtController := event_controller.NewEventController(evtSvc)

	evdRepo := mysql_evd_repo.NewEvidenceRepository(db)
	evdSvc := evd_svc.NewEvidenceService(evdRepo, storageService)
	evdController := evd_controller.NewEvidenceController(evdSvc)

	evidence := e.Group("/evidences")
	evidence.POST("/create", evdController.CreateEvidence, lib.JWTMiddleware())
	evidence.GET("/get/:event-id", evdController.GetEvidences, lib.JWTMiddleware())

	events := e.Group("/events")
	events.POST("", evtController.CreateEvent, lib.JWTMiddleware())
	events.GET("", evtController.GetEvent)

	admins := e.Group("/admins")
	admins.GET("/requesting-users", userController.GetRequestingUser, lib.JWTMiddleware())
	admins.PUT("/change-verification", userController.ApproveVerification, lib.JWTMiddleware())

	e.POST("/register", userController.Register)
	e.POST("/login", userController.Login)
	e.GET("/verify", userController.Verify, lib.JWTMiddleware())
	e.POST("/request-verified", userController.RequestVerified, lib.JWTMiddleware())
	e.POST("/request-verify-email", userController.RequestVerifyEmail, lib.JWTMiddleware())
	e.GET("/verify-email/:user-id", userController.VerifyEmail)
	e.GET("/verify-email/:user-id", userController.VerifyEmail)

}
