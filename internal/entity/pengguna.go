package entity

import "time"

type Pengguna struct {
	IdPengguna     int64     `gorm:"primaryKey;column:id_pengguna;autoIncrement"`
	IdPeran        *int64    `gorm:"column:id_peran"` // Nullable sementara
	NamaLengkap    string    `gorm:"column:nama_lengkap;type:varchar(100);not null"`
	NamaPengguna   string    `gorm:"column:nama_pengguna;type:varchar(100)"`
	Email          string    `gorm:"column:email;type:varchar(100);not null"`
	Sandi          string    `gorm:"column:sandi;type:varchar(100);not null"`
	NomorTelepon   string    `gorm:"column:nomor_telepon;type:varchar(255)"`
	FotoPengguna   string    `gorm:"column:foto_pengguna;type:varchar(255)"`
	StatusAktif    bool      `gorm:"column:status_aktif;not null;default:true"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada;autoCreateTime"`
	DiperbaruiPada time.Time `gorm:"column:diperbarui_pada;autoUpdateTime"`
}
