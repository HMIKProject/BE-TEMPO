package main

import (
	"log"
	"os"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/seeder"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	_ = godotenv.Load("../../.env")
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
		log.Fatal("Gagal koneksi database:", err)
	}

	log.Println("🧹 Memulai pembersihan data untuk persiapan Production...")

	// 1. Truncate tabel anggota dan relasinya
	if err := db.Exec("TRUNCATE TABLE anggota, keahlian_anggota, minat_riset_anggota, riwayat_jabatan CASCADE;").Error; err != nil {
		log.Printf("Peringatan truncate anggota: %v", err)
	} else {
		log.Println("✅ Tabel anggota & riwayat_jabatan berhasil dibersihkan (TRUNCATE)")
	}

	// 2. Bersihkan akun dummy mahasiswa (non-slot)
	if err := db.Exec("DELETE FROM pengguna WHERE email NOT LIKE 'hmik.%@hmik.com' AND email != 'hmik.up@gmail.com';").Error; err != nil {
		log.Printf("Gagal hapus akun non-slot: %v", err)
	} else {
		log.Println("✅ Akun non-slot (dummy mahasiswa) berhasil dihapus")
	}

	// 3. Reset slot akun struktural resmi ke state default (password: 'qwerqwer')
	if err := seeder.SeedAllRoleAccounts(db); err != nil {
		log.Fatalf("Gagal reset slot accounts: %v", err)
	}
	db.Model(&entity.Pengguna{}).Where("1 = 1").Updates(map[string]interface{}{
		"kode_aktivasi":         nil,
		"kode_aktivasi_expired": nil,
		"is_first_login":        false,
		"foto_pengguna":         "",
	})
	log.Println("✅ Slot akun resmi berhasil di-reset ke kondisi KOSONG")

	// 4. Reset data Pemira
	if err := db.Exec("TRUNCATE TABLE suara, pemilih, paslon CASCADE;").Error; err != nil {
		log.Printf("Peringatan truncate pemira: %v", err)
	}
	if err := db.Exec("DELETE FROM pemira;").Error; err != nil {
		log.Printf("Peringatan delete pemira: %v", err)
	}

	now := time.Now()
	cleanPemira := entity.Pemira{
		TahunPeriode:   "2026/2027",
		JudulPemilihan: "Pemilihan Ketua & Wakil Ketua HMIK",
		TanggalMulai:   now,
		TanggalSelesai: now.Add(7 * 24 * time.Hour),
		Status:         "DRAFT",
	}
	if err := db.Create(&cleanPemira).Error; err != nil {
		log.Printf("Gagal buat Pemira draft: %v", err)
	} else {
		log.Printf("✅ Pemira baru dibuat dalam status DRAFT (ID: %d)", cleanPemira.IdPemira)
	}

	// 5. Bersihkan data kabinet pengujian dan set kabinet standar
	db.Exec("DELETE FROM kabinet WHERE nama_kabinet ILIKE '%malam%';")
	db.Model(&entity.Kabinet{}).Where("1 = 1").Update("is_aktif", false)
	db.Model(&entity.Kabinet{}).Where("id_kabinet = ?", 1).Update("is_aktif", true)
	log.Println("✅ Data kabinet di-reset (Kabinet 1 aktif sebagai fallback)")

	// 6. Bersihkan divisi dummy jika ada
	db.Exec("DELETE FROM divisi_departemen WHERE nama_divisi = 'ea';")

	log.Println("🎉 Database telah bersih (Clean Slate) dan siap untuk pengujian Pemira & Input Data Mahasiswa Baru!")
}
