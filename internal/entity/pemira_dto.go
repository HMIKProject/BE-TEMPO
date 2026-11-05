package entity

import "time"

// ==========================================
// VOTER AUTH DTOs
// ==========================================

// RequestOTPRequest merepresentasikan payload saat mahasiswa meminta OTP di form login
type RequestOTPRequest struct {
	NIM         string `json:"nim" binding:"required" example:"10123001"`
	EmailKampus string `json:"email_kampus" binding:"required,email" example:"10123001@student.universitaspertamina.ac.id"`
}

// RequestOTPResponse merepresentasikan kembalian setelah OTP berhasil di-generate & dikirim
type RequestOTPResponse struct {
	NIM              string `json:"nim" example:"10123001"`
	EmailTujuan      string `json:"email_tujuan" example:"10123001@student.universitaspertamina.ac.id"`
	MasaBerlakuMenit int    `json:"masa_berlaku_menit" example:"5"`
	CooldownDetik    int    `json:"cooldown_detik" example:"60"`
}

// VerifyOTPRequest merepresentasikan input kode 6-digit OTP oleh pemilih
type VerifyOTPRequest struct {
	NIM     string `json:"nim" binding:"required" example:"10123001"`
	KodeOTP string `json:"kode_otp" binding:"required,len=6" example:"748291"`
}

// BallotVoterInfo informasi ringkas pemilih di dalam respons verifikasi
type BallotVoterInfo struct {
	NIM           string `json:"nim" example:"10123001"`
	StatusMemilih bool   `json:"status_memilih" example:"false"`
}

// VerifyOTPResponse merepresentasikan token bilik suara (Ballot Token) yang diterbitkan
type VerifyOTPResponse struct {
	BallotToken string          `json:"ballot_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int             `json:"expires_in" example:"900"` // 15 menit (dalam detik)
	Pemilih     BallotVoterInfo `json:"pemilih"`
}

// ResendOTPRequest payload saat pemilih meminta pengiriman ulang OTP
type ResendOTPRequest struct {
	NIM string `json:"nim" binding:"required" example:"10123001"`
}

// ==========================================
// VOTING & BALLOT DTOs
// ==========================================

// CastVoteRequest payload pencoblosan paslon oleh pemilih
type CastVoteRequest struct {
	IdPaslon int64 `json:"id_paslon" binding:"required" example:"1"`
}

// CastVoteResponse hasil konfirmasi pencoblosan suara
type CastVoteResponse struct {
	StatusVote    string    `json:"status_vote" example:"COMPLETED"`
	WaktuTercatat time.Time `json:"waktu_tercatat" example:"2026-11-04T15:30:00Z"`
}

// ==========================================
// ADMIN MONITORING DTOs
// ==========================================

// ElectionStatsResponse ringkasan 4 kartu metrik statistik pemilu
type ElectionStatsResponse struct {
	TotalSuaraMasuk       int64   `json:"total_suara_masuk" example:"327"`
	TotalDPT              int64   `json:"total_dpt" example:"450"`
	PersentasePartisipasi float64 `json:"persentase_partisipasi" example:"72.67"`
	TotalPaslon           int64   `json:"total_paslon" example:"2"`
}

// ElectionChartCandidate rincian perolehan suara per paslon untuk chart
type ElectionChartCandidate struct {
	IdPaslon   int64   `json:"id_paslon" example:"1"`
	NomorUrut  int     `json:"nomor_urut" example:"1"`
	NamaPaslon string  `json:"nama_paslon" example:"Nanami Kento & Loid Forger"`
	FotoPaslon string  `json:"foto_paslon" example:"https://res.cloudinary.com/.../paslon1.png"`
	TotalSuara int64   `json:"total_suara" example:"222"`
	Persentase float64 `json:"persentase" example:"67.89"`
}

// ElectionChartResponse data grafik batang pleno pemilu
type ElectionChartResponse struct {
	TotalSuara int64                    `json:"total_suara" example:"327"`
	IsLive     bool                     `json:"is_live" example:"true"`
	Candidates []ElectionChartCandidate `json:"candidates"`
}

// VoterListResponse data tabel pemilih DPT untuk admin
type VoterListResponse struct {
	Voters    []Pemilih `json:"voters"`
	TotalData int64     `json:"total_data" example:"450"`
}

// ==========================================
// ADMIN CANDIDATE (PASLON) CRUD DTOs
// ==========================================

// CreatePaslonRequest payload penambahan paslon baru oleh admin
type CreatePaslonRequest struct {
	NomorUrut   int    `json:"nomor_urut" binding:"required" example:"3"`
	NimKetua    string `json:"nim_ketua" binding:"required" example:"105223001"`
	NamaKetua   string `json:"nama_ketua" binding:"required" example:"Calon Ketua"`
	NimWakil    string `json:"nim_wakil" binding:"required" example:"105223002"`
	NamaWakil   string `json:"nama_wakil" binding:"required" example:"Calon Wakil"`
	NamaKabinet string `json:"nama_kabinet" binding:"required" example:"Kabinet Bersama"`
	FotoPaslon  string `json:"foto_paslon" example:"https://..."`
	Visi        string `json:"visi" example:"Mewujudkan HMIK unggul"`
	Misi        string `json:"misi" example:"1. Kolaborasi riset"`
}

// UpdatePaslonRequest payload edit paslon oleh admin
type UpdatePaslonRequest struct {
	NomorUrut   int    `json:"nomor_urut" binding:"required" example:"3"`
	NimKetua    string `json:"nim_ketua" binding:"required" example:"105223001"`
	NamaKetua   string `json:"nama_ketua" binding:"required" example:"Calon Ketua"`
	NimWakil    string `json:"nim_wakil" binding:"required" example:"105223002"`
	NamaWakil   string `json:"nama_wakil" binding:"required" example:"Calon Wakil"`
	NamaKabinet string `json:"nama_kabinet" binding:"required" example:"Kabinet Bersama"`
	FotoPaslon  string `json:"foto_paslon" example:"https://..."`
	Visi        string `json:"visi" example:"Mewujudkan HMIK unggul"`
	Misi        string `json:"misi" example:"1. Kolaborasi riset"`
}

// ==========================================
// ADMIN ELECTION CONTROL & MPIK DTOs
// ==========================================

// ElectionSettingsResponse data parameter & kontrol acara pemira
type ElectionSettingsResponse struct {
	IdPemira       int64     `json:"id_pemira" example:"1"`
	TahunPeriode   string    `json:"tahun_periode" example:"2026/2027"`
	Status         string    `json:"status" example:"BERJALAN"`
	TanggalMulai   time.Time `json:"tanggal_mulai"`
	TanggalSelesai time.Time `json:"tanggal_selesai"`
	TotalPaslon    int64     `json:"total_paslon" example:"2"`
	TotalSuara     int64     `json:"total_suara" example:"327"`
	MpikStatus     string    `json:"mpik_status" example:"PENDING"`
	WinningPaslon  *Paslon   `json:"winning_paslon,omitempty"`
}

// UpdatePemiraStatusRequest payload saklar pemicu status pemira
type UpdatePemiraStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=DRAFT BERJALAN SELESAI" example:"BERJALAN"`
}

// UpdatePemiraScheduleRequest payload pengaturan jadwal acara
type UpdatePemiraScheduleRequest struct {
	TahunPeriode   string    `json:"tahun_periode" binding:"required" example:"2026/2027"`
	TanggalMulai   time.Time `json:"tanggal_mulai" binding:"required"`
	TanggalSelesai time.Time `json:"tanggal_selesai" binding:"required"`
}

// PelantikanMPIKResponse konfirmasi pelantikan kabinet dan BPH baru pasca sidang MPIK
type PelantikanMPIKResponse struct {
	StatusPelantikan string   `json:"status_pelantikan" example:"APPROVED"`
	Pesan            string   `json:"pesan" example:"Pelantikan berhasil! Kabinet dan BPH baru telah aktif di Manajemen Himpunan."`
	KabinetBaru      Kabinet  `json:"kabinet_baru"`
	KetuaBaru        *Anggota `json:"ketua_baru,omitempty"`
	WakilBaru        *Anggota `json:"wakil_baru,omitempty"`
}
