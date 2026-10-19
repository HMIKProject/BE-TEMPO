package entity

import "time"

type Pemira struct {
	IdPemira       int64     `gorm:"primaryKey;column:id_pemira;autoIncrement" json:"id_pemira"`
	TahunPeriode   string    `gorm:"column:tahun_periode;type:varchar(50);not null" json:"tahun_periode"`
	TanggalMulai   time.Time `gorm:"column:tanggal_mulai;not null" json:"tanggal_mulai"`
	TanggalSelesai time.Time `gorm:"column:tanggal_selesai;not null" json:"tanggal_selesai"`
	Status         string    `gorm:"column:status;type:varchar(20);default:'DRAFT'" json:"status"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`

	// Relasi
	Paslons  []Paslon  `gorm:"foreignKey:IdPemira" json:"paslon,omitempty"`
	Pemilihs []Pemilih `gorm:"foreignKey:IdPemira" json:"pemilih,omitempty"`
	Suaras   []Suara   `gorm:"foreignKey:IdPemira" json:"suara,omitempty"`
}
