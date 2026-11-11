package routes

import (
	"github.com/HMIKProject/hmik-corex-backend/internal/handler"
	"github.com/HMIKProject/hmik-corex-backend/internal/middleware"
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

	// Pemira (E-Voting)
	pemiraRepo := repository.NewPemiraRepository(db)
	pemiraUsecase := usecase.NewPemiraUsecase(pemiraRepo)
	pemiraHandler := handler.NewPemiraHandler(pemiraUsecase)

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
			authGroup.POST("/set-initial-password", authHandler.SetInitialPassword)
		}

		// Pemira (E-Voting) Voter Routes
		pemiraGroup := v1.Group("/pemira")
		{
			pemiraAuthGroup := pemiraGroup.Group("/auth")
			{
				pemiraAuthGroup.POST("/request-otp", pemiraHandler.RequestOTP)
				pemiraAuthGroup.POST("/verify-otp", pemiraHandler.VerifyOTP)
				pemiraAuthGroup.POST("/resend-otp", pemiraHandler.ResendOTP)
			}

			pemiraGroup.GET("/kandidat", pemiraHandler.GetActiveCandidates)
			pemiraGroup.GET("/jadwal", pemiraHandler.GetPublicSchedule)
			pemiraGroup.POST("/vote", middleware.VoterMiddleware(), pemiraHandler.CastVote)
		}

		// Endpoint untuk halaman tim Company Profile
		v1.GET("/company-profile/team", anggotaHandler.GetAllAnggota)

		// Endpoint untuk Departemen
		v1.GET("/departemen", departemenHandler.GetAllDepartemen)

		// Endpoint untuk Program Kerja (Global, filterable by status/is_unggulan)
		v1.GET("/program-kerja", prokerHandler.GetAllProgramKerja)

		// Jalur Rute Admin yang Dilindungi Middleware Token JWT (Standar Obsidian)
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware())
		{
			dashboardGroup := adminGroup.Group("/dashboard")
			{
				dashboardGroup.GET("/stats", dashboardHandler.GetDashboardStats)
			}

			// Pemira Monitoring & Management
			adminPemiraGroup := adminGroup.Group("/pemira")
			{
				adminPemiraGroup.GET("/stats", pemiraHandler.GetElectionStats)
				adminPemiraGroup.GET("/chart", pemiraHandler.GetElectionChart)
				adminPemiraGroup.GET("/pemilih", pemiraHandler.GetVoters)
				adminPemiraGroup.POST("/pemilih/:id/reset", pemiraHandler.ResetVoter)
				adminPemiraGroup.POST("/reset-kotak-suara", pemiraHandler.ResetKotakSuara)

				// Manajemen Paslon (Kandidat) CRUD
				adminPemiraGroup.POST("/kandidat", pemiraHandler.CreateCandidate)
				adminPemiraGroup.PUT("/kandidat/:id", pemiraHandler.UpdateCandidate)
				adminPemiraGroup.DELETE("/kandidat/:id", pemiraHandler.DeleteCandidate)

				// Pengaturan Acara & Status Trigger
				adminPemiraGroup.GET("/settings", pemiraHandler.GetElectionSettings)
				adminPemiraGroup.PUT("/settings", pemiraHandler.UpdateElectionSchedule)
				adminPemiraGroup.POST("/status", pemiraHandler.UpdateElectionStatus)

				// Pelantikan MPIK & Transisi Manajemen Himpunan
				adminPemiraGroup.POST("/pelantikan-mpik", pemiraHandler.PelantikanMPIK)
				adminPemiraGroup.POST("/new-period", pemiraHandler.CreateNewPeriod)
			}

			// Endpoint unggah file gambar admin
			adminGroup.POST("/upload", uploadHandler.UploadImage)
		}

		// Rute upload terproteksi JWT (Backward Compatibility)
		v1.POST("/upload", middleware.AuthMiddleware(), uploadHandler.UploadImage)

		// Rute Kompatibilitas Khusus Dashboard Admin (Terproteksi Middleware JWT)
		dashboardAdminGroup := v1.Group("/dashboard-admin-halaman-utama")
		dashboardAdminGroup.Use(middleware.AuthMiddleware())
		{
			dashboardAdminGroup.GET("/ringkasan", dashboardHandler.GetDashboardStats)
			dashboardAdminGroup.GET("/total-anggota-aktif", dashboardHandler.GetTotalAnggotaAktif)
			dashboardAdminGroup.GET("/total-modul", dashboardHandler.GetTotalModul)
			dashboardAdminGroup.GET("/total-program-kerja", dashboardHandler.GetTotalProgramKerja)
			dashboardAdminGroup.GET("/voting-aktif", dashboardHandler.GetVotingAktif)
		}

		// Alias untuk fleksibilitas (Terproteksi Middleware JWT)
		v1.GET("/dashboard", middleware.AuthMiddleware(), dashboardHandler.GetDashboardStats)
	}
}
