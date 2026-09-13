package handler

import (
	"net/http"

	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DepartemenHandler struct {
	usecase usecase.DepartemenUsecase
}

func NewDepartemenHandler(usecase usecase.DepartemenUsecase) *DepartemenHandler {
	return &DepartemenHandler{usecase}
}

// GetAllDepartemen godoc
// @Summary      Ambil daftar departemen HMIK
// @Description  Mengembalikan semua data departemen beserta anggota dan program kerjanya
// @Tags         Departemen
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.JSendResponse  "Berhasil mengambil data departemen"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengambil data dari server"
// @Router       /departemen [get]
func (h *DepartemenHandler) GetAllDepartemen(c *gin.Context) {
	depts, err := h.usecase.GetAllDepartemen()
	if err != nil {
		response := utils.BuildErrorResponse("Gagal mengambil data departemen dari server")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse("Data departemen berhasil diambil", depts)
	c.JSON(http.StatusOK, response)
}
