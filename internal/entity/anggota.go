package entity

import "time"

type RiwayatJabatan struct {
	IdRiwayat      int64     `gorm:"primaryKey;column:id_riwayat;autoIncrement" json:"id_riwayat"`
	IdAnggota      int64     `gorm:"column:id_anggota;not null" json:"id_anggota"`
	Nim            string    `gorm:"column:nim;not null" json:"nim"`
	Periode        string    `gorm:"column:periode;not null" json:"periode"`
	Jabatan        string    `gorm:"column:jabatan;not null" json:"jabatan"`
	NamaDepartemen string    `gorm:"column:nama_departemen" json:"nama_departemen,omitempty"`
	NamaDivisi     string    `gorm:"column:nama_divisi" json:"nama_divisi,omitempty"`
	Status         string    `gorm:"column:status;not null;default:'Demisioner'" json:"status"`
	DibuatPada     time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
}

func (RiwayatJabatan) TableName() string {
	return "riwayat_jabatan"
}

type Anggota struct {
	IdAnggota    int64     `gorm:"primaryKey;column:id_anggota;autoIncrement" json:"id_anggota"`
	PenggunaID   *int64    `gorm:"column:id_pengguna" json:"id_pengguna,omitempty"` // Nullable untuk Demisioner
	Pengguna     *Pengguna `gorm:"foreignKey:PenggunaID;references:IdPengguna" json:"pengguna,omitempty"`
	Nim          string    `gorm:"column:nim;type:varchar(20);unique;not null" json:"nim"`
	NamaLengkap  string    `gorm:"column:nama_lengkap;type:varchar(100);not null" json:"nama_lengkap"`
	FotoProfil   string    `gorm:"column:foto_profil;type:varchar(255)" json:"foto_profil,omitempty"`
	ProgramStudi string    `gorm:"column:program_studi;type:varchar(100);not null" json:"program_studi"`
	Angkatan     int       `gorm:"column:angkatan;not null" json:"angkatan"`
	Jabatan      string    `gorm:"column:jabatan;type:varchar(100)" json:"jabatan,omitempty"`
	Status       string    `gorm:"column:status;type:varchar(50);not null;default:'Pengurus'" json:"status"` // Pengurus, Magang, Demisioner
	Tautan       string    `gorm:"column:tautan;type:varchar(255)" json:"tautan,omitempty"`
	Files        string    `gorm:"column:files;type:varchar(255)" json:"files,omitempty"`
	Periode      string    `gorm:"column:periode;type:varchar(20)" json:"periode,omitempty"`
	CvUrl        string    `gorm:"column:cv_url;type:varchar(255)" json:"cv_url,omitempty"`
	DibuatPada   time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
	DiubahPada   time.Time `gorm:"column:diubah_pada;autoUpdateTime" json:"diubah_pada"`

	// Relasi One-to-Many (Opsional, karena BPH mungkin tidak ada di departemen spesifik)
	DepartemenID *int64      `gorm:"column:id_departemen" json:"id_departemen,omitempty"`
	Departemen   *Departemen `gorm:"foreignKey:DepartemenID;references:IdDepartemen" json:"departemen,omitempty"`

	DivisiID *int64            `gorm:"column:id_divisi" json:"id_divisi,omitempty"`
	Divisi   *DivisiDepartemen `gorm:"foreignKey:DivisiID;references:IdDivisi" json:"divisi,omitempty"`

	// Relasi One-to-Many ke Tabel Perantara (Pivot)
	KeahlianAnggota   []KeahlianAnggota   `gorm:"foreignKey:IdAnggota" json:"keahlian,omitempty"`
	MinatRisetAnggota []MinatRisetAnggota `gorm:"foreignKey:IdAnggota" json:"minat_riset,omitempty"`

	// Riwayat Jabatan Rekam Jejak
	RiwayatJabatan []RiwayatJabatan `gorm:"foreignKey:IdAnggota" json:"riwayat_jabatan,omitempty"`
}
