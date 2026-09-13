package handler

import (
	"net/http"
	"strconv"

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
