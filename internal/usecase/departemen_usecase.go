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
	depts, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Sanitasi kredensial pengguna agar tidak bocor
	for i := range depts {
		for j := range depts[i].Anggota {
			if depts[i].Anggota[j].Pengguna != nil {
				depts[i].Anggota[j].Pengguna.Sandi = ""
				depts[i].Anggota[j].Pengguna.NomorTelepon = ""
			}
		}
	}

	return depts, nil
}
