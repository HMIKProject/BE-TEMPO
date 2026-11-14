package repository

import (
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

type ModulRepository interface {
	GetAll(filter entity.ModulFilterDTO) ([]entity.ModulPembelajaran, int64, error)
	GetByID(id int64) (*entity.ModulPembelajaran, error)
	Create(modul *entity.ModulPembelajaran) error
	Update(modul *entity.ModulPembelajaran) error
	Delete(id int64) error
	CountAll() (int64, error)
}

type modulRepository struct {
	db *gorm.DB
}

func NewModulRepository(db *gorm.DB) ModulRepository {
	return &modulRepository{db: db}
}

func (r *modulRepository) GetAll(filter entity.ModulFilterDTO) ([]entity.ModulPembelajaran, int64, error) {
	var list []entity.ModulPembelajaran
	var total int64

	query := r.db.Model(&entity.ModulPembelajaran{})

	if filter.Search != "" {
		s := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(judul) LIKE ? OR LOWER(mata_kuliah) LIKE ? OR LOWER(deskripsi) LIKE ?", s, s, s)
	}

	if filter.Kategori != "" && filter.Kategori != "Semua" {
		query = query.Where("kategori = ?", filter.Kategori)
	}

	if filter.Semester != "" && filter.Semester != "Semua" {
		query = query.Where("semester = ?", filter.Semester)
	}

	if filter.MataKuliah != "" && filter.MataKuliah != "Semua" {
		query = query.Where("mata_kuliah = ?", filter.MataKuliah)
	}

	if filter.TipeFile != "" && filter.TipeFile != "Semua" {
		query = query.Where("LOWER(tipe_file) = ?", strings.ToLower(filter.TipeFile))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 50
	}
	offset := (page - 1) * limit

	if err := query.Order("dibuat_pada DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *modulRepository) GetByID(id int64) (*entity.ModulPembelajaran, error) {
	var modul entity.ModulPembelajaran
	if err := r.db.First(&modul, "id_modul = ?", id).Error; err != nil {
		return nil, err
	}
	return &modul, nil
}

func (r *modulRepository) Create(modul *entity.ModulPembelajaran) error {
	return r.db.Create(modul).Error
}

func (r *modulRepository) Update(modul *entity.ModulPembelajaran) error {
	return r.db.Save(modul).Error
}

func (r *modulRepository) Delete(id int64) error {
	return r.db.Delete(&entity.ModulPembelajaran{}, "id_modul = ?", id).Error
}

func (r *modulRepository) CountAll() (int64, error) {
	var total int64
	err := r.db.Model(&entity.ModulPembelajaran{}).Count(&total).Error
	return total, err
}
