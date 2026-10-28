package entity

import "time"

type Pengguna struct {
	IdPengguna     int64     `gorm:"primaryKey;column:id_pengguna;autoIncrement" json:"id_pengguna"`
	IdPeran        *int64    `gorm:"column:id_peran" json:"id_peran,omitempty"` // Nullable sementara
	NamaPengguna   string    `gorm:"column:nama_pengguna;type:varchar(100)" json:"nama_pengguna"`
	Email          string    `gorm:"column:email;type:varchar(100);not null" json:"email"`
	Sandi          string    `gorm:"column:sandi;type:varchar(100);not null" json:"-"`
	NomorTelepon   string    `gorm:"column:nomor_telepon;type:varchar(255)" json:"-"`
	StatusAktif    bool      `gorm:"column:status_aktif;not null;default:true" json:"status_aktif"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
