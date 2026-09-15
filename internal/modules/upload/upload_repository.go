package upload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UploadRepository interface {
	CreateUpload(ctx context.Context, model *UploadModel) error
	GetUploadsByUserID(ctx context.Context, userID int) ([]UploadModel, error)
	GetUploadByID(ctx context.Context, id int) (*UploadModel, error)
}

type uploadRepository struct {
	db *sql.DB
}

func NewUploadRepository(db *sql.DB) UploadRepository {
	return &uploadRepository{db: db}
}

func (r *uploadRepository) CreateUpload(ctx context.Context, model *UploadModel) error {
	query := `INSERT INTO uploads (user_id, file_name, system_filename, file_path, file_type, file_size, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query,
		model.UserID,
		model.FileName,
		model.SystemFilename,
		model.FilePath,
		model.FileType,
		model.FileSize,
		model.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository CreateUpload: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil {
		model.ID = id
	}

	return nil
}

func (r *uploadRepository) GetUploadsByUserID(ctx context.Context, userID int) ([]UploadModel, error) {
	query := `SELECT id, user_id, file_name, system_filename, file_path, file_type, file_size, created_at
	          FROM uploads
	          WHERE user_id = ?
	          ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository GetUploadsByUserID: %w", err)
	}
	defer rows.Close()

	var uploads []UploadModel
	for rows.Next() {
		var item UploadModel
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.FileName,
			&item.SystemFilename,
			&item.FilePath,
			&item.FileType,
			&item.FileSize,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository scan upload: %w", err)
		}
		uploads = append(uploads, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository rows err: %w", err)
	}

	return uploads, nil
}

func (r *uploadRepository) GetUploadByID(ctx context.Context, id int) (*UploadModel, error) {
	query := `SELECT id, user_id, file_name, system_filename, file_path, file_type, file_size, created_at
	          FROM uploads
	          WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var item UploadModel
	if err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.FileName,
		&item.SystemFilename,
		&item.FilePath,
		&item.FileType,
		&item.FileSize,
		&item.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository GetUploadByID: %w", err)
	}

	return &item, nil
}
