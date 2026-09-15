package upload

import "time"

type UploadModel struct {
	ID             int64     `column:"id"`
	UserID         int64     `column:"user_id"`
	FileName       string    `column:"file_name"`
	SystemFilename string    `column:"system_filename"`
	FilePath       string    `column:"file_path"`
	FileType       string    `column:"file_type"`
	FileSize       int64     `column:"file_size"`
	CreatedAt      time.Time `column:"created_at"`
}

type UploadResponse struct {
	ID             int64     `json:"id"`
	FileName       string    `json:"file_name"`
	SystemFilename string    `json:"system_filename"`
	FilePath       string    `json:"file_path"`
	FileType       string    `json:"file_type"`
	FileSize       int64     `json:"file_size"`
	CreatedAt      time.Time `json:"created_at"`
}
