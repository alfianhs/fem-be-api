package main

import (
	admin_http "app/app/delivery/http/admin"
	member_http "app/app/delivery/http/member"
	"app/app/delivery/http/middleware"
	superadmin_http "app/app/delivery/http/superadmin"
	mongo_repository "app/app/repository/mongo"
	s3_repository "app/app/repository/s3"
	admin_usecase "app/app/usecase/admin"
	member_usecase "app/app/usecase/member"
	superadmin_usecase "app/app/usecase/superadmin"
	"app/docs"
	"app/helpers"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func init() {
	_ = godotenv.Load()
}

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	// programmatically set swagger info
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "Futsal Event Management"
	}
	docs.SwaggerInfo.Title = appName
	docs.SwaggerInfo.Description = "API Documentations"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "5"
	}
	timeout, _ := strconv.Atoi(timeoutStr)
	timeoutContext := time.Duration(timeout) * time.Second

	// logger
	writers := make([]io.Writer, 0)
	if logSTDOUT, _ := strconv.ParseBool(os.Getenv("LOG_TO_STDOUT")); logSTDOUT {
		writers = append(writers, os.Stdout)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(writers...))

	// set gin writer to logrus
	gin.DefaultWriter = logrus.StandardLogger().Writer()

	// init mongo db
	mongo := helpers.ConnectMongoDB(timeoutContext, os.Getenv("MONGO_URL"))

	// init mongo repository
	mongoDbRepo := mongo_repository.NewMongoDbRepo(mongo)

	// init s3 repository
	s3Repo := s3_repository.NewS3Repository(timeoutContext)

	// init superadmin usecase
	superadminUsecase := superadmin_usecase.NewSuperadminAppUsecase(superadmin_usecase.RepoInjection{
		MongoDbRepo: mongoDbRepo,
		S3Repo:      s3Repo,
	}, timeoutContext)

	// init admin usecase
	adminUsecase := admin_usecase.NewAdminAppUsecase(admin_usecase.RepoInjection{
		MongoDbRepo: mongoDbRepo,
	}, timeoutContext)

	// init member usecase
	memberUsecase := member_usecase.NewMemberAppUsecase(member_usecase.RepoInjection{
		MongoDbRepo: mongoDbRepo,
	}, timeoutContext)

	// init middleware
	middleware := middleware.NewAppMiddleware()

	// gin mode realease when go env is production
	if os.Getenv("GO_ENV") == "production" || os.Getenv("GO_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin
	ginEngine := gin.New()

	// panic recovery
	ginEngine.Use(middleware.Recovery())

	// logger
	ginEngine.Use(middleware.Logger(io.MultiWriter(writers...)))

	// cors
	ginEngine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Ticket-Token"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	// init route
	superadmin_http.NewSuperadminRouteHandler(superadminUsecase, ginEngine, middleware)
	admin_http.NewAdminRouteHandler(adminUsecase, ginEngine, middleware)
	member_http.NewMemberRouteHandler(memberUsecase, ginEngine, middleware)

	// default route
	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"message": "Welcome",
		})
	})

	// swagger route
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")

	logrus.Infof("Service running on port %s", port)
	ginEngine.Run(":" + port)
}
