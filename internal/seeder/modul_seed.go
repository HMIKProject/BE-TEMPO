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
			Judul:       "Kalkulus 1 - Diktat Lengkap & Latihan Soal UAS",
			Deskripsi:   "Diktat komprehensif Kalkulus 1 mencakup limit, turunan, integral fungsi, serta arsip review soal ujian.",
			Semester:    "Semester 1",
			MataKuliah:  "Kalkulus 1",
			TipeFile:    "pdf",
			UkuranFile:  "7.1 MB",
			UkuranBytes: 7444889,
			FileURL:     "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654581/hmik/learning/Semester%201/Kalkulus%201/Kalkulus_1.pdf",
			DibuatOleh:  "Admin HMIK",
		},
		{
			Judul:       "Matematika Diskrit - Logika Proposisi & Teori Graf",
			Deskripsi:   "Slide materi logika matematika, himpunan, relasi fungsi, kombinatorika, dan teori graf terapan.",
			Semester:    "Semester 1",
			MataKuliah:  "Matematika Diskrit",
			TipeFile:    "pdf",
			UkuranFile:  "3.8 MB",
			UkuranBytes: 3984588,
			FileURL:     "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654588/hmik/learning/Semester%201/Matematika%20Diskrit/10-13.TeoriHimpunanElementer.pdf",
			DibuatOleh:  "Admin HMIK",
		},
		{
			Judul:       "Organisasi Komputer - Arsitektur CPU & MIPS",
			Deskripsi:   "Materi arsitektur sistem komputer, representasi biner, gerbang logika, dan register instruksi MIPS.",
			Semester:    "Semester 1",
			MataKuliah:  "Organisasi Komputer",
			TipeFile:    "pptx",
			UkuranFile:  "12.4 MB",
			UkuranBytes: 13002342,
			FileURL:     "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654637/hmik/learning/Semester%201/Organisasi%20Komputer/Pertemuan%201%20-%20Pengantar%20Organisasi%20Komputer.pdf",
			DibuatOleh:  "Admin HMIK",
		},
		{
			Judul:       "Dasar Pemrograman - Pointer, Array & Memory Management",
			Deskripsi:   "Modul praktikum dasar algoritma, kontrol alur, penanganan memori dan pointer bahasa C/C++.",
			Semester:    "Semester 2",
			MataKuliah:  "Dasar Pemrograman",
			TipeFile:    "pdf",
			UkuranFile:  "4.5 MB",
			UkuranBytes: 4718592,
			FileURL:     "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654662/hmik/learning/Semester%202/Dasar%20Pemrograman/UAS_DasPro_2024.pdf",
			DibuatOleh:  "Admin HMIK",
		},
		{
			Judul:       "Algoritma dan Struktur Data - Tree, BST & Graph Traversal",
			Deskripsi:   "Panduan implementasi struktur data non-linear Tree, Binary Search Tree, Graph, dan BFS/DFS traversal.",
			Semester:    "Semester 3",
			MataKuliah:  "Algoritma dan Struktur Data",
			TipeFile:    "pdf",
			UkuranFile:  "3.1 MB",
			UkuranBytes: 3250585,
			FileURL:     "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654669/hmik/learning/Semester%203/Algoritma%20dan%20Struktur%20Data/W1_-_Array_and_Struct.pdf",
			DibuatOleh:  "Admin HMIK",
		},
		{
			Judul:       "PBO - Design Patterns, Encapsulation & Polymorphism Java",
			Deskripsi:   "Slide kuliah konsep OOP, pewarisan sifat, interface polymorphism, dan penanganan exception.",
			Semester:    "Semester 4",
			MataKuliah:  "PBO (Pemrograman Berorientasi Objek)",
			TipeFile:    "pdf",
			UkuranFile:  "6.2 MB",
			UkuranBytes: 6501171,
			FileURL:     "https://res.cloudinary.com/du9sbnbx9/image/upload/v1794654795/hmik/learning/Semester%204/PBO%20%28Pemrograman%20Berorientasi%20Objek%29/11-Polymorphism.pdf",
			DibuatOleh:  "Admin HMIK",
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
