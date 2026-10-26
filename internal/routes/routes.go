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
	
	// Dashboard 
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)
	
	// Anggota (Company Profile)
	anggotaRepo := repository.NewAnggotaRepository(db)
	anggotaUsecase := usecase.NewAnggotaUsecase(anggotaRepo)
	anggotaHandler := handler.NewAnggotaHandler(anggotaUsecase)

	// Departemen
	departemenRepo := repository.NewDepartemenRepository(db)
	departemenUsecase := usecase.NewDepartemenUsecase(departemenRepo)
	departemenHandler := handler.NewDepartemenHandler(departemenUsecase)

	// Program Kerja
	prokerRepo := repository.NewProgramKerjaRepository(db)
	prokerUsecase := usecase.NewProgramKerjaUsecase(prokerRepo)
	prokerHandler := handler.NewProgramKerjaHandler(prokerUsecase)

	// Upload (Cloudinary)
	uploadHandler := handler.NewUploadHandler()

	// Auth
	penggunaRepo := repository.NewPenggunaRepository(db)
	authUsecase := usecase.NewAuthUsecase(penggunaRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	// --- Mendaftarkan Routes ---

	// Endpoint khusus untuk memunculkan halaman website Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Grup API Versi 1
	v1 := router.Group("/api/v1")
	{
		// Auth Routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
		}

		// Endpoint untuk halaman tim Company Profile
		v1.GET("/company-profile/team", anggotaHandler.GetAllAnggota)

		// Endpoint untuk Departemen
		v1.GET("/departemen", departemenHandler.GetAllDepartemen)

		// Endpoint untuk Program Kerja (Global, filterable by status/is_unggulan)
		v1.GET("/program-kerja", prokerHandler.GetAllProgramKerja)

		// Endpoint untuk unggah file gambar
		v1.POST("/upload", uploadHandler.UploadImage)
		
		// Endpoint Khusus Dashboard Admin - Halaman Utama (Spesifik per Card)
		dashboardAdminGroup := v1.Group("/dashboard-admin-halaman-utama")
		{
			dashboardAdminGroup.GET("/ringkasan", dashboardHandler.GetDashboardStats)
			dashboardAdminGroup.GET("/total-anggota-aktif", dashboardHandler.GetTotalAnggotaAktif)
			dashboardAdminGroup.GET("/total-modul", dashboardHandler.GetTotalModul)
			dashboardAdminGroup.GET("/total-program-kerja", dashboardHandler.GetTotalProgramKerja)
			dashboardAdminGroup.GET("/voting-aktif", dashboardHandler.GetVotingAktif)
		}

		// Alias untuk fleksibilitas & kompatibilitas
		v1.GET("/dashboard", dashboardHandler.GetDashboardStats)
		v1.GET("/admin/dashboard/stats", dashboardHandler.GetDashboardStats)
	}
}
