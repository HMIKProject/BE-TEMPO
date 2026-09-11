package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// Memuat file .env jika ada
	godotenv.Load()

	// 1. Koneksi ke Database menggunakan port 5433 (seperti instruksi ke user) atau DATABASE_URL dari Cloud
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Wajib agar tidak error plural (e.g. anggotas)
		},
	})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database. Pastikan Docker sudah menyala di port 5433.\nError:", err)
	}

	fmt.Println("🚀 Database terkoneksi!")

	// 2. Eksekusi File SQL (Migrasi DDL Manual) karena dilarang pakai AutoMigrate
	filesToMigrate := []string{
		"migration/000000_init.up.sql",
		"migration/000002_roles.up.sql",
		"migration/000001_users.up.sql",
		"migration/000005_anggota.up.sql",
		"migration/000009_skills.up.sql",
		"migration/000008_anggota_skills.up.sql",
		"migration/000007.research.interests.up.sql",
		"migration/000006_anggota_research_interests.up.sql",
		"migration/000023_add_jabatan_to_anggota.up.sql",
	}

	fmt.Println("🏗️ Menjalankan skrip SQL untuk membangun tabel...")
	for _, file := range filesToMigrate {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Gagal membaca file %s: %v", file, err)
		}
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			log.Fatalf("Gagal mengeksekusi migrasi %s: %v", file, err)
		}
	}

	// 3. Menjalankan Seeder (Insert Data Dummy dengan GORM)
	fmt.Println("🌱 Memasukkan data dummy ke database...")

	// Hapus isi data sebelumnya agar tidak conflict (Opsional, untuk seeder lokal)
	db.Exec("TRUNCATE TABLE pengguna, anggota, keahlian, minat_riset, keahlian_anggota, minat_riset_anggota RESTART IDENTITY CASCADE;")

	// --- Seeder Keahlian ---
	keahlianGolang := entity.Keahlian{NamaKeahlian: "Golang", Deskripsi: "Bahasa Pemrograman Backend"}
	keahlianReact := entity.Keahlian{NamaKeahlian: "React JS", Deskripsi: "Library Frontend Modern"}
	db.Create(&keahlianGolang)
	db.Create(&keahlianReact)

	// --- Seeder Minat Riset ---
	risetAI := entity.MinatRiset{NamaMinat: "Artificial Intelligence", Deskripsi: "AI & Machine Learning"}
	risetWeb := entity.MinatRiset{NamaMinat: "Web Development", Deskripsi: "Pengembangan Website"}
	db.Create(&risetAI)
	db.Create(&risetWeb)

	// --- Seeder Pengguna & Anggota dari JSON ---
	type AnggotaData struct {
		Nim          string `json:"nim"`
		NamaLengkap  string `json:"nama_lengkap"`
		Email        string `json:"email"`
		ProgramStudi string `json:"program_studi"`
		Angkatan     int    `json:"angkatan"`
		Jabatan      string `json:"jabatan"`
		Foto         string `json:"foto"`
	}

	jsonFile, err := os.ReadFile("cmd/seeder/data_anggota.json")
	if err != nil {
		log.Fatalf("Gagal membaca file JSON data anggota: %v", err)
	}

	var dataAnggota []AnggotaData
	if err := json.Unmarshal(jsonFile, &dataAnggota); err != nil {
		log.Fatalf("Gagal parse JSON data anggota: %v", err)
	}

	for _, data := range dataAnggota {
		pengguna := entity.Pengguna{
			NamaLengkap:  data.NamaLengkap,
			Email:        data.Email,
			Sandi:        "HMIK" + data.Nim,
			FotoPengguna: data.Foto,
		}
		
		if pengguna.FotoPengguna == "" {
			pengguna.FotoPengguna = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + data.NamaLengkap
		}
		
		db.Create(&pengguna)

		anggota := entity.Anggota{
			IdUser:       pengguna.IdPengguna,
			Nim:          data.Nim,
			ProgramStudi: data.ProgramStudi,
			Angkatan:     data.Angkatan,
			Jabatan:      data.Jabatan,
		}
		db.Create(&anggota)

		// Dummy Keahlian dan Minat Riset agar profile tidak kosong
		db.Create(&entity.KeahlianAnggota{IdAnggota: anggota.IdAnggota, IdKeahlian: keahlianGolang.IdKeahlian, TingkatPenguasaan: 3})
		db.Create(&entity.MinatRisetAnggota{IdAnggota: anggota.IdAnggota, IdMinat: risetWeb.IdMinat})
	}

	fmt.Println("✅ Proses Seeding berhasil! Database siap digunakan.")
}
