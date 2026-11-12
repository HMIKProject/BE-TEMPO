package handler

import (
	"net/http"
	"strconv"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminOrganisasiHandler struct {
	orgUsecase usecase.AdminOrganisasiUsecase
}

func NewAdminOrganisasiHandler(u usecase.AdminOrganisasiUsecase) *AdminOrganisasiHandler {
	return &AdminOrganisasiHandler{orgUsecase: u}
}

// GetStrukturOrganisasi godoc
// @Summary      Susunan Struktur Organisasi Lengkap (BPH, Departemen, Divisi)
// @Description  Mengambil bagan struktur organisasi HMIK termasuk data kabinet aktif, jajaran BPH, dan departemen beserta divisi & kepala departemen
// @Tags         Admin - Struktur Organisasi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.StrukturOrganisasiResponse}
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/organisasi [get]
func (h *AdminOrganisasiHandler) GetStrukturOrganisasi(c *gin.Context) {
	res, err := h.orgUsecase.GetStrukturOrganisasi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil struktur organisasi: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Struktur organisasi berhasil dimuat", res))
}

// CreateDepartemen godoc
// @Summary      Tambah Departemen Baru
// @Description  Mendaftarkan departemen baru ke dalam bagan himpunan (Khusus BPH)
// @Tags         Admin - Struktur Organisasi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body entity.CreateDeptRequest true "Data Departemen Baru"
// @Success      201 {object} utils.JSendResponse{data=entity.Departemen}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/departemen [post]
func (h *AdminOrganisasiHandler) CreateDepartemen(c *gin.Context) {
	var req entity.CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload tidak valid", []string{err.Error()}))
		return
	}

	dept, err := h.orgUsecase.CreateDepartemen(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Departemen baru berhasil ditambahkan", dept))
}

// UpdateDepartemen godoc
// @Summary      Perbarui Data Departemen
// @Description  Mengubah nama, deskripsi, logo, atau status aktif departemen (Khusus BPH)
// @Tags         Admin - Struktur Organisasi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Departemen"
// @Param        payload body entity.UpdateDeptRequest true "Data Perubahan Departemen"
// @Success      200 {object} utils.JSendResponse{data=entity.Departemen}
// @Failure      400 {object} utils.JSendResponse
// @Failure      404 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/departemen/{id} [put]
func (h *AdminOrganisasiHandler) UpdateDepartemen(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID departemen tidak valid", nil))
		return
	}

	var req entity.UpdateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload tidak valid", []string{err.Error()}))
		return
	}

	dept, err := h.orgUsecase.UpdateDepartemen(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Departemen berhasil diperbarui", dept))
}

// ArchiveDepartemen godoc
// @Summary      Arsipkan Departemen (Soft Delete)
// @Description  Menonaktifkan departemen (status_aktif = false) agar riwayat tidak terputus (Khusus BPH)
// @Tags         Admin - Struktur Organisasi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Departemen"
// @Success      200 {object} utils.JSendResponse
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/departemen/{id} [delete]
func (h *AdminOrganisasiHandler) ArchiveDepartemen(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID departemen tidak valid", nil))
		return
	}

	err = h.orgUsecase.ArchiveDepartemen(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Departemen berhasil diarsipkan", nil))
}

// CreateDivisi godoc
// @Summary      Tambah Divisi Baru
// @Description  Mendaftarkan divisi baru di bawah departemen tertentu (Khusus BPH)
// @Tags         Admin - Struktur Organisasi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body entity.CreateDivisiRequest true "Data Divisi Baru"
// @Success      201 {object} utils.JSendResponse{data=entity.DivisiDepartemen}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/divisi [post]
func (h *AdminOrganisasiHandler) CreateDivisi(c *gin.Context) {
	var req entity.CreateDivisiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload tidak valid", []string{err.Error()}))
		return
	}

	divisi, err := h.orgUsecase.CreateDivisi(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Divisi baru berhasil ditambahkan", divisi))
}

// UpdateDivisi godoc
// @Summary      Perbarui Data Divisi
// @Description  Mengubah nama divisi, deskripsi, atau status aktif divisi (Khusus BPH)
// @Tags         Admin - Struktur Organisasi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Divisi"
// @Param        payload body entity.UpdateDivisiRequest true "Data Perubahan Divisi"
// @Success      200 {object} utils.JSendResponse{data=entity.DivisiDepartemen}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/divisi/{id} [put]
func (h *AdminOrganisasiHandler) UpdateDivisi(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID divisi tidak valid", nil))
		return
	}

	var req entity.UpdateDivisiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload tidak valid", []string{err.Error()}))
		return
	}

	divisi, err := h.orgUsecase.UpdateDivisi(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Divisi berhasil diperbarui", divisi))
}
