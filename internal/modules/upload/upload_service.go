package upload

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
)

const (
	MaxUploadSize = 5 * 1024 * 1024 // 5MB
	UploadDir     = "./uploads"
)

var (
	ErrFileRequired  = errors.New("file is required (key: 'file')")
	ErrFileTooLarge  = errors.New("file size exceeds 5MB limit")
	ErrInvalidFormat = errors.New("invalid file format. Allowed formats: jpg, jpeg, png, webp, gif")
)

var allowedMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type UploadService interface {
	UploadImage(ctx context.Context, userID int, fileHeader *multipart.FileHeader, scheme, host string) (*UploadResponse, error)
	GetUserUploads(ctx context.Context, userID int) ([]UploadResponse, error)
}

type uploadService struct {
	cfg        *configs.Config
	uploadRepo UploadRepository
}

func NewUploadService(cfg *configs.Config, uploadRepo UploadRepository) UploadService {
	return &uploadService{
		cfg:        cfg,
		uploadRepo: uploadRepo,
	}
}

func (s *uploadService) UploadImage(ctx context.Context, userID int, fileHeader *multipart.FileHeader, scheme, host string) (*UploadResponse, error) {
	if fileHeader == nil {
		return nil, ErrFileRequired
	}

	if fileHeader.Size > MaxUploadSize {
		return nil, ErrFileTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Sniff 512 bytes pertama untuk menentukan MIME type secara akurat (security best practice)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])
	expectedExt, isAllowed := allowedMIMETypes[contentType]
	if !isAllowed {
		return nil, ErrInvalidFormat
	}

	// Kembalikan pointer pembaca ke awal file
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	clientExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	ext := expectedExt
	if clientExt == ".jpeg" && expectedExt == ".jpg" {
		ext = ".jpeg"
	}

	randomBytes := make([]byte, 4)
	_, _ = rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	uniqueFilename := fmt.Sprintf("%d-%s%s", time.Now().Unix(), randomHex, ext)
	dstPath := filepath.Join(UploadDir, uniqueFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Normalisasi path ke forward slash untuk kompatibilitas web (cross-platform)
	normalizedPath := filepath.ToSlash(dstPath)

	// Simpan metadata ke tabel uploads melalui repository
	model := UploadModel{
		UserID:         int64(userID),
		FileName:       fileHeader.Filename,
		SystemFilename: uniqueFilename,
		FilePath:       normalizedPath,
		FileType:       contentType,
		FileSize:       fileHeader.Size,
		CreatedAt:      time.Now(),
	}

	if err := s.uploadRepo.CreateUpload(ctx, &model); err != nil {
		return nil, fmt.Errorf("failed to record upload metadata: %w", err)
	}

	return &UploadResponse{
		ID:             model.ID,
		FileName:       model.FileName,
		SystemFilename: model.SystemFilename,
		FilePath:       model.FilePath,
		FileType:       model.FileType,
		FileSize:       model.FileSize,
		CreatedAt:      model.CreatedAt,
	}, nil
}

func (s *uploadService) GetUserUploads(ctx context.Context, userID int) ([]UploadResponse, error) {
	records, err := s.uploadRepo.GetUploadsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]UploadResponse, 0, len(records))
	for _, r := range records {
		responses = append(responses, UploadResponse{
			ID:             r.ID,
			FileName:       r.FileName,
			SystemFilename: r.SystemFilename,
			FilePath:       r.FilePath,
			FileType:       r.FileType,
			FileSize:       r.FileSize,
			CreatedAt:      r.CreatedAt,
		})
	}

	return responses, nil
}
