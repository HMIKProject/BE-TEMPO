package entity

// KeahlianAnggota merupakan tabel perantara (pivot) antara tabel anggota dan keahlian
type KeahlianAnggota struct {
	IdAnggota         int64    `gorm:"primaryKey;column:id_anggota"`
	IdKeahlian        int64    `gorm:"primaryKey;column:id_keahlian"`
	TingkatPenguasaan int      `gorm:"column:tingkat_penguasaan;not null;default:1"`
	Anggota           Anggota  `gorm:"foreignKey:IdAnggota;references:IdAnggota"`
	Keahlian          Keahlian `gorm:"foreignKey:IdKeahlian;references:IdKeahlian"`
}
