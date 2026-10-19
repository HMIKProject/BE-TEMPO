package entity

import "time"

type Kabinet struct {
	IdKabinet      int64      `gorm:"primaryKey;column:id_kabinet;autoIncrement" json:"id_kabinet"`
	NamaKabinet    string     `gorm:"column:nama_kabinet;type:varchar(100);not null" json:"nama_kabinet"`
	TahunMulai     int        `gorm:"column:tahun_mulai;not null" json:"tahun_mulai"`
	TahunSelesai   int        `gorm:"column:tahun_selesai;not null" json:"tahun_selesai"`
	LogoKabinet    string     `gorm:"column:logo_kabinet;type:varchar(255)" json:"logo_kabinet"`
	Visi           string     `gorm:"column:visi;type:text" json:"visi"`
	Misi           string     `gorm:"column:misi;type:text" json:"misi"`
	IsAktif        bool       `gorm:"column:is_aktif;default:false" json:"is_aktif"`
	TanggalMulai   *time.Time `gorm:"column:tanggal_mulai;type:date" json:"tanggal_mulai"`
	TanggalDemis   *time.Time `gorm:"column:tanggal_demis;type:date" json:"tanggal_demis"`
	DibuatPada     time.Time  `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time  `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
