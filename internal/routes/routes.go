package routes

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/handler"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_ "github.com/HMIKProject/hmik-corex-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes mengatur semua pendaftaran jalur API (endpoints)
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// --- Inisialisasi Layer Clean Architecture ---
	
	// Anggota (Company Profile)
	anggotaRepo := repository.NewAnggotaRepository(db)
	anggotaUsecase := usecase.NewAnggotaUsecase(anggotaRepo)
	anggotaHandler := handler.NewAnggotaHandler(anggotaUsecase)

	// Upload (Cloudinary)
	uploadHandler := handler.NewUploadHandler()

	// --- Mendaftarkan Routes ---

	// Endpoint khusus untuk memunculkan halaman website Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Grup API Versi 1
	v1 := router.Group("/api/v1")
	{
		// Endpoint untuk halaman tim Company Profile
		v1.GET("/company-profile/team", anggotaHandler.GetAllAnggota)

		// Endpoint untuk unggah file gambar
		v1.POST("/upload", uploadHandler.UploadImage)
	}
}
