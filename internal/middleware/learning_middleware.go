package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// LearningAccessMiddleware melindungi akses materi modul perkuliahan:
// Mengizinkan user dengan JWT Pengurus HMIK (Login Tipe 1) ATAU JWT Mahasiswa (Login Tipe 2).
func LearningAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Modul pembelajaran dosen dilindungi dan khusus untuk Mahasiswa Ilmu Komputer Universitas Pertamina. Silakan login terlebih dahulu.", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "rahasia_untuk_local_development_saja"
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Sesi login Anda telah berakhir atau tidak valid. Silakan login kembali.", nil))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Klaim sesi tidak valid", nil))
			c.Abort()
			return
		}

		if nim, ok := claims["nim"].(string); ok && nim != "" {
			c.Set("nim", nim)
		}
		if idPengguna, ok := claims["id_pengguna"]; ok {
			c.Set("id_pengguna", idPengguna)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}

		c.Next()
	}
}
