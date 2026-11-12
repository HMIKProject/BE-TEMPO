package usecase_test

import (
	"errors"
	"testing"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
)

type mockOrganisasiRepo struct {
	depts       map[int64]*entity.Departemen
	divs        map[int64]*entity.DivisiDepartemen
	nextDeptID  int64
	nextDivID   int64
	mockBPH     []entity.BPHMemberInfo
	mockKabinet *entity.Kabinet
}

func newMockOrganisasiRepo() *mockOrganisasiRepo {
	return &mockOrganisasiRepo{
		depts:      make(map[int64]*entity.Departemen),
		divs:       make(map[int64]*entity.DivisiDepartemen),
		nextDeptID: 1,
		nextDivID:  1,
		mockKabinet: &entity.Kabinet{
			IdKabinet:   1,
			NamaKabinet: "Kabinet Sinergi Progresif",
			TahunMulai:  2026,
			TahunSelesai: 2027,
			IsAktif:     true,
		},
		mockBPH: []entity.BPHMemberInfo{
			{
				IdAnggota:   1,
				NamaLengkap: "Ketua Himpunan Mock",
				Nim:         "105223001",
				Jabatan:     "Ketua Himpunan",
				Status:      "Pengurus",
			},
		},
	}
}

func (m *mockOrganisasiRepo) GetStrukturOrganisasi() (*entity.StrukturOrganisasiResponse, error) {
	var deptList []entity.OrganisasiDeptItem
	for _, d := range m.depts {
		var divList []entity.DivisiDepartemen
		for _, div := range m.divs {
			if div.IdDepartemen == d.IdDepartemen {
				divList = append(divList, *div)
			}
		}
		deptList = append(deptList, entity.OrganisasiDeptItem{
			IdDepartemen:   d.IdDepartemen,
			NamaDepartemen: d.NamaDepartemen,
			Deskripsi:      d.Deskripsi,
			Logo:           d.Logo,
			StatusAktif:    d.StatusAktif,
			Divisi:         divList,
			TotalDivisi:    int64(len(divList)),
		})
	}

	return &entity.StrukturOrganisasiResponse{
		Kabinet:    m.mockKabinet,
		BPH:        m.mockBPH,
		Departemen: deptList,
	}, nil
}

func (m *mockOrganisasiRepo) FindDepartemenByID(id int64) (*entity.Departemen, error) {
	dept, ok := m.depts[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return dept, nil
}

func (m *mockOrganisasiRepo) CreateDepartemen(dept *entity.Departemen) error {
	dept.IdDepartemen = m.nextDeptID
	m.nextDeptID++
	m.depts[dept.IdDepartemen] = dept
	return nil
}

func (m *mockOrganisasiRepo) UpdateDepartemen(dept *entity.Departemen) error {
	m.depts[dept.IdDepartemen] = dept
	return nil
}

func (m *mockOrganisasiRepo) ArchiveDepartemen(id int64) error {
	dept, ok := m.depts[id]
	if !ok {
		return errors.New("record not found")
	}
	dept.StatusAktif = false
	return nil
}

func (m *mockOrganisasiRepo) FindDivisiByID(id int64) (*entity.DivisiDepartemen, error) {
	div, ok := m.divs[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return div, nil
}

func (m *mockOrganisasiRepo) CreateDivisi(divisi *entity.DivisiDepartemen) error {
	divisi.IdDivisi = m.nextDivID
	m.nextDivID++
	m.divs[divisi.IdDivisi] = divisi
	return nil
}

func (m *mockOrganisasiRepo) UpdateDivisi(divisi *entity.DivisiDepartemen) error {
	m.divs[divisi.IdDivisi] = divisi
	return nil
}

func (m *mockOrganisasiRepo) ArchiveDivisi(id int64) error {
	div, ok := m.divs[id]
	if !ok {
		return errors.New("record not found")
	}
	div.StatusAktif = false
	return nil
}

func TestAdminOrganisasiUsecase(t *testing.T) {
	repo := newMockOrganisasiRepo()
	uc := usecase.NewAdminOrganisasiUsecase(repo)

	// 1. Test Create Departemen
	createDeptReq := entity.CreateDeptRequest{
		NamaDepartemen: "Riset dan Teknologi",
		Deskripsi:      "Departemen teknologi HMIK",
		Logo:           "https://cloudinary.com/ristek.png",
	}
	dept, err := uc.CreateDepartemen(createDeptReq)
	if err != nil {
		t.Fatalf("Gagal membuat departemen: %v", err)
	}
	if dept.NamaDepartemen != "Riset dan Teknologi" || !dept.StatusAktif {
		t.Errorf("Data departemen tidak sesuai: %+v", dept)
	}

	// Test Create Departemen with empty name (must fail)
	_, err = uc.CreateDepartemen(entity.CreateDeptRequest{NamaDepartemen: "  "})
	if err == nil {
		t.Errorf("Ekspektasi error nama departemen kosong, namun sukses")
	}

	// 2. Test Create Divisi under the department
	createDivReq := entity.CreateDivisiRequest{
		IdDepartemen: dept.IdDepartemen,
		NamaDivisi:   "Riset & Inovasi",
		Deskripsi:    "Divisi penelitian ilmiah",
	}
	div, err := uc.CreateDivisi(createDivReq)
	if err != nil {
		t.Fatalf("Gagal membuat divisi: %v", err)
	}
	if div.NamaDivisi != "Riset & Inovasi" || div.IdDepartemen != dept.IdDepartemen || !div.StatusAktif {
		t.Errorf("Data divisi tidak sesuai: %+v", div)
	}

	// Test Create Divisi with non-existent department (must fail)
	_, err = uc.CreateDivisi(entity.CreateDivisiRequest{IdDepartemen: 999, NamaDivisi: "Test"})
	if err == nil {
		t.Errorf("Ekspektasi error departemen tidak ada, namun sukses")
	}

	// 3. Test GetStrukturOrganisasi
	org, err := uc.GetStrukturOrganisasi()
	if err != nil {
		t.Fatalf("Gagal mengambil struktur organisasi: %v", err)
	}
	if org.Kabinet == nil || len(org.BPH) == 0 || len(org.Departemen) != 1 {
		t.Errorf("Struktur organisasi tidak lengkap: %+v", org)
	}
	if len(org.Departemen[0].Divisi) != 1 {
		t.Errorf("Divisi pada departemen tidak terdaftar: %+v", org.Departemen[0].Divisi)
	}

	// 4. Test Update Departemen
	statusFalse := false
	updatedDept, err := uc.UpdateDepartemen(dept.IdDepartemen, entity.UpdateDeptRequest{
		NamaDepartemen: "Riset dan Teknologi Baru",
		Deskripsi:      "Deskripsi baru",
		StatusAktif:    &statusFalse,
	})
	if err != nil {
		t.Fatalf("Gagal update departemen: %v", err)
	}
	if updatedDept.NamaDepartemen != "Riset dan Teknologi Baru" || updatedDept.StatusAktif {
		t.Errorf("Update departemen gagal: %+v", updatedDept)
	}

	// 5. Test Update Divisi
	updatedDiv, err := uc.UpdateDivisi(div.IdDivisi, entity.UpdateDivisiRequest{
		NamaDivisi:  "Riset & Inovasi Terapan",
		Deskripsi:   "Deskripsi terapan",
		StatusAktif: &statusFalse,
	})
	if err != nil {
		t.Fatalf("Gagal update divisi: %v", err)
	}
	if updatedDiv.NamaDivisi != "Riset & Inovasi Terapan" || updatedDiv.StatusAktif {
		t.Errorf("Update divisi gagal: %+v", updatedDiv)
	}

	// 6. Test Archive Departemen
	err = uc.ArchiveDepartemen(dept.IdDepartemen)
	if err != nil {
		t.Fatalf("Gagal arsipkan departemen: %v", err)
	}
	archivedDept, _ := repo.FindDepartemenByID(dept.IdDepartemen)
	if archivedDept.StatusAktif {
		t.Errorf("Status aktif departemen harusnya false")
	}
}
