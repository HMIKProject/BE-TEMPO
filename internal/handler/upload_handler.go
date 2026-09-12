package handler

import (
	"net/http"

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
// @Description  Menerima file form-data gambar dan mengunggahnya ke Cloudinary, mengembalikan URL gambar
// @Tags         Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        image   formData  file    true  "File gambar yang akan diunggah"
// @Param        folder  query     string  false "Folder tujuan di Cloudinary (opsional)"
// @Success      200  {object}  utils.JSendResponse  "Gambar berhasil diunggah"
// @Failure      400  {object}  utils.JSendResponse  "Input file tidak valid"
// @Failure      500  {object}  utils.JSendResponse  "Gagal mengunggah ke Cloudinary"
// @Router       /upload [post]
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 1. Set maksimal ukuran memori untuk file (contoh: 5MB)
	c.Request.ParseMultipartForm(5 << 20)

	// 2. Mengambil file dari field bernama "image" di HTML form-data
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("Gagal menemukan file gambar (key form-data harus 'image')"))
		return
	}
	defer file.Close()

	// Menentukan folder tujuan (Bisa dari parameter URL, atau default ke hmik/cp/profile)
	folder := c.Query("folder")
	if folder == "" {
		folder = "hmik/cp/profile" // Default sesuai arsitektur Fase 1
	}

	// 3. Serahkan file ke utilitas Cloudinary (Koki Eksternal) untuk diunggah
	imageUrl, err := utils.UploadImageToCloudinary(file, fileHeader.Filename, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal mengunggah gambar ke Cloudinary: "+err.Error()))
		return
	}

	// 4. Kembalikan URL hasil unggahan ke Frontend
	data := map[string]string{
		"image_url": imageUrl,
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Gambar berhasil diunggah", data))
}
