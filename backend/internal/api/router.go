package api

import (
	"github.com/wenaixi/wenxi-cloud/backend/config"
	"github.com/wenaixi/wenxi-cloud/backend/internal/api/handlers"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/middleware"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/jwt"
	"github.com/wenaixi/wenxi-cloud/backend/internal/repository"
	"github.com/wenaixi/wenxi-cloud/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// JWT Manager
	jwtManager := jwt.NewJWTManager(cfg.JWTsecret, cfg.JWTExpire)

	// Initialize DB (called in main.go, but we accept it here)
	db := getDB()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)
	shareRepo := repository.NewShareRepository(db)
	lanzouRepo := repository.NewLanZouTokenRepository(db)
	uploadRepo := repository.NewUploadSessionRepository(db)
	folderRepo := repository.NewFolderRepository(db)
	recycleRepo := repository.NewRecycleBinRepository(db)
	versionRepo := repository.NewFileVersionRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, jwtManager)
	fileSvc := service.NewFileService(fileRepo)
	folderSvc := service.NewFolderService(folderRepo)
	shareSvc := service.NewShareService(shareRepo, fileRepo)
	lanzouSvc := service.NewLanZouService(lanzouRepo)
	uploadSvc := service.NewUploadService(uploadRepo, lanzouSvc, fileSvc)
	downloadSvc := service.NewDownloadService(fileSvc, lanzouSvc)
	recycleSvc := service.NewRecycleBinService(recycleRepo, fileRepo, folderRepo)
	versionSvc := service.NewFileVersionService(versionRepo, fileRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	fileHandler := handlers.NewFileHandler(fileSvc)
	folderHandler := handlers.NewFolderHandler(folderSvc)
	shareHandler := handlers.NewShareHandler(shareSvc)
	lanzouHandler := handlers.NewLanZouHandler(lanzouSvc, uploadSvc)
	downloadHandler := handlers.NewDownloadHandler(downloadSvc)
	recycleHandler := handlers.NewRecycleHandler(recycleSvc)
	uploadHandler := handlers.NewUploadHandler(fileSvc, versionSvc)
	shareParseHandler := handlers.NewShareParseHandler(service.NewShareParseService(lanzou.NewClient("")))

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
			files.PUT("/:id", fileHandler.RenameFile)
			files.PUT("/:id/move", fileHandler.MoveFile)
			files.PUT("/:id/description", fileHandler.UpdateFileDescription)
			files.DELETE("/:id", fileHandler.DeleteFile)
			files.POST("/:id/share", shareHandler.CreateShare)
			files.GET("/:id/download", downloadHandler.GetDownloadURL)
			files.POST("/upload", uploadHandler.UploadFile)
			files.POST("/upload-url", uploadHandler.GetUploadURL)
			files.GET("/:id/versions", uploadHandler.ListVersions)
			files.POST("/:id/versions/:version_id/restore", uploadHandler.RestoreVersion)
		}

		// Folder routes (protected)
		folders := api.Group("/folders")
		folders.Use(middleware.AuthRequired(jwtManager))
		{
			folders.GET("", folderHandler.ListFolders)
			folders.POST("", folderHandler.CreateFolder)
			folders.GET("/:id", folderHandler.GetFolder)
			folders.PUT("/:id", folderHandler.UpdateFolder)
			folders.PUT("/:id/description", folderHandler.UpdateFolderDescription)
			folders.DELETE("/:id", folderHandler.DeleteFolder)
			folders.PUT("/:id/move", folderHandler.MoveFolder)
		}

		// Share routes
		shares := api.Group("/shares")
		{
			shares.GET("/:token", shareHandler.GetShare)
			shares.POST("/:token/validate", shareHandler.ValidateShare)
			shares.POST("", middleware.AuthRequired(jwtManager), shareHandler.CreateShareViaBody)
			shares.GET("", middleware.AuthRequired(jwtManager), shareHandler.ListShares)
			shares.DELETE("/:id", middleware.AuthRequired(jwtManager), shareHandler.DeleteShare)
		}

		// LanZou routes (protected)
		lanzou := api.Group("/lanzou")
		lanzou.Use(middleware.AuthRequired(jwtManager))
		{
			lanzou.POST("/connect", lanzouHandler.Connect)
			lanzou.GET("/status", lanzouHandler.GetStatus)
			lanzou.DELETE("/connect", lanzouHandler.Disconnect)
			lanzou.GET("/files", lanzouHandler.ListFiles)
			lanzou.GET("/folders", lanzouHandler.ListFolders)
			lanzou.POST("/folders", lanzouHandler.CreateFolder)
			lanzou.POST("/share", lanzouHandler.CreateShare)
			lanzou.GET("/files/:id/url", lanzouHandler.GetFileURL)
			lanzou.PUT("/access", lanzouHandler.SetAccess)
			lanzou.PUT("/rename/:id", lanzouHandler.Rename)
			lanzou.PUT("/move/:id", lanzouHandler.Move)
			lanzou.POST("/batch/delete", lanzouHandler.BatchDelete)
			lanzou.POST("/batch/move", lanzouHandler.BatchMove)
			lanzou.GET("/files/:id/description", lanzouHandler.GetFileDescription)
			lanzou.PUT("/files/:id/description", lanzouHandler.SetFileDescription)
			lanzou.POST("/upload/init", lanzouHandler.InitializeUpload)
			lanzou.POST("/upload/complete/:id", lanzouHandler.CompleteUpload)
			lanzou.GET("/upload/status/:id", lanzouHandler.UploadStatus)
			lanzou.GET("/profile", lanzouHandler.GetProfile)
		}

		// Recycle bin routes (protected)
		recycle := api.Group("/recycle")
		recycle.Use(middleware.AuthRequired(jwtManager))
		{
			recycle.GET("", recycleHandler.List)
			recycle.POST("/:id/restore", recycleHandler.Restore)
			recycle.DELETE("/:id", recycleHandler.Delete)
			recycle.DELETE("/clear", recycleHandler.Clear)
		}

		// Share parse routes (protected)
		shareParse := api.Group("/lanzou/share")
		shareParse.Use(middleware.AuthRequired(jwtManager))
		{
			shareParse.POST("/parse", shareParseHandler.ParseShare)
			shareParse.POST("/batch-parse", shareParseHandler.BatchParseShare)
			shareParse.POST("/validate", shareParseHandler.ValidateShareURL)
			shareParse.POST("/download", shareParseHandler.GetShareDownloadURL)
		}
	}

	return r
}

// Package-level DB variable, set by main.go
var globalDB *gorm.DB

func SetDB(db *gorm.DB) {
	globalDB = db
}

func getDB() *gorm.DB {
	return globalDB
}
