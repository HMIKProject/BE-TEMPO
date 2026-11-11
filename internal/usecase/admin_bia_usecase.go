package usecase

import (
	"errors"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

// AdminBIAUsecase adalah interface logika bisnis pengelolaan Buku Induk Anggota (BIA)
type AdminBIAUsecase interface {
	GetAnggotaList(filter entity.AnggotaFilter) (*entity.AnggotaListResponse, error)
	GetAnggotaByID(id int64) (*entity.Anggota, error)
	CreateAnggota(req entity.CreateAnggotaRequest) (*entity.Anggota, error)
	UpdateAnggota(id int64, req entity.UpdateAnggotaRequest) (*entity.Anggota, error)
	DeleteAnggota(id int64) error
}

type adminBIAUsecase struct {
	anggotaRepo repository.AnggotaRepository
}

// NewAdminBIAUsecase menginisiasi usecase BIA baru
func NewAdminBIAUsecase(anggotaRepo repository.AnggotaRepository) AdminBIAUsecase {
	return &adminBIAUsecase{anggotaRepo: anggotaRepo}
}

// GetAnggotaList mengambil daftar anggota terfilter dan terpaginasi
func (u *adminBIAUsecase) GetAnggotaList(filter entity.AnggotaFilter) (*entity.AnggotaListResponse, error) {
	anggotas, totalData, err := u.anggotaRepo.FindWithFilter(filter)
	if err != nil {
		return nil, err
	}

	return &entity.AnggotaListResponse{
		TotalData: totalData,
		Anggota:   anggotas,
	}, nil
}

// GetAnggotaByID mengambil detail profil lengkap seorang anggota
func (u *adminBIAUsecase) GetAnggotaByID(id int64) (*entity.Anggota, error) {
	anggota, err := u.anggotaRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("data anggota tidak ditemukan")
	}
	return anggota, nil
}

// CreateAnggota mendaftarkan anggota baru secara manual ke Buku Induk
func (u *adminBIAUsecase) CreateAnggota(req entity.CreateAnggotaRequest) (*entity.Anggota, error) {
	// 1. Cek apakah NIM sudah terdaftar
	existing, _ := u.anggotaRepo.FindByNIM(req.Nim)
	if existing != nil {
		return nil, errors.New("NIM " + req.Nim + " sudah terdaftar dalam Buku Induk Anggota")
	}

	// 2. Default foto jika kosong
	foto := req.FotoProfil
	if foto == "" {
		foto = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + req.NamaLengkap
	}

	status := req.Status
	if status == "" {
		status = "Pengurus"
	}

	newAnggota := &entity.Anggota{
		Nim:          req.Nim,
		NamaLengkap:  req.NamaLengkap,
		FotoProfil:   foto,
		ProgramStudi: req.ProgramStudi,
		Angkatan:     req.Angkatan,
		Jabatan:      req.Jabatan,
		Status:       status,
		Tautan:       req.Tautan,
		CvUrl:        req.CvUrl,
		IdDepartemen: req.IdDepartemen,
		IdDivisi:     req.IdDivisi,
		Periode:      req.Periode,
	}

	if err := u.anggotaRepo.Create(newAnggota); err != nil {
		return nil, err
	}

	// Ambil kembali data lengkap dengan eager loading
	return u.anggotaRepo.FindByID(newAnggota.IdAnggota)
}

// UpdateAnggota memperbarui profil mahasiswa pada Buku Induk
func (u *adminBIAUsecase) UpdateAnggota(id int64, req entity.UpdateAnggotaRequest) (*entity.Anggota, error) {
	anggota, err := u.anggotaRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("data anggota tidak ditemukan")
	}

	anggota.NamaLengkap = req.NamaLengkap
	if req.FotoProfil != "" {
		anggota.FotoProfil = req.FotoProfil
	}
	anggota.ProgramStudi = req.ProgramStudi
	anggota.Angkatan = req.Angkatan
	anggota.Jabatan = req.Jabatan
	anggota.Status = req.Status
	anggota.Tautan = req.Tautan
	anggota.CvUrl = req.CvUrl
	anggota.IdDepartemen = req.IdDepartemen
	anggota.IdDivisi = req.IdDivisi
	anggota.Periode = req.Periode

	if err := u.anggotaRepo.Update(anggota); err != nil {
		return nil, err
	}

	return u.anggotaRepo.FindByID(id)
}

// DeleteAnggota menghapus entri anggota (hanya jika tidak memegang slot aktif)
func (u *adminBIAUsecase) DeleteAnggota(id int64) error {
	anggota, err := u.anggotaRepo.FindByID(id)
	if err != nil {
		return errors.New("data anggota tidak ditemukan")
	}

	// Proteksi integritas: jika sedang memegang slot akun aktif, tidak boleh langsung dihapus
	if anggota.IdPengguna != nil {
		return errors.New("tidak dapat menghapus anggota yang sedang memegang slot akun aktif; silakan lepas tautan slot akun terlebih dahulu")
	}

	return u.anggotaRepo.Delete(id)
}
