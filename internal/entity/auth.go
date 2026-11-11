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
	RequirePasswordChange bool   `json:"require_password_change" example:"false"`
	SetupToken            string `json:"setup_token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Email                 string `json:"email,omitempty" example:"hmik.humas.01@hmik.com"`
	Jabatan               string `json:"jabatan,omitempty" example:"Ketua Departemen Hubungan External"`
}

// SetInitialPasswordRequest payload aktivasi mandiri pengurus baru dari popup /login
type SetInitialPasswordRequest struct {
	SetupToken      string `json:"setup_token" binding:"required" example:"eyJhbGci..."`
	PasswordBaru    string `json:"password_baru" binding:"required,min=8" example:"PasswordPribadi123!"`
	KonfirmasiSandi string `json:"konfirmasi_sandi" binding:"required,eqfield=PasswordBaru" example:"PasswordPribadi123!"`
}