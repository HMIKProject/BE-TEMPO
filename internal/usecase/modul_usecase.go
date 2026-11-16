package usecase

import (
	"errors"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

type ModulUsecase interface {
	GetAll(filter entity.ModulFilterDTO) ([]entity.ModulPembelajaran, int64, error)
	GetByID(id int64) (*entity.ModulPembelajaran, error)
	Create(modul *entity.ModulPembelajaran) error
	Update(id int64, req *entity.ModulPembelajaran) error
	Delete(id int64) error
	RenameFolder(semester, oldMataKuliah, newMataKuliah string) error
}

type modulUsecase struct {
	modulRepo repository.ModulRepository
}

func NewModulUsecase(modulRepo repository.ModulRepository) ModulUsecase {
	return &modulUsecase{modulRepo: modulRepo}
}

func (u *modulUsecase) GetAll(filter entity.ModulFilterDTO) ([]entity.ModulPembelajaran, int64, error) {
	return u.modulRepo.GetAll(filter)
}

func (u *modulUsecase) GetByID(id int64) (*entity.ModulPembelajaran, error) {
	if id <= 0 {
		return nil, errors.New("id modul tidak valid")
	}
	return u.modulRepo.GetByID(id)
}

func (u *modulUsecase) Create(modul *entity.ModulPembelajaran) error {
	if strings.TrimSpace(modul.Judul) == "" {
		return errors.New("judul modul wajib diisi")
	}
	if strings.TrimSpace(modul.FileURL) == "" {
		return errors.New("file URL modul wajib ada")
	}

	// Batasan Maksimal 50 MB sesuai Figma dan BP-03 (50 * 1024 * 1024 bytes)
	const maxBytes = 50 * 1024 * 1024
	if modul.UkuranBytes > maxBytes {
		return errors.New("ukuran file melebihi batas maksimal 50 MB")
	}

	if modul.Semester == "" {
		modul.Semester = "Semester 1"
	}
	if modul.MataKuliah == "" {
		modul.MataKuliah = "Umum"
	}

	// Normalisasi TipeFile
	modul.TipeFile = strings.ToLower(strings.TrimPrefix(modul.TipeFile, "."))

	return u.modulRepo.Create(modul)
}

func (u *modulUsecase) Update(id int64, req *entity.ModulPembelajaran) error {
	existing, err := u.modulRepo.GetByID(id)
	if err != nil {
		return errors.New("modul tidak ditemukan")
	}

	if strings.TrimSpace(req.Judul) != "" {
		existing.Judul = req.Judul
	}
	if req.Deskripsi != "" {
		existing.Deskripsi = req.Deskripsi
	}
	if req.Kategori != "" {
		existing.Kategori = req.Kategori
	}
	if req.Semester != "" {
		existing.Semester = req.Semester
	}
	if req.MataKuliah != "" {
		existing.MataKuliah = req.MataKuliah
	}
	if req.FileURL != "" {
		existing.FileURL = req.FileURL
	}
	if req.TipeFile != "" {
		existing.TipeFile = strings.ToLower(strings.TrimPrefix(req.TipeFile, "."))
	}
	if req.UkuranFile != "" {
		existing.UkuranFile = req.UkuranFile
	}
	if req.UkuranBytes > 0 {
		existing.UkuranBytes = req.UkuranBytes
	}

	return u.modulRepo.Update(existing)
}

func (u *modulUsecase) Delete(id int64) error {
	_, err := u.modulRepo.GetByID(id)
	if err != nil {
		return errors.New("modul tidak ditemukan")
	}
	return u.modulRepo.Delete(id)
}

func (u *modulUsecase) RenameFolder(semester, oldMataKuliah, newMataKuliah string) error {
	sem := strings.TrimSpace(semester)
	oldM := strings.TrimSpace(oldMataKuliah)
	newM := strings.TrimSpace(newMataKuliah)
	if sem == "" || oldM == "" || newM == "" {
		return errors.New("semester, nama folder lama, dan nama folder baru wajib diisi")
	}
	return u.modulRepo.RenameFolder(sem, oldM, newM)
}
