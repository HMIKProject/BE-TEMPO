package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AdminSlotUsecase adalah interface logika bisnis manajemen slot akun pinjam pakai
type AdminSlotUsecase interface {
	GetAccountSlots() ([]entity.AccountSlotResponse, error)
	AssignSlot(req entity.AssignSlotRequest) (*entity.AssignSlotResponse, error)
	UnassignSlot(idPengguna int64) error
	DemisionerMassal() (*entity.DemisionerSummary, error)
	ResetSlotPassword(idPengguna int64) (*entity.ResetSlotResponse, error)
	UpdateSlotRole(idPengguna int64, idRole int64) error
}

type adminSlotUsecase struct {
	slotRepo repository.SlotRepository
}

// NewAdminSlotUsecase menginisiasi usecase slot akun baru
func NewAdminSlotUsecase(slotRepo repository.SlotRepository) AdminSlotUsecase {
	return &adminSlotUsecase{slotRepo: slotRepo}
}

// generateActivationCode membuat kode unik sekali pakai format ACT-XXXX-XXXX (alfanumerik kapital)
func generateActivationCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Hindari karakter ambigu O/0, I/1
	part1 := make([]byte, 4)
	part2 := make([]byte, 4)

	for i := range part1 {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		part1[i] = charset[num.Int64()]
	}

	for i := range part2 {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		part2[i] = charset[num.Int64()]
	}

	return fmt.Sprintf("ACT-%s-%s", string(part1), string(part2)), nil
}

// generateRandomPasswordHash menghasilkan hash bcrypt dari random string untuk mengamankan akun kosong
func generateRandomPasswordHash() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomStr := hex.EncodeToString(randomBytes)
	hash, err := bcrypt.GenerateFromPassword([]byte(randomStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GetAccountSlots mengambil seluruh slot akun beserta info keterisiannya
func (u *adminSlotUsecase) GetAccountSlots() ([]entity.AccountSlotResponse, error) {
	return u.slotRepo.FindAllSlots()
}

// AssignSlot menautkan slot kosong ke mahasiswa terpilih dan menerbitkan kode aktivasi 7 hari
func (u *adminSlotUsecase) AssignSlot(req entity.AssignSlotRequest) (*entity.AssignSlotResponse, error) {
	// 1. Cek slot akun
	slot, err := u.slotRepo.FindSlotByID(req.IdPengguna)
	if err != nil {
		return nil, errors.New("slot akun tidak ditemukan")
	}

	// 2. Buat kode aktivasi 8 karakter dengan TTL 7 hari
	code, err := generateActivationCode()
	if err != nil {
		return nil, errors.New("gagal menerbitkan kode aktivasi")
	}
	expiredAt := time.Now().Add(7 * 24 * time.Hour)

	// 3. Tautkan via repository (Atomic Transaction)
	if err := u.slotRepo.AssignSlot(req.IdPengguna, req.IdAnggota, code, expiredAt); err != nil {
		return nil, err
	}

	// 4. Ambil data peminjam untuk informasi kembalian
	borrower, err := u.slotRepo.FindBorrowerBySlotID(req.IdPengguna)
	if err != nil {
		return nil, errors.New("slot berhasil ditautkan tetapi gagal memuat profil anggota")
	}

	return &entity.AssignSlotResponse{
		IdPengguna:          slot.IdPengguna,
		Email:               slot.Email,
		NamaPengguna:        slot.NamaPengguna,
		Jabatan:             slot.NamaPengguna,
		NamaAnggota:         borrower.NamaLengkap,
		NimAnggota:          borrower.Nim,
		KodeAktivasi:        code,
		KodeAktivasiExpired: expiredAt,
	}, nil
}

// UnassignSlot melepas kaitan akun dari anggota, mendemisionerkan anggota, dan mereset password acak
func (u *adminSlotUsecase) UnassignSlot(idPengguna int64) error {
	// Pastikan slot ada
	_, err := u.slotRepo.FindSlotByID(idPengguna)
	if err != nil {
		return errors.New("slot akun tidak ditemukan")
	}

	// Buat hash sandi acak baru (Zero-Knowledge)
	randomHash, err := generateRandomPasswordHash()
	if err != nil {
		return errors.New("gagal mereset keamanan kredensial slot")
	}

	return u.slotRepo.UnassignSlot(idPengguna, randomHash)
}

// DemisionerMassal mengeksekusi transisi serentak akhir periode kabinet
func (u *adminSlotUsecase) DemisionerMassal() (*entity.DemisionerSummary, error) {
	randomHash, err := generateRandomPasswordHash()
	if err != nil {
		return nil, errors.New("gagal mereset keamanan kredensial massal")
	}

	totalAkun, totalAnggota, err := u.slotRepo.DemisionerMassal(randomHash)
	if err != nil {
		return nil, err
	}

	return &entity.DemisionerSummary{
		TotalAkunDireset:  totalAkun,
		TotalAnggotaDemis: totalAnggota,
		WaktuEksekusi:     time.Now(),
	}, nil
}

// ResetSlotPassword menerbitkan ulang kode aktivasi 7 hari bagi pengurus yang lupa kata sandi
func (u *adminSlotUsecase) ResetSlotPassword(idPengguna int64) (*entity.ResetSlotResponse, error) {
	slot, err := u.slotRepo.FindSlotByID(idPengguna)
	if err != nil {
		return nil, errors.New("slot akun tidak ditemukan")
	}

	// Terbitkan kode aktivasi baru
	code, err := generateActivationCode()
	if err != nil {
		return nil, errors.New("gagal menerbitkan kode aktivasi baru")
	}
	expiredAt := time.Now().Add(7 * 24 * time.Hour)

	if err := u.slotRepo.ResetSlotPassword(idPengguna, code, expiredAt); err != nil {
		return nil, err
	}

	return &entity.ResetSlotResponse{
		IdPengguna:          slot.IdPengguna,
		Email:               slot.Email,
		KodeAktivasi:        code,
		KodeAktivasiExpired: expiredAt,
	}, nil
}

// UpdateSlotRole mengubah peran hak akses slot akun
func (u *adminSlotUsecase) UpdateSlotRole(idPengguna int64, idRole int64) error {
	_, err := u.slotRepo.FindSlotByID(idPengguna)
	if err != nil {
		return errors.New("slot akun tidak ditemukan")
	}

	return u.slotRepo.UpdateSlotRole(idPengguna, idRole)
}
