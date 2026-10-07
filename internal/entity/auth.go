package entity

// LoginRequest merepresentasikan data yang dikirim oleh Frontend saat login
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"hmik.human.01@hmik.com"` // Bisa berupa Email atau Username
	Password   string `json:"password" binding:"required" example:"Human102230462026"`
}

// LoginResponse merepresentasikan kembalian token JWT ke Frontend
type LoginResponse struct {
	Token       string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	IDPengguna  int64  `json:"id_pengguna" example:"1"`
	NamaLengkap string `json:"nama_lengkap" example:"Budi Human Resource"`
}
