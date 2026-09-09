package handler

import (
	"net/http"

	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AnggotaHandler mengatur jalur komunikasi HTTP untuk entitas anggota
type AnggotaHandler struct {
	usecase usecase.AnggotaUsecase
}

// NewAnggotaHandler menginisiasi handler baru
func NewAnggotaHandler(usecase usecase.AnggotaUsecase) *AnggotaHandler {
	return &AnggotaHandler{usecase}
}

// GetAllAnggota adalah endpoint untuk mengambil seluruh daftar tim company profile
func (h *AnggotaHandler) GetAllAnggota(c *gin.Context) {
	// Meminta Koki (Usecase) untuk mengambilkan data
	anggotas, err := h.usecase.GetAllAnggota()
	if err != nil {
		// Menggunakan utils JSend untuk merapikan error
		response := utils.BuildErrorResponse("Gagal mengambil data anggota dari server")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Membungkus hasil sukses dengan utils JSend
	response := utils.BuildSuccessResponse("Data anggota tim berhasil diambil", anggotas)
	c.JSON(http.StatusOK, response)
}
