package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware adalah satpam penjaga gerbang yang memastikan user memiliki token valid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Token JWT tidak ditemukan", nil))
			c.Abort()
			return
		}

		// 2. Ekstrak string token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "rahasia_untuk_local_development_saja"
		}

		// 3. Parsing dan verifikasi token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Pastikan algoritma yang digunakan adalah HMAC (HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		// 4. Jika error atau tidak valid (termasuk jika sudah expired > 24 jam)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Token tidak valid atau sudah kadaluarsa", nil))
			c.Abort()
			return
		}

		// 5. Ekstrak Payload dan tempelkan ke Gin Context agar bisa dibaca oleh Handler
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			c.Set("id_pengguna", claims["id_pengguna"])
			c.Set("id_peran", claims["id_peran"])
		}

		// 6. Lanjut ke proses Handler
		c.Next()
	}
}
