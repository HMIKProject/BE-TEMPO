package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(req entity.LoginRequest) (entity.LoginResponse, error)
	SetInitialPassword(req entity.SetInitialPasswordRequest) (entity.LoginResponse, error)
}

type authUsecase struct {
	penggunaRepo repository.PenggunaRepository
}

func NewAuthUsecase(repo repository.PenggunaRepository) AuthUsecase {
	return &authUsecase{penggunaRepo: repo}
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_untuk_local_development_saja"
	}
	return secret
}

func (u *authUsecase) Login(req entity.LoginRequest) (entity.LoginResponse, error) {
	// 1. Cari pengguna di Database berdasarkan Email atau Username
	pengguna, err := u.penggunaRepo.FindByEmailOrUsername(req.Identifier)
	if err != nil {
		// Pesan error generik agar tidak membocorkan keberadaan akun
		return entity.LoginResponse{}, errors.New("Email atau kata sandi tidak valid")
	}

	if !pengguna.StatusAktif {
		return entity.LoginResponse{}, errors.New("Akun telah dinonaktifkan, silakan hubungi pengurus BPH")
	}

	// 2. Cek apakah akun memiliki Kode Aktivasi aktif (Akun baru / baru di-reset)
	if pengguna.KodeAktivasi != nil && *pengguna.KodeAktivasi != "" {
		// Validasi batas kedaluwarsa (7 hari)
		if pengguna.KodeAktivasiExpired != nil && time.Now().After(*pengguna.KodeAktivasiExpired) {
			return entity.LoginResponse{}, errors.New("Email atau kata sandi tidak valid")
		}

		// Cocokkan input password dengan kode aktivasi
		if req.Password == *pengguna.KodeAktivasi {
			// Terbitkan tiket sementara (setup token) untuk ganti sandi awal
			setupToken, err := u.generateSetupJWT(pengguna)
			if err != nil {
				return entity.LoginResponse{}, errors.New("Gagal memproses tiket aktivasi akun")
			}

			return entity.LoginResponse{
				IDPengguna:            pengguna.IdPengguna,
				NamaPengguna:          pengguna.NamaPengguna,
				RequirePasswordChange: true,
				SetupToken:            setupToken,
				Email:                 pengguna.Email,
				Jabatan:               pengguna.NamaLengkap,
			}, nil
		}

		// Jika tidak cocok dengan kode aktivasi
		return entity.LoginResponse{}, errors.New("Email atau kata sandi tidak valid")
	}

	// 3. Login Normal: Verifikasi Password menggunakan Bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(pengguna.Sandi), []byte(req.Password))
	if err != nil {
		return entity.LoginResponse{}, errors.New("Email atau kata sandi tidak valid")
	}

	// 4. Buat JWT Token Resmi
	token, err := u.generateJWT(pengguna)
	if err != nil {
		return entity.LoginResponse{}, errors.New("Gagal membuat sesi login")
	}

	// 5. Kembalikan Response Sukses Login Penuh
	return entity.LoginResponse{
		Token:                 token,
		IDPengguna:            pengguna.IdPengguna,
		NamaPengguna:          pengguna.NamaPengguna,
		RequirePasswordChange: false,
	}, nil
}

func (u *authUsecase) SetInitialPassword(req entity.SetInitialPasswordRequest) (entity.LoginResponse, error) {
	if req.PasswordBaru != req.KonfirmasiSandi {
		return entity.LoginResponse{}, errors.New("Konfirmasi kata sandi tidak cocok")
	}

	// 1. Verifikasi Setup Token
	idPengguna, err := u.validateSetupJWT(req.SetupToken)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	// 2. Ambil data pengguna
	pengguna, err := u.penggunaRepo.FindByID(idPengguna)
	if err != nil {
		return entity.LoginResponse{}, errors.New("Pengguna tidak ditemukan")
	}

	// 3. Hash sandi baru mandiri dengan Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		return entity.LoginResponse{}, errors.New("Gagal mengenkripsi kata sandi baru")
	}

	// 4. Update data pengguna: hanguskan kode aktivasi & aktifkan password baru
	pengguna.Sandi = string(hashedPassword)
	pengguna.KodeAktivasi = nil
	pengguna.KodeAktivasiExpired = nil
	pengguna.IsFirstLogin = false
	pengguna.TokenVersion++

	if err := u.penggunaRepo.Update(pengguna); err != nil {
		return entity.LoginResponse{}, errors.New("Gagal menyimpan kata sandi baru ke database")
	}

	// 5. Terbitkan Token JWT Resmi agar pengurus langsung masuk ke dashboard
	token, err := u.generateJWT(pengguna)
	if err != nil {
		return entity.LoginResponse{}, errors.New("Gagal membuat sesi login setelah aktivasi")
	}

	return entity.LoginResponse{
		Token:                 token,
		IDPengguna:            pengguna.IdPengguna,
		NamaPengguna:          pengguna.NamaPengguna,
		RequirePasswordChange: false,
	}, nil
}

func (u *authUsecase) generateJWT(p *entity.Pengguna) (string, error) {
	secret := getJWTSecret()

	claims := jwt.MapClaims{
		"id_pengguna":   p.IdPengguna,
		"id_peran":      p.IdPeran,
		"token_version": p.TokenVersion,
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *authUsecase) generateSetupJWT(p *entity.Pengguna) (string, error) {
	secret := getJWTSecret()

	claims := jwt.MapClaims{
		"id_pengguna": p.IdPengguna,
		"purpose":     "initial_setup",
		"exp":         time.Now().Add(15 * time.Minute).Unix(), // 15 menit untuk menyelesaikan form popup
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *authUsecase) validateSetupJWT(tokenString string) (int64, error) {
	secret := getJWTSecret()

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Metode penandatanganan token tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("Tiket aktivasi tidak valid atau telah kedaluwarsa, silakan ulangi login")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["purpose"] != "initial_setup" {
		return 0, errors.New("Tiket tidak valid untuk pengaturan sandi awal")
	}

	idFloat, ok := claims["id_pengguna"].(float64)
	if !ok {
		return 0, errors.New("Data pengguna tidak valid pada tiket")
	}

	return int64(idFloat), nil
}
