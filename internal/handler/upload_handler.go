package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UploadHandler mengatur jalur komunikasi HTTP untuk fitur unggah file
type UploadHandler struct{}

// NewUploadHandler menginisiasi handler upload baru
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage godoc
// @Summary      Unggah gambar ke Cloudinary
// @Description  Menerima file form-data gambar terverifikasi (JPEG, PNG, WebP, GIF max 5MB) dan mengunggahnya ke Cloudinary
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image   formData  file    true  "File gambar yang akan diunggah (maksimal 5MB, format: jpg, jpeg, png, webp, gif)"
// @Param        folder  query     string  false "Folder tujuan di Cloudinary (opsional)"
// @Success      200  {object}  utils.JSendResponse  "Gambar berhasil diunggah"
// @Failure      400  {object}  utils.JSendResponse  "Input file tidak valid atau ukuran/format ditolak"
// @Failure      401  {object}  utils.JSendResponse  "Akses Ditolak: Token JWT tidak ditemukan atau tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengunggah ke Cloudinary"
// @Router       /upload [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 1. Set maksimal ukuran memori untuk parsing form
	c.Request.ParseMultipartForm(5 << 20)

	// 2. Mengambil file dari field bernama "image" di HTML form-data
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal menemukan file gambar (key form-data harus 'image')", nil))
		return
	}
	defer file.Close()

	// 3. Validasi Ukuran File (Maksimal 5MB)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Ukuran file terlalu besar. Maksimal 5MB", nil))
		return
	}

	// 4. Validasi Ekstensi File
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format ekstensi file tidak diizinkan. Hanya .jpg, .jpeg, .png, .webp, atau .gif", nil))
		return
	}

	// 5. Validasi Magic Bytes (Deteksi MIME-type sebenarnya dari konten file)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Gagal membaca header konten file"))
		return
	}

	// Kembalikan kursor file ke awal agar Cloudinary dapat membaca seluruh isi file
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal memproses file gambar"))
			return
		}
	}

	detectedMime := http.DetectContentType(buffer[:n])
	allowedMimes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	if !allowedMimes[detectedMime] {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Konten file bukan gambar yang valid (terdeteksi: "+detectedMime+")", nil))
		return
	}

	// 6. Sanitasi folder tujuan (Mencegah Directory Traversal)
	folder := c.Query("folder")
	if folder == "" {
		folder = "hmik/cp/profile" // Default sesuai arsitektur Fase 1
	} else {
		folder = filepath.Clean(folder)
		folder = strings.TrimPrefix(folder, "/")
		folder = strings.TrimPrefix(folder, "../")
		if !strings.HasPrefix(folder, "hmik/") {
			folder = "hmik/" + folder
		}
	}

	// 7. Serahkan file ke utilitas Cloudinary untuk diunggah
	imageUrl, err := utils.UploadImageToCloudinary(file, fileHeader.Filename, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengunggah gambar ke Cloudinary: "+err.Error()))
		return
	}

	// 8. Kembalikan URL hasil unggahan ke Frontend
	data := map[string]string{
		"image_url": imageUrl,
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Gambar berhasil diunggah", data))
}
