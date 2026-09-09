package entity

type MinatRiset struct {
	IdMinat   int64  `gorm:"primaryKey;column:id_minat;autoIncrement"`
	NamaMinat string `gorm:"column:nama_minat;type:varchar(100);not null"`
	Deskripsi string `gorm:"column:deskripsi;type:varchar(255)"`
}
