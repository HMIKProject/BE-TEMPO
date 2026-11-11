package entity

import "time"

type Pengguna struct {
	IdPengguna     int64     `gorm:"primaryKey;column:id_pengguna;autoIncrement" json:"id_pengguna"`
	IdPeran        *int64    `gorm:"column:id_peran" json:"id_peran,omitempty"` // Nullable sementara
	NamaLengkap    string    `gorm:"column:nama_lengkap;type:varchar(255);not null" json:"nama_lengkap"`
	NamaPengguna   string    `gorm:"column:nama_pengguna;type:varchar(100)" json:"nama_pengguna"`
	Email          string    `gorm:"column:email;type:varchar(100);not null" json:"email"`
	Sandi          string    `gorm:"column:sandi;type:varchar(100);not null" json:"-"`
	NomorTelepon   string    `gorm:"column:nomor_telepon;type:varchar(255)" json:"-"`
	FotoPengguna        string     `gorm:"column:foto_pengguna;type:varchar(255)" json:"foto_pengguna,omitempty"`
	StatusAktif         bool       `gorm:"column:status_aktif;not null;default:true" json:"status_aktif"`
	TokenVersion        int        `gorm:"column:token_version;not null;default:1" json:"token_version"`
	KodeAktivasi        *string    `gorm:"column:kode_aktivasi;type:varchar(50)" json:"-"`
	KodeAktivasiExpired *time.Time `gorm:"column:kode_aktivasi_expired" json:"-"`
	IsFirstLogin        bool       `gorm:"column:is_first_login;not null;default:false" json:"is_first_login"`
	DibuatPada          time.Time  `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada      time.Time  `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
