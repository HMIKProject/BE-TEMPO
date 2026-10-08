package entity

import "time"

type DivisiDepartemen struct {
	IdDivisi       int64          `gorm:"primaryKey;column:id_divisi;autoIncrement" json:"id_divisi"`
	IdDepartemen   int64          `gorm:"column:id_departemen;not null" json:"id_departemen"`
	NamaDivisi     string         `gorm:"column:nama_divisi;type:varchar(255);not null" json:"nama_divisi"`
	Deskripsi      string         `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	DibuatPada     time.Time      `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time      `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`

	// Relasi One-to-Many
	ProgramKerja []ProgramKerja `gorm:"foreignKey:IdDivisi" json:"program_kerja,omitempty"`
	Anggota      []Anggota      `gorm:"foreignKey:IdDivisi" json:"anggota,omitempty"`
}
