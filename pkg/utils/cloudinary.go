package utils

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadImageToCloudinary menerima file gambar (multipart.File) dan mengunggahnya ke Cloudinary
func UploadImageToCloudinary(file multipart.File, fileName string, folderName string) (string, error) {
	// 1. Ambil URL rahasia dari environment (.env)
	cldUrl := os.Getenv("CLOUDINARY_URL")
	if cldUrl == "" {
		return "", errors.New("konfigurasi CLOUDINARY_URL belum diatur di sistem atau file .env")
	}

	// 2. Buat instance Cloudinary dari URL tersebut
	cld, err := cloudinary.NewFromURL(cldUrl)
	if err != nil {
		return "", err
	}

	// 3. Mulai proses upload ke folder dinamis
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folderName, // Semua foto akan tersimpan rapi di dalam folder dinamis ini
	})

	if err != nil {
		return "", err
	}

	// 4. Kembalikan URL gambar yang sudah aman (HTTPS)
	return uploadResult.SecureURL, nil
}
