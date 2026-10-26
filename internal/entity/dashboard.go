package entity

// CardDetail merepresentasikan data lengkap sebuah kartu metrik di Dashboard Admin (Gambar 2 UI)
type CardDetail struct {
	CardKey             string `json:"card_key"`
	Judul               string `json:"judul"`
	Total               int64  `json:"total"`
	PertumbuhanBulanIni int64  `json:"pertumbuhan_bulan_ini"`
	Keterangan          string `json:"keterangan"`
	WarnaTema           string `json:"warna_tema"`   // Contoh: "blue", "yellow", "green", "orange"
	TautanAksi          string `json:"tautan_aksi"`   // Rute navigasi: "/admin/buku-induk-anggota", dll.
}

// StatItem merepresentasikan data sebuah kartu metrik dengan angka total dan pertumbuhannya
type StatItem struct {
	Total               int64  `json:"total"`
	PertumbuhanBulanIni int64  `json:"pertumbuhan_bulan_ini"`
	Keterangan          string `json:"keterangan,omitempty"`
}

// VotingStatItem merepresentasikan metrik voting pada dashboard
type VotingStatItem struct {
	Total          int64  `json:"total"`
	SedangBerjalan int64  `json:"sedang_berjalan"`
	Keterangan     string `json:"keterangan,omitempty"`
}

// DashboardStatsResponse adalah DTO agregasi metrik untuk Halaman Utama Dashboard Admin (Gambar 2 UI)
type DashboardStatsResponse struct {
	TotalMahasiswaAktif int64          `json:"total_mahasiswa_aktif"`
	TotalModul          int64          `json:"total_modul"`
	TotalProgramKerja   int64          `json:"total_program_kerja"`
	VotingAktif         int64          `json:"voting_aktif"`
	MahasiswaAktif      StatItem       `json:"mahasiswa_aktif"`
	Modul               StatItem       `json:"modul"`
	ProgramKerja        StatItem       `json:"program_kerja"`
	Voting              VotingStatItem `json:"voting"`

	// Detail visual kartu spesifik
	CardTotalAnggotaAktif CardDetail `json:"card_total_anggota_aktif"`
	CardTotalModul        CardDetail `json:"card_total_modul"`
	CardTotalProgramKerja CardDetail `json:"card_total_program_kerja"`
	CardVotingAktif       CardDetail `json:"card_voting_aktif"`
}

// Dashboard alias untuk backwards-compatibility
type Dashboard = DashboardStatsResponse
