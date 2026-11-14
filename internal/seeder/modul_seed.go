package seeder

import (
	"log"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

// SeedSampleModul mengisi data sampel modul pembelajaran jika tabel masih kosong
func SeedSampleModul(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.ModulPembelajaran{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sampleModules := []entity.ModulPembelajaran{
		{
			Judul:       "Modul 1 - Pengantar Kotlin & OOP Dasar",
			Deskripsi:   "Materi pengantar dasar sintaksis bahasa Kotlin, variabel, control flow, dan object-oriented programming.",
			Kategori:    "Akademik",
			Semester:    "Semester 3",
			MataKuliah:  "Pemrograman Perangkat Bergerak",
			TipeFile:    "pdf",
			UkuranFile:  "4.2 MB",
			UkuranBytes: 4404019,
			FileURL:     "/uploads/modul/sample_modul1_kotlin.pdf",
			DibuatOleh:  "Divisi Akademik HMIK",
		},
		{
			Judul:       "Modul 2 - Layout, UI & Jetpack Compose",
			Deskripsi:   "Slide presentasi komprehensif tentang arsitektur declarative UI menggunakan Android Jetpack Compose.",
			Kategori:    "Akademik",
			Semester:    "Semester 3",
			MataKuliah:  "Pemrograman Perangkat Bergerak",
			TipeFile:    "pptx",
			UkuranFile:  "12.8 MB",
			UkuranBytes: 13421772,
			FileURL:     "/uploads/modul/sample_modul2_compose.pptx",
			DibuatOleh:  "Divisi Akademik HMIK",
		},
		{
			Judul:       "Latihan Soal & Starter Code: Firebase & Room DB",
			Deskripsi:   "Kumpulan proyek starter dan arsip soal latihan integrasi database lokal Room dan cloud storage Firebase.",
			Kategori:    "Akademik",
			Semester:    "Semester 3",
			MataKuliah:  "Pemrograman Perangkat Bergerak",
			TipeFile:    "zip",
			UkuranFile:  "8.5 MB",
			UkuranBytes: 8912896,
			FileURL:     "/uploads/modul/sample_starter_firebase.zip",
			DibuatOleh:  "Divisi Akademik HMIK",
		},
		{
			Judul:       "Modul Struktur Data: Implementasi Tree & Graph",
			Deskripsi:   "Panduan praktikum implementasi struktur data non-linear Tree, Binary Search Tree, Graph, dan BFS/DFS traversal.",
			Kategori:    "Akademik",
			Semester:    "Semester 2",
			MataKuliah:  "Algoritma & Struktur Data",
			TipeFile:    "pdf",
			UkuranFile:  "3.1 MB",
			UkuranBytes: 3250585,
			FileURL:     "/uploads/modul/sample_modul_struktur_data.pdf",
			DibuatOleh:  "Divisi Akademik HMIK",
		},
		{
			Judul:       "Panduan UI/UX & Design System Neobrutalism HMIK",
			Deskripsi:   "Dokumentasi panduan desain, komponen Neobrutalism, palet warna, tipografi, dan aset visual HMIK.",
			Kategori:    "Riset dan Teknologi",
			Semester:    "Semester 4",
			MataKuliah:  "Interaksi Manusia dan Komputer",
			TipeFile:    "pdf",
			UkuranFile:  "15.4 MB",
			UkuranBytes: 16148070,
			FileURL:     "/uploads/modul/sample_design_system_hmik.pdf",
			DibuatOleh:  "Divisi Medkre & Riset",
		},
	}

	for _, m := range sampleModules {
		if err := db.Create(&m).Error; err != nil {
			log.Printf("⚠️ Gagal seeding modul %s: %v", m.Judul, err)
		}
	}

	log.Printf("✅ [Seeder] Berhasil menginisialisasi %d sampel modul pembelajaran", len(sampleModules))
	return nil
}
