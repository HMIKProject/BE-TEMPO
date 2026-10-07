package main

import (
	"log"
	"net/http"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/config"
	"github.com/HMIKProject/hmik-corex-backend/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Membaca file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment bawaan")
	}

	// Inisialisasi Database
	db := config.InitDB()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://hmik.kampus.ac.id"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "HMIK CoreX Backend is running",
		})
	})

	// PENTING: Mendaftarkan seluruh rute API dan Swagger ke dalam Gin
	routes.SetupRoutes(router, db)

	// Jalankan server
	log.Println("Server berjalan di port :8080")
	router.Run(":8080")
}
