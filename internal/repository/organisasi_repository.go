package repository

import (
	"errors"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

type OrganisasiRepository interface {
	GetStrukturOrganisasi() (*entity.StrukturOrganisasiResponse, error)
	FindDepartemenByID(id int64) (*entity.Departemen, error)
	CreateDepartemen(dept *entity.Departemen) error
	UpdateDepartemen(dept *entity.Departemen) error
	ArchiveDepartemen(id int64) error

	FindDivisiByID(id int64) (*entity.DivisiDepartemen, error)
	CreateDivisi(divisi *entity.DivisiDepartemen) error
	UpdateDivisi(divisi *entity.DivisiDepartemen) error
	ArchiveDivisi(id int64) error
}

type organisasiRepository struct {
	db *gorm.DB
}

func NewOrganisasiRepository(db *gorm.DB) OrganisasiRepository {
	return &organisasiRepository{db: db}
}

func (r *organisasiRepository) GetStrukturOrganisasi() (*entity.StrukturOrganisasiResponse, error) {
	// 1. Ambil data Kabinet Aktif (Prioritaskan is_aktif = true, fallback ke id_kabinet terbaru)
	var kabinet entity.Kabinet
	err := r.db.Where("is_aktif = ?", true).Order("id_kabinet DESC").First(&kabinet).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var kabinetPtr *entity.Kabinet
	if kabinet.IdKabinet > 0 {
		kabinetPtr = &kabinet
	}

	// 2. Ambil data Pengurus Inti BPH (Kahim, Wakahim, Sekretaris, Bendahara)
	// BPH dicirikan dengan id_departemen IS NULL atau jabatan ketua/wakil himpunan dengan status Pengurus
	var bphAnggota []entity.Anggota
	err = r.db.Where("status = ? AND (id_departemen IS NULL OR jabatan ILIKE ? OR jabatan ILIKE ?)",
		"Pengurus", "%ketua himpunan%", "%wakil ketua himpunan%").
		Order("id_anggota ASC").
		Find(&bphAnggota).Error
	if err != nil {
		return nil, err
	}

	// Sorting urutan BPH hierarkis: Kahim -> Wakahim -> Sekre 1 -> Sekre 2 -> Benda 1 -> Benda 2
	bphList := make([]entity.BPHMemberInfo, 0, len(bphAnggota))
	for _, a := range bphAnggota {
		bphList = append(bphList, entity.BPHMemberInfo{
			IdAnggota:    a.IdAnggota,
			NamaLengkap:  a.NamaLengkap,
			Nim:          a.Nim,
			Jabatan:      a.Jabatan,
			FotoProfil:   a.FotoProfil,
			Angkatan:     a.Angkatan,
			ProgramStudi: a.ProgramStudi,
			Status:       a.Status,
			Periode:      a.Periode,
		})
	}

	// 3. Ambil seluruh data Departemen beserta Divisi-nya
	var rawDepts []entity.Departemen
	err = r.db.Order("id_departemen ASC").Find(&rawDepts).Error
	if err != nil {
		return nil, err
	}

	deptItems := make([]entity.OrganisasiDeptItem, 0, len(rawDepts))
	for _, d := range rawDepts {
		// Ambil divisi di bawah departemen ini
		var divisiList []entity.DivisiDepartemen
		_ = r.db.Where("id_departemen = ?", d.IdDepartemen).Order("id_divisi ASC").Find(&divisiList).Error

		// Cari Kepala Departemen (Kadep)
		var kadepAnggota entity.Anggota
		errKadep := r.db.Where("id_departemen = ? AND status = ? AND (jabatan ILIKE ? OR jabatan ILIKE ?)",
			d.IdDepartemen, "Pengurus", "%ketua departemen%", "%kadep%").
			First(&kadepAnggota).Error

		var kadepInfo *entity.DeptKadepInfo
		if errKadep == nil && kadepAnggota.IdAnggota > 0 {
			kadepInfo = &entity.DeptKadepInfo{
				IdAnggota:   kadepAnggota.IdAnggota,
				NamaLengkap: kadepAnggota.NamaLengkap,
				Nim:         kadepAnggota.Nim,
				Jabatan:     kadepAnggota.Jabatan,
				FotoProfil:  kadepAnggota.FotoProfil,
				Angkatan:    kadepAnggota.Angkatan,
			}
		}

		// Hitung total anggota pengurus aktif di departemen ini
		var totalAnggota int64
		_ = r.db.Model(&entity.Anggota{}).Where("id_departemen = ? AND status = ?", d.IdDepartemen, "Pengurus").Count(&totalAnggota).Error

		deptItems = append(deptItems, entity.OrganisasiDeptItem{
			IdDepartemen:   d.IdDepartemen,
			NamaDepartemen: d.NamaDepartemen,
			Deskripsi:      d.Deskripsi,
			Logo:           d.Logo,
			StatusAktif:    d.StatusAktif,
			Kadep:          kadepInfo,
			Divisi:         divisiList,
			TotalAnggota:   totalAnggota,
			TotalDivisi:    int64(len(divisiList)),
		})
	}

	return &entity.StrukturOrganisasiResponse{
		Kabinet:    kabinetPtr,
		BPH:        bphList,
		Departemen: deptItems,
	}, nil
}

func (r *organisasiRepository) FindDepartemenByID(id int64) (*entity.Departemen, error) {
	var dept entity.Departemen
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *organisasiRepository) CreateDepartemen(dept *entity.Departemen) error {
	return r.db.Create(dept).Error
}

func (r *organisasiRepository) UpdateDepartemen(dept *entity.Departemen) error {
	return r.db.Save(dept).Error
}

func (r *organisasiRepository) ArchiveDepartemen(id int64) error {
	return r.db.Model(&entity.Departemen{}).Where("id_departemen = ?", id).Update("status_aktif", false).Error
}

func (r *organisasiRepository) FindDivisiByID(id int64) (*entity.DivisiDepartemen, error) {
	var div entity.DivisiDepartemen
	err := r.db.First(&div, id).Error
	if err != nil {
		return nil, err
	}
	return &div, nil
}

func (r *organisasiRepository) CreateDivisi(divisi *entity.DivisiDepartemen) error {
	return r.db.Create(divisi).Error
}

func (r *organisasiRepository) UpdateDivisi(divisi *entity.DivisiDepartemen) error {
	return r.db.Save(divisi).Error
}

func (r *organisasiRepository) ArchiveDivisi(id int64) error {
	return r.db.Model(&entity.DivisiDepartemen{}).Where("id_divisi = ?", id).Update("status_aktif", false).Error
}
