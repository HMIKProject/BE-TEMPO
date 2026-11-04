package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// VoterMiddleware melindungi endpoint bilik suara pemilih menggunakan Ballot Token (15 menit)
func VoterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Token bilik suara tidak ditemukan", nil))
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
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		// 4. Jika error atau kadaluarsa (> 15 menit)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Sesi bilik suara telah berakhir atau token tidak valid", nil))
			c.Abort()
			return
		}

		// 5. Ekstrak Claims dan pastikan role adalah "voter"
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Klaim token tidak valid", nil))
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "voter" {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: Token ini bukan token pemilih bilik suara", nil))
			c.Abort()
			return
		}

		nim, ok := claims["nim"].(string)
		if !ok || nim == "" {
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: NIM pemilih tidak ditemukan pada token", nil))
			c.Abort()
			return
		}

		var idPemira int64
		switch v := claims["id_pemira"].(type) {
		case float64:
			idPemira = int64(v)
		case int64:
			idPemira = v
		case int:
			idPemira = int64(v)
		default:
			c.JSON(http.StatusUnauthorized, utils.BuildFailResponse("Akses Ditolak: ID Pemira pada token tidak valid", nil))
			c.Abort()
			return
		}

		// 6. Simpan identitas pemilih ke dalam Gin context
		c.Set("nim", nim)
		c.Set("id_pemira", idPemira)

		c.Next()
	}
}
