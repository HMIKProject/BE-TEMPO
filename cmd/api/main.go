package main

import (
	"log"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// 1. Koneksi ke Database (Port 5433 sesuai setup seeder)
	dsn := "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}

	// 2. Inisiasi Mesin Gin
	router := gin.Default()

	// 3. Konfigurasi Middleware CORS (Sesuai Panduan HMIK_note)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://hmik.kampus.ac.id"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 4. Mendaftarkan Seluruh Endpoint API
	routes.SetupRoutes(router, db)

	// 5. Menjalankan Server di port 8080
	log.Println("🚀 Server Backend berjalan di http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}