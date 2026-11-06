package main

import (
	"log"
	"os"

	"github.com/HMIKProject/hmik-corex-backend/internal/seeder"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	godotenv.Load("../../.env")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}

	if err := seeder.SeedAllRoleAccounts(db); err != nil {
		log.Fatal("Gagal seed accounts:", err)
	}
}
