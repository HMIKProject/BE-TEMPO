package entity

import "time"

// --- DTO BUKU INDUK ANGGOTA (BIA) ---

// AnggotaFilter parameter query untuk menyaring daftar anggota
type AnggotaFilter struct {
	Search   string `form:"search" example:"Anton"`
	Status   string `form:"status" example:"Pengurus"` // Pengurus, Magang, Demisioner, Semua
	Angkatan int    `form:"angkatan" example:"2023"`
	Page     int    `form:"page,default=1" example:"1"`
	Limit    int    `form:"limit,default=50" example:"50"`
}

// CreateAnggotaRequest payload tambah anggota baru manual
type CreateAnggotaRequest struct {
	Nim          string `json:"nim" binding:"required" example:"105223046"`
	NamaLengkap  string `json:"nama_lengkap" binding:"required" example:"Anton Dicaprio"`
	FotoProfil   string `json:"foto_profil" example:"https://res.cloudinary.com/.../profile.png"`
	ProgramStudi string `json:"program_studi" binding:"required" example:"Ilmu Komputer"`
	Angkatan     int    `json:"angkatan" binding:"required" example:"2023"`
	Jabatan      string `json:"jabatan" example:"Ketua Divisi Media Kreatif"`
	Status       string `json:"status" binding:"required,oneof=Pengurus Magang Demisioner" example:"Pengurus"`
	Tautan       string `json:"tautan" example:"https://www.linkedin.com/in/antondicaprio"`
	CvUrl        string `json:"cv_url" example:"https://res.cloudinary.com/.../cv.pdf"`
	IdDepartemen *int64 `json:"id_departemen" example:"1"`
	IdDivisi     *int64 `json:"id_divisi" example:"2"`
	Periode      string `json:"periode" example:"2026/2027"`
}

// UpdateAnggotaRequest payload perbarui data anggota
type UpdateAnggotaRequest struct {
	NamaLengkap  string `json:"nama_lengkap" binding:"required" example:"Anton Dicaprio"`
	FotoProfil   string `json:"foto_profil" example:"https://res.cloudinary.com/.../profile.png"`
	ProgramStudi string `json:"program_studi" binding:"required" example:"Ilmu Komputer"`
	Angkatan     int    `json:"angkatan" binding:"required" example:"2023"`
	Jabatan      string `json:"jabatan" example:"Ketua Divisi Media Kreatif"`
	Status       string `json:"status" binding:"required,oneof=Pengurus Magang Demisioner" example:"Pengurus"`
	Tautan       string `json:"tautan" example:"https://www.linkedin.com/in/antondicaprio"`
	CvUrl        string `json:"cv_url" example:"https://res.cloudinary.com/.../cv.pdf"`
	IdDepartemen *int64 `json:"id_departemen" example:"1"`
	IdDivisi     *int64 `json:"id_divisi" example:"2"`
	Periode      string `json:"periode" example:"2026/2027"`
}

// AnggotaListResponse pembungkus list anggota beserta total count pagination
type AnggotaListResponse struct {
	TotalData int64     `json:"total_data"`
	Anggota   []Anggota `json:"anggota"`
}

// --- DTO SLOT AKUN PINJAM PAKAI ---

// SlotPeminjamInfo data ringkas anggota yang sedang meminjam slot
type SlotPeminjamInfo struct {
	IdAnggota   int64  `json:"id_anggota"`
	Nim         string `json:"nim"`
	NamaLengkap string `json:"nama_lengkap"`
	Angkatan    int    `json:"angkatan"`
	Periode     string `json:"periode"`
	FotoProfil  string `json:"foto_profil"`
}

// AccountSlotResponse representasi kartu slot akun di dashboard admin
type AccountSlotResponse struct {
	IdPengguna          int64             `json:"id_pengguna"`
	NamaPengguna        string            `json:"nama_pengguna"`
	Email               string            `json:"email"`
	Jabatan             string            `json:"jabatan"`
	StatusSlot          string            `json:"status_slot"` // "TERPAKAI" atau "KOSONG"
	Peminjam            *SlotPeminjamInfo `json:"peminjam,omitempty"`
	KodeAktivasi        *string           `json:"kode_aktivasi,omitempty"`
	KodeAktivasiExpired *time.Time        `json:"kode_aktivasi_expired,omitempty"`
	IsFirstLogin        bool              `json:"is_first_login"`
}

// AssignSlotRequest payload menautkan slot kosong ke mahasiswa BIA
type AssignSlotRequest struct {
	IdPengguna int64 `json:"id_pengguna" binding:"required" example:"1"`
	IdAnggota  int64 `json:"id_anggota" binding:"required" example:"5"`
}

// AssignSlotResponse kembalian sukses penautan akun dengan Kode Aktivasi
type AssignSlotResponse struct {
	IdPengguna          int64     `json:"id_pengguna"`
	Email               string    `json:"email"`
	NamaPengguna        string    `json:"nama_pengguna"`
	Jabatan             string    `json:"jabatan"`
	NamaAnggota         string    `json:"nama_anggota"`
	NimAnggota          string    `json:"nim_anggota"`
	KodeAktivasi        string    `json:"kode_aktivasi"`
	KodeAktivasiExpired time.Time `json:"kode_aktivasi_expired"`
}

// UnassignSlotRequest payload melepas kaitan akun dari anggota
type UnassignSlotRequest struct {
	IdPengguna int64 `json:"id_pengguna" binding:"required" example:"1"`
}

// ResetSlotResponse kembalian terbit ulang kode aktivasi jika lupa sandi
type ResetSlotResponse struct {
	IdPengguna          int64     `json:"id_pengguna"`
	Email               string    `json:"email"`
	KodeAktivasi        string    `json:"kode_aktivasi"`
	KodeAktivasiExpired time.Time `json:"kode_aktivasi_expired"`
}

// DemisionerSummary ringkasan eksekusi demisioner massal akhir periode
type DemisionerSummary struct {
	TotalAkunDireset     int64     `json:"total_akun_direset"`
	TotalAnggotaDemis    int64     `json:"total_anggota_demis"`
	WaktuEksekusi        time.Time `json:"waktu_eksekusi"`
}

// UpdateSlotRoleRequest payload ubah peran hak akses slot akun
type UpdateSlotRoleRequest struct {
	IdPeran int64 `json:"id_peran" binding:"required" example:"1"`
}

// --- DTO STRUKTUR ORGANISASI & DEPARTEMEN ---

// BPHMemberInfo data ringkas pengurus inti BPH (Kahim, Wakahim, Sekre, Benda)
type BPHMemberInfo struct {
	IdAnggota    int64  `json:"id_anggota"`
	NamaLengkap  string `json:"nama_lengkap"`
	Nim          string `json:"nim"`
	Jabatan      string `json:"jabatan"`
	FotoProfil   string `json:"foto_profil"`
	Angkatan     int    `json:"angkatan"`
	ProgramStudi string `json:"program_studi"`
	Status       string `json:"status"`
	Periode      string `json:"periode"`
}

// DeptKadepInfo data kepala departemen
type DeptKadepInfo struct {
	IdAnggota   int64  `json:"id_anggota"`
	NamaLengkap string `json:"nama_lengkap"`
	Nim         string `json:"nim"`
	Jabatan     string `json:"jabatan"`
	FotoProfil  string `json:"foto_profil"`
	Angkatan    int    `json:"angkatan"`
}

// OrganisasiDeptItem representasi departemen dalam struktur organisasi
type OrganisasiDeptItem struct {
	IdDepartemen   int64              `json:"id_departemen"`
	NamaDepartemen string             `json:"nama_departemen"`
	Deskripsi      string             `json:"deskripsi"`
	Logo           string             `json:"logo"`
	StatusAktif    bool               `json:"status_aktif"`
	Kadep          *DeptKadepInfo     `json:"kadep,omitempty"`
	Divisi         []DivisiDepartemen `json:"divisi"`
	TotalAnggota   int64              `json:"total_anggota"`
	TotalDivisi    int64              `json:"total_divisi"`
}

// StrukturOrganisasiResponse pembungkus susunan organisasi lengkap
type StrukturOrganisasiResponse struct {
	Kabinet    *Kabinet             `json:"kabinet,omitempty"`
	BPH        []BPHMemberInfo      `json:"bph"`
	Departemen []OrganisasiDeptItem `json:"departemen"`
}

// CreateDeptRequest payload tambah departemen baru
type CreateDeptRequest struct {
	NamaDepartemen string `json:"nama_departemen" binding:"required" example:"Riset dan Teknologi"`
	Deskripsi      string `json:"deskripsi" example:"Departemen yang membawahi bidang teknologi dan inovasi"`
	Logo           string `json:"logo" example:"https://res.cloudinary.com/.../logo.png"`
}

// UpdateDeptRequest payload perbarui departemen
type UpdateDeptRequest struct {
	NamaDepartemen string `json:"nama_departemen" binding:"required" example:"Riset dan Teknologi"`
	Deskripsi      string `json:"deskripsi" example:"Departemen yang membawahi bidang teknologi dan inovasi"`
	Logo           string `json:"logo" example:"https://res.cloudinary.com/.../logo.png"`
	StatusAktif    *bool  `json:"status_aktif" example:"true"`
}

// CreateDivisiRequest payload tambah divisi baru
type CreateDivisiRequest struct {
	IdDepartemen int64  `json:"id_departemen" binding:"required" example:"1"`
	NamaDivisi   string `json:"nama_divisi" binding:"required" example:"Riset & Inovasi"`
	Deskripsi    string `json:"deskripsi" example:"Divisi penelitian dan eksplorasi ilmiah"`
}

// UpdateDivisiRequest payload perbarui divisi
type UpdateDivisiRequest struct {
	NamaDivisi  string `json:"nama_divisi" binding:"required" example:"Riset & Inovasi"`
	Deskripsi   string `json:"deskripsi" example:"Divisi penelitian dan eksplorasi ilmiah"`
	StatusAktif *bool  `json:"status_aktif" example:"true"`
}

