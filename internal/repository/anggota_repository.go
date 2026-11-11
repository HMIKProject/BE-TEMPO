package repository

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

// AnggotaRepository adalah interface untuk interaksi database anggota
type AnggotaRepository interface {
	FindAll() ([]entity.Anggota, error)
	FindWithFilter(filter entity.AnggotaFilter) ([]entity.Anggota, int64, error)
	FindByID(id int64) (*entity.Anggota, error)
	FindByNIM(nim string) (*entity.Anggota, error)
	Create(anggota *entity.Anggota) error
	Update(anggota *entity.Anggota) error
	Delete(id int64) error
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
		Preload("Departemen").
		Preload("Divisi").
		Preload("KeahlianAnggota.Keahlian").
		Preload("MinatRisetAnggota.MinatRiset").
		Find(&anggotas).Error
		
	return anggotas, err
}

// FindWithFilter mengambil daftar anggota berdasarkan filter pencarian, status, angkatan, serta pagination
func (r *anggotaRepository) FindWithFilter(filter entity.AnggotaFilter) ([]entity.Anggota, int64, error) {
	var anggotas []entity.Anggota
	var totalData int64

	query := r.db.Model(&entity.Anggota{})

	// Filter Search (Nama Lengkap atau NIM)
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("nama_lengkap ILIKE ? OR nim ILIKE ?", searchPattern, searchPattern)
	}

	// Filter Status
	if filter.Status != "" && filter.Status != "Semua" {
		query = query.Where("status = ?", filter.Status)
	}

	// Filter Angkatan
	if filter.Angkatan > 0 {
		query = query.Where("angkatan = ?", filter.Angkatan)
	}

	// Hitung total sebelum limit/offset
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Pagination setup
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 50
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Pengguna").
		Preload("Departemen").
		Preload("Divisi").
		Preload("KeahlianAnggota.Keahlian").
		Preload("MinatRisetAnggota.MinatRiset").
		Order("angkatan DESC, nama_lengkap ASC").
		Offset(offset).
		Limit(limit).
		Find(&anggotas).Error

	return anggotas, totalData, err
}

// FindByID mencari anggota berdasarkan ID spesifik
func (r *anggotaRepository) FindByID(id int64) (*entity.Anggota, error) {
	var anggota entity.Anggota
	err := r.db.
		Preload("Pengguna").
		Preload("Departemen").
		Preload("Divisi").
		Preload("KeahlianAnggota.Keahlian").
		Preload("MinatRisetAnggota.MinatRiset").
		First(&anggota, id).Error
	if err != nil {
		return nil, err
	}
	return &anggota, nil
}

// FindByNIM mencari anggota berdasarkan NIM unik
func (r *anggotaRepository) FindByNIM(nim string) (*entity.Anggota, error) {
	var anggota entity.Anggota
	err := r.db.Where("nim = ?", nim).First(&anggota).Error
	if err != nil {
		return nil, err
	}
	return &anggota, nil
}

// Create menyimpan anggota baru ke database
func (r *anggotaRepository) Create(anggota *entity.Anggota) error {
	return r.db.Create(anggota).Error
}

// Update memperbarui data profil anggota
func (r *anggotaRepository) Update(anggota *entity.Anggota) error {
	return r.db.Save(anggota).Error
}

// Delete menghapus entri anggota
func (r *anggotaRepository) Delete(id int64) error {
	return r.db.Delete(&entity.Anggota{}, id).Error
}

