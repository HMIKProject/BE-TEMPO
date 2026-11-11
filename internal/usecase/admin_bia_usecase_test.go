package usecase_test

import (
	"errors"
	"testing"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
)

type mockAnggotaRepo struct {
	data   map[int64]*entity.Anggota
	byNIM  map[string]*entity.Anggota
	nextID int64
}

func newMockAnggotaRepo() *mockAnggotaRepo {
	return &mockAnggotaRepo{
		data:   make(map[int64]*entity.Anggota),
		byNIM:  make(map[string]*entity.Anggota),
		nextID: 1,
	}
}

func (m *mockAnggotaRepo) FindAll() ([]entity.Anggota, error) {
	var res []entity.Anggota
	for _, a := range m.data {
		res = append(res, *a)
	}
	return res, nil
}

func (m *mockAnggotaRepo) FindWithFilter(filter entity.AnggotaFilter) ([]entity.Anggota, int64, error) {
	var filtered []entity.Anggota
	for _, a := range m.data {
		if filter.Status != "" && filter.Status != "Semua" && a.Status != filter.Status {
			continue
		}
		if filter.Angkatan > 0 && a.Angkatan != filter.Angkatan {
			continue
		}
		filtered = append(filtered, *a)
	}
	return filtered, int64(len(filtered)), nil
}

func (m *mockAnggotaRepo) FindByID(id int64) (*entity.Anggota, error) {
	a, ok := m.data[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return a, nil
}

func (m *mockAnggotaRepo) FindByNIM(nim string) (*entity.Anggota, error) {
	a, ok := m.byNIM[nim]
	if !ok {
		return nil, errors.New("record not found")
	}
	return a, nil
}

func (m *mockAnggotaRepo) Create(anggota *entity.Anggota) error {
	anggota.IdAnggota = m.nextID
	m.nextID++
	m.data[anggota.IdAnggota] = anggota
	m.byNIM[anggota.Nim] = anggota
	return nil
}

func (m *mockAnggotaRepo) Update(anggota *entity.Anggota) error {
	m.data[anggota.IdAnggota] = anggota
	m.byNIM[anggota.Nim] = anggota
	return nil
}

func (m *mockAnggotaRepo) Delete(id int64) error {
	a, ok := m.data[id]
	if !ok {
		return errors.New("record not found")
	}
	delete(m.byNIM, a.Nim)
	delete(m.data, id)
	return nil
}

func TestAdminBIAUsecase_CRUD(t *testing.T) {
	repo := newMockAnggotaRepo()
	uc := usecase.NewAdminBIAUsecase(repo)

	// 1. Test CreateAnggota
	createReq := entity.CreateAnggotaRequest{
		Nim:          "105223046",
		NamaLengkap:  "Anton Dicaprio",
		ProgramStudi: "Ilmu Komputer",
		Angkatan:     2023,
		Jabatan:      "Ketua Divisi Media Kreatif",
		Status:       "Pengurus",
		Tautan:       "https://linkedin.com/in/antondicaprio",
	}

	created, err := uc.CreateAnggota(createReq)
	if err != nil {
		t.Fatalf("Gagal membuat anggota: %v", err)
	}
	if created.Nim != "105223046" || created.NamaLengkap != "Anton Dicaprio" {
		t.Errorf("Data anggota tidak sesuai: %+v", created)
	}

	// 2. Test Duplicate NIM Error
	_, err = uc.CreateAnggota(createReq)
	if err == nil {
		t.Errorf("Ekspektasi error duplikasi NIM, tetapi berhasil dibuat")
	}

	// 3. Test GetAnggotaByID
	found, err := uc.GetAnggotaByID(created.IdAnggota)
	if err != nil {
		t.Fatalf("Gagal mengambil anggota by ID: %v", err)
	}
	if found.IdAnggota != created.IdAnggota {
		t.Errorf("ID anggota tidak cocok: %d vs %d", found.IdAnggota, created.IdAnggota)
	}

	// 4. Test UpdateAnggota
	updateReq := entity.UpdateAnggotaRequest{
		NamaLengkap:  "Anton Dicaprio Updated",
		ProgramStudi: "Ilmu Komputer",
		Angkatan:     2023,
		Jabatan:      "Kepala Departemen Media Informasi",
		Status:       "Pengurus",
	}
	updated, err := uc.UpdateAnggota(created.IdAnggota, updateReq)
	if err != nil {
		t.Fatalf("Gagal update anggota: %v", err)
	}
	if updated.NamaLengkap != "Anton Dicaprio Updated" {
		t.Errorf("Nama lengkap tidak terupdate: %s", updated.NamaLengkap)
	}

	// 5. Test DeleteAnggota with active slot assigned (must fail)
	slotID := int64(10)
	updated.IdPengguna = &slotID
	repo.Update(updated)

	err = uc.DeleteAnggota(created.IdAnggota)
	if err == nil {
		t.Errorf("Ekspektasi error hapus anggota yang memegang slot aktif, namun sukses")
	}

	// 6. Test DeleteAnggota after unassigning slot (must succeed)
	updated.IdPengguna = nil
	repo.Update(updated)

	err = uc.DeleteAnggota(created.IdAnggota)
	if err != nil {
		t.Fatalf("Gagal menghapus anggota bebas slot: %v", err)
	}
}
