package entity

// LoginRequest merepresentasikan data yang dikirim oleh Frontend saat login
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"hmik.humas.01@hmik.com"` // Bisa berupa Email atau Username
	Password   string `json:"password" binding:"required" example:"qwerqwer"`
}

// LoginResponse merepresentasikan kembalian login atau pemicu ganti sandi awal ke Frontend
type LoginResponse struct {
	Token                 string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	IDPengguna            int64  `json:"id_pengguna" example:"1"`
	NamaPengguna          string `json:"nama_pengguna" example:"hmik.humas.01"`
	NamaLengkap           string `json:"nama_lengkap,omitempty" example:"Super Admin"`
	RequirePasswordChange bool   `json:"require_password_change" example:"false"`
	SetupToken            string `json:"setup_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Email                 string `json:"email,omitempty" example:"hmik.humas.01@hmik.com"`
	Jabatan               string `json:"jabatan,omitempty" example:"Ketua Departemen Hubungan External"`
	Role                  string `json:"role,omitempty" example:"Administrator"`
	Initials              string `json:"initials,omitempty" example:"SA"`
}

// UserProfileResponse merepresentasikan profil akun pengguna yang sedang aktif/login
type UserProfileResponse struct {
	IDPengguna   int64  `json:"id_pengguna"`
	NamaPengguna string `json:"nama_pengguna"`
	NamaLengkap  string `json:"nama_lengkap"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Jabatan      string `json:"jabatan"`
	Initials     string `json:"initials"`
	Foto         string `json:"foto,omitempty"`
}

// SetInitialPasswordRequest payload aktivasi mandiri pengurus baru dari popup /login
type SetInitialPasswordRequest struct {
	SetupToken      string `json:"setup_token" binding:"required" example:"eyJhbGci..."`
	PasswordBaru    string `json:"password_baru" binding:"required,min=8" example:"PasswordPribadi123!"`
	KonfirmasiSandi string `json:"konfirmasi_sandi" binding:"required,eqfield=PasswordBaru" example:"PasswordPribadi123!"`
}

// StudentRequestOTPRequest payload permohonan kode OTP mahasiswa umum (Login Tipe 2)
type StudentRequestOTPRequest struct {
	NIM         string `json:"nim" binding:"required" example:"105223046"`
	EmailKampus string `json:"email_kampus" binding:"required,email" example:"105223046@student.universitaspertamina.ac.id"`
}

// StudentRequestOTPResponse kembalian request OTP mahasiswa
type StudentRequestOTPResponse struct {
	EmailKampus  string `json:"email_kampus"`
	NIM          string `json:"nim"`
	SisaCooldown int    `json:"sisa_cooldown"`
	KodeOTP      string `json:"kode_otp,omitempty"` // Ditampilkan khusus di mode dev/testing
}

// StudentVerifyOTPRequest payload verifikasi OTP mahasiswa umum (Login Tipe 2)
type StudentVerifyOTPRequest struct {
	NIM         string `json:"nim" binding:"required" example:"105223046"`
	EmailKampus string `json:"email_kampus" binding:"required,email" example:"105223046@student.universitaspertamina.ac.id"`
	KodeOTP     string `json:"kode_otp" binding:"required,len=6" example:"123456"`
}

// StudentVerifyOTPResponse kembalian sukses login mahasiswa umum
type StudentVerifyOTPResponse struct {
	Token       string `json:"token" example:"eyJhbGci..."`
	NIM         string `json:"nim" example:"105223046"`
	EmailKampus string `json:"email_kampus" example:"105223046@student.universitaspertamina.ac.id"`
	Role        string `json:"role" example:"Mahasiswa"`
	NamaLengkap string `json:"nama_lengkap" example:"Mahasiswa Ilmu Komputer"`
}