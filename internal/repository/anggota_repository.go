package repository

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

// AnggotaRepository adalah interface untuk interaksi database anggota
type AnggotaRepository interface {
	FindAll() ([]entity.Anggota, error)
}

type anggotaRepository struct {
	db *gorm.DB
}

// NewAnggotaRepository menginisiasi repository anggota baru
func NewAnggotaRepository(db *gorm.DB) AnggotaRepository {
	return &anggotaRepository{db}
}

// FindAll mengambil semua anggota beserta relasinya (Keahlian & Minat)
func (r *anggotaRepository) FindAll() ([]entity.Anggota, error) {
	var anggotas []entity.Anggota
	
	// Menggunakan Preload untuk menarik data dari tabel relasi (Eager Loading)
	err := r.db.
		Preload("Pengguna").
		Preload("KeahlianAnggota.Keahlian").
		Preload("MinatRisetAnggota.MinatRiset").
		Find(&anggotas).Error
		
	return anggotas, err
}
