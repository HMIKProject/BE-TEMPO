package entity

import "time"

type Departemen struct {
	IdDepartemen   int64          `gorm:"primaryKey;column:id_departemen;autoIncrement" json:"id_departemen"`
	NamaDepartemen string         `gorm:"column:nama_departemen;type:varchar(255 जीता);not null" json:"nama_departemen"`
	Deskripsi      string         `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	Logo           string         `gorm:"column:logo;type:varchar(255)" json:"logo"`
	DibuatPada     time.Time      `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time      `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`

	// Relasi One-to-Many
	ProgramKerja []ProgramKerja `gorm:"foreignKey:IdDepartemen" json:"program_kerja,omitempty"`
	Anggota      []Anggota      `gorm:"foreignKey:IdDepartemen" json:"anggota,omitempty"`
}
