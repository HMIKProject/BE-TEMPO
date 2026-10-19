package entity

import "time"

type Paslon struct {
	IdPaslon       int64     `gorm:"primaryKey;column:id_paslon;autoIncrement" json:"id_paslon"`
	IdPemira       int64     `gorm:"column:id_pemira;not null" json:"id_pemira"`
	NomorUrut      int       `gorm:"column:nomor_urut;not null" json:"nomor_urut"`
	NimKetua       string    `gorm:"column:nim_ketua;type:varchar(20);not null" json:"nim_ketua"`
	NamaKetua      string    `gorm:"column:nama_ketua;type:varchar(100);not null" json:"nama_ketua"`
	NimWakil       string    `gorm:"column:nim_wakil;type:varchar(20);not null" json:"nim_wakil"`
	NamaWakil      string    `gorm:"column:nama_wakil;type:varchar(100);not null" json:"nama_wakil"`
	NamaKabinet    string    `gorm:"column:nama_kabinet;type:varchar(100);not null" json:"nama_kabinet"`
	FotoPaslon     string    `gorm:"column:foto_paslon;type:varchar(255)" json:"foto_paslon"`
	Visi           string    `gorm:"column:visi;type:text" json:"visi"`
	Misi           string    `gorm:"column:misi;type:text" json:"misi"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiperbaruiPada time.Time `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
