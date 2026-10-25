package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Upload mengalirkan (stream) data file ke bucket Cloudflare R2 secara hemat memori
// dan mengembalikan URL publik file yang tersimpan.
func (r *r2Storage) Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) (string, error) {
	if body == nil {
		return "", fmt.Errorf("upload body cannot be nil")
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	}

	if size > 0 {
		input.ContentLength = aws.Int64(size)
	}

	_, err := r.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("upload object to R2: %w", err)
	}

	return r.GetPublicURL(key), nil
}

// Delete menghapus objek dari Cloudflare R2 berdasarkan nama key atau full URL-nya.
func (r *r2Storage) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}

	cleanKey := key
	if r.publicURL != "" && strings.HasPrefix(key, r.publicURL) {
		cleanKey = strings.TrimPrefix(key, r.publicURL)
		cleanKey = strings.TrimPrefix(cleanKey, "/")
	}

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(cleanKey),
	})
	if err != nil {
		return fmt.Errorf("delete object from R2: %w", err)
	}

	return nil
}

// GetPublicURL merakit URL publik/CDN untuk akses gambar.
func (r *r2Storage) GetPublicURL(key string) string {
	cleanKey := strings.TrimLeft(key, "/")
	if r.publicURL != "" {
		return fmt.Sprintf("%s/%s", r.publicURL, cleanKey)
	}
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", r.accountID, r.bucket, cleanKey)
}
