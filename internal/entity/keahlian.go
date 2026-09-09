package entity

type Keahlian struct {
	IdKeahlian   int64  `gorm:"primaryKey;column:id_keahlian;autoIncrement"`
	NamaKeahlian string `gorm:"column:nama_keahlian;type:varchar(100);not null"`
	Deskripsi    string `gorm:"column:deskripsi;type:varchar(255)"`
}
