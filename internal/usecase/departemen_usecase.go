package usecase

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

type DepartemenUsecase interface {
	GetAllDepartemen() ([]entity.Departemen, error)
}

type departemenUsecase struct {
	repo repository.DepartemenRepository
}

func NewDepartemenUsecase(repo repository.DepartemenRepository) DepartemenUsecase {
	return &departemenUsecase{repo}
}

func (u *departemenUsecase) GetAllDepartemen() ([]entity.Departemen, error) {
	return u.repo.FindAll()
}
