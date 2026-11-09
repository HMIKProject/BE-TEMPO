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
		senderName = "Panitia Pemira HMIK"
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	subject := "Kode Akses Masuk Bilik Suara Pemira HMIK"
	boundary := fmt.Sprintf("hmik_pemira_%d", time.Now().UnixNano())
	userDomain := "gmail.com"
	if parts := strings.Split(smtpUser, "@"); len(parts) > 1 {
		userDomain = parts[1]
	}
	messageID := fmt.Sprintf("<%d.%d@%s>", time.Now().Unix(), time.Now().UnixNano()%1000000, userDomain)
	dateHeader := time.Now().Format(time.RFC1123Z)

	textBody := fmt.Sprintf("PEMILIHAN RAYA HMIK\r\n\r\n"+
		"Halo Mahasiswa (NIM: %s),\r\n\r\n"+
		"Berikut kode akses verifikasi untuk masuk ke Bilik Suara Virtual Pemira:\r\n"+
		">>> %s <<<\r\n\r\n"+
		"Kode verifikasi ini hanya berlaku selama 5 menit.\r\n"+
		"Jangan bagikan kode ini kepada siapa pun untuk menjaga kerahasiaan suara Anda.\r\n\r\n"+
		"Pesan ini dikirimkan otomatis kepada pemilih yang terdaftar di DPT Pemira HMIK.\r\n"+
		"© Panitia Pemilihan Raya HMIK", nim, otpCode)

	htmlBody := buildOTPEmailHTML(nim, otpCode)

	message := fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Reply-To: %s\r\n"+
		"Subject: %s\r\n"+
		"Date: %s\r\n"+
		"Message-ID: %s\r\n"+
		"X-Priority: 3\r\n"+
		"Importance: normal\r\n"+
		"Auto-Submitted: auto-generated\r\n"+
		"X-Auto-Response-Suppress: All\r\n"+
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
		senderName, smtpUser, toEmail, smtpUser, subject, dateHeader, messageID, boundary, boundary, textBody, boundary, htmlBody, boundary)

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
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #F8FAFC; margin: 0; padding: 20px; color: #1E293B; }
        .container { max-width: 480px; margin: 0 auto; background-color: #FFFFFF; border: 2px solid #CBD5E1; border-radius: 12px; padding: 32px 24px; }
        .header { text-align: center; margin-bottom: 24px; }
        .badge { display: inline-block; background-color: #EEF2FF; color: #4059AD; border: 1px solid #C7D2FE; border-radius: 16px; padding: 4px 12px; font-size: 11px; font-weight: 700; letter-spacing: 0.5px; text-transform: uppercase; margin-bottom: 12px; }
        .title { font-size: 18px; font-weight: 700; color: #0F172A; margin: 0 0 8px 0; }
        .greeting { color: #475569; font-size: 14px; margin: 0; line-height: 1.5; }
        .otp-box { background-color: #F1F5F9; border-radius: 8px; padding: 16px; text-align: center; margin: 24px 0; }
        .otp-number { font-size: 32px; font-weight: 800; letter-spacing: 6px; color: #4059AD; margin: 0; font-family: monospace, sans-serif; }
        .warning { font-size: 13px; color: #64748B; text-align: center; line-height: 1.5; margin: 16px 0; }
        .footer { text-align: center; font-size: 11px; color: #94A3B8; margin-top: 24px; border-top: 1px solid #F1F5F9; padding-top: 14px; line-height: 1.5; }
    </style>
</head>
<body>
    <!-- Hidden preheader text for Outlook inbox preview -->
    <div style="display:none;font-size:1px;color:#333333;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;">
        Kode akses bilik suara virtual Pemira HMIK Anda adalah %s. Berlaku selama 5 menit.
    </div>

    <div class="container">
        <div class="header">
            <div class="badge">Pemilihan Raya HMIK</div>
            <h2 class="title">Verifikasi Bilik Suara Virtual</h2>
            <p class="greeting">Halo Mahasiswa (NIM: %s), berikut kode verifikasi untuk masuk ke bilik suara:</p>
        </div>
        <div class="otp-box">
            <h1 class="otp-number">%s</h1>
        </div>
        <p class="warning">
            Kode ini hanya berlaku selama <strong>5 menit</strong>.<br>
            Demi menjaga kerahasiaan suara Anda, mohon tidak membagikan kode ini kepada siapa pun.
        </p>
        <div class="footer">
            Email ini dikirimkan otomatis karena NIM Anda terdaftar dalam DPT Pemira HMIK.<br>
            &copy; Panitia Pemilihan Raya HMIK
        </div>
    </div>
</body>
</html>`, otpCode, nim, otpCode)
}
