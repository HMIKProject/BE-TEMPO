package usecase

import (
	"errors"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

type ProgramKerjaUsecase interface {
	GetAllProgramKerja(status string, isUnggulan *bool) ([]entity.ProgramKerja, error)
	GetProgramKerjaByID(id int64) (*entity.ProgramKerja, error)
	CreateProgramKerja(proker *entity.ProgramKerja) error
	UpdateProgramKerja(id int64, req *entity.ProgramKerja) error
	DeleteProgramKerja(id int64) error
}

type programKerjaUsecase struct {
	repo repository.ProgramKerjaRepository
}

func NewProgramKerjaUsecase(repo repository.ProgramKerjaRepository) ProgramKerjaUsecase {
	return &programKerjaUsecase{repo}
}

func (u *programKerjaUsecase) GetAllProgramKerja(status string, isUnggulan *bool) ([]entity.ProgramKerja, error) {
	return u.repo.FindAll(status, isUnggulan)
}

func (u *programKerjaUsecase) GetProgramKerjaByID(id int64) (*entity.ProgramKerja, error) {
	return u.repo.FindByID(id)
}

func (u *programKerjaUsecase) CreateProgramKerja(proker *entity.ProgramKerja) error {
	if proker.NamaProker == "" {
		return errors.New("nama program kerja wajib diisi")
	}
	if proker.IdDepartemen <= 0 {
		return errors.New("departemen penaung program kerja wajib dipilih")
	}
	if proker.Status == "" {
		proker.Status = "Mendatang"
	}
	if proker.IdDivisi != nil && *proker.IdDivisi <= 0 {
		proker.IdDivisi = nil
	}
	return u.repo.Create(proker)
}

func (u *programKerjaUsecase) UpdateProgramKerja(id int64, req *entity.ProgramKerja) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("program kerja tidak ditemukan")
	}

	if req.NamaProker != "" {
		existing.NamaProker = req.NamaProker
	}
	if req.IdDepartemen > 0 {
		existing.IdDepartemen = req.IdDepartemen
	}
	if req.IdDivisi != nil && *req.IdDivisi <= 0 {
		existing.IdDivisi = nil
	} else {
		existing.IdDivisi = req.IdDivisi
	}
	existing.Deskripsi = req.Deskripsi
	existing.Foto = req.Foto
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.IsUnggulan = req.IsUnggulan

	return u.repo.Update(existing)
}

func (u *programKerjaUsecase) DeleteProgramKerja(id int64) error {
	_, err := u.repo.FindByID(id)
	if err != nil {
		return errors.New("program kerja tidak ditemukan")
	}
	return u.repo.Delete(id)
}
