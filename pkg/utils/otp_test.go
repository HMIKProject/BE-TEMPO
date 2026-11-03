package utils

import (
	"regexp"
	"testing"
)

func TestGenerateCryptoOTP(t *testing.T) {
	// Test panjang 6 digit
	otp, err := GenerateCryptoOTP(6)
	if err != nil {
		t.Fatalf("GenerateCryptoOTP gagal: %v", err)
	}

	if len(otp) != 6 {
		t.Errorf("Panjang OTP harus 6, didapat: %d (otp: %s)", len(otp), otp)
	}

	// Validasi bahwa seluruh karakter adalah digit angka (0-9)
	match, err := regexp.MatchString(`^\d{6}$`, otp)
	if err != nil || !match {
		t.Errorf("OTP harus berupa 6 digit angka numerik murni, didapat: %s", otp)
	}

	// Test keacakan (dua OTP berurutan tidak boleh sama)
	otp2, err := GenerateCryptoOTP(6)
	if err != nil {
		t.Fatalf("GenerateCryptoOTP kedua gagal: %v", err)
	}

	if otp == otp2 {
		t.Errorf("Dua OTP berurutan tidak boleh identik: %s vs %s", otp, otp2)
	}
}

func TestSendOTPEmail_DevelopmentMode(t *testing.T) {
	// Test bahwa pada development mode (tanpa SMTP), fungsi mengembalikan nil tanpa error
	t.Setenv("APP_ENV", "development")
	t.Setenv("SMTP_HOST", "")

	err := SendOTPEmail("10123001@student.universitaspertamina.ac.id", "10123001", "123456")
	if err != nil {
		t.Errorf("SendOTPEmail development mode gagal: %v", err)
	}
}
