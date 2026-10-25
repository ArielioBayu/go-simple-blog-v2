package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ClientAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// R2Storage mendefinisikan kontrak operasi penyimpanan objek ke Cloudflare R2.
type R2Storage interface {
	Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetPublicURL(key string) string
}

type r2Storage struct {
	client    S3ClientAPI
	bucket    string
	publicURL string
	accountID string
}

// NewR2Storage menginisialisasi client Cloudflare R2 menggunakan konfigurasi terpusat.
func NewR2Storage(ctx context.Context, r2Cfg configs.CloudflareR2) (R2Storage, error) {
	if r2Cfg.AccountID == "" || r2Cfg.AccessKeyID == "" ||
		r2Cfg.SecretAccessKey == "" || r2Cfg.BucketName == "" {
		return nil, fmt.Errorf("R2 configuration is incomplete: account_id, access_key_id, secret_access_key, and bucket_name are required")
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				r2Cfg.AccessKeyID,
				r2Cfg.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load R2 aws config: %w", err)
	}

	endpoint := fmt.Sprintf(
		"https://%s.r2.cloudflarestorage.com",
		r2Cfg.AccountID,
	)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &r2Storage{
		client:    s3Client,
		bucket:    r2Cfg.BucketName,
		publicURL: strings.TrimRight(r2Cfg.PublicURL, "/"),
		accountID: r2Cfg.AccountID,
	}, nil
}

func NewR2StorageWithClient(client S3ClientAPI, bucket, publicURL, accountID string) R2Storage {
	return &r2Storage{
		client:    client,
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
		accountID: accountID,
	}
}
