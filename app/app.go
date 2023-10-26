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
	global_repo "miniproject-alterra/module/global/repository"
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
	"github.com/sashabaranov/go-openai"
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

func Bootstrap(db *gorm.DB, echo *echo.Group, config *config.AppConfig) {

	sess, err := newAWSSession(config)
	if err != nil {
		panic("Failed to create AWS session")
	}
	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)
	storageService := global_service.NewStorageService(uploader, downloader, s3Client)

	emailService := global_service.NewEmailService(config)
	globalRepo := global_repo.NewGlobalRepo(db)
	openaiClient := openai.NewClient(config.OPENAPI_KEY)
	openaiSvc := global_service.NewOpenAIService(openaiClient, openai.GPT3Dot5Turbo)

	userRepository := mysql_user_repository.NewUserRepository(db)
	userService := user_service.NewUserService(userRepository, emailService, storageService, config)
	userController := user_controller.NewUserController(userService, emailService, storageService)

	eventRepo := mysql_event_repository.NewEventRepository(db)
	eventSvc := event_service.NewEventService(eventRepo, storageService, globalRepo, openaiSvc)
	eventController := event_controller.NewEventController(eventSvc)

	evidenceRepo := mysql_evd_repo.NewEvidenceRepository(db)
	evidenceSvc := evd_svc.NewEvidenceService(evidenceRepo, storageService)
	evidenceController := evd_controller.NewEvidenceController(evidenceSvc)

	// Route

	evidence := echo.Group("/evidences")
	evidence.POST("/create", evidenceController.CreateEvidence, lib.JWTMiddleware())
	evidence.GET("/get/:event-id", evidenceController.GetEvidences, lib.JWTMiddleware())

	events := echo.Group("/events")
	events.GET("", eventController.GetEvent)
	events.POST("", eventController.CreateEvent, lib.JWTMiddleware())
	events.PUT("/update/:event-id", eventController.UpdateEvent, lib.JWTMiddleware())
	events.PUT("/update-image/:event-id", eventController.UpdateImage, lib.JWTMiddleware())
	events.DELETE("delete/:event-id", eventController.DeleteEvent, lib.JWTMiddleware())

	admin := echo.Group("/admin", lib.JWTMiddleware())
	admin.GET("/events/waiting", eventController.GetWaitingEvents)
	admin.PUT("/events/publish/:event-id", eventController.PublishEvent)
	admin.PUT("/events/takedown/:event-id", eventController.TakedownEvent)
	admin.GET("/users", userController.GetAllUser)
	admin.GET("/users/requesting-users", userController.GetRequestingUser)
	admin.PUT("/users/change-role/:user-id", userController.ChangeUserRole)
	admin.DELETE("/users/delete/:user-id", userController.DeleteUser)

	users := echo.Group("/users")
	users.GET("/profile", userController.GetUserProfile, lib.JWTMiddleware())
	users.GET("/verify-email/:user-id", userController.VerifyEmail)
	users.GET("/request-verify-email", userController.RequestVerifyEmail, lib.JWTMiddleware())
	users.GET("/request-verify-user", userController.RequestVerified, lib.JWTMiddleware())
	users.PUT("", userController.UpdateUser, lib.JWTMiddleware())
	users.PUT("/photo", userController.UpdatePhoto, lib.JWTMiddleware())
	users.DELETE("/self-delete", userController.UserSelfDelete, lib.JWTMiddleware())

	// Auth
	echo.POST("/register", userController.Register)
	echo.POST("/login", userController.Login)

}
