package main

import (
	"fmt"
	"log"

	"github.com/HMIKProject/hmik-corex-backend/config"
	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Membaca file .env
	godotenv.Load("../../.env") // Sesuaikan path relatif dari cmd/seeder/auth_seed

	// Connect to Database
	db := config.InitDB()

	// Password yang akan kita hash
	plainPassword := "Human102230462026"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal melakukan hashing:", err)
	}

	// Data Dummy Akun Slot
	akun := entity.Pengguna{
		NamaLengkap:  "Budi Human Resource",
		NamaPengguna: "hmik.human.01",
		Email:        "hmik.human.01@hmik.com",
		Sandi:        string(hashedPassword),
		StatusAktif:  true,
	}

	// Insert ke database (Abaikan error kalau email sudah ada)
	result := db.Where(entity.Pengguna{Email: akun.Email}).FirstOrCreate(&akun)
	
	if result.Error != nil {
		log.Fatal("Gagal insert akun:", result.Error)
	}

	fmt.Println("✅ Akun dummy berhasil dibuat!")
	fmt.Printf("Email: %s\nPassword: %s\n", akun.Email, plainPassword)
}
