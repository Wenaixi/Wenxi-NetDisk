package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/model"
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

// mockUploadRepo 模拟上传仓库
type mockUploadRepo struct {
	sessions []*model.UploadSession
	nextID   uint
}

func newMockUploadRepo() *mockUploadRepo {
	return &mockUploadRepo{
		sessions: []*model.UploadSession{},
		nextID:   1,
	}
}

func (m *mockUploadRepo) Create(session *model.UploadSession) error {
	session.ID = m.nextID
	m.nextID++
	m.sessions = append(m.sessions, session)
	return nil
}

func (m *mockUploadRepo) FindByID(id uint) (*model.UploadSession, error) {
	for _, s := range m.sessions {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("session not found")
}

func (m *mockUploadRepo) FindByUserIDAndHash(userID uint, fileHash string) (*model.UploadSession, error) {
	for _, s := range m.sessions {
		if s.UserID == userID && s.FileHash == fileHash {
			return s, nil
		}
	}
	return nil, errors.New("session not found")
}

func (m *mockUploadRepo) Update(session *model.UploadSession) error {
	for i, s := range m.sessions {
		if s.ID == session.ID {
			m.sessions[i] = session
			return nil
		}
	}
	return errors.New("session not found")
}

func (m *mockUploadRepo) Delete(id uint) error {
	return nil
}

func (m *mockUploadRepo) DeleteByUserID(userID uint) error {
	return nil
}

// mockLanzouService 模拟蓝奏云服务
type mockLanzouService struct {
	connected bool
	client    *lanzou.Client
}

func (m *mockLanzouService) IsConnected(userID uint) bool {
	return m.connected
}

func (m *mockLanzouService) GetClient(userID uint) (*lanzou.Client, error) {
	if !m.connected {
		return nil, errors.New("lanzou not connected")
	}
	if m.client != nil {
		return m.client, nil
	}
	return lanzou.NewClient("test-cookie"), nil
}

// setupMockUploadServer 创建模拟上传服务器并返回配置好的客户端
func setupMockUploadServer(t *testing.T) (*httptest.Server, *lanzou.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/html5up.php" {
			w.Write([]byte(`{"zt":1,"info":"上传成功","text":[{"f_id":"123","is_newd":"https://wwn.lanzouf.com","downs":"0","icon":"txt","id":"456","name":"test.txt","size":"1024","time":"2024-01-01","onof":"0"}]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := lanzou.NewClient("test-cookie")
	client.SetBaseURL(server.URL)
	client.SetUploadBaseURL(server.URL)
	return server, client
}

// mockFileService 模拟文件服务
type mockFileService struct{}

func (m *mockFileService) CreateMetadata(userID uint, req *CreateFileRequest) (*model.File, error) {
	return &model.File{
		ID:     1,
		UserID: userID,
		Name:   req.Name,
		Size:   req.Size,
		LanZouFileID: req.LanZouFileID,
	}, nil
}

func (m *mockFileService) ListFiles(userID uint) ([]model.File, error) {
	return nil, nil
}

func (m *mockFileService) GetFile(userID, fileID uint) (*model.File, error) {
	return nil, nil
}

func (m *mockFileService) DeleteFile(userID, fileID uint) error {
	return nil
}

// TestUploadService_InitializeUpload tests upload session initialization
func TestUploadService_InitializeUpload(t *testing.T) {
	t.Run("should create new upload session", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
			MimeType: "text/plain",
		}

		resp, err := svc.InitializeUpload(1, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Error("expected response, got nil")
			return
		}
		if resp.SessionID == 0 {
			t.Error("expected session ID > 0")
		}
		if resp.TotalChunks != 1 {
			t.Errorf("expected 1 chunk, got %d", resp.TotalChunks)
		}
		if resp.ChunkSize != ChunkSize {
			t.Errorf("expected chunk size %d, got %d", ChunkSize, resp.ChunkSize)
		}
	})

	t.Run("should return existing session if not completed", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		resp1, err1 := svc.InitializeUpload(1, req)
		if err1 != nil {
			t.Errorf("expected no error on first call, got %v", err1)
			return
		}

		resp2, err2 := svc.InitializeUpload(1, req)
		if err2 != nil {
			t.Errorf("expected no error on second call, got %v", err2)
			return
		}

		if resp1.SessionID != resp2.SessionID {
			t.Errorf("expected same session ID, got %d and %d", resp1.SessionID, resp2.SessionID)
		}
	})

	t.Run("should fail if lanzou not connected", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: false}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		resp, err := svc.InitializeUpload(1, req)

		if err == nil {
			t.Error("expected error when lanzou not connected")
		}
		if resp != nil {
			t.Error("expected nil response when lanzou not connected")
		}
	})
}

// TestUploadService_GetUploadStatus tests upload status retrieval
func TestUploadService_GetUploadStatus(t *testing.T) {
	t.Run("should return session for valid user", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		resp, _ := svc.InitializeUpload(1, req)

		session, err := svc.GetUploadStatus(1, resp.SessionID)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if session == nil {
			t.Error("expected session, got nil")
			return
		}
		if session.FileName != "test.txt" {
			t.Errorf("expected filename 'test.txt', got '%s'", session.FileName)
		}
	})

	t.Run("should fail for wrong user", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		resp, _ := svc.InitializeUpload(1, req)

		session, err := svc.GetUploadStatus(999, resp.SessionID)

		if err == nil {
			t.Error("expected error for wrong user")
		}
		if session != nil {
			t.Error("expected nil session for wrong user")
		}
	})
}

// TestUploadService_ResumeUpload tests upload resumption
func TestUploadService_ResumeUpload(t *testing.T) {
	t.Run("should resume existing incomplete upload", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		_, _ = svc.InitializeUpload(1, req)

		// 使用相同的文件哈希来恢复上传
		hash := calculateFileHash("test.txt", 1024)
		resumeResp, err := svc.ResumeUpload(1, hash)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resumeResp == nil {
			t.Error("expected response, got nil")
			return
		}
		if resumeResp.SessionID == 0 {
			t.Errorf("expected non-zero session ID")
		}
	})

	t.Run("should fail if no session found", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)

		_, err := svc.ResumeUpload(1, "nonexistent-hash")
		if err == nil {
			t.Error("expected no session found error")
		}
	})

	t.Run("should fail if upload already completed", func(t *testing.T) {
		repo := newMockUploadRepo()
		_, client := setupMockUploadServer(t)
		lanzouSvc := &mockLanzouService{connected: true, client: client}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024, // 1 chunk
		}

		resp, _ := svc.InitializeUpload(1, req)

		// Upload the only chunk to complete
_, _ = svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		hash := calculateFileHash("test.txt", 1024)
		_, err := svc.ResumeUpload(1, hash)
		if err == nil {
			t.Error("expected already completed error")
		}
	})
}

// TestUploadService_CalculateFileHash tests file hash calculation
func TestUploadService_CalculateFileHash(t *testing.T) {
	hash1 := calculateFileHash("test.txt", 1024)
	hash2 := calculateFileHash("test.txt", 1024)
	hash3 := calculateFileHash("different.txt", 1024)

	if hash1 != hash2 {
		t.Error("same file should produce same hash")
	}
	if hash1 == hash3 {
		t.Error("different files should produce different hashes")
	}
	if len(hash1) != 32 {
		t.Errorf("expected 32 char hash, got %d", len(hash1))
	}
}

// TestUploadService_UploadChunk tests chunk upload
func TestUploadService_UploadChunk(t *testing.T) {
	t.Run("should increment chunks uploaded", func(t *testing.T) {
		repo := newMockUploadRepo()
		_, client := setupMockUploadServer(t)
		lanzouSvc := &mockLanzouService{connected: true, client: client}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 10 * 1024 * 1024, // 10MB, 5 chunks
		}

		resp, _ := svc.InitializeUpload(1, req)

_, err := svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("chunk data"),
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		session, _ := repo.FindByID(resp.SessionID)
		if session.ChunksUploaded != 1 {
			t.Errorf("expected 1 chunk uploaded, got %d", session.ChunksUploaded)
		}
	})

	t.Run("should fail for non-existent session", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)

_, err := svc.UploadChunk(1, 999, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		if err == nil {
			t.Error("expected session not found error")
		}
	})

	t.Run("should fail for wrong user", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		resp, _ := svc.InitializeUpload(1, req)

_, err := svc.UploadChunk(999, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		if err == nil {
			t.Error("expected access denied error")
		}
	})

	t.Run("should fail for completed upload", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024, // 1 chunk
		}

		resp, _ := svc.InitializeUpload(1, req)

		// Upload the only chunk to complete
_, _ = svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		// Try to upload again
_, err := svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		if err == nil {
			t.Error("expected already completed error")
		}
	})

	t.Run("should fail for invalid chunk index", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024, // 1 chunk
		}

		resp, _ := svc.InitializeUpload(1, req)

_, err := svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 5,
			Data:       []byte("data"),
		})

		if err == nil {
			t.Error("expected invalid chunk index error")
		}
	})

	t.Run("should mark completed when all chunks uploaded", func(t *testing.T) {
		repo := newMockUploadRepo()
		_, client := setupMockUploadServer(t)
		lanzouSvc := &mockLanzouService{connected: true, client: client}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024, // 1 chunk
		}

		resp, _ := svc.InitializeUpload(1, req)

_, err := svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		session, _ := repo.FindByID(resp.SessionID)
		if session.Status != StatusCompleted {
			t.Errorf("expected status completed, got %s", session.Status)
		}
	})
}

// TestUploadService_CompleteUpload tests completing upload
func TestUploadService_CompleteUpload(t *testing.T) {
	t.Run("should create file metadata after all chunks uploaded", func(t *testing.T) {
		repo := newMockUploadRepo()
		_, client := setupMockUploadServer(t)
		lanzouSvc := &mockLanzouService{connected: true, client: client}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024, // 1 chunk
		}

		resp, _ := svc.InitializeUpload(1, req)

		// Upload the chunk
_, _ = svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
			ChunkIndex: 0,
			Data:       []byte("data"),
		})

		// Complete upload
		completeReq := &CompleteUploadRequest{
			EncryptionKey:   "enc-key-123",
			EncryptionNonce: "nonce-456",
			LanZouFileID:    "lanzou-789",
		}

		file, err := svc.CompleteUpload(1, resp.SessionID, completeReq)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if file == nil {
			t.Error("expected file, got nil")
			return
		}
		if file.Name != "test.txt" {
			t.Errorf("expected filename 'test.txt', got '%s'", file.Name)
		}
		if file.LanZouFileID != "lanzou-789" {
			t.Errorf("expected lanzou file ID 'lanzou-789', got '%s'", file.LanZouFileID)
		}
	})

	t.Run("should fail for non-existent session", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)

		_, err := svc.CompleteUpload(1, 999, &CompleteUploadRequest{
			EncryptionKey:   "key",
			EncryptionNonce: "nonce",
			LanZouFileID:    "id",
		})

		if err == nil {
			t.Error("expected session not found error")
		}
	})

	t.Run("should fail for wrong user", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)
		req := &InitializeUploadRequest{
			FileName: "test.txt",
			FileSize: 1024,
		}

		resp, _ := svc.InitializeUpload(1, req)

		_, err := svc.CompleteUpload(999, resp.SessionID, &CompleteUploadRequest{
			EncryptionKey:   "key",
			EncryptionNonce: "nonce",
			LanZouFileID:    "id",
		})

		if err == nil {
			t.Error("expected access denied error")
		}
	})
}

// TestUploadService_GetLanZouClient tests getting lanzou client
func TestUploadService_GetLanZouClient(t *testing.T) {
	t.Run("should return client when connected", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: true}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)

		client, err := svc.GetLanZouClient(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if client == nil {
			t.Error("expected client, got nil")
		}
	})

	t.Run("should return error when not connected", func(t *testing.T) {
		repo := newMockUploadRepo()
		lanzouSvc := &mockLanzouService{connected: false}
		fileSvc := &mockFileService{}

		svc := NewUploadService(repo, lanzouSvc, fileSvc)

		client, err := svc.GetLanZouClient(1)

		if err == nil {
			t.Error("expected not connected error")
		}
		if client != nil {
			t.Error("expected nil client when not connected")
		}
	})
}

// TestUploadService_InitializeUpload_NotConnected 测试未连接蓝奏云
func TestUploadService_InitializeUpload_NotConnected(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: false}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	req := &InitializeUploadRequest{FileName: "test.txt", FileSize: 1024}

	_, err := svc.InitializeUpload(1, req)
	if err == nil {
		t.Error("expected lanzou not connected error")
	}
}

// TestUploadService_InitializeUpload_LargeFile 测试大文件分块
func TestUploadService_InitializeUpload_LargeFile(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	// 5MB file should be 3 chunks (5MB / 2MB = 2.5 → ceil = 3)
	req := &InitializeUploadRequest{FileName: "large.bin", FileSize: 5 * 1024 * 1024}

	resp, err := svc.InitializeUpload(1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.TotalChunks != 3 {
		t.Errorf("expected 3 chunks, got %d", resp.TotalChunks)
	}
}

// TestUploadService_UploadChunk_NegativeIndex 测试负分块索引
func TestUploadService_UploadChunk_NegativeIndex(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	req := &InitializeUploadRequest{FileName: "test.txt", FileSize: 1024}
	resp, _ := svc.InitializeUpload(1, req)

_, err := svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
		ChunkIndex: -1,
		Data:       []byte("data"),
	})
	if err == nil {
		t.Error("expected invalid chunk index error for negative index")
	}
}

// TestUploadService_CompleteUpload_FileSvcError 测试文件服务创建失败
func TestUploadService_CompleteUpload_FileSvcError(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileServiceWithError{err: errors.New("db error")}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	req := &InitializeUploadRequest{FileName: "test.txt", FileSize: 1024}
	resp, _ := svc.InitializeUpload(1, req)

	// Upload chunk to make it completable
	_, _ = svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{ChunkIndex: 0, Data: []byte("data")})

	_, err := svc.CompleteUpload(1, resp.SessionID, &CompleteUploadRequest{
		EncryptionKey: "key", EncryptionNonce: "nonce", LanZouFileID: "id",
	})
	if err == nil {
		t.Error("expected error from file service")
	}
}

// TestUploadService_GetUploadStatus_SessionNotFound 测试会话不存在
func TestUploadService_GetUploadStatus_SessionNotFound(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	_, err := svc.GetUploadStatus(1, 999)
	if err == nil {
		t.Error("expected session not found error")
	}
}

// TestUploadService_ResumeUpload_CompletedViaChunks 测试通过分块上传完成后恢复
func TestUploadService_ResumeUpload_CompletedViaChunks(t *testing.T) {
	repo := newMockUploadRepo()
	_, client := setupMockUploadServer(t)
	lanzouSvc := &mockLanzouService{connected: true, client: client}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	req := &InitializeUploadRequest{FileName: "test.txt", FileSize: 1024}
	resp, _ := svc.InitializeUpload(1, req)

	// Complete the upload
	_, _ = svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{ChunkIndex: 0, Data: []byte("data")})

	// Try to resume
	hash := calculateFileHash("test.txt", 1024)
	_, err := svc.ResumeUpload(1, hash)
	if err == nil {
		t.Error("expected already completed error")
	}
}

// TestUploadService_ResumeUpload_NotFound 测试恢复不存在的会话
func TestUploadService_ResumeUpload_NotFound(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	_, err := svc.ResumeUpload(1, "nonexistent")
	if err == nil {
		t.Error("expected no upload session found error")
	}
}

// mockFileServiceWithError 模拟返回错误的文件服务
type mockFileServiceWithError struct {
	err error
}

func (m *mockFileServiceWithError) CreateMetadata(userID uint, req *CreateFileRequest) (*model.File, error) {
	return nil, m.err
}

func (m *mockFileServiceWithError) ListFiles(userID uint) ([]model.File, error) {
	return nil, nil
}

func (m *mockFileServiceWithError) GetFile(userID, fileID uint) (*model.File, error) {
	return nil, nil
}

func (m *mockFileServiceWithError) DeleteFile(userID, fileID uint) error {
	return nil
}

// TestUploadService_InitializeUpload_ZeroByteFile 测试零字节文件
func TestUploadService_InitializeUpload_ZeroByteFile(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	req := &InitializeUploadRequest{FileName: "empty.txt", FileSize: 0}

	resp, err := svc.InitializeUpload(1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// 0 byte file should have 0 chunks
	if resp.TotalChunks != 0 {
		t.Errorf("expected 0 chunks for empty file, got %d", resp.TotalChunks)
	}
}

// TestUploadService_UploadChunk_OutOfRange 测试超出范围的分块索引
func TestUploadService_UploadChunk_OutOfRange(t *testing.T) {
	repo := newMockUploadRepo()
	lanzouSvc := &mockLanzouService{connected: true}
	fileSvc := &mockFileService{}

	svc := NewUploadService(repo, lanzouSvc, fileSvc)
	// 3MB file = 2 chunks
	req := &InitializeUploadRequest{FileName: "test.bin", FileSize: 3 * 1024 * 1024}
	resp, _ := svc.InitializeUpload(1, req)

	// Try to upload chunk index 2 (should only be 0, 1)
_, err := svc.UploadChunk(1, resp.SessionID, &UploadChunkRequest{
		ChunkIndex: 2,
		Data:       []byte("data"),
	})
	if err == nil {
		t.Error("expected invalid chunk index error")
	}
}
