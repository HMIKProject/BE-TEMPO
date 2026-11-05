package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	_ = godotenv.Load("../../.env")
	_ = os.Setenv("APP_ENV", "development")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=rahasia_admin dbname=hmik_db port=5433 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		t.Skipf("Database tidak tersedia, lewati integrasi E2E: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, db)

	return router, db
}

func generateAdminTestToken() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "rahasia_untuk_local_development_saja"
	}

	claims := jwt.MapClaims{
		"id_pengguna": 1,
		"id_peran":    1,
		"role":        "admin",
		"exp":         time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestPemiraE2E_FullLifecycle(t *testing.T) {
	router, db := setupTestRouter(t)

	// Persiapkan test voter
	testNIM := "105223099"
	testEmail := "105223099@student.universitaspertamina.ac.id"

	// Pastikan pemira aktif untuk test bilik suara
	db.Exec("UPDATE pemira SET status = 'BERJALAN', tanggal_mulai = NOW() - INTERVAL '1 day', tanggal_selesai = NOW() + INTERVAL '2 day' WHERE id_pemira = 1")

	// Reset status voter di DB sebelum test
	db.Exec("UPDATE pemilih SET status_memilih = false, kode_otp = NULL, otp_expired = NULL WHERE nim = ?", testNIM)

	// --- 1. Request OTP dengan Email Ilegal ---
	t.Run("1. Request OTP Email Ilegal Ditolak", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"nim":          testNIM,
			"email_kampus": "hacker@gmail.com",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/auth/request-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ekspektasi 400 Bad Request, dapat: %d", w.Code)
		}
	})

	// --- 2. Request OTP Berhasil ---
	t.Run("2. Request OTP Berhasil", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"nim":          testNIM,
			"email_kampus": testEmail,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/auth/request-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK, dapat: %d, body: %s", w.Code, w.Body.String())
		}
	})

	// --- 2b. Auto-Register Mahasiswa Baru Belum Ada di DPT ---
	t.Run("2b. Auto-Register Mahasiswa Baru Belum Ada di DPT", func(t *testing.T) {
		newStudentNIM := "105223888"
		newStudentEmail := "105223888@student.universitaspertamina.ac.id"

		// Hapus jika sebelumnya ada
		db.Exec("DELETE FROM pemilih WHERE nim = ?", newStudentNIM)

		body, _ := json.Marshal(map[string]string{
			"nim":          newStudentNIM,
			"email_kampus": newStudentEmail,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/auth/request-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK untuk auto-register, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		// Pastikan mahasiswa langsung tercatat di tabel pemilih
		var check entity.Pemilih
		if err := db.Where("nim = ?", newStudentNIM).First(&check).Error; err != nil {
			t.Fatalf("Mahasiswa baru gagal tercatat di tabel pemilih: %v", err)
		}
		if check.EmailKampus != newStudentEmail || check.StatusMemilih != false {
			t.Errorf("Data pemilih yang dibuat tidak sesuai: %+v", check)
		}

		// Cleanup
		db.Exec("DELETE FROM pemilih WHERE nim = ?", newStudentNIM)
	})

	// Ambil kode OTP yang tersimpan di DB
	var pemilih entity.Pemilih
	if err := db.Where("nim = ?", testNIM).First(&pemilih).Error; err != nil {
		t.Fatalf("Gagal mengambil voter dari DB: %v", err)
	}
	if pemilih.KodeOtp == "" {
		t.Fatalf("Kode OTP kosong di database!")
	}
	correctOTP := pemilih.KodeOtp

	// --- 3. Verifikasi OTP dengan Kode Salah ---
	t.Run("3. Verifikasi OTP Kode Salah Ditolak", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"nim":      testNIM,
			"kode_otp": "000000",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/auth/verify-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ekspektasi 400 Bad Request untuk kode salah, dapat: %d", w.Code)
		}
	})

	// --- 4. Verifikasi OTP Benar -> Dapatkan Ballot Token ---
	var ballotToken string
	t.Run("4. Verifikasi OTP Benar Menerbitkan Ballot Token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"nim":      testNIM,
			"kode_otp": correctOTP,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/auth/verify-otp", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data entity.VerifyOTPResponse `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Gagal unmarshal response verify OTP: %v", err)
		}

		if resp.Data.BallotToken == "" {
			t.Fatalf("Ballot token kosong di response!")
		}
		ballotToken = resp.Data.BallotToken
	})

	// --- 5. Ambil Daftar Kandidat Paslon ---
	var paslonID int64
	t.Run("5. Ambil Daftar Paslon Aktif", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/pemira/kandidat", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data []entity.Paslon `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) == 0 {
			t.Fatalf("Daftar paslon kosong!")
		}
		paslonID = resp.Data[0].IdPaslon
	})

	// --- 6. Pemberian Suara (Cast Vote) ---
	t.Run("6. Pemberian Suara Sah Berhasil", func(t *testing.T) {
		body, _ := json.Marshal(map[string]int64{
			"id_paslon": paslonID,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/vote", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ballotToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Ekspektasi 201 Created, dapat: %d, body: %s", w.Code, w.Body.String())
		}
	})

	// --- 7. Anti Double-Voting: Pencoblosan Ulang Ditolak ---
	t.Run("7. Anti Double-Voting: Pencoblosan Ulang Ditolak", func(t *testing.T) {
		body, _ := json.Marshal(map[string]int64{
			"id_paslon": paslonID,
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pemira/vote", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ballotToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ekspektasi 400 Bad Request untuk double-voting, dapat: %d", w.Code)
		}
	})

	// --- 8. Admin Metrik & Monitoring ---
	adminToken := generateAdminTestToken()

	t.Run("8. Admin Mengambil Ringkasan Statistik Pemilu", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/pemira/stats", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK untuk admin stats, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data entity.ElectionStatsResponse `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Data.TotalSuaraMasuk < 1 {
			t.Errorf("Ekspektasi suara masuk >= 1, dapat: %d", resp.Data.TotalSuaraMasuk)
		}
	})

	// --- 9. Admin Mengambil Data Grafik Pleno ---
	t.Run("9. Admin Mengambil Data Grafik Pleno", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/pemira/chart", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK untuk admin chart, dapat: %d, body: %s", w.Code, w.Body.String())
		}
	})

	// --- 10. Admin Reset Status Pemilih ---
	t.Run("10. Admin Reset Hak Pilih Berhasil", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/pemira/pemilih/1/reset", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK untuk reset pemilih, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		var check entity.Pemilih
		db.First(&check, 1)
		if check.StatusMemilih {
			t.Errorf("Status memilih harusnya false setelah di-reset!")
		}
	})
}

func TestAdminPemiraCRUD_And_PelantikanMPIK(t *testing.T) {
	router, db := setupTestRouter(t)
	adminToken := generateAdminTestToken()

	// 1. Create Paslon Baru
	var newPaslonID int64
	t.Run("1. Admin Tambah Paslon Baru Berhasil", func(t *testing.T) {
		body, _ := json.Marshal(entity.CreatePaslonRequest{
			NomorUrut:   99,
			NimKetua:    "105223091",
			NamaKetua:   "Test Kahim",
			NimWakil:    "105223092",
			NamaWakil:   "Test Wakahim",
			NamaKabinet: "Kabinet Sinergi Uji",
			FotoPaslon:  "https://example.com/test.png",
			Visi:        "Visi Uji Coba",
			Misi:        "Misi Uji Coba",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/pemira/kandidat", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Ekspektasi 201 Created, dapat: %d, body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data entity.Paslon `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Gagal unmarshal paslon: %v", err)
		}
		newPaslonID = resp.Data.IdPaslon
		if newPaslonID == 0 {
			t.Fatalf("ID Paslon baru tidak boleh 0")
		}
	})

	// 2. Update Paslon
	t.Run("2. Admin Update Paslon Berhasil", func(t *testing.T) {
		body, _ := json.Marshal(entity.UpdatePaslonRequest{
			NomorUrut:   99,
			NimKetua:    "105223091",
			NamaKetua:   "Test Kahim Diperbarui",
			NimWakil:    "105223092",
			NamaWakil:   "Test Wakahim Diperbarui",
			NamaKabinet: "Kabinet Sinergi Uji Diperbarui",
			FotoPaslon:  "https://example.com/test_updated.png",
			Visi:        "Visi Terupdate",
			Misi:        "Misi Terupdate",
		})
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/pemira/kandidat/"+strconv.FormatInt(newPaslonID, 10), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK, dapat: %d, body: %s", w.Code, w.Body.String())
		}
	})

	// 3. Settings & Status Control
	t.Run("3. Admin Ambil & Update Pengaturan Acara", func(t *testing.T) {
		// GET Settings
		reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/pemira/settings", nil)
		reqGet.Header.Set("Authorization", "Bearer "+adminToken)
		wGet := httptest.NewRecorder()
		router.ServeHTTP(wGet, reqGet)

		if wGet.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK get settings, dapat: %d", wGet.Code)
		}

		// PUT Schedule
		scheduleBody, _ := json.Marshal(entity.UpdatePemiraScheduleRequest{
			TahunPeriode:   "2026/2027",
			TanggalMulai:   time.Now().Add(-24 * time.Hour),
			TanggalSelesai: time.Now().Add(48 * time.Hour),
		})
		reqPut, _ := http.NewRequest(http.MethodPut, "/api/v1/admin/pemira/settings", bytes.NewBuffer(scheduleBody))
		reqPut.Header.Set("Content-Type", "application/json")
		reqPut.Header.Set("Authorization", "Bearer "+adminToken)
		wPut := httptest.NewRecorder()
		router.ServeHTTP(wPut, reqPut)

		if wPut.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK update schedule, dapat: %d", wPut.Code)
		}

		// POST Status Trigger
		statusBody, _ := json.Marshal(entity.UpdatePemiraStatusRequest{Status: "BERJALAN"})
		reqStatus, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/pemira/status", bytes.NewBuffer(statusBody))
		reqStatus.Header.Set("Content-Type", "application/json")
		reqStatus.Header.Set("Authorization", "Bearer "+adminToken)
		wStatus := httptest.NewRecorder()
		router.ServeHTTP(wStatus, reqStatus)

		if wStatus.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK update status, dapat: %d", wStatus.Code)
		}
	})

	// 4. Strict Deletion Protection (LUBER JURDIL)
	t.Run("4. Proteksi Penghapusan Paslon Berjalan Benar", func(t *testing.T) {
		// Masukkan suara sah untuk paslon 5 agar proteksi penghapusan paslon teruji
		db.Exec("INSERT INTO suara (id_pemira, id_paslon, waktu_memilih) VALUES (1, 5, NOW())")

		reqDelVoted, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/pemira/kandidat/5", nil)
		reqDelVoted.Header.Set("Authorization", "Bearer "+adminToken)
		wDelVoted := httptest.NewRecorder()
		router.ServeHTTP(wDelVoted, reqDelVoted)

		if wDelVoted.Code != http.StatusBadRequest {
			t.Errorf("Ekspektasi 400 Bad Request untuk paslon dengan suara masuk, dapat: %d", wDelVoted.Code)
		}

		// Paslon baru yang belum punya suara -> boleh dihapus
		reqDelNew, _ := http.NewRequest(http.MethodDelete, "/api/v1/admin/pemira/kandidat/"+strconv.FormatInt(newPaslonID, 10), nil)
		reqDelNew.Header.Set("Authorization", "Bearer "+adminToken)
		wDelNew := httptest.NewRecorder()
		router.ServeHTTP(wDelNew, reqDelNew)

		if wDelNew.Code != http.StatusOK {
			t.Errorf("Ekspektasi 200 OK untuk paslon tanpa suara, dapat: %d", wDelNew.Code)
		}
	})

	// 5. Pelantikan MPIK & Transisi Manajemen Himpunan
	t.Run("5. Pengesahan & Pelantikan MPIK", func(t *testing.T) {
		// Ubah status ke SELESAI agar syarat pelantikan terpenuhi
		statusSelesai, _ := json.Marshal(entity.UpdatePemiraStatusRequest{Status: "SELESAI"})
		reqStatus, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/pemira/status", bytes.NewBuffer(statusSelesai))
		reqStatus.Header.Set("Content-Type", "application/json")
		reqStatus.Header.Set("Authorization", "Bearer "+adminToken)
		wStatus := httptest.NewRecorder()
		router.ServeHTTP(wStatus, reqStatus)

		if wStatus.Code != http.StatusOK {
			t.Fatalf("Gagal mengubah status pemira ke SELESAI: %d", wStatus.Code)
		}

		// Eksekusi Pelantikan MPIK
		reqPelantikan, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/pemira/pelantikan-mpik", nil)
		reqPelantikan.Header.Set("Authorization", "Bearer "+adminToken)
		wPelantikan := httptest.NewRecorder()
		router.ServeHTTP(wPelantikan, reqPelantikan)

		if wPelantikan.Code != http.StatusOK {
			t.Fatalf("Ekspektasi 200 OK pelantikan MPIK, dapat: %d, body: %s", wPelantikan.Code, wPelantikan.Body.String())
		}

		var resp struct {
			Data entity.PelantikanMPIKResponse `json:"data"`
		}
		if err := json.Unmarshal(wPelantikan.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Gagal parse response pelantikan: %v", err)
		}

		if resp.Data.StatusPelantikan != "APPROVED" {
			t.Errorf("Ekspektasi status APPROVED, dapat: %s", resp.Data.StatusPelantikan)
		}

		// Verifikasi Kabinet Aktif di DB
		var activeKabinet entity.Kabinet
		if err := db.Where("is_aktif = ?", true).First(&activeKabinet).Error; err != nil {
			t.Fatalf("Gagal menemukan kabinet aktif baru di DB: %v", err)
		}

		// Verifikasi BPH Terdaftar di Anggota dengan Angkatan Ekstraksi
		if resp.Data.KetuaBaru != nil {
			if resp.Data.KetuaBaru.Angkatan < 2000 {
				t.Errorf("Ekstraksi angkatan Ketua tidak valid: %d", resp.Data.KetuaBaru.Angkatan)
			}
		}

		// Reset kembali status Pemira ke BERJALAN dan kosongkan suara/pemilih uji coba
		statusBerjalan, _ := json.Marshal(entity.UpdatePemiraStatusRequest{Status: "BERJALAN"})
		reqResetStatus, _ := http.NewRequest(http.MethodPost, "/api/v1/admin/pemira/status", bytes.NewBuffer(statusBerjalan))
		reqResetStatus.Header.Set("Content-Type", "application/json")
		reqResetStatus.Header.Set("Authorization", "Bearer "+adminToken)
		wReset := httptest.NewRecorder()
		router.ServeHTTP(wReset, reqResetStatus)

		// Bersihkan suara dan pemilih uji coba
		db.Exec("TRUNCATE TABLE suara, pemilih RESTART IDENTITY")
	})
}

