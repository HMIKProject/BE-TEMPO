package handler

import (
	"net/http"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: u}
}

// Login godoc
// @Summary      Login Pengurus HMIK & Deteksi Aktivasi
// @Description  Endpoint terpadu untuk login biasa atau aktivasi akun baru menggunakan Kode Aktivasi pada kolom sandi
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body entity.LoginRequest true "Kredensial Login atau Kode Aktivasi"
// @Success      200 {object} utils.JSendResponse{data=entity.LoginResponse}
// @Failure      400 {object} utils.JSendResponse
// @Failure      401 {object} utils.JSendResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.authUsecase.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	if res.RequirePasswordChange {
		c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kode aktivasi valid, silakan atur sandi baru", res))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Login berhasil", res))
}

// SetInitialPassword godoc
// @Summary      Aktivasi Akun & Pengaturan Kata Sandi Baru
// @Description  Endpoint untuk pengurus baru menetapkan kata sandi pribadi mandiri dari modal popup di halaman login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body entity.SetInitialPasswordRequest true "Payload Sandi Baru dan Tiket Aktivasi"
// @Success      200 {object} utils.JSendResponse{data=entity.LoginResponse}
// @Failure      400 {object} utils.JSendResponse
// @Failure      401 {object} utils.JSendResponse
// @Router       /auth/set-initial-password [post]
func (h *AuthHandler) SetInitialPassword(c *gin.Context) {
	var req entity.SetInitialPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.authUsecase.SetInitialPassword(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kata sandi berhasil diatur, akun Anda kini aktif", res))
}

// GetMe godoc
// @Summary      Ambil Profil Pengguna Saat Ini (Admin)
// @Description  Mengembalikan informasi profil pengguna yang sedang login berdasarkan token JWT
// @Tags         Admin - Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.UserProfileResponse} "Profil pengguna berhasil diambil"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token tidak valid"
// @Router       /admin/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	idVal, exists := c.Get("id_pengguna")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Sesi login tidak ditemukan", nil))
		return
	}

	var idPengguna int64
	switch v := idVal.(type) {
	case float64:
		idPengguna = int64(v)
	case int64:
		idPengguna = v
	case int:
		idPengguna = int64(v)
	default:
		c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Format ID pengguna tidak valid", nil))
		return
	}

	profile, err := h.authUsecase.GetMe(c.Request.Context(), idPengguna)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Profil pengguna berhasil diambil", profile))
}

// StudentRequestOTP godoc
// @Summary      Permintaan Kode OTP Mahasiswa (Login Tipe 2)
// @Description  Memvalidasi NIM (diawali 1052) & email kampus (@student.universitaspertamina.ac.id) lalu menerbitkan OTP
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body entity.StudentRequestOTPRequest true "Payload NIM dan Email Kampus"
// @Success      200 {object} utils.JSendResponse{data=entity.StudentRequestOTPResponse}
// @Failure      400 {object} utils.JSendResponse
// @Router       /auth/student/request-otp [post]
func (h *AuthHandler) StudentRequestOTP(c *gin.Context) {
	var req entity.StudentRequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format data tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.authUsecase.StudentRequestOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kode OTP berhasil dikirimkan ke email kampus Anda", res))
}

// StudentVerifyOTP godoc
// @Summary      Verifikasi Kode OTP Mahasiswa (Login Tipe 2)
// @Description  Memverifikasi 6 digit kode OTP dan menerbitkan token sesi JWT mahasiswa
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body entity.StudentVerifyOTPRequest true "Payload Verifikasi OTP"
// @Success      200 {object} utils.JSendResponse{data=entity.StudentVerifyOTPResponse}
// @Failure      400 {object} utils.JSendResponse
// @Router       /auth/student/verify-otp [post]
func (h *AuthHandler) StudentVerifyOTP(c *gin.Context) {
	var req entity.StudentVerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format data tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.authUsecase.StudentVerifyOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Verifikasi berhasil. Selamat datang di Portal Pembelajaran HMIK!", res))
}

// StudentResendOTP godoc
// @Summary      Kirim Ulang Kode OTP Mahasiswa
// @Description  Mengirim ulang OTP ke email mahasiswa jika belum kedaluwarsa
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body entity.StudentRequestOTPRequest true "Payload Kirim Ulang OTP"
// @Success      200 {object} utils.JSendResponse{data=entity.StudentRequestOTPResponse}
// @Failure      400 {object} utils.JSendResponse
// @Router       /auth/student/resend-otp [post]
func (h *AuthHandler) StudentResendOTP(c *gin.Context) {
	h.StudentRequestOTP(c)
}
