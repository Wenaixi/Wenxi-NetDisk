package api

import (
	"github.com/wenaixi/wenxi-cloud/backend/config"
	"github.com/wenaixi/wenxi-cloud/backend/internal/api/handlers"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/middleware"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/jwt"
	"github.com/wenaixi/wenxi-cloud/backend/internal/repository"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// JWT Manager
	jwtManager := jwt.NewJWTManager(cfg.JWTsecret, cfg.JWTExpire)

	// Repositories
	userRepo := repository.NewUserRepository(nil) // Will be set after DB init
	fileRepo := repository.NewFileRepository(nil)
	shareRepo := repository.NewShareRepository(nil)
	lanzouRepo := repository.NewLanZouTokenRepository(nil)

	// Services
	authSvc := service.NewAuthService(userRepo, jwtManager)
	fileSvc := service.NewFileService(fileRepo)
	shareSvc := service.NewShareService(shareRepo, fileRepo)
	lanzouSvc := service.NewLanZouService(lanzouRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	fileHandler := handlers.NewFileHandler(fileSvc)
	shareHandler := handlers.NewShareHandler(shareSvc)
	lanzouHandler := handlers.NewLanZouHandler(lanzouSvc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthRequired(jwtManager), authHandler.GetCurrentUser)
		}

		// File routes (protected)
		files := api.Group("/files")
		files.Use(middleware.AuthRequired(jwtManager))
		{
			files.GET("", fileHandler.ListFiles)
			files.POST("", fileHandler.CreateFileMetadata)
			files.GET("/:id", fileHandler.GetFile)
			files.DELETE("/:id", fileHandler.DeleteFile)
			files.POST("/:id/share", shareHandler.CreateShare)
		}

		// Share routes
		shares := api.Group("/shares")
		{
			shares.GET("/:token", shareHandler.GetShare)
			shares.DELETE("/:id", middleware.AuthRequired(jwtManager), shareHandler.DeleteShare)
		}

		// LanZou routes (protected)
		lanzou := api.Group("/lanzou")
		lanzou.Use(middleware.AuthRequired(jwtManager))
		{
			lanzou.POST("/connect", lanzouHandler.Connect)
			lanzou.GET("/status", lanzouHandler.GetStatus)
		}
	}

	return r
}
