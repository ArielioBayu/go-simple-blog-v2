package upload_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type mockUploadRepo struct {
	mock.Mock
}

func (m *mockUploadRepo) CreateUpload(ctx context.Context, model *upload.UploadModel) error {
	args := m.Called(ctx, model)
	if args.Get(0) != nil {
		model.ID = 123
	}
	return args.Error(0)
}

func (m *mockUploadRepo) GetUploadsByUserID(ctx context.Context, userID int) ([]upload.UploadModel, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]upload.UploadModel), args.Error(1)
}

func (m *mockUploadRepo) GetUploadByID(ctx context.Context, id int) (*upload.UploadModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*upload.UploadModel), args.Error(1)
}

func (m *mockUploadRepo) GetUploadByPathOrFilename(ctx context.Context, identifier string) (*upload.UploadModel, error) {
	args := m.Called(ctx, identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*upload.UploadModel), args.Error(1)
}

func (m *mockUploadRepo) DeleteUpload(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock R2 storage
type mockR2Storage struct {
	mock.Mock
}

func (m *mockR2Storage) Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) (string, error) {
	args := m.Called(ctx, key, body, size, contentType)
	return args.String(0), args.Error(1)
}

func (m *mockR2Storage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockR2Storage) GetPublicURL(key string) string {
	args := m.Called(key)
	return args.String(0)
}

// Helper to create valid multipart file header
func createMultipartFileHeader(filename string, content []byte) (*multipart.FileHeader, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(content); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}
	return req.MultipartForm.File["file"][0], nil
}

// Minimal valid PNG header bytes
var validPNGBytes = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
	0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	0x42, 0x60, 0x82,
}

func TestUploadService_UploadImage_R2Success(t *testing.T) {
	fileHeader, err := createMultipartFileHeader("test.png", validPNGBytes)
	assert.NoError(t, err)

	mockRepo := new(mockUploadRepo)
	mockStorage := new(mockR2Storage)

	expectedURL := "https://pub-mock.r2.dev/unique-test.png"

	mockStorage.On("Upload", mock.Anything, mock.AnythingOfType("string"), mock.Anything, int64(len(validPNGBytes)), "image/png").
		Return(expectedURL, nil).
		Once()

	mockRepo.On("CreateUpload", mock.Anything, mock.MatchedBy(func(m *upload.UploadModel) bool {
		return m.UserID == 1 && m.FileName == "test.png" && m.FilePath == expectedURL && m.FileType == "image/png"
	})).Return(nil).Once()

	cfg := &configs.Config{}
	svc := upload.NewUploadService(cfg, mockRepo, mockStorage)

	res, err := svc.UploadImage(context.Background(), 1, fileHeader, "http", "localhost:9888")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test.png", res.FileName)
	assert.Equal(t, expectedURL, res.FilePath)
	assert.Equal(t, "image/png", res.FileType)

	mockStorage.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestUploadService_UploadImage_R2Error(t *testing.T) {
	fileHeader, err := createMultipartFileHeader("test.png", validPNGBytes)
	assert.NoError(t, err)

	mockRepo := new(mockUploadRepo)
	mockStorage := new(mockR2Storage)

	mockStorage.On("Upload", mock.Anything, mock.AnythingOfType("string"), mock.Anything, int64(len(validPNGBytes)), "image/png").
		Return("", errors.New("s3 connection failed")).
		Once()

	cfg := &configs.Config{}
	svc := upload.NewUploadService(cfg, mockRepo, mockStorage)

	res, err := svc.UploadImage(context.Background(), 1, fileHeader, "http", "localhost:9888")
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "s3 connection failed")

	mockStorage.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "CreateUpload")
}

func TestUploadService_UploadImage_InvalidFormat(t *testing.T) {
	// Plain text content is not an allowed image format
	fileHeader, err := createMultipartFileHeader("test.txt", []byte("plain text content here"))
	assert.NoError(t, err)

	mockRepo := new(mockUploadRepo)
	mockStorage := new(mockR2Storage)

	cfg := &configs.Config{}
	svc := upload.NewUploadService(cfg, mockRepo, mockStorage)

	res, err := svc.UploadImage(context.Background(), 1, fileHeader, "http", "localhost:9888")
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, constants.ErrInvalidFormat, err)

	mockStorage.AssertNotCalled(t, "Upload")
	mockRepo.AssertNotCalled(t, "CreateUpload")
}

func TestUploadService_UploadImage_LocalFallback(t *testing.T) {
	fileHeader, err := createMultipartFileHeader("fallback.png", validPNGBytes)
	assert.NoError(t, err)

	mockRepo := new(mockUploadRepo)

	mockRepo.On("CreateUpload", mock.Anything, mock.MatchedBy(func(m *upload.UploadModel) bool {
		return m.UserID == 2 && m.FileName == "fallback.png" && m.FileType == "image/png"
	})).Return(nil).Once()

	cfg := &configs.Config{}
	// r2Storage is nil -> should fallback to local
	svc := upload.NewUploadService(cfg, mockRepo, nil)

	res, err := svc.UploadImage(context.Background(), 2, fileHeader, "http", "localhost:9888")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "fallback.png", res.FileName)

	// Clean up generated local file
	if res.FilePath != "" {
		_ = os.Remove(filepath.FromSlash(res.FilePath))
	}

	mockRepo.AssertExpectations(t)
}
