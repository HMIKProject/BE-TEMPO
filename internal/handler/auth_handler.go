package handler

import (
	"net/http"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: u}
}

// Login godoc
// @Summary Login Pengurus HMIK
// @Description Endpoint untuk mendapatkan JWT Token menggunakan email/username dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Kredensial Login"
// @Success 200 {object} utils.JSendResponse{data=entity.LoginResponse}
// @Failure 401 {object} utils.JSendResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	
	// Bind JSON dari request body ke struct LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	// Panggil Usecase untuk proses Login
	res, err := h.authUsecase.Login(req)
	if err != nil {
		// HTTP 401 Unauthorized jika kredensial salah
		c.JSON(http.StatusUnauthorized, utils.BuildFailResponse(err.Error(), nil))
		return
	}

	// Sukses Login, kembalikan JWT (JSend format)
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Login berhasil", res))
}
