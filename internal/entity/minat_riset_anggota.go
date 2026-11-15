package entity

// MinatRisetAnggota merupakan tabel perantara (pivot) antara tabel anggota dan minat_riset
type MinatRisetAnggota struct {
	IdAnggota  int64      `gorm:"primaryKey;column:id_anggota"`
	IdMinat    int64      `gorm:"primaryKey;column:id_minat"`
	MinatRiset MinatRiset `gorm:"foreignKey:IdMinat;references:IdMinat" json:"minat_riset,omitempty"`
}
