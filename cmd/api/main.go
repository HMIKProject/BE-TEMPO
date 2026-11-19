package main

import (
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/HMIKProject/hmik-corex-backend/docs"
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/routes"
	"github.com/HMIKProject/hmik-corex-backend/internal/seeder"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// @title           HMIK CoreX API
// @version         1.0
// @description     API Dokumentasi untuk Backend Sistem Himpunan Mahasiswa Ilmu Komputer
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	// Memuat file .env jika ada (Berguna untuk lokal, akan diabaikan jika di Cloud tanpa .env)
	godotenv.Load()

	// 1. Koneksi ke Database
	// Render akan menyediakan DATABASE_URL, jika kosong kita pakai fallback ke localhost
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
		log.Fatal("Gagal terkoneksi ke database:", err)
	}

	// Memastikan kolom-kolom baru (token_version, kode_aktivasi, status_aktif, dll.) selalu termigrasi
	_ = db.AutoMigrate(
		&entity.Pengguna{},
		&entity.Anggota{},
		&entity.Departemen{},
		&entity.DivisiDepartemen{},
		&entity.ModulPembelajaran{},
		&entity.Pemira{},
		&entity.Paslon{},
		&entity.Pemilih{},
		&entity.Suara{},
		&entity.Kabinet{},
		&entity.RiwayatJabatan{},
	)

	// Bersihkan data dummy anggota pengujian (Dankon & Anton) agar bagan organisasi bersih
	db.Exec("DELETE FROM anggota WHERE nama_lengkap IN ('Dankon', 'Anton Chandrawansah') OR nim IN ('105224003', '105224000')")

	// Inisialisasi otomatis akun slot jabatan dengan password 'qwerqwer'
	_ = seeder.SeedAllRoleAccounts(db)
	// _ = seeder.SeedSampleModul(db) // Ditiadakan agar modul tetap bersih untuk input manual pengurus

	// Memperbaiki URL file modul lama dari /uploads/... ke Cloudinary CDN jika ada
	db.Model(&entity.ModulPembelajaran{}).
		Where("file_url LIKE ?", "%/uploads/modul/semester_1/kalkulus_1%").
		Update("file_url", "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654581/hmik/learning/Semester%201/Kalkulus%201/Kalkulus_1.pdf")
	db.Model(&entity.ModulPembelajaran{}).
		Where("file_url LIKE ?", "%/uploads/modul/semester_1/matematika_diskrit%").
		Update("file_url", "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654588/hmik/learning/Semester%201/Matematika%20Diskrit/10-13.TeoriHimpunanElementer.pdf")
	db.Model(&entity.ModulPembelajaran{}).
		Where("file_url LIKE ?", "%/uploads/modul/semester_1/organisasi_komputer%").
		Update("file_url", "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654637/hmik/learning/Semester%201/Organisasi%20Komputer/Pertemuan%201%20-%20Pengantar%20Organisasi%20Komputer.pdf")
	db.Model(&entity.ModulPembelajaran{}).
		Where("file_url LIKE ?", "%/uploads/modul/semester_2/dasar_pemrograman%").
		Update("file_url", "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654662/hmik/learning/Semester%202/Dasar%20Pemrograman/UAS_DasPro_2024.pdf")
	db.Model(&entity.ModulPembelajaran{}).
		Where("file_url LIKE ?", "%/uploads/modul/semester_3/algoritma_dan_struktur_data%").
		Update("file_url", "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654669/hmik/learning/Semester%203/Algoritma%20dan%20Struktur%20Data/W1_-_Array_and_Struct.pdf")
	db.Model(&entity.ModulPembelajaran{}).
		Where("file_url LIKE ?", "%/uploads/modul/semester_4/pbo%").
		Update("file_url", "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654795/hmik/learning/Semester%204/PBO%20%28Pemrograman%20Berorientasi%20Objek%29/11-Polymorphism.pdf")

	// Memastikan selalu ada setidaknya satu record sesi Pemira berstatus DRAFT jika tabel kosong
	var pemiraCount int64
	db.Model(&entity.Pemira{}).Count(&pemiraCount)
	if pemiraCount == 0 {
		now := time.Now().UTC()
		defaultPemira := entity.Pemira{
			TahunPeriode:   "2026/2027",
			JudulPemilihan: "Pemilihan Ketua & Wakil Ketua HMIK",
			TanggalMulai:   now,
			TanggalSelesai: now.Add(7 * 24 * time.Hour),
			Status:         "DRAFT",
		}
		if err := db.Create(&defaultPemira).Error; err != nil {
			log.Printf("Peringatan: Gagal auto-seed DRAFT pemira: %v", err)
		} else {
			log.Printf("✅ Auto-seed DRAFT Pemira berhasil dibuat dengan ID: %d", defaultPemira.IdPemira)
		}
	}

	// 2. Inisiasi Mesin Gin
	router := gin.Default()

	// 3. Konfigurasi Middleware CORS Terproteksi (Sesuai Panduan HMIK_note)
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://localhost:5173" || origin == "http://localhost:3000" || origin == "https://hmik.kampus.ac.id" {
				return true
			}
			if strings.HasSuffix(origin, ".netlify.app") || 
			   strings.HasSuffix(origin, ".pages.dev") || 
			   strings.HasSuffix(origin, ".workers.dev") || 
			   strings.Contains(origin, "hmikup.id") || 
			   strings.Contains(origin, "gowebdigital.agency") {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static file serving untuk berkas unggahan lokal (uploads/modul, dll.)
	router.Static("/uploads", "./uploads")

	// Middleware Anti-Stale Cache: Cegah browser menahan data usang dari API publik & admin
	router.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	})

	// 4. Mendaftarkan Seluruh Endpoint API
	routes.SetupRoutes(router, db)

	// Endpoint Diagnostik Versi Live & Seeder On-Demand
	router.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app":     "HMIK CoreX Backend",
			"version": "v16-fix-riwayat-modul-cache",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	router.GET("/api/v1/auth/debug-accounts", func(c *gin.Context) {
		var users []struct {
			IdPengguna   int64  `json:"id_pengguna"`
			NamaPengguna string `json:"nama_pengguna"`
			Email        string `json:"email"`
			StatusAktif  bool   `json:"status_aktif"`
		}
		db.Table("pengguna").Select("id_pengguna, nama_pengguna, email, status_aktif").Order("id_pengguna desc").Limit(35).Scan(&users)
		c.JSON(200, gin.H{
			"total_accounts_sample": len(users),
			"accounts":              users,
		})
	})

	router.GET("/api/v1/auth/force-seed", func(c *gin.Context) {
		err := seeder.SeedAllRoleAccounts(db)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"status":           "success",
			"message":          "Seeder akun role berhasil dieksekusi secara live!",
			"default_password": "qwerqwer",
		})
	})

	// 5. Menjalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server Backend berjalan di port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}