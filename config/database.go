package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=rahasia dbname=hmik_db port=5432"

	// Perhatikan bagian NamingStrategy di bawah ini!
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // WAJIB: Memaksa GORM tidak menambahkan huruf 's'
		},
	})

	if err != nil {
		panic("Gagal terkoneksi ke database")
	}
	return db
}
