package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://adminhmik:Hmik_DB123@hmik-db-server.postgres.database.azure.com:5432/postgres?sslmode=require"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}

	fmt.Println("🚀 Menjalankan migrasi tabel Pemira (32 s/d 36)...")
	migrations := []string{
		"migration/000032_create_kabinet.up.sql",
		"migration/000033_create_pemira.up.sql",
		"migration/000034_create_paslon.up.sql",
		"migration/000035_create_pemilih.up.sql",
		"migration/000036_create_suara.up.sql",
	}

	for _, file := range migrations {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Gagal membaca migrasi %s: %v", file, err)
		}
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			log.Printf("Peringatan migrasi %s: %v", file, err)
		} else {
			fmt.Printf("✅ Migrasi sukses: %s\n", file)
		}
	}

	// 1. Cek atau Buat Acara Pemira
	var activePemira entity.Pemira
	err = db.Where("status = ?", "BERJALAN").First(&activePemira).Error
	if err != nil {
		now := time.Now().UTC()
		activePemira = entity.Pemira{
			TahunPeriode:   "2026/2027",
			TanggalMulai:   now.Add(-24 * time.Hour),
			TanggalSelesai: now.Add(7 * 24 * time.Hour),
			Status:         "BERJALAN",
		}
		if err := db.Create(&activePemira).Error; err != nil {
			log.Fatalf("Gagal membuat data pemira: %v", err)
		}
		fmt.Printf("✅ Pemira baru dibuat dengan ID: %d (Status: BERJALAN)\n", activePemira.IdPemira)
	} else {
		fmt.Printf("ℹ️ Pemira aktif sudah ada dengan ID: %d (Status: %s)\n", activePemira.IdPemira, activePemira.Status)
	}

	// 2. Cek atau Buat Kandidat Paslon
	var countPaslon int64
	db.Model(&entity.Paslon{}).Where("id_pemira = ?", activePemira.IdPemira).Count(&countPaslon)
	if countPaslon == 0 {
		paslons := []entity.Paslon{
			{
				IdPemira:    activePemira.IdPemira,
				NomorUrut:   1,
				NimKetua:    "105223001",
				NamaKetua:   "Nanami Kento",
				NimWakil:    "105223002",
				NamaWakil:   "Loid Forger",
				NamaKabinet: "Kabinet Sinergi Progresif",
				FotoPaslon:  "https://res.cloudinary.com/du9sbnbx9/image/upload/v1789247377/hmik/cp/profile/debpdp2ask5fnagcwcvt.png",
				Visi:        "Mewujudkan HMIK sebagai wadah kolaborasi inovatif, inklusif, dan berdaya saing global.",
				Misi:        "1. Meningkatkan kapabilitas teknis dan riset mahasiswa Ilmu Komputer.\n2. Memperluas jejaring profesional dengan industri dan alumni.\n3. Mengembangkan tata kelola himpunan yang transparan dan akuntabel.",
			},
			{
				IdPemira:    activePemira.IdPemira,
				NomorUrut:   2,
				NimKetua:    "105223003",
				NamaKetua:   "Gojo Satoru",
				NimWakil:    "105223004",
				NamaWakil:   "Suguru Geto",
				NamaKabinet: "Kabinet Harmoni Berkelanjutan",
				FotoPaslon:  "https://res.cloudinary.com/du9sbnbx9/image/upload/v1789147773/hmik/cp/profile/rpmvpb4uje8ekb0us1ru.png",
				Visi:        "HMIK Berdaya, Berprestasi, dan Berdampak Nyata bagi Almamater dan Bangsa.",
				Misi:        "1. Optimalisasi potensi akademik dan keprofesian mahasiswa.\n2. Pengabdian masyarakat berbasis teknologi tepat guna.\n3. Membangun ruang apresiasi dan kreativitas mahasiswa Ilkom.",
			},
		}

		for _, p := range paslons {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Gagal buat paslon %d: %v", p.NomorUrut, err)
			} else {
				fmt.Printf("✅ Berhasil input Paslon No %d: %s & %s (%s)\n", p.NomorUrut, p.NamaKetua, p.NamaWakil, p.NamaKabinet)
			}
		}
	} else {
		fmt.Printf("ℹ️ Paslon untuk Pemira ID %d sudah ada (%d paslon)\n", activePemira.IdPemira, countPaslon)
	}

	// 3. Masukkan 3 pemilih contoh (termasuk NIM mahasiswa yang ada di db pengguna)
	sampleVoters := []struct {
		Nim   string
		Email string
	}{
		{"105223046", "105223046@student.universitaspertamina.ac.id"},
		{"105223001", "105223001@student.universitaspertamina.ac.id"},
		{"105223049", "105223049@student.universitaspertamina.ac.id"},
	}

	for _, v := range sampleVoters {
		var existing entity.Pemilih
		if err := db.Where("id_pemira = ? AND nim = ?", activePemira.IdPemira, v.Nim).First(&existing).Error; err != nil {
			newVoter := entity.Pemilih{
				IdPemira:      activePemira.IdPemira,
				Nim:           v.Nim,
				EmailKampus:   v.Email,
				StatusMemilih: false,
			}
			db.Create(&newVoter)
			fmt.Printf("✅ Terdaftar Pemilih DPT Uji Coba: %s (%s)\n", v.Nim, v.Email)
		}
	}

	fmt.Println("\n🎉 Inisialisasi Fitur Pemira (Database & Data Awal) 100% Selesai!")
}
