package handler

import (
	"net/http"
	"strconv"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PemiraHandler struct {
	usecase usecase.PemiraUsecase
}

func NewPemiraHandler(u usecase.PemiraUsecase) *PemiraHandler {
	return &PemiraHandler{usecase: u}
}

// RequestOTP godoc
// @Summary      Permintaan Kode OTP Pemilih
// @Description  Memvalidasi NIM dan email kampus mahasiswa terdaftar di DPT, kemudian mengirimkan kode 6-digit OTP
// @Tags         Pemira
// @Accept       json
// @Produce      json
// @Param        request body entity.RequestOTPRequest true "Kredensial NIM & Email Kampus"
// @Success      200 {object} utils.JSendResponse{data=entity.RequestOTPResponse} "OTP berhasil dikirim"
// @Failure      400 {object} utils.JSendResponse "Validasi gagal atau akun tidak memenuhi syarat"
// @Router       /pemira/auth/request-otp [post]
func (h *PemiraHandler) RequestOTP(c *gin.Context) {
	var req entity.RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.RequestOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kode OTP berhasil dikirimkan ke email kampus Anda", res))
}

// VerifyOTP godoc
// @Summary      Verifikasi Kode OTP & Terbitkan Ballot Token
// @Description  Memvalidasi kode OTP 6-digit pemilih dan mengembalikan Ballot JWT Token (15 menit) untuk mengakses bilik suara
// @Tags         Pemira
// @Accept       json
// @Produce      json
// @Param        request body entity.VerifyOTPRequest true "NIM dan Kode OTP"
// @Success      200 {object} utils.JSendResponse{data=entity.VerifyOTPResponse} "Verifikasi berhasil, token diterbitkan"
// @Failure      400 {object} utils.JSendResponse "Kode OTP salah, expired, atau format tidak valid"
// @Router       /pemira/auth/verify-otp [post]
func (h *PemiraHandler) VerifyOTP(c *gin.Context) {
	var req entity.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Verifikasi OTP berhasil, sesi bilik suara aktif", res))
}

// ResendOTP godoc
// @Summary      Kirim Ulang Kode OTP Pemilih
// @Description  Mengirimkan kembali kode OTP baru jika jeda cooldown 60 detik telah terpenuhi
// @Tags         Pemira
// @Accept       json
// @Produce      json
// @Param        request body entity.ResendOTPRequest true "NIM Mahasiswa"
// @Success      200 {object} utils.JSendResponse{data=entity.RequestOTPResponse} "Kode OTP baru berhasil dikirim"
// @Failure      400 {object} utils.JSendResponse "Masih dalam masa cooldown atau data tidak valid"
// @Router       /pemira/auth/resend-otp [post]
func (h *PemiraHandler) ResendOTP(c *gin.Context) {
	var req entity.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.ResendOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kode OTP baru berhasil dikirimkan ke email kampus Anda", res))
}

// GetActiveCandidates godoc
// @Summary      Ambil Daftar Paslon Aktif Pemira
// @Description  Mengambil seluruh kandidat pasangan calon ketua & wakil ketua pada periode pemira aktif
// @Tags         Pemira
// @Produce      json
// @Success      200 {object} utils.JSendResponse{data=[]entity.Paslon} "Daftar paslon aktif berhasil diambil"
// @Failure      500 {object} utils.JSendResponse "Gagal mengambil data paslon"
// @Router       /pemira/kandidat [get]
func (h *PemiraHandler) GetActiveCandidates(c *gin.Context) {
	candidates, err := h.usecase.GetActiveCandidates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Daftar paslon aktif berhasil diambil", candidates))
}

// CastVote godoc
// @Summary      Pemberian Suara di Bilik Suara (Cast Vote)
// @Description  Mencatat suara pemilih secara atomik dan anonim (Anti Double-Voting) menggunakan Ballot Token
// @Tags         Pemira
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.CastVoteRequest true "ID Paslon yang dipilih"
// @Success      201 {object} utils.JSendResponse{data=entity.CastVoteResponse} "Suara berhasil dicatat secara sah dan anonim"
// @Failure      400 {object} utils.JSendResponse "Hak suara sudah pernah digunakan atau paslon tidak valid"
// @Failure      401 {object} utils.JSendResponse "Ballot token tidak valid atau sesi kadaluarsa"
// @Failure      500 {object} utils.JSendResponse "Gagal memproses transaksi suara"
// @Router       /pemira/vote [post]
func (h *PemiraHandler) CastVote(c *gin.Context) {
	nimVal, exists := c.Get("nim")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("NIM tidak teridentifikasi pada token", nil))
		return
	}
	nim, _ := nimVal.(string)

	idPemiraVal, exists := c.Get("id_pemira")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("ID Pemira tidak teridentifikasi pada token", nil))
		return
	}
	idPemira, _ := idPemiraVal.(int64)

	var req entity.CastVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.CastVote(c.Request.Context(), nim, idPemira, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Suara Anda berhasil dicatat secara sah dan anonim", res))
}

// GetElectionStats godoc
// @Summary      Ambil Ringkasan Metrik Pemilu (Admin)
// @Description  Mengembalikan agregat 4 kartu metrik: Suara Masuk, Total DPT, Persentase Partisipasi, Total Paslon
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.ElectionStatsResponse} "Statistik pemilu berhasil diambil"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Failure      500 {object} utils.JSendResponse "Gagal mengambil data statistik"
// @Router       /admin/pemira/stats [get]
func (h *PemiraHandler) GetElectionStats(c *gin.Context) {
	stats, err := h.usecase.GetElectionStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Statistik pemilu berhasil diambil", stats))
}

// GetElectionChart godoc
// @Summary      Ambil Data Grafik Perolehan Suara Pemilu (Admin)
// @Description  Mengembalikan data perolehan suara per paslon secara realtime untuk visualisasi grafik batang pleno
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.ElectionChartResponse} "Data grafik pemilu berhasil diambil"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Failure      500 {object} utils.JSendResponse "Gagal mengambil data grafik pemilu"
// @Router       /admin/pemira/chart [get]
func (h *PemiraHandler) GetElectionChart(c *gin.Context) {
	chart, err := h.usecase.GetElectionChart(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data grafik pleno pemilu berhasil diambil", chart))
}

// GetVoters godoc
// @Summary      Ambil Daftar Pemilih DPT (Admin)
// @Description  Mengembalikan daftar pemilih tetap dengan pagination, filter pencarian (nama/nim), dan filter status memilih
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Param        search query string false "Cari berdasarkan nama atau NIM"
// @Param        status_memilih query bool false "Filter status memilih (true/false)"
// @Param        limit query int false "Batas jumlah data per halaman (default: 20)"
// @Param        offset query int false "Offset pagination data"
// @Success      200 {object} utils.JSendResponse{data=entity.VoterListResponse} "Daftar pemilih DPT berhasil diambil"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Failure      500 {object} utils.JSendResponse "Gagal mengambil data pemilih"
// @Router       /admin/pemira/pemilih [get]
func (h *PemiraHandler) GetVoters(c *gin.Context) {
	search := c.Query("search")
	var statusPtr *bool
	if statusStr := c.Query("status_memilih"); statusStr != "" {
		b := statusStr == "true" || statusStr == "1"
		statusPtr = &b
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	voters, err := h.usecase.GetVoters(c.Request.Context(), search, statusPtr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Daftar pemilih DPT berhasil diambil", voters))
}

// ResetVoter godoc
// @Summary      Reset Status Hak Pilih Pemilih (Admin)
// @Description  Mengembalikan status hak pilih mahasiswa menjadi belum memilih jika terjadi kendala force majeure
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Pemilih di tabel DPT"
// @Success      200 {object} utils.JSendResponse "Status hak pilih pemilih berhasil di-reset"
// @Failure      400 {object} utils.JSendResponse "ID pemilih tidak valid atau gagal melakukan reset"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Failure      500 {object} utils.JSendResponse "Gagal mereset status pemilih"
// @Router       /admin/pemira/pemilih/{id}/reset [post]
func (h *PemiraHandler) ResetVoter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID pemilih tidak valid", nil))
		return
	}

	if err := h.usecase.ResetVoter(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Status hak pilih berhasil di-reset menjadi belum memilih", nil))
}
