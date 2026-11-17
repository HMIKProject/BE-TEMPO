package usecase

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(req entity.LoginRequest) (entity.LoginResponse, error)
	SetInitialPassword(req entity.SetInitialPasswordRequest) (entity.LoginResponse, error)
	GetMe(ctx context.Context, idPengguna int64) (*entity.UserProfileResponse, error)
	StudentRequestOTP(ctx context.Context, req entity.StudentRequestOTPRequest) (*entity.StudentRequestOTPResponse, error)
	StudentVerifyOTP(ctx context.Context, req entity.StudentVerifyOTPRequest) (*entity.StudentVerifyOTPResponse, error)
	StudentResendOTP(ctx context.Context, req entity.StudentRequestOTPRequest) (*entity.StudentRequestOTPResponse, error)
}

type authUsecase struct {
	penggunaRepo repository.PenggunaRepository
	slotRepo     repository.SlotRepository
}

func NewAuthUsecase(repo repository.PenggunaRepository, slotRepo repository.SlotRepository) AuthUsecase {
	return &authUsecase{
		penggunaRepo: repo,
		slotRepo:     slotRepo,
	}
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
	displayName, role, initials, _ := u.resolveProfileInfo(pengguna)
	return entity.LoginResponse{
		Token:                 token,
		IDPengguna:            pengguna.IdPengguna,
		NamaPengguna:          pengguna.NamaPengguna,
		NamaLengkap:           displayName,
		Email:                 pengguna.Email,
		Jabatan:               role,
		Role:                  role,
		Initials:              initials,
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

	displayName, role, initials, _ := u.resolveProfileInfo(pengguna)
	return entity.LoginResponse{
		Token:                 token,
		IDPengguna:            pengguna.IdPengguna,
		NamaPengguna:          pengguna.NamaPengguna,
		NamaLengkap:           displayName,
		Email:                 pengguna.Email,
		Jabatan:               role,
		Role:                  role,
		Initials:              initials,
		RequirePasswordChange: false,
	}, nil
}

func (u *authUsecase) resolveProfileInfo(pengguna *entity.Pengguna) (displayName string, role string, initials string, foto string) {
	displayName = pengguna.NamaLengkap
	if displayName == "" {
		displayName = pengguna.NamaPengguna
	}
	role = "Administrator"
	foto = pengguna.FotoPengguna

	if u.slotRepo != nil {
		borrower, err := u.slotRepo.FindBorrowerBySlotID(pengguna.IdPengguna)
		if err == nil && borrower != nil {
			if borrower.NamaLengkap != "" {
				displayName = borrower.NamaLengkap
			}
			if borrower.Jabatan != "" {
				role = borrower.Jabatan
			}
			if borrower.FotoProfil != "" {
				foto = borrower.FotoProfil
			}
		}
	}

	initials = utils.GenerateInitials(displayName)
	return displayName, role, initials, foto
}

func (u *authUsecase) GetMe(ctx context.Context, idPengguna int64) (*entity.UserProfileResponse, error) {
	pengguna, err := u.penggunaRepo.FindByID(idPengguna)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	displayName, role, initials, foto := u.resolveProfileInfo(pengguna)

	return &entity.UserProfileResponse{
		IDPengguna:   pengguna.IdPengguna,
		NamaPengguna: pengguna.NamaPengguna,
		NamaLengkap:  displayName,
		Email:        pengguna.Email,
		Role:         role,
		Jabatan:      role,
		Initials:     initials,
		Foto:         foto,
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

// StudentRequestOTP memproses permohonan OTP untuk login mahasiswa umum (Login Tipe 2)
func (u *authUsecase) StudentRequestOTP(ctx context.Context, req entity.StudentRequestOTPRequest) (*entity.StudentRequestOTPResponse, error) {
	nim := strings.TrimSpace(req.NIM)
	if !strings.HasPrefix(nim, "1052") {
		return nil, errors.New("NIM tidak valid. Hanya mahasiswa Program Studi Ilmu Komputer (diawali 1052) yang berhak mengakses layanan ini")
	}

	email := strings.ToLower(strings.TrimSpace(req.EmailKampus))
	if !strings.HasSuffix(email, "@student.universitaspertamina.ac.id") {
		return nil, errors.New("Gunakan alamat email resmi mahasiswa (@student.universitaspertamina.ac.id)")
	}

	otp, err := utils.GenerateCryptoOTP(6)
	if err != nil {
		otp = "123456"
	}
	expiredAt := time.Now().Add(5 * time.Minute)

	pengguna, err := u.penggunaRepo.FindByEmailOrUsername(nim)
	if err != nil {
		pengguna, err = u.penggunaRepo.FindByEmailOrUsername(email)
	}

	roleID := int64(3)
	if err != nil || pengguna == nil {
		newPengguna := &entity.Pengguna{
			NamaPengguna:        nim,
			Email:               email,
			NamaLengkap:         "Mahasiswa Ilmu Komputer",
			Sandi:               "otp_student_login",
			StatusAktif:         true,
			IdPeran:             &roleID,
			KodeAktivasi:        &otp,
			KodeAktivasiExpired: &expiredAt,
		}
		if err := u.penggunaRepo.Create(newPengguna); err != nil {
			return nil, errors.New("Gagal mendaftarkan sesi mahasiswa")
		}
	} else {
		// Cek cooldown 60 detik
		if pengguna.KodeAktivasiExpired != nil {
			createdAt := pengguna.KodeAktivasiExpired.Add(-5 * time.Minute)
			elapsed := time.Since(createdAt)
			if elapsed >= 0 && elapsed < 60*time.Second {
				sisa := int((60*time.Second - elapsed).Seconds()) + 1
				return &entity.StudentRequestOTPResponse{
					EmailKampus:  email,
					NIM:          nim,
					SisaCooldown: sisa,
					KodeOTP:      *pengguna.KodeAktivasi,
				}, nil
			}
		}

		pengguna.KodeAktivasi = &otp
		pengguna.KodeAktivasiExpired = &expiredAt
		if err := u.penggunaRepo.Update(pengguna); err != nil {
			return nil, errors.New("Gagal memperbarui kode OTP")
		}
	}

	// Kirim OTP via SMTP ke email kampus mahasiswa (Outlook)
	_ = utils.SendOTPEmailWithPurpose(email, nim, otp, "learning")

	log.Printf("📩 [STUDENT OTP] NIM: %s | Email: %s | Kode OTP: %s", nim, email, otp)

	return &entity.StudentRequestOTPResponse{
		EmailKampus:  email,
		NIM:          nim,
		SisaCooldown: 60,
		KodeOTP:      otp,
	}, nil
}

// StudentVerifyOTP memvalidasi kode OTP dan menerbitkan token sesi JWT mahasiswa (Login Tipe 2)
func (u *authUsecase) StudentVerifyOTP(ctx context.Context, req entity.StudentVerifyOTPRequest) (*entity.StudentVerifyOTPResponse, error) {
	nim := strings.TrimSpace(req.NIM)
	email := strings.ToLower(strings.TrimSpace(req.EmailKampus))

	if !strings.HasPrefix(nim, "1052") {
		return nil, errors.New("NIM tidak valid. Hanya mahasiswa Program Studi Ilmu Komputer (diawali 1052) yang diizinkan")
	}

	pengguna, err := u.penggunaRepo.FindByEmailOrUsername(nim)
	if err != nil {
		pengguna, err = u.penggunaRepo.FindByEmailOrUsername(email)
		if err != nil {
			return nil, errors.New("Data mahasiswa belum terdaftar. Silakan minta kode OTP terlebih dahulu")
		}
	}

	// Master OTP "000000" selalu diizinkan untuk kemudahan dev & testing
	isValid := false
	if req.KodeOTP == "000000" {
		isValid = true
	} else if pengguna.KodeAktivasi != nil && *pengguna.KodeAktivasi == req.KodeOTP {
		if pengguna.KodeAktivasiExpired == nil || time.Now().Before(*pengguna.KodeAktivasiExpired) {
			isValid = true
		}
	}

	if !isValid {
		return nil, errors.New("Kode OTP tidak valid atau sudah kedaluwarsa")
	}

	// Hanguskan OTP
	pengguna.KodeAktivasi = nil
	pengguna.KodeAktivasiExpired = nil
	_ = u.penggunaRepo.Update(pengguna)

	token, err := u.generateStudentJWT(pengguna, nim)
	if err != nil {
		return nil, errors.New("Gagal menerbitkan sesi mahasiswa")
	}

	displayName := pengguna.NamaLengkap
	if displayName == "" {
		displayName = "Mahasiswa Ilmu Komputer"
	}

	return &entity.StudentVerifyOTPResponse{
		Token:       token,
		NIM:         nim,
		EmailKampus: pengguna.Email,
		Role:        "Mahasiswa",
		NamaLengkap: displayName,
	}, nil
}

func (u *authUsecase) StudentResendOTP(ctx context.Context, req entity.StudentRequestOTPRequest) (*entity.StudentRequestOTPResponse, error) {
	return u.StudentRequestOTP(ctx, req)
}

func (u *authUsecase) generateStudentJWT(p *entity.Pengguna, nim string) (string, error) {
	claims := jwt.MapClaims{
		"id_pengguna": p.IdPengguna,
		"nim":         nim,
		"email":       p.Email,
		"role":        "Mahasiswa",
		"exp":         time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getJWTSecret()))
}
