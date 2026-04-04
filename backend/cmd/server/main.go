package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtSecret = []byte("your-secret-key-change-in-production")

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type File struct {
	ID               int64
	UserID           int64
	Name             string
	Size             int64
	LanZouFileID     string
	LanZouFolderID  string
	EncryptionKey    string
	EncryptionNonce  string
	MimeType         string
	CreatedAt        time.Time
}

type Config struct {
	ServerPort string
	DBPath     string
}

func LoadConfig() *Config {
	godotenv.Load()
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/wenxi.db"
	}
	return &Config{
		ServerPort: port,
		DBPath:     dbPath,
	}
}

func InitDB(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			size INTEGER NOT NULL,
			lanzou_file_id TEXT,
			lanzou_folder_id TEXT,
			encryption_key TEXT NOT NULL,
			encryption_nonce TEXT NOT NULL,
			mime_type TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS shares (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_id INTEGER NOT NULL,
			share_token TEXT UNIQUE NOT NULL,
			password_hash TEXT,
			expires_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS lanzou_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL UNIQUE,
			cookie TEXT NOT NULL,
			token_value TEXT,
			expires_at DATETIME,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS upload_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			file_name TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			file_hash TEXT NOT NULL,
			chunks_total INTEGER NOT NULL,
			chunks_uploaded INTEGER DEFAULT 0,
			lanzou_upload_url TEXT,
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", int64(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Next()
	}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": data})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": msg})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": msg})
}

// Handlers
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		BadRequest(c, "failed to hash password")
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		req.Username, req.Email, hashedPassword,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			BadRequest(c, "username or email already exists")
			return
		}
		BadRequest(c, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	Success(c, gin.H{"id": id, "username": req.Username, "email": req.Email})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	var user User
	err := db.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		Unauthorized(c, "invalid credentials")
		return
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		Unauthorized(c, "invalid credentials")
		return
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		BadRequest(c, "failed to generate token")
		return
	}

	Success(c, gin.H{"access_token": token, "token_type": "bearer"})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var user User
	err := db.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		uid,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	Success(c, gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "created_at": user.CreatedAt})
}

func ListFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	rows, err := db.Query(
		"SELECT id, name, size, lanzou_file_id, created_at FROM files WHERE user_id = ? ORDER BY created_at DESC",
		uid,
	)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.Name, &f.Size, &f.LanZouFileID, &f.CreatedAt); err != nil {
			continue
		}
		files = append(files, f)
	}

	Success(c, files)
}

func CreateFileMetadata(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var req struct {
		Name            string `json:"name" binding:"required"`
		Size            int64  `json:"size" binding:"required"`
		LanZouFileID    string `json:"lanzou_file_id" binding:"required"`
		LanZouFolderID  string `json:"lanzou_folder_id"`
		EncryptionKey   string `json:"encryption_key" binding:"required"`
		EncryptionNonce string `json:"encryption_nonce" binding:"required"`
		MimeType        string `json:"mime_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	result, err := db.Exec(
		`INSERT INTO files (user_id, name, size, lanzou_file_id, lanzou_folder_id, encryption_key, encryption_nonce, mime_type)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		uid, req.Name, req.Size, req.LanZouFileID, req.LanZouFolderID, req.EncryptionKey, req.EncryptionNonce, req.MimeType,
	)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	Success(c, gin.H{"id": id, "name": req.Name})
}

func GetFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	fileID := c.Param("id")

	var f File
	err := db.QueryRow(
		"SELECT id, name, size, lanzou_file_id, encryption_key, created_at FROM files WHERE id = ? AND user_id = ?",
		fileID, uid,
	).Scan(&f.ID, &f.Name, &f.Size, &f.LanZouFileID, &f.EncryptionKey, &f.CreatedAt)
	if err != nil {
		NotFound(c, "file not found")
		return
	}

	Success(c, f)
}

func DeleteFile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	fileID := c.Param("id")

	result, err := db.Exec("DELETE FROM files WHERE id = ? AND user_id = ?", fileID, uid)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		NotFound(c, "file not found")
		return
	}

	Success(c, gin.H{"message": "file deleted"})
}

func CreateShare(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	fileID := c.Param("id")

	// Verify file ownership
	var ownerID int64
	err := db.QueryRow("SELECT user_id FROM files WHERE id = ?", fileID).Scan(&ownerID)
	if err != nil {
		NotFound(c, "file not found")
		return
	}
	if ownerID != uid {
		Unauthorized(c, "access denied")
		return
	}

	var req struct {
		Password  *string `json:"password"`
		ExpiresAt *string `json:"expires_at"`
	}
	c.ShouldBindJSON(&req)

	// Generate share token
	tokenBytes := make([]byte, 16)
	for i := range tokenBytes {
		tokenBytes[i] = byte(time.Now().UnixNano() % 256)
	}
	shareToken := fmt.Sprintf("%x", tokenBytes)[:32]

	var passwordHash *string
	if req.Password != nil && *req.Password != "" {
		h, _ := HashPassword(*req.Password)
		passwordHash = &h
	}

	result, err := db.Exec(
		"INSERT INTO shares (file_id, share_token, password_hash) VALUES (?, ?, ?)",
		fileID, shareToken, passwordHash,
	)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	Success(c, gin.H{"id": id, "share_token": shareToken, "share_url": "/api/shares/" + shareToken})
}

func GetShare(c *gin.Context) {
	token := c.Param("token")

	var share struct {
		ID           int64
		FileID       int64
		PasswordHash *string
		FileName     string
		FileSize     int64
	}
	err := db.QueryRow(
		`SELECT s.id, s.file_id, s.password_hash, f.name, f.size
		 FROM shares s JOIN files f ON s.file_id = f.id
		 WHERE s.share_token = ?`,
		token,
	).Scan(&share.ID, &share.FileID, &share.PasswordHash, &share.FileName, &share.FileSize)
	if err != nil {
		NotFound(c, "share not found")
		return
	}

	Success(c, gin.H{
		"file_id":          share.FileID,
		"file_name":        share.FileName,
		"file_size":        share.FileSize,
		"requires_password": share.PasswordHash != nil,
	})
}

func DeleteShare(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	shareID := c.Param("id")

	// Verify ownership through file
	var ownerID int64
	err := db.QueryRow(
		`SELECT f.user_id FROM shares s JOIN files f ON s.file_id = f.id WHERE s.id = ?`,
		shareID,
	).Scan(&ownerID)
	if err != nil {
		NotFound(c, "share not found")
		return
	}
	if ownerID != uid {
		Unauthorized(c, "access denied")
		return
	}

	db.Exec("DELETE FROM shares WHERE id = ?", shareID)
	Success(c, gin.H{"message": "share deleted"})
}

func ConnectLanZou(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var req struct {
		Cookie     string `json:"cookie" binding:"required"`
		TokenValue string `json:"token_value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	_, err := db.Exec(
		`INSERT INTO lanzou_tokens (user_id, cookie, token_value, updated_at)
		 VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(user_id) DO UPDATE SET cookie = ?, token_value = ?, updated_at = CURRENT_TIMESTAMP`,
		uid, req.Cookie, req.TokenValue, req.Cookie, req.TokenValue,
	)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"status": "connected"})
}

func GetLanZouStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var token struct {
		Cookie    string
		ExpiresAt sql.NullTime
	}
	err := db.QueryRow("SELECT cookie, expires_at FROM lanzou_tokens WHERE user_id = ?", uid).Scan(&token.Cookie, &token.ExpiresAt)
	if err != nil {
		Success(c, gin.H{"connected": false})
		return
	}

	Success(c, gin.H{"connected": true, "expires_at": token.ExpiresAt})
}

func main() {
	cfg := LoadConfig()

	if err := InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:   []string{"Origin", "Content-Type", "Authorization"},
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
			auth.GET("/me", AuthRequired(), GetCurrentUser)
		}

		files := api.Group("/files")
		files.Use(AuthRequired())
		{
			files.GET("", ListFiles)
			files.POST("", CreateFileMetadata)
			files.GET("/:id", GetFile)
			files.DELETE("/:id", DeleteFile)
			files.POST("/:id/share", CreateShare)
		}

		shares := api.Group("/shares")
		{
			shares.GET("/:token", GetShare)
			shares.DELETE("/:id", AuthRequired(), DeleteShare)
		}

		lanzou := api.Group("/lanzou")
		lanzou.Use(AuthRequired())
		{
			lanzou.POST("/connect", ConnectLanZou)
			lanzou.GET("/status", GetLanZouStatus)
		}
	}

	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
