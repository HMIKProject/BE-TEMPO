package repository

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

type DepartemenRepository interface {
	FindAll() ([]entity.Departemen, error)
}

type departemenRepository struct {
	db *gorm.DB
}

func NewDepartemenRepository(db *gorm.DB) DepartemenRepository {
	return &departemenRepository{db}
}

func (r *departemenRepository) FindAll() ([]entity.Departemen, error) {
	var depts []entity.Departemen
	err := r.db.
		Preload("ProgramKerja").
		Preload("Anggota").
		Preload("Anggota.Pengguna").
		Find(&depts).Error
	return depts, err
}
