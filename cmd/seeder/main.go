package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database. Pastikan Docker sudah menyala di port 5433.\nError:", err)
	}

	fmt.Println("🚀 Database terkoneksi!")

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
		"migration/000028_add_is_unggulan_to_proker.up.sql",
		"migration/000029_create_divisi_departemen.up.sql",
		"migration/000030_update_program_kerja_add_divisi.up.sql",
		"migration/000031_update_anggota_add_divisi.up.sql",
	}

	fmt.Println("🏗️ Menjalankan skrip SQL untuk membangun tabel...")
	for _, file := range filesToMigrate {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Gagal membaca file %s: %v", file, err)
		}
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			log.Printf("Peringatan: Gagal mengeksekusi migrasi %s: %v", file, err)
		}
	}

	fmt.Println("🌱 Memasukkan data dummy ke database...")

	db.Exec("TRUNCATE TABLE divisi_departemen, departemen, program_kerja, pengguna, anggota, keahlian, minat_riset, keahlian_anggota, minat_riset_anggota RESTART IDENTITY CASCADE;")

	// Seed Departments and Divisions
	departments := []struct {
		Nama      string
		Deskripsi string
		Logo      string
		Divisions []string
		Prokers   []struct {
			Nama string
			Foto string
		}
	}{
		{
			Nama:      "Riset dan Teknologi",
			Deskripsi: "Departemen yang berfokus pada pengembangan kompetensi mahasiswa dalam bidang penelitian, teknologi, dan rekayasa perangkat lunak.",
			Logo:      "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365010/hmik/cp/departemen/hwwczwa85n3pqapfla60.png",
			Divisions: []string{"Riset & Inovasi", "Pengembangan Produk"},
			Prokers: []struct{Nama string; Foto string}{
				{Nama: "HMIK-CoreX", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Tech Talk", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Hackathon Internal", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
			},
		},
		{
			Nama:      "Media Informasi",
			Deskripsi: "Departemen yang mengelola publikasi, dokumentasi, dan penyebaran informasi kepada seluruh anggota HMIK serta pihak eksternal.",
			Logo:      "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365010/hmik/cp/departemen/hwwczwa85n3pqapfla60.png",
			Divisions: []string{"Media Kreatif", "Publikasi & Dokumentasi"},
			Prokers: []struct{Nama string; Foto string}{
				{Nama: "HMIK Magazine", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Podcast HMIK", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
			},
		},
		{
			Nama:      "Hubungan External",
			Deskripsi: "Departemen yang bertugas membangun dan menjaga hubungan baik dengan pihak di luar himpunan, seperti organisasi mahasiswa lain dan lembaga profesional.",
			Logo:      "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365010/hmik/cp/departemen/hwwczwa85n3pqapfla60.png",
			Divisions: []string{"Kemitraan Strategis", "Pengabdian Masyarakat"},
			Prokers: []struct{Nama string; Foto string}{
				{Nama: "Studi Banding", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Sosial Mengajar", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Alumni Gathering", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
			},
		},
		{
			Nama:      "Manajemen Internal",
			Deskripsi: "Departemen yang mengatur urusan rumah tangga himpunan, memastikan kesekretariatan dan inventaris berjalan dengan baik.",
			Logo:      "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365010/hmik/cp/departemen/hwwczwa85n3pqapfla60.png",
			Divisions: []string{"MSDA", "Aset & Perlengkapan", "Minat & Bakat"},
			Prokers: []struct{Nama string; Foto string}{
				{Nama: "Upgrading Pengurus", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Makrab HMIK", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
			},
		},
		{
			Nama:      "Academy UP",
			Deskripsi: "Departemen yang bertanggung jawab atas pengembangan akademik, keprofesian, serta menampung minat dan bakat mahasiswa.",
			Logo:      "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365010/hmik/cp/departemen/hwwczwa85n3pqapfla60.png",
			Divisions: []string{"Akademik", "Keprofesian", "Minat & Bakat"},
			Prokers: []struct{Nama string; Foto string}{
				{Nama: "Study Club", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Seminar Nasional", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
				{Nama: "Workshop IT", Foto: "https://res.cloudinary.com/du9sbnbx9/image/upload/v1791365011/hmik/cp/proker/aibxheziqovxrvtwq0dj.png"},
			},
		},
	}

	for _, d := range departments {
		dept := entity.Departemen{
			NamaDepartemen: d.Nama,
			Deskripsi:      d.Deskripsi,
			Logo:           d.Logo,
		}
		db.Create(&dept)

		for _, divName := range d.Divisions {
			div := entity.DivisiDepartemen{
				IdDepartemen: dept.IdDepartemen,
				NamaDivisi:   divName,
				Deskripsi:    "Divisi " + divName + " pada departemen " + d.Nama,
			}
			db.Create(&div)
		}

		for _, p := range d.Prokers {
			proker := entity.ProgramKerja{
				IdDepartemen: dept.IdDepartemen,
				NamaProker:   p.Nama,
				Deskripsi:    "Program kerja unggulan: " + p.Nama,
				Foto:         p.Foto,
				IsUnggulan:   true,
				Status:       "Mendatang",
			}
			db.Create(&proker)
		}
	}

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
		hashedSandi, _ := bcrypt.GenerateFromPassword([]byte("HMIK"+data.Nim), bcrypt.DefaultCost)
		pengguna := entity.Pengguna{
			NamaLengkap:  data.NamaLengkap,
			Email:        data.Email,
			Sandi:        string(hashedSandi),
			FotoPengguna: data.Foto,
		}
		
		if pengguna.FotoPengguna == "" {
			pengguna.FotoPengguna = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + data.NamaLengkap
		}
		db.Create(&pengguna)
		
		var idDepartemen *int64
		var idDivisi *int64

		if data.Jabatan != "" && data.Jabatan != "Ketua Himpunan" && data.Jabatan != "Wakil Ketua Himpunan" && data.Jabatan != "Sekretaris 1" && data.Jabatan != "Sekretaris 2" && data.Jabatan != "Bendahara 1" && data.Jabatan != "Bendahara 2" {
			jabatanLower := strings.ToLower(data.Jabatan)
			var deptName string
			
			if strings.Contains(jabatanLower, "riset dan teknologi") || strings.Contains(jabatanLower, "riset & inovasi") || strings.Contains(jabatanLower, "pengembangan produk") {
				deptName = "Riset dan Teknologi"
			} else if strings.Contains(jabatanLower, "media informasi") || strings.Contains(jabatanLower, "media kreatif") || strings.Contains(jabatanLower, "publikasi & dokumentasi") {
				deptName = "Media Informasi"
			} else if strings.Contains(jabatanLower, "hubungan external") || strings.Contains(jabatanLower, "kemitraan strategis") || strings.Contains(jabatanLower, "pengabdian masyarakat") {
				deptName = "Hubungan External"
			} else if strings.Contains(jabatanLower, "manajemen internal") || strings.Contains(jabatanLower, "msda") || strings.Contains(jabatanLower, "aset & perlengkapan") || strings.Contains(jabatanLower, "minat & bakat") {
				deptName = "Manajemen Internal"
			} else if strings.Contains(jabatanLower, "academy up") || strings.Contains(jabatanLower, "akademik") || strings.Contains(jabatanLower, "keprofesian") {
				deptName = "Academy UP"
			}

			if deptName != "" {
				var dept entity.Departemen
				db.Where("nama_departemen = ?", deptName).First(&dept)
				if dept.IdDepartemen != 0 {
					idDepartemen = &dept.IdDepartemen
				}
			}

			var divName string
			if strings.Contains(jabatanLower, "media kreatif") { divName = "Media Kreatif" }
			if strings.Contains(jabatanLower, "publikasi & dokumentasi") { divName = "Publikasi & Dokumentasi" }
			if strings.Contains(jabatanLower, "riset & inovasi") { divName = "Riset & Inovasi" }
			if strings.Contains(jabatanLower, "pengembangan produk") { divName = "Pengembangan Produk" }
			if strings.Contains(jabatanLower, "kemitraan strategis") { divName = "Kemitraan Strategis" }
			if strings.Contains(jabatanLower, "pengabdian masyarakat") { divName = "Pengabdian Masyarakat" }
			if strings.Contains(jabatanLower, "msda") { divName = "MSDA" }
			if strings.Contains(jabatanLower, "aset & perlengkapan") { divName = "Aset & Perlengkapan" }
			if strings.Contains(jabatanLower, "akademik") { divName = "Akademik" }
			if strings.Contains(jabatanLower, "keprofesian") { divName = "Keprofesian" }
			if strings.Contains(jabatanLower, "minat & bakat") { divName = "Minat & Bakat" }

			if divName != "" && idDepartemen != nil {
				var div entity.DivisiDepartemen
				db.Where("nama_divisi = ? AND id_departemen = ?", divName, *idDepartemen).First(&div)
				if div.IdDivisi != 0 {
					idDivisi = &div.IdDivisi
				}
			}
		}

		anggota := entity.Anggota{
			IdUser:       pengguna.IdPengguna,
			Nim:          data.Nim,
			ProgramStudi: data.ProgramStudi,
			Angkatan:     data.Angkatan,
			Jabatan:      data.Jabatan,
			IdDepartemen: idDepartemen,
			IdDivisi:     idDivisi,
		}
		db.Create(&anggota)

		for _, namaKeahlian := range data.Keahlian {
			var k entity.Keahlian
			db.Where(entity.Keahlian{NamaKeahlian: namaKeahlian}).FirstOrCreate(&k, entity.Keahlian{NamaKeahlian: namaKeahlian})
			db.Create(&entity.KeahlianAnggota{
				IdAnggota:         anggota.IdAnggota,
				IdKeahlian:        k.IdKeahlian,
				TingkatPenguasaan: 3,
			})
		}
	}

	fmt.Println("✅ Proses Seeding berhasil! Database siap digunakan dengan Divisi dan Proker terhubung.")
}
