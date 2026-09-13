package usecase

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

type ProgramKerjaUsecase interface {
	GetAllProgramKerja(status string, isUnggulan *bool) ([]entity.ProgramKerja, error)
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
