package entity

import "time"

type Pemilih struct {
	IdPemilih      int64      `gorm:"primaryKey;column:id_pemilih;autoIncrement" json:"id_pemilih"`
	IdPemira       int64      `gorm:"column:id_pemira;not null" json:"id_pemira"`
	Nim            string     `gorm:"column:nim;type:varchar(20);not null" json:"nim"`
	EmailKampus    string     `gorm:"column:email_kampus;type:varchar(100);not null" json:"email_kampus"`
	StatusMemilih  bool       `gorm:"column:status_memilih;default:false" json:"status_memilih"`
	KodeOtp        string     `gorm:"column:kode_otp;type:varchar(10)" json:"-"` // Hidden dari JSON agar aman
	OtpExpired     *time.Time `gorm:"column:otp_expired" json:"-"`
	DibuatPada     time.Time  `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time  `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
