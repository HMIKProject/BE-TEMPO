package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// 1. Koneksi ke Database menggunakan port 5433 (seperti instruksi ke user)
	dsn := "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable"
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

	// --- Seeder Pengguna & Anggota ---
	// Data 1
	pengguna1 := entity.Pengguna{
		NamaLengkap:  "Damos Santoso",
		Email:        "damos@student.ac.id",
		Sandi:        "hashed_password_123",
		FotoPengguna: "https://api.dicebear.com/7.x/avataaars/svg?seed=Damos",
	}
	db.Create(&pengguna1)

	anggota1 := entity.Anggota{
		IdUser:       pengguna1.IdPengguna,
		Nim:          "1900018001",
		ProgramStudi: "Informatika",
		Angkatan:     2019,
	}
	db.Create(&anggota1)

	// Insert Pivot Anggota 1
	db.Create(&entity.KeahlianAnggota{IdAnggota: anggota1.IdAnggota, IdKeahlian: keahlianGolang.IdKeahlian, TingkatPenguasaan: 5})
	db.Create(&entity.KeahlianAnggota{IdAnggota: anggota1.IdAnggota, IdKeahlian: keahlianReact.IdKeahlian, TingkatPenguasaan: 4})
	db.Create(&entity.MinatRisetAnggota{IdAnggota: anggota1.IdAnggota, IdMinat: risetWeb.IdMinat})

	// Data 2
	pengguna2 := entity.Pengguna{
		NamaLengkap:  "Siti Aminah",
		Email:        "siti@student.ac.id",
		Sandi:        "hashed_password_456",
		FotoPengguna: "https://api.dicebear.com/7.x/avataaars/svg?seed=Siti",
	}
	db.Create(&pengguna2)

	anggota2 := entity.Anggota{
		IdUser:       pengguna2.IdPengguna,
		Nim:          "1900018002",
		ProgramStudi: "Sistem Informasi",
		Angkatan:     2019,
	}
	db.Create(&anggota2)

	// Insert Pivot Anggota 2
	db.Create(&entity.KeahlianAnggota{IdAnggota: anggota2.IdAnggota, IdKeahlian: keahlianReact.IdKeahlian, TingkatPenguasaan: 3})
	db.Create(&entity.MinatRisetAnggota{IdAnggota: anggota2.IdAnggota, IdMinat: risetAI.IdMinat})

	fmt.Println("✅ Proses Seeding berhasil! Database siap digunakan.")
}
