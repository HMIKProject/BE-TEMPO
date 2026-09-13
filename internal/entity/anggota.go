package entity

import "time"

type Anggota struct {
	IdAnggota    int64     `gorm:"primaryKey;column:id_anggota;autoIncrement"`
	IdUser       int64     `gorm:"column:id_user;not null"`
	Pengguna     Pengguna  `gorm:"foreignKey:IdUser;references:IdPengguna"`
	Nim          string    `gorm:"column:nim;type:varchar(20);unique;not null"`
	ProgramStudi string    `gorm:"column:program_studi;type:varchar(100);not null"`
	Angkatan     int       `gorm:"column:angkatan;not null"`
	Jabatan      string    `gorm:"column:jabatan;type:varchar(100)"`
	Status       bool      `gorm:"column:status;not null;default:true"`
	Tautan       string    `gorm:"column:tautan;type:varchar(255)"`
	Files        string    `gorm:"column:files;type:varchar(255)"`
	DibuatPada   time.Time `gorm:"column:dibuat_pada;autoCreateTime"`
	DiubahPada   time.Time `gorm:"column:diubah_pada;autoUpdateTime"`

	// Relasi One-to-Many (Opsional, karena BPH mungkin tidak ada di departemen spesifik)
	IdDepartemen *int64      `gorm:"column:id_departemen"`
	Departemen   *Departemen `gorm:"foreignKey:IdDepartemen;references:IdDepartemen"`

	// Relasi One-to-Many ke Tabel Perantara (Pivot)
	KeahlianAnggota   []KeahlianAnggota   `gorm:"foreignKey:IdAnggota"`
	MinatRisetAnggota []MinatRisetAnggota `gorm:"foreignKey:IdAnggota"`
}
