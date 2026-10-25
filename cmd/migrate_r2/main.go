package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/storage"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/internalsql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("==========================================================")
	log.Println("  MIGRASI DATA UPLOADS LOKAL KE CLOUDFLARE R2 & SINKRON DB")
	log.Println("==========================================================")

	// 1. Inisialisasi Konfigurasi
	err := configs.Init(
		configs.WithConfigFolder([]string{"./internal/configs/"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatalf("[FATAL] Gagal membaca konfigurasi: %v", err)
	}

	cfg := configs.Get()
	r2Cfg := cfg.GetR2Config()

	if r2Cfg.AccountID == "" || r2Cfg.AccessKeyID == "" || r2Cfg.SecretAccessKey == "" || r2Cfg.BucketName == "" {
		log.Fatalf("[FATAL] Kredensial Cloudflare R2 belum lengkap di config.yaml!")
	}

	log.Printf("[INFO] Bucket Target: %s", r2Cfg.BucketName)
	log.Printf("[INFO] Account ID:    %s", r2Cfg.AccountID)

	// 2. Hubungkan ke Database MySQL
	db, err := internalsql.Connect(cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("[FATAL] Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()
	log.Println("[INFO] Terhubung ke MySQL Database.")

	// 3. Inisialisasi Cloudflare R2 Client
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	r2Client, err := storage.NewR2Storage(ctx, r2Cfg)
	if err != nil {
		log.Fatalf("[FATAL] Gagal inisialisasi Cloudflare R2 client: %v", err)
	}
	log.Println("[INFO] Cloudflare R2 client siap.")

	// 4. Baca direktori uploads
	uploadDir := "./uploads"
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Fatalf("[FATAL] Gagal membaca folder %s: %v", uploadDir, err)
	}

	var files []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			files = append(files, entry)
		}
	}

	totalFiles := len(files)
	log.Printf("[INFO] Menemukan %d file untuk dimigrasikan dari folder %s\n\n", totalFiles, uploadDir)

	if totalFiles == 0 {
		log.Println("[INFO] Tidak ada file yang perlu dimigrasi.")
		return
	}

	successCount := 0
	failedCount := 0
	totalDbUpdated := int64(0)

	// 5. Eksekusi Migrasi Per File
	for i, entry := range files {
		filename := entry.Name()
		filePath := filepath.Join(uploadDir, filename)

		fileInfo, err := entry.Info()
		if err != nil {
			log.Printf("[%d/%d] [ERROR] Gagal mendapatkan info file %s: %v", i+1, totalFiles, filename, err)
			failedCount++
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("[%d/%d] [ERROR] Gagal membuka file %s: %v", i+1, totalFiles, filename, err)
			failedCount++
			continue
		}

		// Sniff MIME type
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType := http.DetectContentType(buf[:n])
		_, _ = file.Seek(0, io.SeekStart)

		// Upload stream ke Cloudflare R2
		uploadStartTime := time.Now()
		fileURL, err := r2Client.Upload(ctx, filename, file, fileInfo.Size(), contentType)
		file.Close()

		if err != nil {
			log.Printf("[%d/%d] [ERROR] Gagal upload %s ke R2: %v", i+1, totalFiles, filename, err)
			failedCount++
			continue
		}

		duration := time.Since(uploadStartTime)
		log.Printf("[%d/%d] [UPLOAD OK] %s (%.2f KB, %v) -> %s",
			i+1, totalFiles, filename, float64(fileInfo.Size())/1024.0, duration.Round(time.Millisecond), fileURL)

		// 6. Sinkronisasi Record ke Database MySQL
		dbUpdatedThisFile := int64(0)

		// A. Tabel uploads: update file_path
		res, err := db.ExecContext(ctx,
			`UPDATE uploads 
			 SET file_path = ? 
			 WHERE system_filename = ? OR file_path LIKE ? OR file_path LIKE ?`,
			fileURL, filename, "%"+filename, "%"+filename+"%",
		)
		if err != nil {
			log.Printf("     [WARN] Gagal update tabel uploads untuk %s: %v", filename, err)
		} else if affected, _ := res.RowsAffected(); affected > 0 {
			dbUpdatedThisFile += affected
		}

		// B. Tabel users: avatar_url & banner_url
		resUserAvatar, _ := db.ExecContext(ctx,
			`UPDATE users SET avatar_url = ? WHERE avatar_url LIKE ?`,
			fileURL, "%"+filename,
		)
		if affected, _ := resUserAvatar.RowsAffected(); affected > 0 {
			dbUpdatedThisFile += affected
		}

		resUserBanner, _ := db.ExecContext(ctx,
			`UPDATE users SET banner_url = ? WHERE banner_url LIKE ?`,
			fileURL, "%"+filename,
		)
		if affected, _ := resUserBanner.RowsAffected(); affected > 0 {
			dbUpdatedThisFile += affected
		}

		// C. Tabel posts: post_content jika berisi path lokal
		resPost, _ := db.ExecContext(ctx,
			`UPDATE posts SET post_content = REPLACE(post_content, ?, ?) WHERE post_content LIKE ?`,
			filename, fileURL, "%"+filename+"%",
		)
		if affected, _ := resPost.RowsAffected(); affected > 0 {
			dbUpdatedThisFile += affected
		}

		totalDbUpdated += dbUpdatedThisFile
		log.Printf("     [DB SYNC OK] %d baris diperbarui di database.", dbUpdatedThisFile)
		successCount++
	}

	log.Println("\n==========================================================")
	log.Println("                  RINGKASAN MIGRASI DATA                  ")
	log.Println("==========================================================")
	log.Printf("  Total File Ditemukan : %d\n", totalFiles)
	log.Printf("  Berhasil Diunggah    : %d\n", successCount)
	log.Printf("  Gagal Diunggah       : %d\n", failedCount)
	log.Printf("  Total Baris DB Di-update : %d\n", totalDbUpdated)
	log.Println("==========================================================")
	if failedCount == 0 {
		log.Println("[SELESAI] Seluruh file lokal berhasil dibackup ke Cloudflare R2 dan database telah tersinkronisasi!")
	} else {
		log.Printf("[PERINGATAN] Selesai dengan %d kegagalan. Silahkan periksa log di atas.\n", failedCount)
	}
}
