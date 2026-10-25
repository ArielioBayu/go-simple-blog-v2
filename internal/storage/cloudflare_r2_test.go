package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	putObjectFunc    func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	deleteObjectFunc func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if m.putObjectFunc != nil {
		return m.putObjectFunc(ctx, params, optFns...)
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if m.deleteObjectFunc != nil {
		return m.deleteObjectFunc(ctx, params, optFns...)
	}
	return &s3.DeleteObjectOutput{}, nil
}

func TestR2Storage_Upload_Success(t *testing.T) {
	mockClient := &mockS3Client{
		putObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			if *params.Bucket != "test-bucket" {
				t.Errorf("expected bucket 'test-bucket', got %s", *params.Bucket)
			}
			if *params.Key != "posts/image.jpg" {
				t.Errorf("expected key 'posts/image.jpg', got %s", *params.Key)
			}
			if *params.ContentType != "image/jpeg" {
				t.Errorf("expected content-type 'image/jpeg', got %s", *params.ContentType)
			}
			content, _ := io.ReadAll(params.Body)
			if string(content) != "fake-image-data" {
				t.Errorf("unexpected body content: %s", string(content))
			}
			return &s3.PutObjectOutput{}, nil
		},
	}

	storage := NewR2StorageWithClient(mockClient, "test-bucket", "https://cdn.example.com", "acc123")
	data := bytes.NewReader([]byte("fake-image-data"))

	url, err := storage.Upload(context.Background(), "posts/image.jpg", data, int64(data.Len()), "image/jpeg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedURL := "https://cdn.example.com/posts/image.jpg"
	if url != expectedURL {
		t.Errorf("expected URL %s, got %s", expectedURL, url)
	}
}

func TestR2Storage_Upload_NilBodyError(t *testing.T) {
	mockClient := &mockS3Client{}
	storage := NewR2StorageWithClient(mockClient, "test-bucket", "", "acc123")

	_, err := storage.Upload(context.Background(), "posts/image.jpg", nil, 0, "image/jpeg")
	if err == nil {
		t.Fatal("expected error for nil body, got nil")
	}
}

func TestR2Storage_Upload_S3Error(t *testing.T) {
	mockClient := &mockS3Client{
		putObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			return nil, errors.New("network failure to R2")
		},
	}

	storage := NewR2StorageWithClient(mockClient, "test-bucket", "", "acc123")
	data := bytes.NewReader([]byte("data"))

	_, err := storage.Upload(context.Background(), "image.png", data, int64(data.Len()), "image/png")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "network failure to R2") {
		t.Errorf("expected error to contain 'network failure to R2', got %v", err)
	}
}

func TestR2Storage_Delete_Success(t *testing.T) {
	deletedKey := ""
	mockClient := &mockS3Client{
		deleteObjectFunc: func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			deletedKey = *params.Key
			return &s3.DeleteObjectOutput{}, nil
		},
	}

	storage := NewR2StorageWithClient(mockClient, "test-bucket", "https://cdn.example.com", "acc123")

	// Test deleting by raw key
	err := storage.Delete(context.Background(), "avatar/pic.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedKey != "avatar/pic.jpg" {
		t.Errorf("expected deletedKey 'avatar/pic.jpg', got %s", deletedKey)
	}

	// Test deleting by full URL
	err = storage.Delete(context.Background(), "https://cdn.example.com/posts/photo.webp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deletedKey != "posts/photo.webp" {
		t.Errorf("expected deletedKey 'posts/photo.webp', got %s", deletedKey)
	}
}

func TestR2Storage_GetPublicURL_Fallback(t *testing.T) {
	mockClient := &mockS3Client{}
	storage := NewR2StorageWithClient(mockClient, "my-bucket", "", "my-account")

	url := storage.GetPublicURL("folder/image.png")
	expected := "https://my-account.r2.cloudflarestorage.com/my-bucket/folder/image.png"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}
