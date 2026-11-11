package handler

import (
	"net/http"
	"strconv"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminBIAHandler struct {
	biaUsecase usecase.AdminBIAUsecase
}

func NewAdminBIAHandler(u usecase.AdminBIAUsecase) *AdminBIAHandler {
	return &AdminBIAHandler{biaUsecase: u}
}

// GetAnggotaList godoc
// @Summary      Daftar Anggota Buku Induk (BIA)
// @Description  Mengambil daftar anggota Buku Induk dengan penyaringan nama/NIM, status, angkatan, serta pagination
// @Tags         Admin - Buku Induk Anggota
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        search query string false "Cari berdasarkan nama atau NIM"
// @Param        status query string false "Filter status (Pengurus, Magang, Demisioner, Semua)"
// @Param        angkatan query int false "Filter tahun angkatan"
// @Param        page query int false "Nomor halaman (default 1)"
// @Param        limit query int false "Jumlah data per halaman (default 50)"
// @Success      200 {object} utils.JSendResponse{data=entity.AnggotaListResponse}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/anggota [get]
func (h *AdminBIAHandler) GetAnggotaList(c *gin.Context) {
	var filter entity.AnggotaFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Parameter query filter tidak valid", []string{err.Error()}))
		return
	}

	res, err := h.biaUsecase.GetAnggotaList(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil daftar anggota: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Daftar anggota berhasil dimuat", res))
}

// GetAnggotaByID godoc
// @Summary      Detail Anggota Buku Induk
// @Description  Mengambil detail profil lengkap seorang anggota berdasarkan ID primary key
// @Tags         Admin - Buku Induk Anggota
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Anggota"
// @Success      200 {object} utils.JSendResponse{data=entity.Anggota}
// @Failure      400 {object} utils.JSendResponse
// @Failure      404 {object} utils.JSendResponse
// @Router       /admin/anggota/{id} [get]
func (h *AdminBIAHandler) GetAnggotaByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID anggota tidak valid", nil))
		return
	}

	res, err := h.biaUsecase.GetAnggotaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Detail anggota berhasil dimuat", res))
}

// CreateAnggota godoc
// @Summary      Tambah Anggota Baru
// @Description  Menambahkan data mahasiswa baru secara manual ke Buku Induk Anggota
// @Tags         Admin - Buku Induk Anggota
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.CreateAnggotaRequest true "Payload Tambah Anggota"
// @Success      201 {object} utils.JSendResponse{data=entity.Anggota}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/anggota [post]
func (h *AdminBIAHandler) CreateAnggota(c *gin.Context) {
	var req entity.CreateAnggotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload tambah anggota tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.biaUsecase.CreateAnggota(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Anggota baru berhasil ditambahkan", res))
}

// UpdateAnggota godoc
// @Summary      Perbarui Data Anggota
// @Description  Memperbarui data profil, jabatan, status, atau departemen anggota
// @Tags         Admin - Buku Induk Anggota
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Anggota"
// @Param        request body entity.UpdateAnggotaRequest true "Payload Perbarui Data Anggota"
// @Success      200 {object} utils.JSendResponse{data=entity.Anggota}
// @Failure      400 {object} utils.JSendResponse
// @Failure      404 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/anggota/{id} [put]
func (h *AdminBIAHandler) UpdateAnggota(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID anggota tidak valid", nil))
		return
	}

	var req entity.UpdateAnggotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload update anggota tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.biaUsecase.UpdateAnggota(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data anggota berhasil diperbarui", res))
}

// DeleteAnggota godoc
// @Summary      Hapus Data Anggota
// @Description  Menghapus data anggota dari Buku Induk (hanya jika anggota tidak memegang slot akun aktif)
// @Tags         Admin - Buku Induk Anggota
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Anggota"
// @Success      200 {object} utils.JSendResponse
// @Failure      400 {object} utils.JSendResponse
// @Failure      404 {object} utils.JSendResponse
// @Router       /admin/anggota/{id} [delete]
func (h *AdminBIAHandler) DeleteAnggota(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID anggota tidak valid", nil))
		return
	}

	if err := h.biaUsecase.DeleteAnggota(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Data anggota berhasil dihapus", nil))
}
