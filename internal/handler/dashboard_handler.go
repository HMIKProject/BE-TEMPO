package handler

import (
	"net/http"

	_ "github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	usecase usecase.DashboardUsecase
}

func NewDashboardHandler(u usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{usecase: u}
}

// GetDashboardStats godoc
// @Summary      Ambil seluruh ringkasan kartu Dashboard Admin Halaman Utama
// @Description  Mengembalikan agregasi 4 kartu metrik (Total Anggota Aktif, Total Modul, Total Proker, Voting Aktif)
// @Tags         Dashboard Admin - Halaman Utama
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSendResponse{data=entity.DashboardStatsResponse}  "Berhasil mengambil data ringkasan dashboard"
// @Failure      401  {object}  utils.JSendResponse  "Akses Ditolak: Token JWT tidak ditemukan atau tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /admin/dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.usecase.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil data ringkasan dashboard"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Berhasil mengambil data ringkasan dashboard", stats))
}

// GetTotalAnggotaAktif godoc
// @Summary      Ambil detail kartu Total Anggota Aktif
// @Description  Mengembalikan data kartu biru: total anggota berstatus pengurus/magang beserta pertumbuhannya
// @Tags         Dashboard Admin - Halaman Utama
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSendResponse{data=entity.CardDetail}  "Berhasil mengambil detail kartu Total Anggota Aktif"
// @Failure      401  {object}  utils.JSendResponse  "Akses Ditolak: Token JWT tidak ditemukan atau tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /dashboard-admin-halaman-utama/total-anggota-aktif [get]
func (h *DashboardHandler) GetTotalAnggotaAktif(c *gin.Context) {
	card, err := h.usecase.GetTotalAnggotaAktif(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil data kartu Total Anggota Aktif"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data kartu Total Anggota Aktif berhasil diambil", card))
}

// GetTotalModul godoc
// @Summary      Ambil detail kartu Total Modul
// @Description  Mengembalikan data kartu kuning: total modul pembelajaran beserta pertumbuhannya
// @Tags         Dashboard Admin - Halaman Utama
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSendResponse{data=entity.CardDetail}  "Berhasil mengambil detail kartu Total Modul"
// @Failure      401  {object}  utils.JSendResponse  "Akses Ditolak: Token JWT tidak ditemukan atau tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /dashboard-admin-halaman-utama/total-modul [get]
func (h *DashboardHandler) GetTotalModul(c *gin.Context) {
	card, err := h.usecase.GetTotalModul(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil data kartu Total Modul"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data kartu Total Modul berhasil diambil", card))
}

// GetTotalProgramKerja godoc
// @Summary      Ambil detail kartu Total Program Kerja
// @Description  Mengembalikan data kartu hijau: total program kerja himpunan beserta pertumbuhannya
// @Tags         Dashboard Admin - Halaman Utama
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSendResponse{data=entity.CardDetail}  "Berhasil mengambil detail kartu Total Program Kerja"
// @Failure      401  {object}  utils.JSendResponse  "Akses Ditolak: Token JWT tidak ditemukan atau tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /dashboard-admin-halaman-utama/total-program-kerja [get]
func (h *DashboardHandler) GetTotalProgramKerja(c *gin.Context) {
	card, err := h.usecase.GetTotalProgramKerja(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil data kartu Total Program Kerja"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data kartu Total Program Kerja berhasil diambil", card))
}

// GetVotingAktif godoc
// @Summary      Ambil detail kartu Voting Aktif
// @Description  Mengembalikan data kartu oranye: total voting aktif dan kegiatan pemilihan yang sedang berjalan
// @Tags         Dashboard Admin - Halaman Utama
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.JSendResponse{data=entity.CardDetail}  "Berhasil mengambil detail kartu Voting Aktif"
// @Failure      401  {object}  utils.JSendResponse  "Akses Ditolak: Token JWT tidak ditemukan atau tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /dashboard-admin-halaman-utama/voting-aktif [get]
func (h *DashboardHandler) GetVotingAktif(c *gin.Context) {
	card, err := h.usecase.GetVotingAktif(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil data kartu Voting Aktif"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data kartu Voting Aktif berhasil diambil", card))
}