package entity

import "time"

type ProgramKerja struct {
	IdProker       int64      `gorm:"primaryKey;column:id_proker;autoIncrement" json:"id_proker"`
	IdDepartemen   int64      `gorm:"column:id_departemen;not null" json:"id_departemen"`
	NamaProker     string     `gorm:"column:nama_proker;type:varchar(255);not null" json:"nama_proker"`
	Deskripsi      string     `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	Foto           string     `gorm:"column:foto;type:varchar(255)" json:"foto"`
	Status         string     `gorm:"column:status;type:varchar(50);default:'Belum Terlaksana'" json:"status"`
	IsUnggulan     bool       `gorm:"column:is_unggulan;default:false" json:"is_unggulan"`
	DibuatPada     time.Time  `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time  `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
