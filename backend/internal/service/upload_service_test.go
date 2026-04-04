package service

import (
	"errors"
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
}

func (m *mockLanzouService) IsConnected(userID uint) bool {
	return m.connected
}

func (m *mockLanzouService) GetClient(userID uint) (*lanzou.Client, error) {
	if !m.connected {
		return nil, errors.New("lanzou not connected")
	}
	return lanzou.NewClient("test-cookie"), nil
}

// mockFileService 模拟文件服务
type mockFileService struct{}

func (m *mockFileService) CreateMetadata(userID uint, req *CreateFileRequest) (*model.File, error) {
	return &model.File{
		ID:     1,
		UserID: userID,
		Name:   req.Name,
		Size:   req.Size,
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
