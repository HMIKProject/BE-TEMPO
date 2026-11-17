package usecase

import (
	"errors"
	"time"

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
	slotRepo    repository.SlotRepository
}

// NewAdminBIAUsecase menginisiasi usecase BIA baru
func NewAdminBIAUsecase(anggotaRepo repository.AnggotaRepository, slotRepo repository.SlotRepository) AdminBIAUsecase {
	return &adminBIAUsecase{anggotaRepo: anggotaRepo, slotRepo: slotRepo}
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
		if existing.Status == "Demisioner" {
			return nil, errors.New("NIM " + req.Nim + " sudah tercatat di sistem sebagai anggota Demisioner (" + existing.NamaLengkap + "). Silakan cari di daftar anggota lalu edit datanya untuk mengaktifkan kembali pada periode baru.")
		}
		return nil, errors.New("NIM " + req.Nim + " sudah terdaftar dalam Buku Induk Anggota")
	}

	// 2. Default foto jika kosong
	foto := req.FotoProfil
	if foto == "" {
		foto = "https://api.dicebear.com/7.x/avataaars/svg?seed=" + req.NamaLengkap
	}

	status := req.Status
	if status == "" {
		status = "Anggota"
	}

	// 3. Jika id_pengguna disediakan, sinkronkan jabatan & struktur dari slot akun resmi
	var slotPengguna *entity.Pengguna
	if req.IdPengguna != nil && *req.IdPengguna > 0 && u.slotRepo != nil {
		slot, err := u.slotRepo.FindSlotByID(*req.IdPengguna)
		if err != nil || slot == nil {
			return nil, errors.New("slot akun resmi tidak ditemukan")
		}
		// Cek apakah slot sudah terpakai oleh pengurus lain
		borrower, err := u.slotRepo.FindBorrowerBySlotID(*req.IdPengguna)
		if err == nil && borrower != nil {
			return nil, errors.New("slot akun tersebut sudah terpakai oleh pengurus lain: " + borrower.NamaLengkap)
		}
		slotPengguna = slot
		mapping := repository.GetSlotOrgMapping(slot.NamaPengguna, slot.NamaLengkap)
		req.Jabatan = mapping.Jabatan
		req.IdDepartemen = mapping.IdDepartemen
		req.IdDivisi = mapping.IdDivisi
		status = "Pengurus"
	} else if req.Jabatan == "" {
		if status == "Magang" {
			req.Jabatan = "Magang"
		} else if status == "Demisioner" {
			req.Jabatan = "Demisioner"
		} else {
			req.Jabatan = "Anggota"
		}
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
		DepartemenID: req.IdDepartemen,
		DivisiID:     req.IdDivisi,
		Periode:      req.Periode,
		PenggunaID:   req.IdPengguna,
	}

	if err := u.anggotaRepo.Create(newAnggota); err != nil {
		return nil, err
	}

	// Jika memilih slot akun, terbitkan kode aktivasi 7 hari
	if slotPengguna != nil && u.slotRepo != nil {
		code, _ := GenerateActivationCode()
		expiredAt := time.Now().Add(7 * 24 * time.Hour)
		_ = u.slotRepo.AssignSlot(slotPengguna.IdPengguna, newAnggota.IdAnggota, code, expiredAt)
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

	// Penanganan sinkronisasi slot akun
	if req.IdPengguna != nil && *req.IdPengguna > 0 && u.slotRepo != nil {
		if anggota.PenggunaID == nil || *anggota.PenggunaID != *req.IdPengguna {
			// Jika anggota sebelumnya memegang slot lain, lepas dulu
			if anggota.PenggunaID != nil {
				hash, _ := GenerateRandomPasswordHash()
				_ = u.slotRepo.UnassignSlot(*anggota.PenggunaID, hash)
			}
			code, _ := GenerateActivationCode()
			expiredAt := time.Now().Add(7 * 24 * time.Hour)
			_ = u.slotRepo.AssignSlot(*req.IdPengguna, id, code, expiredAt)
		}
		// Selaraskan jabatan & struktur dengan slot yang dipegang
		slot, err := u.slotRepo.FindSlotByID(*req.IdPengguna)
		if err == nil && slot != nil {
			mapping := repository.GetSlotOrgMapping(slot.NamaPengguna, slot.NamaLengkap)
			req.Jabatan = mapping.Jabatan
			req.IdDepartemen = mapping.IdDepartemen
			req.IdDivisi = mapping.IdDivisi
			req.Status = "Pengurus"
		}
	} else if (req.IdPengguna == nil || *req.IdPengguna == 0) && anggota.PenggunaID != nil && u.slotRepo != nil {
		// Jika slot dilepas dari anggota
		hash, _ := GenerateRandomPasswordHash()
		_ = u.slotRepo.UnassignSlot(*anggota.PenggunaID, hash)
		anggota.PenggunaID = nil
	}

	// Rekam riwayat jabatan jika terjadi transisi kepengurusan (periode/jabatan berubah atau demisioner aktif kembali)
	isPeriodChange := anggota.Periode != "" && req.Periode != "" && anggota.Periode != req.Periode
	isJabatanChange := anggota.Jabatan != "" && req.Jabatan != "" && anggota.Jabatan != req.Jabatan
	isReactivation := anggota.Status == "Demisioner" && req.Status == "Pengurus"

	if isPeriodChange || isJabatanChange || isReactivation {
		deptName := ""
		if anggota.Departemen != nil {
			deptName = anggota.Departemen.NamaDepartemen
		}
		divName := ""
		if anggota.Divisi != nil {
			divName = anggota.Divisi.NamaDivisi
		}
		riwayat := &entity.RiwayatJabatan{
			IdAnggota:      anggota.IdAnggota,
			Nim:            anggota.Nim,
			Periode:        anggota.Periode,
			Jabatan:        anggota.Jabatan,
			NamaDepartemen: deptName,
			NamaDivisi:     divName,
			Status:         "Demisioner",
		}
		_ = u.anggotaRepo.AddRiwayatJabatan(riwayat)
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
	anggota.DepartemenID = req.IdDepartemen
	anggota.DivisiID = req.IdDivisi
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
	if anggota.PenggunaID != nil {
		return errors.New("tidak dapat menghapus anggota yang sedang memegang slot akun aktif; silakan lepas tautan slot akun terlebih dahulu")
	}

	return u.anggotaRepo.Delete(id)
}
