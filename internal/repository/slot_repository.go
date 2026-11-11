package repository

import (
	"errors"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

// SlotRepository adalah interface interaksi database untuk manajemen slot akun pinjam pakai
type SlotRepository interface {
	FindAllSlots() ([]entity.AccountSlotResponse, error)
	FindSlotByID(idPengguna int64) (*entity.Pengguna, error)
	FindBorrowerBySlotID(idPengguna int64) (*entity.Anggota, error)
	FindSlotByAnggotaID(idAnggota int64) (*entity.Pengguna, error)
	AssignSlot(idPengguna int64, idAnggota int64, kodeAktivasi string, expiredAt time.Time) error
	UnassignSlot(idPengguna int64, randomPasswordHash string) error
	DemisionerMassal(randomPasswordHash string) (int64, int64, error)
	ResetSlotPassword(idPengguna int64, newCode string, expiredAt time.Time) error
	UpdateSlotRole(idPengguna int64, idRole int64) error
}

type slotRepository struct {
	db *gorm.DB
}

// NewSlotRepository menginisiasi repository slot akun baru
func NewSlotRepository(db *gorm.DB) SlotRepository {
	return &slotRepository{db: db}
}

// FindAllSlots mengambil seluruh slot akun struktural beserta info peminjamnya
func (r *slotRepository) FindAllSlots() ([]entity.AccountSlotResponse, error) {
	var penggunaList []entity.Pengguna
	// Query semua akun slot struktural (diawali 'hmik.' dan berakhiran '@hmik.com')
	err := r.db.Where("email LIKE ?", "hmik.%@hmik.com").
		Order("id_pengguna ASC").
		Find(&penggunaList).Error
	if err != nil {
		return nil, err
	}

	// Ambil semua anggota yang terikat dengan akun-akun tersebut
	var anggotaList []entity.Anggota
	err = r.db.Where("id_pengguna IS NOT NULL").Find(&anggotaList).Error
	if err != nil {
		return nil, err
	}

	// Petakan id_pengguna -> anggota
	borrowerMap := make(map[int64]entity.Anggota)
	for _, a := range anggotaList {
		if a.IdPengguna != nil {
			borrowerMap[*a.IdPengguna] = a
		}
	}

	var responses []entity.AccountSlotResponse
	for _, p := range penggunaList {
		resp := entity.AccountSlotResponse{
			IdPengguna:          p.IdPengguna,
			NamaPengguna:        p.NamaPengguna,
			Email:               p.Email,
			Jabatan:             p.NamaPengguna,
			KodeAktivasi:        p.KodeAktivasi,
			KodeAktivasiExpired: p.KodeAktivasiExpired,
			IsFirstLogin:        p.IsFirstLogin,
		}

		if borrower, ok := borrowerMap[p.IdPengguna]; ok {
			resp.StatusSlot = "TERPAKAI"
			resp.Peminjam = &entity.SlotPeminjamInfo{
				IdAnggota:   borrower.IdAnggota,
				Nim:         borrower.Nim,
				NamaLengkap: borrower.NamaLengkap,
				Angkatan:    borrower.Angkatan,
				Periode:     borrower.Periode,
				FotoProfil:  borrower.FotoProfil,
			}
		} else {
			resp.StatusSlot = "KOSONG"
			resp.Peminjam = nil
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

// FindSlotByID mengambil data pengguna slot berdasarkan ID
func (r *slotRepository) FindSlotByID(idPengguna int64) (*entity.Pengguna, error) {
	var p entity.Pengguna
	err := r.db.First(&p, idPengguna).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBorrowerBySlotID mencari anggota yang saat ini sedang meminjam slot
func (r *slotRepository) FindBorrowerBySlotID(idPengguna int64) (*entity.Anggota, error) {
	var a entity.Anggota
	err := r.db.Where("id_pengguna = ?", idPengguna).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// FindSlotByAnggotaID mencari apakah seorang anggota sudah memegang slot lain
func (r *slotRepository) FindSlotByAnggotaID(idAnggota int64) (*entity.Pengguna, error) {
	var a entity.Anggota
	err := r.db.First(&a, idAnggota).Error
	if err != nil {
		return nil, err
	}
	if a.IdPengguna == nil {
		return nil, nil
	}
	return r.FindSlotByID(*a.IdPengguna)
}

// AssignSlot menautkan slot kosong ke anggota dan memasang kode aktivasi 7 hari
func (r *slotRepository) AssignSlot(idPengguna int64, idAnggota int64, kodeAktivasi string, expiredAt time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Cek apakah slot sudah terisi
		var existingBorrower entity.Anggota
		err := tx.Where("id_pengguna = ?", idPengguna).First(&existingBorrower).Error
		if err == nil {
			return errors.New("slot akun sudah terpakai oleh pengurus lain, lepas tautan terlebih dahulu")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 2. Cek apakah anggota sudah punya slot lain
		var targetAnggota entity.Anggota
		if err := tx.First(&targetAnggota, idAnggota).Error; err != nil {
			return errors.New("data anggota tidak ditemukan")
		}
		if targetAnggota.IdPengguna != nil {
			return errors.New("anggota tersebut sudah terdaftar memegang slot akun lain")
		}

		// 3. Update tabel anggota: tautkan id_pengguna & pastikan status 'Pengurus'
		if err := tx.Model(&targetAnggota).Updates(map[string]interface{}{
			"id_pengguna": idPengguna,
			"status":      "Pengurus",
		}).Error; err != nil {
			return err
		}

		// 4. Update tabel pengguna: pasang kode aktivasi & aktifkan is_first_login
		if err := tx.Model(&entity.Pengguna{}).Where("id_pengguna = ?", idPengguna).Updates(map[string]interface{}{
			"kode_aktivasi":         kodeAktivasi,
			"kode_aktivasi_expired": expiredAt,
			"is_first_login":        true,
			"token_version":         gorm.Expr("token_version + 1"),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// UnassignSlot melepas kaitan anggota dari slot, mendemisionerkan anggota, dan mereset password acak
func (r *slotRepository) UnassignSlot(idPengguna int64, randomPasswordHash string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lepas foreign key pada anggota peminjam dan ubah status ke Demisioner
		if err := tx.Model(&entity.Anggota{}).
			Where("id_pengguna = ?", idPengguna).
			Updates(map[string]interface{}{
				"id_pengguna": nil,
				"status":      "Demisioner",
			}).Error; err != nil {
			return err
		}

		// 2. Reset kredensial slot ke random password, hanguskan kode aktivasi, cabut sesi JWT
		if err := tx.Model(&entity.Pengguna{}).
			Where("id_pengguna = ?", idPengguna).
			Updates(map[string]interface{}{
				"sandi":                 randomPasswordHash,
				"kode_aktivasi":         nil,
				"kode_aktivasi_expired": nil,
				"is_first_login":        false,
				"token_version":         gorm.Expr("token_version + 1"),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// DemisionerMassal memutus serentak seluruh slot terpakai, mendemisionerkan anggota, dan mereset sandi massal
func (r *slotRepository) DemisionerMassal(randomPasswordHash string) (int64, int64, error) {
	var totalAkunDireset int64
	var totalAnggotaDemis int64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update anggota yang sedang memegang slot menjadi Demisioner dan unbind id_pengguna
		resAnggota := tx.Model(&entity.Anggota{}).
			Where("id_pengguna IS NOT NULL").
			Updates(map[string]interface{}{
				"status":      "Demisioner",
				"id_pengguna": nil,
			})
		if resAnggota.Error != nil {
			return resAnggota.Error
		}
		totalAnggotaDemis = resAnggota.RowsAffected

		// 2. Reset semua akun slot struktural
		resPengguna := tx.Model(&entity.Pengguna{}).
			Where("email LIKE ?", "hmik.%@hmik.com").
			Updates(map[string]interface{}{
				"sandi":                 randomPasswordHash,
				"kode_aktivasi":         nil,
				"kode_aktivasi_expired": nil,
				"is_first_login":        false,
				"token_version":         gorm.Expr("token_version + 1"),
			})
		if resPengguna.Error != nil {
			return resPengguna.Error
		}
		totalAkunDireset = resPengguna.RowsAffected

		return nil
	})

	return totalAkunDireset, totalAnggotaDemis, err
}

// ResetSlotPassword menerbitkan ulang kode aktivasi baru untuk pemegang akun (kasus lupa sandi)
func (r *slotRepository) ResetSlotPassword(idPengguna int64, newCode string, expiredAt time.Time) error {
	return r.db.Model(&entity.Pengguna{}).
		Where("id_pengguna = ?", idPengguna).
		Updates(map[string]interface{}{
			"kode_aktivasi":         newCode,
			"kode_aktivasi_expired": expiredAt,
			"is_first_login":        true,
			"token_version":         gorm.Expr("token_version + 1"),
		}).Error
}

// UpdateSlotRole mengubah peran / role hak akses slot akun
func (r *slotRepository) UpdateSlotRole(idPengguna int64, idRole int64) error {
	return r.db.Model(&entity.Pengguna{}).
		Where("id_pengguna = ?", idPengguna).
		Update("id_peran", idRole).Error
}
