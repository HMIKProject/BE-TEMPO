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
}

type authUsecase struct {
	penggunaRepo repository.PenggunaRepository
}

func NewAuthUsecase(repo repository.PenggunaRepository) AuthUsecase {
	return &authUsecase{penggunaRepo: repo}
}

func (u *authUsecase) Login(req entity.LoginRequest) (entity.LoginResponse, error) {
	// 1. Cari pengguna di Database berdasarkan Email atau Username
	pengguna, err := u.penggunaRepo.FindByEmailOrUsername(req.Identifier)
	if err != nil {
		// Kita berikan pesan error generik agar keamanan terjaga (tidak memberi tahu email valid tapi password salah)
		return entity.LoginResponse{}, errors.New("Kredensial tidak valid")
	}

	// 2. Verifikasi Password menggunakan Bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(pengguna.Sandi), []byte(req.Password))
	if err != nil {
		return entity.LoginResponse{}, errors.New("Kredensial tidak valid")
	}

	// 3. Buat JWT Token
	token, err := u.generateJWT(pengguna)
	if err != nil {
		return entity.LoginResponse{}, errors.New("Gagal membuat sesi login")
	}

	// 4. Kembalikan Response Sukses
	return entity.LoginResponse{
		Token:       token,
		IDPengguna:  pengguna.IdPengguna,
		NamaLengkap: pengguna.NamaLengkap,
	}, nil
}

func (u *authUsecase) generateJWT(p *entity.Pengguna) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_untuk_local_development_saja" 
	}

	// Payload data yang dimasukkan ke dalam JWT
	claims := jwt.MapClaims{
		"id_pengguna": p.IdPengguna,
		"id_peran":    p.IdPeran,
		"exp":         time.Now().Add(24 * time.Hour).Unix(), // Masa berlaku 24 jam sesuai SOP
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
