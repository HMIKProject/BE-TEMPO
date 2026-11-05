package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

// PemiraUsecase mendefinisikan seluruh logika bisnis untuk alur Pemilihan Raya (Pemilih & Admin)
type PemiraUsecase interface {
	// 1. Voter Authentication
	RequestOTP(ctx context.Context, req entity.RequestOTPRequest) (*entity.RequestOTPResponse, error)
	VerifyOTP(ctx context.Context, req entity.VerifyOTPRequest) (*entity.VerifyOTPResponse, error)
	ResendOTP(ctx context.Context, req entity.ResendOTPRequest) (*entity.RequestOTPResponse, error)

	// 2. Voting Booth (Bilik Suara)
	GetActiveCandidates(ctx context.Context) ([]entity.Paslon, error)
	CastVote(ctx context.Context, nim string, idPemira int64, req entity.CastVoteRequest) (*entity.CastVoteResponse, error)

	// 3. Admin CMS Monitoring
	GetElectionStats(ctx context.Context) (*entity.ElectionStatsResponse, error)
	GetElectionChart(ctx context.Context) (*entity.ElectionChartResponse, error)
	GetVoters(ctx context.Context, search string, status *bool, limit, offset int) (*entity.VoterListResponse, error)
	ResetVoter(ctx context.Context, id int64) error

	// 4. Admin Candidate (Paslon) CRUD
	CreateCandidate(ctx context.Context, req entity.CreatePaslonRequest) (*entity.Paslon, error)
	UpdateCandidate(ctx context.Context, id int64, req entity.UpdatePaslonRequest) (*entity.Paslon, error)
	DeleteCandidate(ctx context.Context, id int64) error

	// 5. Admin Settings & Status Trigger
	GetElectionSettings(ctx context.Context) (*entity.ElectionSettingsResponse, error)
	UpdateElectionSchedule(ctx context.Context, req entity.UpdatePemiraScheduleRequest) error
	UpdateElectionStatus(ctx context.Context, req entity.UpdatePemiraStatusRequest) error

	// 6. Pelantikan MPIK & Transisi Manajemen Himpunan
	PelantikanMPIK(ctx context.Context) (*entity.PelantikanMPIKResponse, error)
}

type pemiraUsecase struct {
	pemiraRepo repository.PemiraRepository
}

// NewPemiraUsecase membuat instance baru PemiraUsecase
func NewPemiraUsecase(repo repository.PemiraRepository) PemiraUsecase {
	return &pemiraUsecase{pemiraRepo: repo}
}

// ==========================================
// 1. VOTER AUTHENTICATION USECASES
// ==========================================

func (u *pemiraUsecase) RequestOTP(ctx context.Context, req entity.RequestOTPRequest) (*entity.RequestOTPResponse, error) {
	// 1. Validasi Domain Email Institusi Mahasiswa
	email := strings.TrimSpace(strings.ToLower(req.EmailKampus))
	if !strings.HasSuffix(email, "@student.universitaspertamina.ac.id") {
		return nil, errors.New("hanya email resmi mahasiswa (@student.universitaspertamina.ac.id) yang diizinkan masuk")
	}

	// 2. Cek Periode Pemira yang Sedang Aktif
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		return nil, errors.New("tidak ada pemilihan raya yang sedang aktif saat ini")
	}

	now := time.Now().UTC()
	if now.Before(pemira.TanggalMulai.UTC()) || now.After(pemira.TanggalSelesai.UTC()) {
		return nil, errors.New("bilik suara sedang berada di luar jadwal pemungutan suara resmi")
	}

	// 3. Cari Data Pemilih di DPT (Real-Time Auto-Registration)
	pemilih, err := u.pemiraRepo.FindPemilihByNimAndPemira(ctx, req.NIM, pemira.IdPemira)
	if err != nil {
		// Mahasiswa baru pertama kali meminta OTP -> Daftarkan secara otomatis ke DPT Real-Time
		newPemilih := &entity.Pemilih{
			IdPemira:      pemira.IdPemira,
			Nim:           req.NIM,
			EmailKampus:   email,
			StatusMemilih: false,
		}
		if err := u.pemiraRepo.CreatePemilih(ctx, newPemilih); err != nil {
			return nil, errors.New("gagal mendaftarkan data pemilih")
		}
		pemilih = newPemilih
	} else {
		// 4. Jika sudah pernah terdaftar, validasi kesesuaian email
		if strings.ToLower(pemilih.EmailKampus) != email {
			return nil, errors.New("email yang Anda masukkan tidak sesuai dengan email yang terdaftar untuk NIM ini")
		}
	}

	// 5. Validasi Anti Double-Voting
	if pemilih.StatusMemilih {
		return nil, errors.New("hak suara Anda sudah digunakan sebelumnya. Akses ke bilik suara ditutup")
	}

	// 6. Validasi Cooldown Resend (60 Detik)
	if pemilih.OtpExpired != nil {
		expiredAt := pemilih.OtpExpired.UTC()
		// OTP berlaku 5 menit, berarti waktu pembuatan = expiredAt - 5 menit
		createdAt := expiredAt.Add(-5 * time.Minute)

		elapsed := now.Sub(createdAt)
		if elapsed >= 0 && elapsed < 60*time.Second {
			sisaDetik := int((60*time.Second - elapsed).Seconds()) + 1
			if sisaDetik > 60 {
				sisaDetik = 60
			}
			return nil, fmt.Errorf("harap tunggu %d detik sebelum meminta kode OTP kembali", sisaDetik)
		}
	}

	// 7. Generate OTP Kriptografis 6 Digit
	otp, err := utils.GenerateCryptoOTP(6)
	if err != nil {
		return nil, errors.New("gagal memproses pembuatan kode verifikasi")
	}

	// 8. Simpan OTP ke Database dengan Masa Berlaku 5 Menit (UTC)
	expiredAt := now.Add(5 * time.Minute)
	if err := u.pemiraRepo.UpdateOTP(ctx, pemilih.IdPemilih, otp, expiredAt); err != nil {
		return nil, errors.New("gagal memperbarui kode verifikasi pemilih")
	}

	// 9. Kirim OTP via Mailer (Dual-Mode: SMTP / Console Fallback)
	if err := utils.SendOTPEmail(pemilih.EmailKampus, pemilih.Nim, otp); err != nil {
		return nil, errors.New("gagal mengirimkan email verifikasi, silakan coba beberapa saat lagi")
	}

	return &entity.RequestOTPResponse{
		NIM:              pemilih.Nim,
		EmailTujuan:      pemilih.EmailKampus,
		MasaBerlakuMenit: 5,
		CooldownDetik:    60,
	}, nil
}

func (u *pemiraUsecase) VerifyOTP(ctx context.Context, req entity.VerifyOTPRequest) (*entity.VerifyOTPResponse, error) {
	// 1. Cek Periode Pemira Aktif
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		return nil, errors.New("tidak ada pemilihan raya yang sedang aktif saat ini")
	}

	// 2. Cari Data Pemilih di DPT
	pemilih, err := u.pemiraRepo.FindPemilihByNimAndPemira(ctx, req.NIM, pemira.IdPemira)
	if err != nil {
		return nil, errors.New("NIM tidak ditemukan dalam DPT")
	}

	// 3. Validasi Hak Suara
	if pemilih.StatusMemilih {
		return nil, errors.New("hak suara Anda sudah digunakan sebelumnya")
	}

	// 4. Validasi Keberadaan & Masa Berlaku OTP
	if pemilih.KodeOtp == "" || pemilih.OtpExpired == nil {
		return nil, errors.New("kode OTP belum pernah diminta, silakan lakukan login kembali")
	}

	now := time.Now().UTC()
	if now.After(pemilih.OtpExpired.UTC()) {
		return nil, errors.New("kode OTP telah kedaluwarsa (melebihi 5 menit), silakan minta kode baru")
	}

	// 5. Cocokkan Kode OTP
	if pemilih.KodeOtp != req.KodeOTP {
		return nil, errors.New("kode OTP yang Anda masukkan salah")
	}

	// 6. Hanguskan Kode OTP Seketika (One-Time Invalidation)
	_ = u.pemiraRepo.ClearOTP(ctx, pemilih.IdPemilih)

	// 7. Terbitkan Ballot JWT Token (Masa Aktif 15 Menit)
	ballotToken, err := u.generateBallotJWT(pemilih.Nim, pemira.IdPemira)
	if err != nil {
		return nil, errors.New("gagal menerbitkan tiket masuk bilik suara")
	}

	return &entity.VerifyOTPResponse{
		BallotToken: ballotToken,
		ExpiresIn:   900, // 15 menit dalam detik
		Pemilih: entity.BallotVoterInfo{
			NIM:           pemilih.Nim,
			StatusMemilih: false,
		},
	}, nil
}

func (u *pemiraUsecase) ResendOTP(ctx context.Context, req entity.ResendOTPRequest) (*entity.RequestOTPResponse, error) {
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		return nil, errors.New("tidak ada pemilihan raya yang sedang aktif saat ini")
	}

	pemilih, err := u.pemiraRepo.FindPemilihByNimAndPemira(ctx, req.NIM, pemira.IdPemira)
	if err != nil {
		return nil, errors.New("NIM Anda belum pernah meminta kode OTP sebelumnya")
	}

	return u.RequestOTP(ctx, entity.RequestOTPRequest{
		NIM:         pemilih.Nim,
		EmailKampus: pemilih.EmailKampus,
	})
}

// ==========================================
// 2. VOTING BOOTH (BILIK SUARA) USECASES
// ==========================================

func (u *pemiraUsecase) GetActiveCandidates(ctx context.Context) ([]entity.Paslon, error) {
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		return nil, errors.New("tidak ada pemilihan raya yang sedang aktif saat ini")
	}

	paslons, err := u.pemiraRepo.GetPaslonByPemiraID(ctx, pemira.IdPemira)
	if err != nil {
		return nil, errors.New("gagal memuat data pasangan calon")
	}

	return paslons, nil
}

func (u *pemiraUsecase) CastVote(ctx context.Context, nim string, idPemira int64, req entity.CastVoteRequest) (*entity.CastVoteResponse, error) {
	// 1. Validasi Keberadaan Paslon
	paslon, err := u.pemiraRepo.GetPaslonByID(ctx, req.IdPaslon)
	if err != nil || paslon.IdPemira != idPemira {
		return nil, errors.New("pasangan calon yang Anda pilih tidak valid untuk periode pemilu ini")
	}

	// 2. Eksekusi Transaksi Atomik (SELECT FOR UPDATE + LUBER JURDIL INSERT)
	if err := u.pemiraRepo.CastVote(ctx, idPemira, nim, req.IdPaslon); err != nil {
		return nil, err
	}

	return &entity.CastVoteResponse{
		StatusVote:    "COMPLETED",
		WaktuTercatat: time.Now(),
	}, nil
}

// ==========================================
// 3. ADMIN CMS MONITORING USECASES
// ==========================================

func (u *pemiraUsecase) GetElectionStats(ctx context.Context) (*entity.ElectionStatsResponse, error) {
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		// Jika belum ada pemira berjalan, kembalikan data 0 tanpa crash
		return &entity.ElectionStatsResponse{
			TotalSuaraMasuk:       0,
			TotalDPT:              0,
			PersentasePartisipasi: 0,
			TotalPaslon:           0,
		}, nil
	}

	totalSuara, _ := u.pemiraRepo.CountTotalSuara(ctx, pemira.IdPemira)
	totalDPT, _ := u.pemiraRepo.CountTotalDPT(ctx, pemira.IdPemira)
	paslons, _ := u.pemiraRepo.GetPaslonByPemiraID(ctx, pemira.IdPemira)

	var partisipasi float64 = 0
	if totalDPT > 0 {
		partisipasi = math.Round((float64(totalSuara)/float64(totalDPT)*100)*100) / 100
	}

	return &entity.ElectionStatsResponse{
		TotalSuaraMasuk:       totalSuara,
		TotalDPT:              totalDPT,
		PersentasePartisipasi: partisipasi,
		TotalPaslon:           int64(len(paslons)),
	}, nil
}

func (u *pemiraUsecase) GetElectionChart(ctx context.Context) (*entity.ElectionChartResponse, error) {
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		return &entity.ElectionChartResponse{
			TotalSuara: 0,
			IsLive:     false,
			Candidates: []entity.ElectionChartCandidate{},
		}, nil
	}

	totalSuara, _ := u.pemiraRepo.CountTotalSuara(ctx, pemira.IdPemira)
	suaraMap, _ := u.pemiraRepo.CountSuaraPerPaslon(ctx, pemira.IdPemira)
	paslons, _ := u.pemiraRepo.GetPaslonByPemiraID(ctx, pemira.IdPemira)

	candidates := make([]entity.ElectionChartCandidate, 0, len(paslons))
	for _, p := range paslons {
		count := suaraMap[p.IdPaslon]
		var pct float64 = 0
		if totalSuara > 0 {
			pct = math.Round((float64(count)/float64(totalSuara)*100)*100) / 100
		}

		nama := p.NamaKetua
		if p.NamaWakil != "" {
			nama = p.NamaKetua + " & " + p.NamaWakil
		}

		candidates = append(candidates, entity.ElectionChartCandidate{
			IdPaslon:   p.IdPaslon,
			NomorUrut:  p.NomorUrut,
			NamaPaslon: nama,
			FotoPaslon: p.FotoPaslon,
			TotalSuara: count,
			Persentase: pct,
		})
	}

	return &entity.ElectionChartResponse{
		TotalSuara: totalSuara,
		IsLive:     pemira.Status == "BERJALAN",
		Candidates: candidates,
	}, nil
}

func (u *pemiraUsecase) GetVoters(ctx context.Context, search string, status *bool, limit, offset int) (*entity.VoterListResponse, error) {
	pemira, err := u.pemiraRepo.GetActivePemira(ctx)
	if err != nil {
		return &entity.VoterListResponse{
			Voters:    []entity.Pemilih{},
			TotalData: 0,
		}, nil
	}

	if limit <= 0 {
		limit = 10
	}

	voters, total, err := u.pemiraRepo.GetPemilihList(ctx, pemira.IdPemira, search, status, limit, offset)
	if err != nil {
		return nil, errors.New("gagal memuat data DPT")
	}

	return &entity.VoterListResponse{
		Voters:    voters,
		TotalData: total,
	}, nil
}

func (u *pemiraUsecase) ResetVoter(ctx context.Context, id int64) error {
	if err := u.pemiraRepo.ResetPemilihStatus(ctx, id); err != nil {
		return errors.New("gagal mereset status pemilih")
	}
	return nil
}

// ==========================================
// TOKEN HELPER (BALLOT JWT)
// ==========================================

// generateBallotJWT menerbitkan tiket masuk bilik suara khusus pemilih berdurasi 15 menit
func (u *pemiraUsecase) generateBallotJWT(nim string, idPemira int64) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_untuk_local_development_saja"
	}

	claims := jwt.MapClaims{
		"nim":       nim,
		"id_pemira": idPemira,
		"role":      "voter",
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ==========================================
// 4. ADMIN CANDIDATE (PASLON) CRUD USECASES
// ==========================================

func (u *pemiraUsecase) CreateCandidate(ctx context.Context, req entity.CreatePaslonRequest) (*entity.Paslon, error) {
	pemira, err := u.pemiraRepo.GetLatestPemira(ctx)
	if err != nil {
		return nil, errors.New("belum ada periode pemilihan raya yang terdaftar")
	}

	paslon := entity.Paslon{
		IdPemira:    pemira.IdPemira,
		NomorUrut:   req.NomorUrut,
		NimKetua:    req.NimKetua,
		NamaKetua:   req.NamaKetua,
		NimWakil:    req.NimWakil,
		NamaWakil:   req.NamaWakil,
		NamaKabinet: req.NamaKabinet,
		FotoPaslon:  req.FotoPaslon,
		Visi:        req.Visi,
		Misi:        req.Misi,
	}

	if err := u.pemiraRepo.CreatePaslon(ctx, &paslon); err != nil {
		return nil, fmt.Errorf("gagal menambahkan pasangan calon: %w", err)
	}

	return &paslon, nil
}

func (u *pemiraUsecase) UpdateCandidate(ctx context.Context, id int64, req entity.UpdatePaslonRequest) (*entity.Paslon, error) {
	paslon, err := u.pemiraRepo.GetPaslonByID(ctx, id)
	if err != nil {
		return nil, errors.New("data pasangan calon tidak ditemukan")
	}

	paslon.NomorUrut = req.NomorUrut
	paslon.NimKetua = req.NimKetua
	paslon.NamaKetua = req.NamaKetua
	paslon.NimWakil = req.NimWakil
	paslon.NamaWakil = req.NamaWakil
	paslon.NamaKabinet = req.NamaKabinet
	if req.FotoPaslon != "" {
		paslon.FotoPaslon = req.FotoPaslon
	}
	paslon.Visi = req.Visi
	paslon.Misi = req.Misi

	if err := u.pemiraRepo.UpdatePaslon(ctx, paslon); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pasangan calon: %w", err)
	}

	return paslon, nil
}

func (u *pemiraUsecase) DeleteCandidate(ctx context.Context, id int64) error {
	// Proteksi Integritas Suara (Asas LUBER JURDIL):
	// Kandidat dilarang dihapus jika sudah ada surat suara sah yang masuk ke bilik suara
	count, err := u.pemiraRepo.CountSuaraByPaslon(ctx, id)
	if err != nil {
		return errors.New("gagal memeriksa integritas suara paslon")
	}
	if count > 0 {
		return errors.New("kandidat tidak dapat dihapus karena sudah memiliki suara tercatat dalam bilik suara (Prinsip LUBER JURDIL)")
	}

	return u.pemiraRepo.DeletePaslon(ctx, id)
}

// ==========================================
// 5. ADMIN SETTINGS & STATUS TRIGGER
// ==========================================

func (u *pemiraUsecase) GetElectionSettings(ctx context.Context) (*entity.ElectionSettingsResponse, error) {
	pemira, err := u.pemiraRepo.GetLatestPemira(ctx)
	if err != nil {
		return nil, errors.New("data pengaturan pemilihan raya belum tersedia")
	}

	paslons, _ := u.pemiraRepo.GetPaslonByPemiraID(ctx, pemira.IdPemira)
	totalSuara, _ := u.pemiraRepo.CountTotalSuara(ctx, pemira.IdPemira)

	var winningPaslon *entity.Paslon
	mpikStatus := "PENDING"
	if pemira.Status == "SELESAI" {
		winner, _, errWinner := u.pemiraRepo.GetWinningPaslon(ctx, pemira.IdPemira)
		if errWinner == nil && winner != nil {
			winningPaslon = winner
			mpikStatus = "READY_FOR_INAUGURATION"
		}
	}

	return &entity.ElectionSettingsResponse{
		IdPemira:       pemira.IdPemira,
		TahunPeriode:   pemira.TahunPeriode,
		Status:         pemira.Status,
		TanggalMulai:   pemira.TanggalMulai,
		TanggalSelesai: pemira.TanggalSelesai,
		TotalPaslon:    int64(len(paslons)),
		TotalSuara:     totalSuara,
		MpikStatus:     mpikStatus,
		WinningPaslon:  winningPaslon,
	}, nil
}

func (u *pemiraUsecase) UpdateElectionSchedule(ctx context.Context, req entity.UpdatePemiraScheduleRequest) error {
	pemira, err := u.pemiraRepo.GetLatestPemira(ctx)
	if err != nil {
		return errors.New("data pemilihan raya tidak ditemukan")
	}

	if req.TanggalSelesai.Before(req.TanggalMulai) {
		return errors.New("tanggal selesai tidak boleh mendahului tanggal mulai")
	}

	return u.pemiraRepo.UpdatePemiraSchedule(ctx, pemira.IdPemira, req.TahunPeriode, req.TanggalMulai, req.TanggalSelesai)
}

func (u *pemiraUsecase) UpdateElectionStatus(ctx context.Context, req entity.UpdatePemiraStatusRequest) error {
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status != "DRAFT" && status != "BERJALAN" && status != "SELESAI" {
		return errors.New("status pemira tidak valid. Pilihan status yang tersedia: DRAFT, BERJALAN, SELESAI")
	}

	pemira, err := u.pemiraRepo.GetLatestPemira(ctx)
	if err != nil {
		return errors.New("data pemilihan raya tidak ditemukan")
	}

	return u.pemiraRepo.UpdatePemiraStatus(ctx, pemira.IdPemira, status)
}

// ==========================================
// 6. PELANTIKAN MPIK & TRANSISI MANAJEMEN HIMPUNAN
// ==========================================

func (u *pemiraUsecase) PelantikanMPIK(ctx context.Context) (*entity.PelantikanMPIKResponse, error) {
	// 1. Validasi Status Pemira Harus SELESAI
	pemira, err := u.pemiraRepo.GetLatestPemira(ctx)
	if err != nil {
		return nil, errors.New("data pemilihan raya tidak ditemukan")
	}
	if pemira.Status != "SELESAI" {
		return nil, errors.New("pelantikan MPIK hanya dapat dilakukan jika status pemilihan raya telah SELESAI")
	}

	// 2. Dapatkan Paslon Pemenang berdasarkan akumulasi suara terbanyak
	winner, _, err := u.pemiraRepo.GetWinningPaslon(ctx, pemira.IdPemira)
	if err != nil || winner == nil {
		return nil, fmt.Errorf("tidak dapat memproses pelantikan: %v", err)
	}

	// 3. Ekstrak Tahun Mulai & Tahun Selesai dari Tahun Periode
	tahunMulai, tahunSelesai := parsePeriodeYears(pemira.TahunPeriode)

	// 4. Promosikan ke Tabel Kabinet (Nonaktifkan Kabinet Lama, Aktifkan Kabinet Baru)
	kabinet := entity.Kabinet{
		NamaKabinet:  winner.NamaKabinet,
		TahunMulai:   tahunMulai,
		TahunSelesai: tahunSelesai,
		LogoKabinet:  winner.FotoPaslon,
		Visi:         winner.Visi,
		Misi:         winner.Misi,
		IsAktif:      true,
	}
	if err := u.pemiraRepo.PromoteKabinet(ctx, &kabinet); err != nil {
		return nil, fmt.Errorf("gagal meresmikan kabinet baru: %w", err)
	}

	// 5. Ekstrak Angkatan Otomatis dari 2 Digit NIM (contoh: 105223049 -> 2023)
	angkatanKetua := parseAngkatanFromNIM(winner.NimKetua)
	angkatanWakil := parseAngkatanFromNIM(winner.NimWakil)

	// 6. Upsert Data Ketua Himpunan ke Tabel Anggota (Status = Pengurus, BPH id_departemen = null)
	ketua, err := u.pemiraRepo.UpsertAnggotaBPH(ctx, winner.NimKetua, winner.NamaKetua, "Ketua Himpunan", angkatanKetua, winner.FotoPaslon)
	if err != nil {
		return nil, fmt.Errorf("gagal mencatat data Ketua Himpunan baru: %w", err)
	}

	// 7. Upsert Data Wakil Ketua Himpunan ke Tabel Anggota
	wakil, err := u.pemiraRepo.UpsertAnggotaBPH(ctx, winner.NimWakil, winner.NamaWakil, "Wakil Ketua Himpunan", angkatanWakil, winner.FotoPaslon)
	if err != nil {
		return nil, fmt.Errorf("gagal mencatat data Wakil Ketua Himpunan baru: %w", err)
	}

	// 8. Demisionerkan BPH Lama (Ubah status Pengurus lama yang bukan Ketua/Wakil baru menjadi Demisioner)
	_ = u.pemiraRepo.DemisionerkanPengurusLama(ctx, []string{winner.NimKetua, winner.NimWakil})

	return &entity.PelantikanMPIKResponse{
		StatusPelantikan: "APPROVED",
		Pesan:            fmt.Sprintf("Pelantikan berhasil! Kabinet %s resmi dilantik oleh MPIK dan BPH baru telah terintegrasi ke Manajemen Himpunan.", winner.NamaKabinet),
		KabinetBaru:      kabinet,
		KetuaBaru:        ketua,
		WakilBaru:        wakil,
	}, nil
}

// Helpers
func parsePeriodeYears(periode string) (int, int) {
	nowYear := time.Now().Year()
	clean := strings.ReplaceAll(periode, "/", "-")
	parts := strings.Split(clean, "-")
	if len(parts) >= 2 {
		y1, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		y2, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 == nil && err2 == nil && y1 > 2000 && y2 > 2000 {
			return y1, y2
		}
	} else if len(parts) == 1 {
		y1, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err == nil && y1 > 2000 {
			return y1, y1 + 1
		}
	}
	return nowYear, nowYear + 1
}

func parseAngkatanFromNIM(nim string) int {
	clean := strings.TrimSpace(nim)
	if len(clean) >= 6 {
		yyStr := clean[4:6]
		if yy, err := strconv.Atoi(yyStr); err == nil && yy >= 0 && yy <= 99 {
			return 2000 + yy
		}
	}
	return time.Now().Year()
}

