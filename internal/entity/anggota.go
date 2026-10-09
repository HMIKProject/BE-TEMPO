package entity

import "time"

type Anggota struct {
	IdAnggota    int64     `gorm:"primaryKey;column:id_anggota;autoIncrement"`
	IdPengguna   *int64    `gorm:"column:id_pengguna"` // Nullable untuk Demisioner
	Pengguna     *Pengguna `gorm:"foreignKey:IdPengguna;references:IdPengguna"`
	Nim          string    `gorm:"column:nim;type:varchar(20);unique;not null"`
	NamaLengkap  string    `gorm:"column:nama_lengkap;type:varchar(100);not null"`
	FotoProfil   string    `gorm:"column:foto_profil;type:varchar(255)"`
	ProgramStudi string    `gorm:"column:program_studi;type:varchar(100);not null"`
	Angkatan     int       `gorm:"column:angkatan;not null"`
	Jabatan      string    `gorm:"column:jabatan;type:varchar(100)"`
	Status       string    `gorm:"column:status;type:varchar(50);not null;default:'Pengurus'"` // Pengurus, Magang, Demisioner
	Tautan       string    `gorm:"column:tautan;type:varchar(255)"`
	Files        string    `gorm:"column:files;type:varchar(255)"`
	DibuatPada   time.Time `gorm:"column:dibuat_pada;autoCreateTime"`
	DiubahPada   time.Time `gorm:"column:diubah_pada;autoUpdateTime"`

	// Relasi One-to-Many (Opsional, karena BPH mungkin tidak ada di departemen spesifik)
	IdDepartemen *int64      `gorm:"column:id_departemen"`
	Departemen   *Departemen `gorm:"foreignKey:IdDepartemen;references:IdDepartemen"`

	IdDivisi     *int64            `gorm:"column:id_divisi"`
	Divisi       *DivisiDepartemen `gorm:"foreignKey:IdDivisi;references:IdDivisi"`

	// Relasi One-to-Many ke Tabel Perantara (Pivot)
	KeahlianAnggota   []KeahlianAnggota   `gorm:"foreignKey:IdAnggota"`
	MinatRisetAnggota []MinatRisetAnggota `gorm:"foreignKey:IdAnggota"`
}
