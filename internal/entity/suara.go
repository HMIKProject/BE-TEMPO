package entity

import "time"

type Suara struct {
	IdSuara      int64     `gorm:"primaryKey;column:id_suara;autoIncrement" json:"id_suara"`
	IdPemira     int64     `gorm:"column:id_pemira;not null" json:"id_pemira"`
	IdPaslon     int64     `gorm:"column:id_paslon;not null" json:"id_paslon"`
	WaktuMemilih time.Time `gorm:"column:waktu_memilih;autoCreateTime" json:"waktu_memilih"`
}
