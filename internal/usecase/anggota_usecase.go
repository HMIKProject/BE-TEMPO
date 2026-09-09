package usecase

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

// AnggotaUsecase adalah interface untuk logika bisnis anggota
type AnggotaUsecase interface {
	GetAllAnggota() ([]entity.Anggota, error)
}

type anggotaUsecase struct {
	repo repository.AnggotaRepository
}

// NewAnggotaUsecase menginisiasi usecase anggota baru
func NewAnggotaUsecase(repo repository.AnggotaRepository) AnggotaUsecase {
	return &anggotaUsecase{repo}
}

// GetAllAnggota memproses permintaan untuk mendapatkan daftar anggota
func (u *anggotaUsecase) GetAllAnggota() ([]entity.Anggota, error) {
	// Di sini Koki bisa menambahkan logika bisnis tambahan jika diperlukan
	// Contoh: memanipulasi URL foto, menyembunyikan password, dll.
	
	anggotas, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	
	// Mengosongkan password agar tidak bocor ke response JSON
	for i := range anggotas {
		anggotas[i].Pengguna.Sandi = ""
	}
	
	return anggotas, nil
}
