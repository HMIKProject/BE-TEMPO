package repository

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

// PenggunaRepository adalah interface untuk interaksi database tabel pengguna
type PenggunaRepository interface {
	FindByEmailOrUsername(identifier string) (*entity.Pengguna, error)
}

type penggunaRepository struct {
	db *gorm.DB
}

// NewPenggunaRepository menginisiasi repository pengguna baru
func NewPenggunaRepository(db *gorm.DB) PenggunaRepository {
	return &penggunaRepository{db}
}

// FindByEmailOrUsername mencari pengguna berdasarkan email atau username
func (r *penggunaRepository) FindByEmailOrUsername(identifier string) (*entity.Pengguna, error) {
	var pengguna entity.Pengguna
	
	// Mencari data pengguna yang memiliki email ATAU nama_pengguna yang cocok dengan input
	err := r.db.Where("email = ? OR nama_pengguna = ?", identifier, identifier).First(&pengguna).Error
	
	if err != nil {
		return nil, err
	}
	
	return &pengguna, nil
}
