package handler

import (
	"net/http"
	"strconv"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminSlotHandler struct {
	slotUsecase usecase.AdminSlotUsecase
}

func NewAdminSlotHandler(u usecase.AdminSlotUsecase) *AdminSlotHandler {
	return &AdminSlotHandler{slotUsecase: u}
}

// GetAccountSlots godoc
// @Summary      Daftar Slot Akun Pinjam Pakai
// @Description  Mengambil seluruh 38 slot akun struktural beserta status keterisian (TERPAKAI/KOSONG), data peminjam, dan status kode aktivasi
// @Tags         Admin - Slot Akun Pinjam Pakai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=[]entity.AccountSlotResponse}
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/akun/slots [get]
func (h *AdminSlotHandler) GetAccountSlots(c *gin.Context) {
	res, err := h.slotUsecase.GetAccountSlots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengambil data slot akun: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Daftar slot akun berhasil dimuat", res))
}

// AssignSlot godoc
// @Summary      Tautkan Slot Akun ke Mahasiswa
// @Description  Menautkan slot akun kosong ke anggota dari Buku Induk dan menerbitkan Kode Aktivasi 7 hari
// @Tags         Admin - Slot Akun Pinjam Pakai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.AssignSlotRequest true "Payload Penautan Slot Akun"
// @Success      200 {object} utils.JSendResponse{data=entity.AssignSlotResponse}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/akun/assign [post]
func (h *AdminSlotHandler) AssignSlot(c *gin.Context) {
	var req entity.AssignSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload penautan akun tidak valid: "+err.Error(), nil))
		return
	}

	res, err := h.slotUsecase.AssignSlot(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Slot akun berhasil ditautkan dan Kode Aktivasi telah diterbitkan", res))
}

// UnassignSlot godoc
// @Summary      Lepas Tautan Slot Akun
// @Description  Melepas kaitan slot akun dari anggota, mendemisionerkan anggota, mereset password acak, dan membatalkan sesi login
// @Tags         Admin - Slot Akun Pinjam Pakai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body entity.UnassignSlotRequest true "Payload Pelepasan Tautan Slot Akun"
// @Success      200 {object} utils.JSendResponse
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/akun/unassign [post]
func (h *AdminSlotHandler) UnassignSlot(c *gin.Context) {
	var req entity.UnassignSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload lepas tautan tidak valid: "+err.Error(), nil))
		return
	}

	if err := h.slotUsecase.UnassignSlot(req.IdPengguna); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Tautan slot akun berhasil dilepas dan kata sandi telah direset secara acak", nil))
}

// DemisionerMassal godoc
// @Summary      Demisioner Massal Akhir Periode
// @Description  Memutus serentak seluruh slot pengurus aktif, mereset kata sandi akun secara acak, dan membatalkan seluruh sesi login
// @Tags         Admin - Slot Akun Pinjam Pakai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.JSendResponse{data=entity.DemisionerSummary}
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/akun/demisioner-massal [post]
func (h *AdminSlotHandler) DemisionerMassal(c *gin.Context) {
	res, err := h.slotUsecase.DemisionerMassal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengeksekusi demisioner massal: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Demisioner massal berhasil dieksekusi", res))
}

// ResetSlotPassword godoc
// @Summary      Reset Sandi Pengurus (Terbit Ulang Kode Aktivasi)
// @Description  Menerbitkan ulang Kode Aktivasi 7 hari bagi pengurus yang lupa sandi dan membatalkan token login aktif
// @Tags         Admin - Slot Akun Pinjam Pakai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Pengguna Slot"
// @Success      200 {object} utils.JSendResponse{data=entity.ResetSlotResponse}
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/akun/{id}/reset-password [post]
func (h *AdminSlotHandler) ResetSlotPassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID pengguna slot tidak valid", nil))
		return
	}

	res, err := h.slotUsecase.ResetSlotPassword(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Kode aktivasi baru berhasil diterbitkan", res))
}

// UpdateSlotRole godoc
// @Summary      Atur Peran / Hak Akses Slot Akun
// @Description  Mengubah peran (id_peran) pada slot akun tertentu
// @Tags         Admin - Slot Akun Pinjam Pakai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID Pengguna Slot"
// @Param        request body entity.UpdateSlotRoleRequest true "Payload Update Role"
// @Success      200 {object} utils.JSendResponse
// @Failure      400 {object} utils.JSendResponse
// @Failure      500 {object} utils.JSendResponse
// @Router       /admin/akun/{id}/role [put]
func (h *AdminSlotHandler) UpdateSlotRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID pengguna slot tidak valid", nil))
		return
	}

	var req entity.UpdateSlotRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Payload update role tidak valid: "+err.Error(), nil))
		return
	}

	if err := h.slotUsecase.UpdateSlotRole(id, req.IdPeran); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Peran slot akun berhasil diperbarui", nil))
}
