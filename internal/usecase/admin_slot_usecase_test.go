package usecase_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
)

type mockSlotRepo struct {
	slots     map[int64]*entity.Pengguna
	anggotas  map[int64]*entity.Anggota
	borrowers map[int64]int64 // slotID -> anggotaID
}

func newMockSlotRepo() *mockSlotRepo {
	return &mockSlotRepo{
		slots:     make(map[int64]*entity.Pengguna),
		anggotas:  make(map[int64]*entity.Anggota),
		borrowers: make(map[int64]int64),
	}
}

func (m *mockSlotRepo) FindAllSlots() ([]entity.AccountSlotResponse, error) {
	var res []entity.AccountSlotResponse
	for _, p := range m.slots {
		resp := entity.AccountSlotResponse{
			IdPengguna:          p.IdPengguna,
			NamaPengguna:        p.NamaPengguna,
			Email:               p.Email,
			Jabatan:             p.NamaPengguna,
			KodeAktivasi:        p.KodeAktivasi,
			KodeAktivasiExpired: p.KodeAktivasiExpired,
			IsFirstLogin:        p.IsFirstLogin,
		}
		if aID, ok := m.borrowers[p.IdPengguna]; ok {
			a := m.anggotas[aID]
			resp.StatusSlot = "TERPAKAI"
			resp.Peminjam = &entity.SlotPeminjamInfo{
				IdAnggota:   a.IdAnggota,
				Nim:         a.Nim,
				NamaLengkap: a.NamaLengkap,
				Angkatan:    a.Angkatan,
				Periode:     a.Periode,
				FotoProfil:  a.FotoProfil,
			}
		} else {
			resp.StatusSlot = "KOSONG"
		}
		res = append(res, resp)
	}
	return res, nil
}

func (m *mockSlotRepo) FindSlotByID(idPengguna int64) (*entity.Pengguna, error) {
	p, ok := m.slots[idPengguna]
	if !ok {
		return nil, errors.New("slot not found")
	}
	return p, nil
}

func (m *mockSlotRepo) FindBorrowerBySlotID(idPengguna int64) (*entity.Anggota, error) {
	aID, ok := m.borrowers[idPengguna]
	if !ok {
		return nil, errors.New("borrower not found")
	}
	a, ok := m.anggotas[aID]
	if !ok {
		return nil, errors.New("anggota not found")
	}
	return a, nil
}

func (m *mockSlotRepo) FindSlotByAnggotaID(idAnggota int64) (*entity.Pengguna, error) {
	for sID, aID := range m.borrowers {
		if aID == idAnggota {
			return m.slots[sID], nil
		}
	}
	return nil, nil
}

func (m *mockSlotRepo) AssignSlot(idPengguna int64, idAnggota int64, kodeAktivasi string, expiredAt time.Time) error {
	if _, ok := m.borrowers[idPengguna]; ok {
		return errors.New("slot akun sudah terpakai")
	}
	for _, aID := range m.borrowers {
		if aID == idAnggota {
			return errors.New("anggota sudah memegang slot lain")
		}
	}
	a, ok := m.anggotas[idAnggota]
	if !ok {
		return errors.New("anggota not found")
	}
	p, ok := m.slots[idPengguna]
	if !ok {
		return errors.New("slot not found")
	}

	m.borrowers[idPengguna] = idAnggota
	a.IdPengguna = &idPengguna
	a.Status = "Pengurus"

	p.KodeAktivasi = &kodeAktivasi
	p.KodeAktivasiExpired = &expiredAt
	p.IsFirstLogin = true
	p.TokenVersion++
	return nil
}

func (m *mockSlotRepo) UnassignSlot(idPengguna int64, randomPasswordHash string) error {
	p, ok := m.slots[idPengguna]
	if !ok {
		return errors.New("slot not found")
	}
	if aID, ok := m.borrowers[idPengguna]; ok {
		if a, ok := m.anggotas[aID]; ok {
			a.IdPengguna = nil
			a.Status = "Demisioner"
		}
		delete(m.borrowers, idPengguna)
	}

	p.Sandi = randomPasswordHash
	p.KodeAktivasi = nil
	p.KodeAktivasiExpired = nil
	p.IsFirstLogin = false
	p.TokenVersion++
	return nil
}

func (m *mockSlotRepo) DemisionerMassal(randomPasswordHash string) (int64, int64, error) {
	totalAnggota := int64(len(m.borrowers))
	totalSlots := int64(len(m.slots))

	for sID, aID := range m.borrowers {
		if a, ok := m.anggotas[aID]; ok {
			a.IdPengguna = nil
			a.Status = "Demisioner"
		}
		delete(m.borrowers, sID)
	}

	for _, p := range m.slots {
		p.Sandi = randomPasswordHash
		p.KodeAktivasi = nil
		p.KodeAktivasiExpired = nil
		p.IsFirstLogin = false
		p.TokenVersion++
	}

	return totalSlots, totalAnggota, nil
}

func (m *mockSlotRepo) ResetSlotPassword(idPengguna int64, newCode string, expiredAt time.Time) error {
	p, ok := m.slots[idPengguna]
	if !ok {
		return errors.New("slot not found")
	}
	p.KodeAktivasi = &newCode
	p.KodeAktivasiExpired = &expiredAt
	p.IsFirstLogin = true
	p.TokenVersion++
	return nil
}

func (m *mockSlotRepo) UpdateSlotRole(idPengguna int64, idRole int64) error {
	p, ok := m.slots[idPengguna]
	if !ok {
		return errors.New("slot not found")
	}
	p.IdPeran = &idRole
	return nil
}

func TestAdminSlotUsecase(t *testing.T) {
	repo := newMockSlotRepo()
	uc := usecase.NewAdminSlotUsecase(repo)

	// Seed 2 slots
	slot1 := &entity.Pengguna{
		IdPengguna:   54,
		Email:        "hmik.kahim.01@hmik.com",
		NamaPengguna: "Ketua Himpunan",
		TokenVersion: 1,
	}
	slot2 := &entity.Pengguna{
		IdPengguna:   55,
		Email:        "hmik.wakahim.01@hmik.com",
		NamaPengguna: "Wakil Ketua Himpunan",
		TokenVersion: 1,
	}
	repo.slots[slot1.IdPengguna] = slot1
	repo.slots[slot2.IdPengguna] = slot2

	// Seed 2 anggota
	ang1 := &entity.Anggota{
		IdAnggota:   1,
		Nim:         "105223001",
		NamaLengkap: "Raihan Putra",
		Status:      "Pengurus",
	}
	ang2 := &entity.Anggota{
		IdAnggota:   2,
		Nim:         "105223049",
		NamaLengkap: "Nabil Hafizuddin",
		Status:      "Pengurus",
	}
	repo.anggotas[ang1.IdAnggota] = ang1
	repo.anggotas[ang2.IdAnggota] = ang2

	// 1. Test GetAccountSlots (initially both KOSONG)
	slots, err := uc.GetAccountSlots()
	if err != nil {
		t.Fatalf("Gagal mengambil slots: %v", err)
	}
	if len(slots) != 2 {
		t.Errorf("Ekspektasi 2 slot, didapat: %d", len(slots))
	}

	// 2. Test AssignSlot
	assignResp, err := uc.AssignSlot(entity.AssignSlotRequest{
		IdPengguna: slot1.IdPengguna,
		IdAnggota:  ang1.IdAnggota,
	})
	if err != nil {
		t.Fatalf("Gagal assign slot: %v", err)
	}
	if !strings.HasPrefix(assignResp.KodeAktivasi, "ACT-") {
		t.Errorf("Format kode aktivasi tidak sesuai: %s", assignResp.KodeAktivasi)
	}
	if assignResp.NamaAnggota != "Raihan Putra" {
		t.Errorf("Nama peminjam tidak sesuai: %s", assignResp.NamaAnggota)
	}

	// 3. Test Assign Already Assigned Slot (must fail)
	_, err = uc.AssignSlot(entity.AssignSlotRequest{
		IdPengguna: slot1.IdPengguna,
		IdAnggota:  ang2.IdAnggota,
	})
	if err == nil {
		t.Errorf("Ekspektasi error slot sudah terpakai, tetapi berhasil")
	}

	// 4. Test ResetSlotPassword
	resetResp, err := uc.ResetSlotPassword(slot1.IdPengguna)
	if err != nil {
		t.Fatalf("Gagal reset password slot: %v", err)
	}
	if !strings.HasPrefix(resetResp.KodeAktivasi, "ACT-") {
		t.Errorf("Kode aktivasi baru tidak berformat ACT-: %s", resetResp.KodeAktivasi)
	}

	// 5. Test UnassignSlot
	err = uc.UnassignSlot(slot1.IdPengguna)
	if err != nil {
		t.Fatalf("Gagal unassign slot: %v", err)
	}
	if ang1.Status != "Demisioner" || ang1.IdPengguna != nil {
		t.Errorf("Anggota belum didemisionerkan dengan benar: status=%s, idPengguna=%v", ang1.Status, ang1.IdPengguna)
	}

	// 6. Test DemisionerMassal
	// Assign slot 2 first
	_, err = uc.AssignSlot(entity.AssignSlotRequest{
		IdPengguna: slot2.IdPengguna,
		IdAnggota:  ang2.IdAnggota,
	})
	if err != nil {
		t.Fatalf("Gagal assign slot 2: %v", err)
	}

	demisSummary, err := uc.DemisionerMassal()
	if err != nil {
		t.Fatalf("Gagal demisioner massal: %v", err)
	}
	if demisSummary.TotalAnggotaDemis != 1 || demisSummary.TotalAkunDireset != 2 {
		t.Errorf("Summary demisioner massal tidak sesuai: %+v", demisSummary)
	}
	if ang2.Status != "Demisioner" {
		t.Errorf("Anggota 2 harusnya Demisioner, tetapi: %s", ang2.Status)
	}

	// 7. Test UpdateSlotRole
	err = uc.UpdateSlotRole(slot1.IdPengguna, 3)
	if err != nil {
		t.Fatalf("Gagal update role slot: %v", err)
	}
	if *slot1.IdPeran != 3 {
		t.Errorf("Role slot belum terupdate: %d", *slot1.IdPeran)
	}
}
