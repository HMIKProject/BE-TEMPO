package routes

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/handler"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes mengatur semua pendaftaran jalur API (endpoints)
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// --- Inisialisasi Layer Clean Architecture ---
	
	// Anggota (Company Profile)
	anggotaRepo := repository.NewAnggotaRepository(db)
	anggotaUsecase := usecase.NewAnggotaUsecase(anggotaRepo)
	anggotaHandler := handler.NewAnggotaHandler(anggotaUsecase)

	// --- Mendaftarkan Routes ---
	
	// Grup API Versi 1
	v1 := router.Group("/api/v1")
	{
		// Endpoint untuk halaman tim Company Profile
		v1.GET("/company-profile/team", anggotaHandler.GetAllAnggota)
	}
}
