package handler

import (
	"net/http"
	"strconv"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ProgramKerjaHandler struct {
	usecase usecase.ProgramKerjaUsecase
}

func NewProgramKerjaHandler(usecase usecase.ProgramKerjaUsecase) *ProgramKerjaHandler {
	return &ProgramKerjaHandler{usecase}
}

// GetAllProgramKerja godoc
// @Summary      Ambil daftar program kerja HMIK
// @Description  Mengembalikan data program kerja, bisa difilter berdasarkan status (Mendatang) atau is_unggulan (true/false)
// @Tags         ProgramKerja
// @Accept       json
// @Produce      json
// @Param        status query string false "Filter berdasarkan status (Contoh: Mendatang, Terlaksana)"
// @Param        is_unggulan query bool false "Filter berdasarkan apakah proker unggulan (true/false)"
// @Success      200  {object}  utils.JSendResponse  "Berhasil mengambil data program kerja"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /program-kerja [get]
func (h *ProgramKerjaHandler) GetAllProgramKerja(c *gin.Context) {
	status := c.Query("status")
	
	var isUnggulanPtr *bool
	isUnggulanStr := c.Query("is_unggulan")
	if isUnggulanStr != "" {
		isUnggulan, err := strconv.ParseBool(isUnggulanStr)
		if err == nil {
			isUnggulanPtr = &isUnggulan
		}
	}

	prokers, err := h.usecase.GetAllProgramKerja(status, isUnggulanPtr)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal mengambil data program kerja dari server")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse("Data program kerja berhasil diambil", prokers)
	c.JSON(http.StatusOK, response)
}

// GetByID mengambil detail program kerja
func (h *ProgramKerjaHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID program kerja tidak valid", []string{"id harus angka"}))
		return
	}

	proker, err := h.usecase.GetProgramKerjaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildFailResponse("Program kerja tidak ditemukan", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Berhasil mengambil detail program kerja", proker))
}

// Create menambahkan program kerja baru (Admin CMS)
func (h *ProgramKerjaHandler) Create(c *gin.Context) {
	var req entity.ProgramKerja
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format JSON tidak valid", []string{err.Error()}))
		return
	}

	if err := h.usecase.CreateProgramKerja(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal menambahkan program kerja", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Program kerja berhasil ditambahkan", req))
}

// Update memperbarui program kerja (Admin CMS)
func (h *ProgramKerjaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID program kerja tidak valid", []string{"id harus angka"}))
		return
	}

	var req entity.ProgramKerja
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format JSON tidak valid", []string{err.Error()}))
		return
	}

	if err := h.usecase.UpdateProgramKerja(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal memperbarui program kerja", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Program kerja berhasil diperbarui", gin.H{"id_proker": id}))
}

// Delete menghapus program kerja (Admin CMS)
func (h *ProgramKerjaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID program kerja tidak valid", []string{"id harus angka"}))
		return
	}

	if err := h.usecase.DeleteProgramKerja(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal menghapus program kerja", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Program kerja berhasil dihapus", gin.H{"id_proker": id}))
}
