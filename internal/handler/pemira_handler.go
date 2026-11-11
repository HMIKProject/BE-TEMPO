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

// ResetKotakSuara godoc
// @Summary      Kosongkan Kotak Suara & Reset Pemilu (Admin)
// @Description  Menghapus seluruh surat suara di kotak suara dan mereset status hak pilih seluruh DPT ke 0 (Uji Coba / Persiapan Pemilu)
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse "Kotak suara berhasil dikosongkan dan data pemilu di-reset ke 0"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Failure      500 {object} utils.JSendResponse "Gagal mengosongkan kotak suara"
// @Router       /admin/pemira/reset-kotak-suara [post]
func (h *PemiraHandler) ResetKotakSuara(c *gin.Context) {
	if err := h.usecase.ResetKotakSuara(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kotak suara berhasil dikosongkan dan seluruh data suara di-reset ke 0", nil))
}

// CreateCandidate godoc
// @Summary      Tambah Pasangan Calon Baru (Admin)
// @Description  Mendaftarkan paslon baru beserta visi, misi, dan informasi ketua/wakil
// @Tags         Admin - Pemira
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.CreatePaslonRequest true "Data Paslon Baru"
// @Success      201 {object} utils.JSendResponse{data=entity.Paslon} "Pasangan calon berhasil ditambahkan"
// @Failure      400 {object} utils.JSendResponse "Format request tidak valid atau data gagal disimpan"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/kandidat [post]
func (h *PemiraHandler) CreateCandidate(c *gin.Context) {
	var req entity.CreatePaslonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.CreateCandidate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Pasangan calon berhasil ditambahkan", res))
}

// UpdateCandidate godoc
// @Summary      Perbarui Pasangan Calon (Admin)
// @Description  Memperbarui data paslon yang sudah ada berdasarkan ID
// @Tags         Admin - Pemira
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Paslon"
// @Param        request body entity.UpdatePaslonRequest true "Data Paslon Diperbarui"
// @Success      200 {object} utils.JSendResponse{data=entity.Paslon} "Data pasangan calon berhasil diperbarui"
// @Failure      400 {object} utils.JSendResponse "ID tidak valid atau update gagal"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/kandidat/{id} [put]
func (h *PemiraHandler) UpdateCandidate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID paslon tidak valid", nil))
		return
	}

	var req entity.UpdatePaslonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.UpdateCandidate(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data pasangan calon berhasil diperbarui", res))
}

// DeleteCandidate godoc
// @Summary      Hapus Pasangan Calon (Admin)
// @Description  Menghapus data paslon (Dilindungi: jika suara sudah masuk, penghapusan ditolak demi LUBER JURDIL)
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Paslon"
// @Success      200 {object} utils.JSendResponse "Pasangan calon berhasil dihapus"
// @Failure      400 {object} utils.JSendResponse "ID tidak valid atau suara sudah masuk sehingga dilarang dihapus"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/kandidat/{id} [delete]
func (h *PemiraHandler) DeleteCandidate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID paslon tidak valid", nil))
		return
	}

	if err := h.usecase.DeleteCandidate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Pasangan calon berhasil dihapus", nil))
}

// GetElectionSettings godoc
// @Summary      Ambil Pengaturan & Parameter Acara Pemira (Admin)
// @Description  Mengembalikan konfigurasi status, periode waktu, jadwal, serta status MPIK
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.ElectionSettingsResponse} "Pengaturan pemira berhasil diambil"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Failure      500 {object} utils.JSendResponse "Gagal memuat pengaturan pemira"
// @Router       /admin/pemira/settings [get]
func (h *PemiraHandler) GetElectionSettings(c *gin.Context) {
	settings, err := h.usecase.GetElectionSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Pengaturan pemira berhasil diambil", settings))
}

// UpdateElectionSchedule godoc
// @Summary      Perbarui Jadwal Acara Pemira (Admin)
// @Description  Memperbarui tanggal mulai, tanggal selesai, dan tahun periode pemilihan raya
// @Tags         Admin - Pemira
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.UpdatePemiraScheduleRequest true "Data Jadwal Baru"
// @Success      200 {object} utils.JSendResponse "Jadwal pemilihan raya berhasil diperbarui"
// @Failure      400 {object} utils.JSendResponse "Format jadwal tidak valid"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/settings [put]
func (h *PemiraHandler) UpdateElectionSchedule(c *gin.Context) {
	var req entity.UpdatePemiraScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	if err := h.usecase.UpdateElectionSchedule(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Jadwal pemilihan raya berhasil diperbarui", nil))
}

// UpdateElectionStatus godoc
// @Summary      Ubah Status Pemira / Saklar Pemicu (Admin)
// @Description  Mengubah status pelaksanaan pemilihan raya (DRAFT, BERJALAN, SELESAI)
// @Tags         Admin - Pemira
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.UpdatePemiraStatusRequest true "Status Baru"
// @Success      200 {object} utils.JSendResponse "Status pemilihan raya berhasil diubah"
// @Failure      400 {object} utils.JSendResponse "Status tidak valid"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/status [post]
func (h *PemiraHandler) UpdateElectionStatus(c *gin.Context) {
	var req entity.UpdatePemiraStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	if err := h.usecase.UpdateElectionStatus(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Status pemilihan raya berhasil diubah", nil))
}

// PelantikanMPIK godoc
// @Summary      Pengesahan & Pelantikan MPIK (Admin)
// @Description  Mengesahkan paslon pemenang pemira, membentuk kabinet aktif baru, dan memperbarui BPH di manajemen himpunan
// @Tags         Admin - Pemira
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.PelantikanMPIKResponse} "Pelantikan MPIK berhasil disahkan"
// @Failure      400 {object} utils.JSendResponse "Pemira belum selesai atau data pemenang tidak lengkap"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/pelantikan-mpik [post]
func (h *PemiraHandler) PelantikanMPIK(c *gin.Context) {
	res, err := h.usecase.PelantikanMPIK(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Pelantikan MPIK berhasil disahkan", res))
}

// GetPublicSchedule godoc
// @Summary      Ambil Jadwal & Status Bilik Suara Publik
// @Description  Mengembalikan informasi jadwal buka/tutup dan status bilik suara untuk ditampilkan di portal pemilih
// @Tags         Pemira
// @Produce      json
// @Success      200 {object} utils.JSendResponse{data=entity.PublicElectionScheduleResponse} "Jadwal pemira publik berhasil diambil"
// @Failure      500 {object} utils.JSendResponse "Gagal memuat jadwal pemira"
// @Router       /pemira/jadwal [get]
func (h *PemiraHandler) GetPublicSchedule(c *gin.Context) {
	schedule, err := h.usecase.GetPublicSchedule(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Jadwal pemira berhasil diambil", schedule))
}

// CreateNewPeriod godoc
// @Summary      Buka Periode Pemira Baru di Masa Depan (Admin)
// @Description  Membuka periode pemilihan raya baru (misal: 2027/2028) saat periode sebelumnya telah selesai
// @Tags         Admin - Pemira
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.CreatePemiraPeriodRequest true "Data Periode Baru"
// @Success      201 {object} utils.JSendResponse{data=entity.Pemira} "Periode pemira baru berhasil dibuat"
// @Failure      400 {object} utils.JSendResponse "Format request tidak valid atau pemira saat ini masih berjalan"
// @Failure      401 {object} utils.JSendResponse "Akses Ditolak: Token admin tidak valid"
// @Router       /admin/pemira/new-period [post]
func (h *PemiraHandler) CreateNewPeriod(c *gin.Context) {
	var req entity.CreatePemiraPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format request tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.usecase.CreateNewPeriod(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Periode pemilihan raya baru berhasil dibuat", res))
}

