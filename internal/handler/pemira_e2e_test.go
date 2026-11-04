package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
