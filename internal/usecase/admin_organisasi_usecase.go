package usecase

import (
	"errors"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

type AdminOrganisasiUsecase interface {
	GetStrukturOrganisasi() (*entity.StrukturOrganisasiResponse, error)
	CreateDepartemen(req entity.CreateDeptRequest) (*entity.Departemen, error)
	UpdateDepartemen(id int64, req entity.UpdateDeptRequest) (*entity.Departemen, error)
	ArchiveDepartemen(id int64) error
	CreateDivisi(req entity.CreateDivisiRequest) (*entity.DivisiDepartemen, error)
	UpdateDivisi(id int64, req entity.UpdateDivisiRequest) (*entity.DivisiDepartemen, error)
}

type adminOrganisasiUsecase struct {
	orgRepo repository.OrganisasiRepository
}

func NewAdminOrganisasiUsecase(orgRepo repository.OrganisasiRepository) AdminOrganisasiUsecase {
	return &adminOrganisasiUsecase{orgRepo: orgRepo}
}

func (u *adminOrganisasiUsecase) GetStrukturOrganisasi() (*entity.StrukturOrganisasiResponse, error) {
	return u.orgRepo.GetStrukturOrganisasi()
}

func (u *adminOrganisasiUsecase) CreateDepartemen(req entity.CreateDeptRequest) (*entity.Departemen, error) {
	nama := strings.TrimSpace(req.NamaDepartemen)
	if nama == "" {
		return nil, errors.New("nama departemen tidak boleh kosong")
	}

	dept := &entity.Departemen{
		NamaDepartemen: nama,
		Deskripsi:      strings.TrimSpace(req.Deskripsi),
		Logo:           strings.TrimSpace(req.Logo),
		StatusAktif:    true,
	}

	err := u.orgRepo.CreateDepartemen(dept)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func (u *adminOrganisasiUsecase) UpdateDepartemen(id int64, req entity.UpdateDeptRequest) (*entity.Departemen, error) {
	dept, err := u.orgRepo.FindDepartemenByID(id)
	if err != nil {
		return nil, errors.New("departemen tidak ditemukan")
	}

	nama := strings.TrimSpace(req.NamaDepartemen)
	if nama == "" {
		return nil, errors.New("nama departemen tidak boleh kosong")
	}

	dept.NamaDepartemen = nama
	dept.Deskripsi = strings.TrimSpace(req.Deskripsi)
	dept.Logo = strings.TrimSpace(req.Logo)
	if req.StatusAktif != nil {
		dept.StatusAktif = *req.StatusAktif
	}

	err = u.orgRepo.UpdateDepartemen(dept)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func (u *adminOrganisasiUsecase) ArchiveDepartemen(id int64) error {
	_, err := u.orgRepo.FindDepartemenByID(id)
	if err != nil {
		return errors.New("departemen tidak ditemukan")
	}
	return u.orgRepo.ArchiveDepartemen(id)
}

func (u *adminOrganisasiUsecase) CreateDivisi(req entity.CreateDivisiRequest) (*entity.DivisiDepartemen, error) {
	_, err := u.orgRepo.FindDepartemenByID(req.IdDepartemen)
	if err != nil {
		return nil, errors.New("departemen induk tidak ditemukan")
	}

	nama := strings.TrimSpace(req.NamaDivisi)
	if nama == "" {
		return nil, errors.New("nama divisi tidak boleh kosong")
	}

	divisi := &entity.DivisiDepartemen{
		IdDepartemen: req.IdDepartemen,
		NamaDivisi:   nama,
		Deskripsi:    strings.TrimSpace(req.Deskripsi),
		StatusAktif:  true,
	}

	err = u.orgRepo.CreateDivisi(divisi)
	if err != nil {
		return nil, err
	}
	return divisi, nil
}

func (u *adminOrganisasiUsecase) UpdateDivisi(id int64, req entity.UpdateDivisiRequest) (*entity.DivisiDepartemen, error) {
	divisi, err := u.orgRepo.FindDivisiByID(id)
	if err != nil {
		return nil, errors.New("divisi tidak ditemukan")
	}

	nama := strings.TrimSpace(req.NamaDivisi)
	if nama == "" {
		return nil, errors.New("nama divisi tidak boleh kosong")
	}

	divisi.NamaDivisi = nama
	divisi.Deskripsi = strings.TrimSpace(req.Deskripsi)
	if req.StatusAktif != nil {
		divisi.StatusAktif = *req.StatusAktif
	}

	err = u.orgRepo.UpdateDivisi(divisi)
	if err != nil {
		return nil, err
	}
	return divisi, nil
}
