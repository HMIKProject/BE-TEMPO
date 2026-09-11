package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/joho/godotenv"
)

type AnggotaData struct {
	Nim          string `json:"nim"`
	NamaLengkap  string `json:"nama_lengkap"`
	Email        string `json:"email"`
	ProgramStudi string `json:"program_studi"`
	Angkatan     int    `json:"angkatan"`
	Jabatan      string `json:"jabatan"`
	Foto         string `json:"foto"`
}

func main() {
	// 1. Muat .env
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan.")
	}

	fmt.Println("🚀 Memulai proses upload masal ke Cloudinary...")

	// 2. Baca JSON
	jsonFile, err := os.ReadFile("cmd/seeder/data_anggota.json")
	if err != nil {
		log.Fatalf("Gagal membaca file JSON: %v", err)
	}

	var dataAnggota []AnggotaData
	if err := json.Unmarshal(jsonFile, &dataAnggota); err != nil {
		log.Fatalf("Gagal parse JSON: %v", err)
	}

	berubah := false

	// 3. Loop dan Upload
	for i, data := range dataAnggota {
		if strings.HasPrefix(data.Foto, "/home/") {
			fmt.Printf("Mencoba upload foto untuk %s (%s)...\n", data.NamaLengkap, data.Nim)
			
			// Buka file lokal
			file, err := os.Open(data.Foto)
			if err != nil {
				fmt.Printf("⚠️ Gagal membuka file %s: %v\n", data.Foto, err)
				continue
			}
			
			// Upload
			url, err := utils.UploadImageToCloudinary(file, data.Nim+".png", "hmik/cp/profile")
			file.Close() // Pastikan file ditutup
			
			if err != nil {
				fmt.Printf("❌ Gagal upload untuk %s: %v\n", data.NamaLengkap, err)
			} else {
				fmt.Printf("✅ Sukses upload %s -> %s\n", data.NamaLengkap, url)
				dataAnggota[i].Foto = url
				berubah = true
			}
		}
	}

	// 4. Simpan kembali ke JSON jika ada perubahan
	if berubah {
		fmt.Println("💾 Menyimpan pembaruan URL ke data_anggota.json...")
		updatedJson, _ := json.MarshalIndent(dataAnggota, "", "  ")
		os.WriteFile("cmd/seeder/data_anggota.json", updatedJson, 0644)
	}

	fmt.Println("🎉 Proses upload selesai!")
}
