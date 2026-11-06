package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// SendOTPEmail mengirimkan kode OTP ke email tujuan dengan deteksi mode otomatis (Dual-Mode Architecture).
// Sesuai standar TS-Infra-Layanan-Mailer-dan-OTP:
// 1. Mode Development (atau jika SMTP_HOST belum diset): Mencetak banner OTP ke terminal tanpa crash/network error.
// 2. Mode Production: Mengirim email HTML resmi via SMTP TLS ke kotak masuk mahasiswa.
func SendOTPEmail(emailTujuan, nim, otpCode string) error {
	appEnv := strings.ToLower(os.Getenv("APP_ENV"))
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = os.Getenv("SMTP_EMAIL")
	}

	// 1. PRODUCTION / SMTP MODE (Aktif jika SMTP_HOST & SMTP_USER terisi dan APP_ENV bukan "development")
	if smtpHost != "" && smtpUser != "" && appEnv != "development" {
		err := sendSMTPEmail(emailTujuan, nim, otpCode)
		if err != nil {
			log.Printf("⚠️ [SMTP ERROR] Gagal mengirim email ke %s: %v", emailTujuan, err)
			return err
		}
		log.Printf("✅ [SMTP SUCCESS] Email OTP berhasil dikirimkan ke %s via %s", emailTujuan, smtpHost)
		return nil
	}

	// 2. FALLBACK / DEVELOPMENT MODE (Cetak ke konsol terminal)
	log.Println("================================================================")
	log.Printf("📧 [DEV MAILER - CONSOLE FALLBACK]")
	log.Printf("   Tujuan   : %s (NIM: %s)", emailTujuan, nim)
	log.Printf("   KODE OTP : >>> %s <<< (Masa Berlaku: 5 Menit)", otpCode)
	log.Println("================================================================")
	return nil
}

// sendSMTPEmail melakukan transmisi paket email HTML via protokol SMTP
func sendSMTPEmail(toEmail, nim, otpCode string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}
	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = os.Getenv("SMTP_EMAIL")
	}
	smtpPass := strings.ReplaceAll(os.Getenv("SMTP_PASSWORD"), " ", "")
	senderName := os.Getenv("SMTP_SENDER_NAME")
	if senderName == "" {
		senderName = "KPU HMIK Universitas Pertamina"
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	subject := fmt.Sprintf("[PEMIRA HMIK] Kode Verifikasi Masuk Bilik Suara: %s", otpCode)
	boundary := fmt.Sprintf("hmik_pemira_%d", time.Now().UnixNano())
	messageID := fmt.Sprintf("<%d.%s@gmail.com>", time.Now().UnixNano(), strings.Split(smtpUser, "@")[0])
	dateHeader := time.Now().Format(time.RFC1123Z)

	textBody := fmt.Sprintf("PEMILIHAN RAYA HMIK UNIVERSITAS PERTAMINA\r\n\r\n"+
		"Halo Mahasiswa (NIM: %s),\r\n\r\n"+
		"Berikut kode verifikasi untuk menggunakan hak suara Anda di Bilik Suara Virtual:\r\n"+
		">>> %s <<<\r\n\r\n"+
		"Kode verifikasi ini hanya berlaku selama 5 menit.\r\n"+
		"Jangan pernah membagikan kode ini kepada siapa pun.\r\n\r\n"+
		"© Komisi Pemilihan Umum HMIK - Universitas Pertamina", nim, otpCode)

	htmlBody := buildOTPEmailHTML(nim, otpCode)

	message := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Date: %s\r\n"+
		"Message-ID: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n"+
		"--%s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: 8bit\r\n\r\n"+
		"%s\r\n\r\n"+
		"--%s\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"Content-Transfer-Encoding: 8bit\r\n\r\n"+
		"%s\r\n\r\n"+
		"--%s--\r\n",
		senderName, smtpUser, toEmail, subject, dateHeader, messageID, boundary, boundary, textBody, boundary, htmlBody, boundary)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	return smtp.SendMail(addr, auth, smtpUser, []string{toEmail}, []byte(message))
}

// buildOTPEmailHTML menghasilkan markup template email responsif bertema HMIK Neo-Brutalisme
func buildOTPEmailHTML(nim, otpCode string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kode Verifikasi Pemira HMIK</title>
    <style>
        body { font-family: 'Helvetica Neue', Arial, sans-serif; background-color: #F8FAFC; margin: 0; padding: 20px; }
        .container { max-width: 500px; margin: 0 auto; background-color: #FFFFFF; border: 3px solid #2B2B2B; border-radius: 16px; padding: 32px 24px; box-shadow: 5px 5px 0px #2B2B2B; }
        .header { text-align: center; margin-bottom: 24px; }
        .badge { display: inline-block; background-color: #F7D070; color: #2B2B2B; border: 1.5px solid #2B2B2B; border-radius: 20px; padding: 4px 14px; font-size: 11px; font-weight: 800; letter-spacing: 1px; text-transform: uppercase; margin-bottom: 12px; }
        .title { font-size: 20px; font-weight: 800; color: #2B2B2B; margin: 0 0 8px 0; }
        .greeting { color: #555555; font-size: 14px; margin: 0; line-height: 1.5; }
        .otp-box { background-color: #F0F2F1; border: 2.5px solid #2B2B2B; border-radius: 12px; padding: 18px; text-align: center; margin: 24px 0; box-shadow: 3px 3px 0px #2B2B2B; }
        .otp-number { font-size: 36px; font-weight: 800; letter-spacing: 8px; color: #4059AD; margin: 0; }
        .warning { font-size: 13px; color: #666666; text-align: center; line-height: 1.5; margin: 16px 0; }
        .warning strong { color: #E74C3C; }
        .footer { text-align: center; font-size: 11px; color: #999999; margin-top: 24px; border-top: 1px solid #E5E7EB; padding-top: 14px; line-height: 1.4; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="badge">PEMILIHAN RAYA HMIK</div>
            <h2 class="title">Verifikasi Bilik Suara Virtual</h2>
            <p class="greeting">Halo Mahasiswa (NIM: %s), berikut kode verifikasi untuk menggunakan hak suara Anda:</p>
        </div>
        <div class="otp-box">
            <h1 class="otp-number">%s</h1>
        </div>
        <p class="warning">
            Kode verifikasi ini hanya berlaku selama <strong>5 menit</strong>.<br>
            Jangan pernah membagikan kode ini kepada siapa pun demi menjaga kedaulatan dan kerahasiaan hak suara Anda.
        </p>
        <div class="footer">
            &copy; Komisi Pemilihan Umum HMIK &bull; Universitas Pertamina<br>
            Pesan ini dikirimkan secara otomatis oleh sistem, mohon tidak membalas email ini.
        </div>
    </div>
</body>
</html>`, nim, otpCode)
}
