package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func formatFileSize(bytes int64) string {
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024.0*1024.0))
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, membaca environment sistem.")
	}

	// 1. Cloudinary Client
	cldUrl := os.Getenv("CLOUDINARY_URL")
	cldUrl = strings.TrimSpace(cldUrl)
	cldUrl = strings.TrimPrefix(cldUrl, "CLOUDINARY_URL=")
	cldUrl = strings.Trim(cldUrl, `"'`)

	cld, err := cloudinary.NewFromURL(cldUrl)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Cloudinary: %v", err)
	}

	// 2. Database Connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Printf("⚠️ Gagal koneksi database: %v. Modul tetap diunggah ke Cloudinary.", err)
	} else {
		_ = db.AutoMigrate(&entity.ModulPembelajaran{})
	}

	// 3. Scan Direktori ModuleWEB
	rootDir := "/home/damos/Downloads/ModuleWEB-20261114T103654Z-1-001/ModuleWEB"
	fmt.Printf("🔍 Memindai berkas dari: %s\n", rootDir)

	type FileTask struct {
		Path       string
		Semester   string
		MataKuliah string
		Filename   string
		Size       int64
		Ext        string
	}

	var tasks []FileTask

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(rootDir, path)
		parts := strings.Split(rel, string(filepath.Separator))

		var semester = "Semester 1"
		var matkul = "Umum"

		for _, p := range parts {
			if strings.HasPrefix(strings.ToLower(p), "semester") {
				semester = p
			}
		}

		if len(parts) >= 3 {
			matkul = parts[len(parts)-2]
		} else if len(parts) == 2 {
			matkul = "Umum"
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		tasks = append(tasks, FileTask{
			Path:       path,
			Semester:   semester,
			MataKuliah: matkul,
			Filename:   info.Name(),
			Size:       info.Size(),
			Ext:        strings.TrimPrefix(ext, "."),
		})
		return nil
	})

	if err != nil {
		log.Fatalf("Gagal membaca folder: %v", err)
	}

	total := len(tasks)
	fmt.Printf("📦 Ditemukan %d berkas untuk diunggah ke Cloudinary...\n", total)

	ctx := context.Background()
	useFilename := true
	uniqueFilename := false

	successCount := 0
	failedCount := 0

	for i, t := range tasks {
		folderName := fmt.Sprintf("hmik/learning/%s/%s", t.Semester, t.MataKuliah)
		fmt.Printf("[%d/%d] 🚀 Uploading: %s (%s - %s)...\n", i+1, total, t.Filename, t.Semester, t.MataKuliah)

		file, err := os.Open(t.Path)
		if err != nil {
			fmt.Printf("   ❌ Gagal buka file: %v\n", err)
			failedCount++
			continue
		}

		res, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
			Folder:         folderName,
			ResourceType:   "auto",
			UseFilename:    &useFilename,
			UniqueFilename: &uniqueFilename,
		})
		file.Close()

		if err != nil {
			fmt.Printf("   ❌ Gagal upload ke Cloudinary: %v\n", err)
			failedCount++
			continue
		}

		successCount++
		fmt.Printf("   ✅ Sukses: %s\n", res.SecureURL)

		// Simpan ke PostgreSQL jika DB terkoneksi
		if db != nil {
			cleanTitle := strings.TrimSuffix(t.Filename, filepath.Ext(t.Filename))
			var existing entity.ModulPembelajaran
			dbCheck := db.Where("judul = ? AND semester = ? AND mata_kuliah = ?", cleanTitle, t.Semester, t.MataKuliah).First(&existing)

			if dbCheck.Error == nil && existing.IDModul > 0 {
				existing.FileURL = res.SecureURL
				existing.UkuranFile = formatFileSize(t.Size)
				existing.UkuranBytes = t.Size
				existing.TipeFile = t.Ext
				existing.DiperbaruiPada = time.Now()
				db.Save(&existing)
			} else {
				mod := entity.ModulPembelajaran{
					Judul:       cleanTitle,
					Deskripsi:   fmt.Sprintf("Materi perkuliahan %s (%s)", t.MataKuliah, t.Semester),
					Semester:    t.Semester,
					MataKuliah:  t.MataKuliah,
					TipeFile:    t.Ext,
					UkuranFile:  formatFileSize(t.Size),
					UkuranBytes: t.Size,
					FileURL:     res.SecureURL,
					DibuatOleh:  "Admin HMIK",
				}
				db.Create(&mod)
			}
		}
	}

	fmt.Printf("\n🎉 SELESAI! Berhasil upload %d/%d berkas ke Cloudinary (Gagal: %d)\n", successCount, total, failedCount)
}
