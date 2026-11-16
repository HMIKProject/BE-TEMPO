package repository

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

type ProgramKerjaRepository interface {
	FindAll(status string, isUnggulan *bool) ([]entity.ProgramKerja, error)
	FindByID(id int64) (*entity.ProgramKerja, error)
	Create(proker *entity.ProgramKerja) error
	Update(proker *entity.ProgramKerja) error
	Delete(id int64) error
}

type programKerjaRepository struct {
	db *gorm.DB
}

func NewProgramKerjaRepository(db *gorm.DB) ProgramKerjaRepository {
	return &programKerjaRepository{db}
}

func (r *programKerjaRepository) FindAll(status string, isUnggulan *bool) ([]entity.ProgramKerja, error) {
	var prokers []entity.ProgramKerja
	query := r.db.Model(&entity.ProgramKerja{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if isUnggulan != nil {
		query = query.Where("is_unggulan = ?", *isUnggulan)
	}

	err := query.Order("id_proker ASC").Find(&prokers).Error
	return prokers, err
}

func (r *programKerjaRepository) FindByID(id int64) (*entity.ProgramKerja, error) {
	var proker entity.ProgramKerja
	if err := r.db.First(&proker, "id_proker = ?", id).Error; err != nil {
		return nil, err
	}
	return &proker, nil
}

func (r *programKerjaRepository) Create(proker *entity.ProgramKerja) error {
	return r.db.Create(proker).Error
}

func (r *programKerjaRepository) Update(proker *entity.ProgramKerja) error {
	return r.db.Save(proker).Error
}

func (r *programKerjaRepository) Delete(id int64) error {
	return r.db.Delete(&entity.ProgramKerja{}, "id_proker = ?", id).Error
}
