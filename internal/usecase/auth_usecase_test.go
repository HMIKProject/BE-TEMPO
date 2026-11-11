package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

type mockPenggunaRepo struct {
	byEmailOrUsername map[string]*entity.Pengguna
	byID              map[int64]*entity.Pengguna
}

func newMockPenggunaRepo() *mockPenggunaRepo {
	return &mockPenggunaRepo{
		byEmailOrUsername: make(map[string]*entity.Pengguna),
		byID:              make(map[int64]*entity.Pengguna),
	}
}

func (m *mockPenggunaRepo) add(p *entity.Pengguna) {
	m.byEmailOrUsername[p.Email] = p
	m.byEmailOrUsername[p.NamaPengguna] = p
	m.byID[p.IdPengguna] = p
}

func (m *mockPenggunaRepo) FindByEmailOrUsername(identifier string) (*entity.Pengguna, error) {
	p, ok := m.byEmailOrUsername[identifier]
	if !ok {
		return nil, errors.New("record not found")
	}
	return p, nil
}

func (m *mockPenggunaRepo) FindByID(id int64) (*entity.Pengguna, error) {
	p, ok := m.byID[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return p, nil
}

func (m *mockPenggunaRepo) Update(pengguna *entity.Pengguna) error {
	m.add(pengguna)
	return nil
}

func TestAuthUsecase_LoginAndActivation(t *testing.T) {
	repo := newMockPenggunaRepo()
	uc := usecase.NewAuthUsecase(repo)

	// Hash password standard
	defaultPass := "qwerqwer"
	hashedDefault, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)

	// 1. Akun Standar Aktif (Normal Pengurus)
	userNormal := &entity.Pengguna{
		IdPengguna:   1,
		NamaPengguna: "hmik.kahim.01",
		Email:        "hmik.kahim.01@hmik.com",
		Sandi:        string(hashedDefault),
		StatusAktif:  true,
		TokenVersion: 1,
	}
	repo.add(userNormal)

	// 2. Akun Baru dengan Kode Aktivasi
	actCode := "ACT-8829-KL90"
	validExp := time.Now().Add(7 * 24 * time.Hour)
	randomHashed, _ := bcrypt.GenerateFromPassword([]byte("random_temp_pass"), bcrypt.DefaultCost)
	userAktivasi := &entity.Pengguna{
		IdPengguna:          2,
		NamaPengguna:        "hmik.humas.01",
		Email:               "hmik.humas.01@hmik.com",
		Sandi:               string(randomHashed),
		KodeAktivasi:        &actCode,
		KodeAktivasiExpired: &validExp,
		StatusAktif:         true,
		TokenVersion:        1,
	}
	repo.add(userAktivasi)

	t.Run("Login Normal: Berhasil", func(t *testing.T) {
		res, err := uc.Login(entity.LoginRequest{
			Identifier: "hmik.kahim.01@hmik.com",
			Password:   "qwerqwer",
		})
		if err != nil {
			t.Fatalf("Ekspektasi login sukses, error: %v", err)
		}
		if res.RequirePasswordChange {
			t.Errorf("Ekspektasi RequirePasswordChange false, dapat true")
		}
		if res.Token == "" {
			t.Errorf("Ekspektasi token terisi")
		}
	})

	t.Run("Login Normal: Password Salah", func(t *testing.T) {
		_, err := uc.Login(entity.LoginRequest{
			Identifier: "hmik.kahim.01@hmik.com",
			Password:   "password_salah",
		})
		if err == nil {
			t.Fatalf("Ekspektasi error kredensial tidak valid")
		}
		if err.Error() != "Email atau kata sandi tidak valid" {
			t.Errorf("Ekspektasi pesan error generik, dapat: %s", err.Error())
		}
	})

	var setupToken string
	t.Run("Aktivasi Akun: Login dengan Kode Aktivasi Berhasil Memicu Popup", func(t *testing.T) {
		res, err := uc.Login(entity.LoginRequest{
			Identifier: "hmik.humas.01@hmik.com",
			Password:   "ACT-8829-KL90",
		})
		if err != nil {
			t.Fatalf("Ekspektasi aktivasi terdeteksi, error: %v", err)
		}
		if !res.RequirePasswordChange {
			t.Errorf("Ekspektasi RequirePasswordChange true")
		}
		if res.SetupToken == "" {
			t.Errorf("Ekspektasi SetupToken terisi")
		}
		setupToken = res.SetupToken
	})

	t.Run("Aktivasi Akun: Kode Aktivasi Salah", func(t *testing.T) {
		_, err := uc.Login(entity.LoginRequest{
			Identifier: "hmik.humas.01@hmik.com",
			Password:   "ACT-SALAH-0000",
		})
		if err == nil {
			t.Fatalf("Ekspektasi error kredensial tidak valid")
		}
		if err.Error() != "Email atau kata sandi tidak valid" {
			t.Errorf("Ekspektasi pesan generik, dapat: %s", err.Error())
		}
	})

	t.Run("SetInitialPassword: Ganti Sandi Mandiri Berhasil", func(t *testing.T) {
		newPassword := "MyPersonalSecret123!"
		res, err := uc.SetInitialPassword(entity.SetInitialPasswordRequest{
			SetupToken:      setupToken,
			PasswordBaru:    newPassword,
			KonfirmasiSandi: newPassword,
		})
		if err != nil {
			t.Fatalf("Ekspektasi ganti sandi awal sukses, error: %v", err)
		}
		if res.Token == "" {
			t.Errorf("Ekspektasi login resmi diterbitkan langsung")
		}

		// Verifikasi di DB: kode aktivasi harus sudah NULL
		updatedUser, _ := repo.FindByID(2)
		if updatedUser.KodeAktivasi != nil {
			t.Errorf("Ekspektasi kode aktivasi di-NULL-kan setelah ganti sandi")
		}

		// Login ulang dengan sandi baru harus berhasil
		loginNew, err := uc.Login(entity.LoginRequest{
			Identifier: "hmik.humas.01@hmik.com",
			Password:   newPassword,
		})
		if err != nil {
			t.Fatalf("Login dengan sandi baru gagal: %v", err)
		}
		if loginNew.RequirePasswordChange {
			t.Errorf("Akun seharusnya sudah aktif permanen")
		}

		// Login ulang dengan kode aktivasi lama harus gagal!
		_, errOld := uc.Login(entity.LoginRequest{
			Identifier: "hmik.humas.01@hmik.com",
			Password:   "ACT-8829-KL90",
		})
		if errOld == nil {
			t.Errorf("Ekspektasi kode aktivasi lama ditolak setelah password baru dibuat")
		}
	})
}
