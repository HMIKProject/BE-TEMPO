package repository

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

type ProgramKerjaRepository interface {
	FindAll(status string, isUnggulan *bool) ([]entity.ProgramKerja, error)
}

type programKerjaRepository struct {
	db *gorm.DB
}

func NewProgramKerjaRepository(db *gorm.DB) ProgramKerjaRepository {
	return &programKerjaRepository{db}
}

func (r *programKerjaRepository) FindAll(status string, isUnggulan *bool) ([]entity.ProgramKerja, error) {
	var prokers []entity.ProgramKerja
	query := r.db

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if isUnggulan != nil {
		query = query.Where("is_unggulan = ?", *isUnggulan)
	}

	err := query.Find(&prokers).Error
	return prokers, err
}
