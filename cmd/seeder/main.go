package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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
		"migration/000024_departemen.up.sql",
		"migration/000025_program_kerja.up.sql",
		"migration/000026_add_departemen_to_anggota.up.sql",
		"migration/000027_add_foto_to_proker.up.sql",
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
	db.Exec("TRUNCATE TABLE departemen, program_kerja, pengguna, anggota, keahlian, minat_riset, keahlian_anggota, minat_riset_anggota RESTART IDENTITY CASCADE;")

	// --- Seeder Pengguna & Anggota dari JSON ---
	type AnggotaData struct {
		Nim          string   `json:"nim"`
		NamaLengkap  string   `json:"nama_lengkap"`
		Email        string   `json:"email"`
		ProgramStudi string   `json:"program_studi"`
		Angkatan     int      `json:"angkatan"`
		Jabatan      string   `json:"jabatan"`
		Foto         string   `json:"foto"`
		Keahlian     []string `json:"keahlian"`
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
		
		var idDepartemen *int64
		if data.Jabatan != "" && data.Jabatan != "Ketua Himpunan" && data.Jabatan != "Wakil Ketua Himpunan" && data.Jabatan != "Sekretaris 1" && data.Jabatan != "Sekretaris 2" && data.Jabatan != "Bendahara 1" && data.Jabatan != "Bendahara 2" {
			namaDept := data.Jabatan
			if strings.Contains(data.Jabatan, "Departemen ") {
				namaDept = strings.SplitN(data.Jabatan, "Departemen ", 2)[1]
			} else if strings.Contains(data.Jabatan, "Divisi ") {
				namaDept = strings.SplitN(data.Jabatan, "Divisi ", 2)[1]
			} else if strings.Contains(data.Jabatan, "Kadiv ") {
				namaDept = strings.SplitN(data.Jabatan, "Kadiv ", 2)[1]
			}
			
			var dept entity.Departemen
			res := db.Where(entity.Departemen{NamaDepartemen: namaDept}).FirstOrCreate(&dept, entity.Departemen{NamaDepartemen: namaDept, Deskripsi: "Deskripsi resmi departemen " + namaDept, Logo: ""})
			if res.RowsAffected > 0 { 
				db.Create(&entity.ProgramKerja{IdDepartemen: dept.IdDepartemen, NamaProker: "Proker Unggulan " + namaDept, Deskripsi: "Contoh deskripsi proker otomatis.", Foto: ""})
			}
			idDepartemen = &dept.IdDepartemen
		}

		anggota := entity.Anggota{
			IdUser:       pengguna.IdPengguna,
			Nim:          data.Nim,
			ProgramStudi: data.ProgramStudi,
			Angkatan:     data.Angkatan,
			Jabatan:      data.Jabatan,
			IdDepartemen: idDepartemen,
		}
		db.Create(&anggota)

		// Memasukkan Keahlian Dinamis dari Array JSON
		for _, namaKeahlian := range data.Keahlian {
			var k entity.Keahlian
			// Cari keahlian, jika tidak ada maka buat baru
			db.Where(entity.Keahlian{NamaKeahlian: namaKeahlian}).FirstOrCreate(&k, entity.Keahlian{NamaKeahlian: namaKeahlian})
			
			// Buat jembatan relasi many-to-many
			db.Create(&entity.KeahlianAnggota{
				IdAnggota:         anggota.IdAnggota,
				IdKeahlian:        k.IdKeahlian,
				TingkatPenguasaan: 3,
			})
		}
	}

	fmt.Println("✅ Proses Seeding berhasil! Database siap digunakan.")
}
