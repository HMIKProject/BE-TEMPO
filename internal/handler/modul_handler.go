package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/usecase"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ModulHandler struct {
	modulUsecase usecase.ModulUsecase
}

func NewModulHandler(modulUsecase usecase.ModulUsecase) *ModulHandler {
	return &ModulHandler{modulUsecase: modulUsecase}
}

// GetAll mengambil daftar modul materi terfilter
func (h *ModulHandler) GetAll(c *gin.Context) {
	var filter entity.ModulFilterDTO
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Parameter query tidak valid", []string{err.Error()}))
		return
	}

	list, total, err := h.modulUsecase.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Berhasil mengambil daftar modul pembelajaran", gin.H{
		"total_data": total,
		"page":       filter.Page,
		"limit":      filter.Limit,
		"modul":      list,
	}))
}

// GetByID mengambil detail modul berdasarkan ID
func (h *ModulHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID modul tidak valid", []string{"id harus angka"}))
		return
	}

	modul, err := h.modulUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildFailResponse("Modul tidak ditemukan", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Berhasil mengambil detail modul", modul))
}

// Create menambahkan modul baru
func (h *ModulHandler) Create(c *gin.Context) {
	var req entity.ModulPembelajaran
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format JSON tidak valid", []string{err.Error()}))
		return
	}

	if err := h.modulUsecase.Create(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal menambahkan modul", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("Modul pembelajaran berhasil ditambahkan", req))
}

// Update memperbarui data modul
func (h *ModulHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID modul tidak valid", []string{"id harus angka"}))
		return
	}

	var req entity.ModulPembelajaran
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format JSON tidak valid", []string{err.Error()}))
		return
	}

	if err := h.modulUsecase.Update(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal memperbarui modul", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Modul pembelajaran berhasil diperbarui", gin.H{"id_modul": id}))
}

// Delete menghapus modul
func (h *ModulHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("ID modul tidak valid", []string{"id harus angka"}))
		return
	}

	if err := h.modulUsecase.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Gagal menghapus modul", []string{err.Error()}))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("Modul pembelajaran berhasil dihapus", gin.H{"id_modul": id}))
}

// UploadFile menangani pengunggahan berkas dokumen/multimedia modul pembelajaran
func (h *ModulHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("File tidak ditemukan pada request", []string{err.Error()}))
		return
	}

	// Batasan Maksimal 50 MB (Figma & BP-03)
	const maxFileSize = 50 * 1024 * 1024
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Ukuran file melebihi batas maksimal 50 MB", []string{"file terlalu besar"}))
		return
	}

	// Ekstensi yang diizinkan: pdf, pptx, ppt, zip, rar, mp4, docx, doc
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".pdf":  true,
		".pptx": true,
		".ppt":  true,
		".zip":  true,
		".rar":  true,
		".mp4":  true,
		".docx": true,
		".doc":  true,
	}

	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, utils.BuildFailResponse("Format file tidak didukung. Harap gunakan PDF, PPTX, DOCX, ZIP, atau MP4", []string{"ekstensi tidak diizinkan"}))
		return
	}

	// Baca parameter semester dan mata kuliah untuk menentukan target folder
	semester := c.DefaultPostForm("semester", "Semester 1")
	mataKuliah := c.DefaultPostForm("mata_kuliah", "Umum")
	cleanSemester := strings.ReplaceAll(strings.ToLower(semester), " ", "_")
	cleanMatkul := strings.ReplaceAll(strings.ToLower(mataKuliah), " ", "_")
	cloudFolder := fmt.Sprintf("hmik/learning/%s/%s", cleanSemester, cleanMatkul)

	var fileURL string
	var uploadErr error

	// 1. Coba upload langsung ke Cloud Storage (Cloudinary)
	srcFile, err := file.Open()
	if err == nil {
		defer srcFile.Close()
		fileURL, uploadErr = utils.UploadFileToCloudinary(srcFile, file.Filename, cloudFolder)
	}

	// 2. Fallback ke penyimpanan lokal jika Cloudinary offline atau tidak terkonfigurasi
	if uploadErr != nil || fileURL == "" {
		uploadDir := filepath.Join("uploads", "modul", cleanSemester, cleanMatkul)
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal menyiapkan direktori penyimpanan"))
			return
		}

		timestamp := time.Now().UnixNano()
		safeFilename := fmt.Sprintf("%d_%s", timestamp, strings.ReplaceAll(file.Filename, " ", "_"))
		dst := filepath.Join(uploadDir, safeFilename)

		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("Gagal menyimpan berkas ke server"))
			return
		}
		fileURL = "/" + dst
	}

	// Format ukuran berkas (contoh: 4.2 MB)
	var formattedSize string
	if file.Size < 1024*1024 {
		formattedSize = fmt.Sprintf("%.1f KB", float64(file.Size)/1024.0)
	} else {
		formattedSize = fmt.Sprintf("%.1f MB", float64(file.Size)/(1024.0*1024.0))
	}

	cleanExt := strings.TrimPrefix(ext, ".")

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("File modul berhasil diunggah", gin.H{
		"file_name":    file.Filename,
		"file_url":     fileURL,
		"tipe_file":    cleanExt,
		"ukuran_file":  formattedSize,
		"ukuran_bytes": file.Size,
	}))
}
