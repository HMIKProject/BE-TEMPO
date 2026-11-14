package entity

import "time"

// ModulPembelajaran merepresentasikan berkas materi/modul pembelajaran (BP-03)
type ModulPembelajaran struct {
	IDModul        int64     `json:"id_modul" gorm:"column:id_modul;primaryKey;autoIncrement"`
	Judul          string    `json:"judul" gorm:"column:judul;type:varchar(255);not null"`
	Deskripsi      string    `json:"deskripsi" gorm:"column:deskripsi;type:text"`
	Kategori       string    `json:"kategori" gorm:"column:kategori;type:varchar(100);default:''"`
	Semester       string    `json:"semester" gorm:"column:semester;type:varchar(50);not null;default:'Semester 1'"` // Semester 1 s/d 8
	MataKuliah     string    `json:"mata_kuliah" gorm:"column:mata_kuliah;type:varchar(150);not null;default:'Umum'"`
	TipeFile       string    `json:"tipe_file" gorm:"column:tipe_file;type:varchar(20);not null;default:'pdf'"`    // pdf, pptx, zip, mp4
	UkuranFile     string    `json:"ukuran_file" gorm:"column:ukuran_file;type:varchar(50);default:'0 KB'"`
	UkuranBytes    int64     `json:"ukuran_bytes" gorm:"column:ukuran_bytes;default:0"`
	FileURL        string    `json:"file_url" gorm:"column:file_url;type:text;not null"`
	DibuatOleh     string    `json:"dibuat_oleh" gorm:"column:dibuat_oleh;type:varchar(150);default:'Admin HMIK'"`
	DibuatPada     time.Time `json:"dibuat_pada" gorm:"column:dibuat_pada;autoCreateTime"`
	DiperbaruiPada time.Time `json:"diperbarui_pada" gorm:"column:diperbarui_pada;autoUpdateTime"`
}

// TableName menentukan nama tabel database secara eksplisit
func (ModulPembelajaran) TableName() string {
	return "modul_pembelajaran"
}

// ModulFilterDTO parameter filter pencarian modul
type ModulFilterDTO struct {
	Search     string `form:"search"`
	Kategori   string `form:"kategori"`
	Semester   string `form:"semester"`
	MataKuliah string `form:"mata_kuliah"`
	TipeFile   string `form:"tipe_file"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=50"`
}

// ModulResponseDTO struktur respons terstandarisasi untuk modul
type ModulResponseDTO struct {
	IDModul    int64  `json:"id_modul"`
	Judul      string `json:"judul"`
	Deskripsi  string `json:"deskripsi"`
	Kategori   string `json:"kategori"`
	Semester   string `json:"semester"`
	MataKuliah string `json:"mata_kuliah"`
	TipeFile   string `json:"tipe_file"`
	UkuranFile string `json:"ukuran_file"`
	FileURL    string `json:"file_url"`
	DibuatOleh string `json:"dibuat_oleh"`
	DibuatPada string `json:"dibuat_pada"`
}
