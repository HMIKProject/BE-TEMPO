package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, continuing with environment variables")
	}
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	cldUrl := os.Getenv("CLOUDINARY_URL")
	if cldUrl == "" {
		log.Fatal("CLOUDINARY_URL is not set")
	}

	cld, err := cloudinary.NewFromURL(cldUrl)
	if err != nil {
		log.Fatal("Failed to init cloudinary:", err)
	}

	ctx := context.Background()
	log.Println("Uploading RISTEK logo...")
	deptLogoResp, err := cld.Upload.Upload(ctx, "/home/damos/Documents/HMIK/asset/departement/RISTEK/fotodept.png", uploader.UploadParams{
		Folder: "hmik/cp/departemen",
	})
	if err != nil {
		log.Fatal("Failed to upload dept logo:", err)
	}
	deptLogoUrl := deptLogoResp.SecureURL
	log.Println("Dept logo uploaded:", deptLogoUrl)

	log.Println("Uploading HMIK logo for proker...")
	prokerLogoResp, err := cld.Upload.Upload(ctx, "/home/damos/Documents/HMIK/HMIK_Code/ModulWeb/public/hmik.png", uploader.UploadParams{
		Folder: "hmik/cp/proker",
	})
	if err != nil {
		log.Fatal("Failed to upload proker logo:", err)
	}
	prokerLogoUrl := prokerLogoResp.SecureURL
	log.Println("Proker logo uploaded:", prokerLogoUrl)

	deskripsiDeptBytes, err := os.ReadFile("/home/damos/Documents/HMIK/asset/departement/RISTEK/deskripsidept.txt")
	if err != nil {
		log.Fatal("Failed to read dept text:", err)
	}
	deskripsiDept := string(deskripsiDeptBytes)

	var dept entity.Departemen
	err = db.Table("departemen").Where("nama_departemen LIKE ?", "%Riset%").First(&dept).Error
	if err != nil {
		log.Fatal("Failed to find department:", err)
	}

	dept.Deskripsi = deskripsiDept
	dept.Logo = deptLogoUrl
	dept.DiperbaruiPada = time.Now()
	db.Table("departemen").Save(&dept)
	log.Println("Department updated successfully!")

	deskripsiProkerBytes, err := os.ReadFile("/home/damos/Documents/HMIK/asset/departement/RISTEK/proker/deskripsiprokerhmikcorex.txt")
	if err != nil {
		log.Fatal("Failed to read proker text:", err)
	}
	deskripsiProker := string(deskripsiProkerBytes)

	var proker entity.ProgramKerja
	err = db.Table("program_kerja").Where("nama_proker = ?", "HMIK-CoreX").First(&proker).Error
	if err != nil {
		proker = entity.ProgramKerja{
			IdDepartemen:   dept.IdDepartemen,
			NamaProker:     "HMIK-CoreX",
			Deskripsi:      deskripsiProker,
			Foto:           prokerLogoUrl,
			Status:         "Mendatang",
			IsUnggulan:     true,
			DibuatPada:     time.Now(),
			DiperbaruiPada: time.Now(),
		}
		db.Table("program_kerja").Create(&proker)
		log.Println("Proker created successfully!")
	} else {
		proker.Deskripsi = deskripsiProker
		proker.Foto = prokerLogoUrl
		proker.IsUnggulan = true
		proker.DiperbaruiPada = time.Now()
		db.Table("program_kerja").Save(&proker)
		log.Println("Proker updated successfully!")
	}
}
