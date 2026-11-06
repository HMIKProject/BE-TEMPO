package seeder

import (
	"log"
	"strings"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAllRoleAccounts menginisialisasi atau memperbarui akun slot jabatan resmi dengan sandi qwerqwer
func SeedAllRoleAccounts(db *gorm.DB) error {
	plainPassword := "qwerqwer"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashStr := string(hashedPassword)

	accounts := []struct {
		Username string
		Email    string
		Jabatan  string
	}{
		// 1. Akun Admin Utama
		{"hmik.admin", "hmik.up@gmail.com", "Super Admin"},
		{"hmik.human.01", "hmik.human.01@hmik.com", "Slot Humas Human"},

		// 2. BPH (Badan Pengurus Harian)
		{"hmik.kahim.01", "hmik.kahim.01@hmik.com", "Ketua Himpunan"},
		{"hmik.wakahim.01", "hmik.wakahim.01@hmik.com", "Wakil Ketua Himpunan"},
		{"hmik.sekretaris.01", "hmik.sekretaris.01@hmik.com", "Sekretaris 1"},
		{"hmik.sekretaris.02", "hmik.sekretaris.02@hmik.com", "Sekretaris 2"},
		{"hmik.bendahara.01", "hmik.bendahara.01@hmik.com", "Bendahara 1"},
		{"hmik.bendahara.02", "hmik.bendahara.02@hmik.com", "Bendahara 2"},

		// 3. Riset dan Teknologi (Ristek)
		{"hmik.ristek.01", "hmik.ristek.01@hmik.com", "Ketua Departemen Riset dan Teknologi"},
		{"hmik.ristek.02", "hmik.ristek.02@hmik.com", "Sekretaris Departemen Riset dan Teknologi"},
		{"hmik.ristek.inovasi.01", "hmik.ristek.inovasi.01@hmik.com", "Ketua Divisi Riset & Inovasi"},
		{"hmik.ristek.inovasi.02", "hmik.ristek.inovasi.02@hmik.com", "Anggota Divisi Riset & Inovasi"},
		{"hmik.ristek.produk.01", "hmik.ristek.produk.01@hmik.com", "Ketua Divisi Pengembangan Produk"},
		{"hmik.ristek.produk.02", "hmik.ristek.produk.02@hmik.com", "Anggota Divisi Pengembangan Produk"},

		// 4. Media Informasi (Medinfo)
		{"hmik.medinfo.01", "hmik.medinfo.01@hmik.com", "Ketua Departemen Media Informasi"},
		{"hmik.medinfo.02", "hmik.medinfo.02@hmik.com", "Sekretaris Departemen Media Informasi"},
		{"hmik.medinfo.kreatif.01", "hmik.medinfo.kreatif.01@hmik.com", "Ketua Divisi Media Kreatif"},
		{"hmik.medinfo.kreatif.02", "hmik.medinfo.kreatif.02@hmik.com", "Anggota Divisi Media Kreatif"},
		{"hmik.medinfo.pubdok.01", "hmik.medinfo.pubdok.01@hmik.com", "Ketua Divisi Publikasi & Dokumentasi"},
		{"hmik.medinfo.pubdok.02", "hmik.medinfo.pubdok.02@hmik.com", "Anggota Divisi Publikasi & Dokumentasi"},

		// 5. Hubungan External (Humas)
		{"hmik.humas.01", "hmik.humas.01@hmik.com", "Ketua Departemen Hubungan External"},
		{"hmik.humas.02", "hmik.humas.02@hmik.com", "Sekretaris Departemen Hubungan External"},
		{"hmik.humas.pengmas.01", "hmik.humas.pengmas.01@hmik.com", "Ketua Divisi Pengabdian Masyarakat"},
		{"hmik.humas.pengmas.02", "hmik.humas.pengmas.02@hmik.com", "Anggota Divisi Pengabdian Masyarakat"},
		{"hmik.humas.kemitraan.01", "hmik.humas.kemitraan.01@hmik.com", "Ketua Divisi Kemitraan Strategis"},
		{"hmik.humas.kemitraan.02", "hmik.humas.kemitraan.02@hmik.com", "Anggota Divisi Kemitraan Strategis"},

		// 6. Manajemen Internal
		{"hmik.internal.01", "hmik.internal.01@hmik.com", "Ketua Departemen Manajemen Internal"},
		{"hmik.internal.02", "hmik.internal.02@hmik.com", "Sekretaris Departemen Manajemen Internal"},
		{"hmik.internal.msda.01", "hmik.internal.msda.01@hmik.com", "Ketua Divisi MSDA"},
		{"hmik.internal.msda.02", "hmik.internal.msda.02@hmik.com", "Anggota Divisi MSDA"},
		{"hmik.internal.aset.01", "hmik.internal.aset.01@hmik.com", "Ketua Divisi Aset & Perlengkapan"},
		{"hmik.internal.aset.02", "hmik.internal.aset.02@hmik.com", "Anggota Divisi Aset & Perlengkapan"},
		{"hmik.internal.minatbakat.01", "hmik.internal.minatbakat.01@hmik.com", "Ketua Divisi Minat & Bakat Internal"},

		// 7. Academy UP
		{"hmik.academy.01", "hmik.academy.01@hmik.com", "Ketua Departemen Academy UP"},
		{"hmik.academy.02", "hmik.academy.02@hmik.com", "Sekretaris Departemen Academy UP"},
		{"hmik.academy.akademik.01", "hmik.academy.akademik.01@hmik.com", "Ketua Divisi Akademik"},
		{"hmik.academy.akademik.02", "hmik.academy.akademik.02@hmik.com", "Anggota Divisi Akademik"},
		{"hmik.academy.keprofesian.01", "hmik.academy.keprofesian.01@hmik.com", "Ketua Divisi Keprofesian"},
		{"hmik.academy.keprofesian.02", "hmik.academy.keprofesian.02@hmik.com", "Anggota Divisi Keprofesian"},
	}

	for _, acc := range accounts {
		var existing entity.Pengguna
		res := db.Where("email = ? OR nama_pengguna = ?", acc.Email, acc.Username).First(&existing)
		if res.Error == nil {
			existing.NamaPengguna = acc.Username
			existing.Email = acc.Email
			existing.Sandi = hashStr
			existing.StatusAktif = true
			db.Save(&existing)
		} else {
			newAcc := entity.Pengguna{
				NamaPengguna: acc.Username,
				Email:        acc.Email,
				Sandi:        hashStr,
				StatusAktif:  true,
			}
			db.Create(&newAcc)
		}
	}

	// Tautkan anggota ke akun slot masing-masing jika jabatan cocok
	for _, acc := range accounts {
		if acc.Jabatan != "" && !strings.Contains(acc.Jabatan, "Admin") && !strings.Contains(acc.Jabatan, "Slot") {
			var p entity.Pengguna
			if err := db.Where("nama_pengguna = ?", acc.Username).First(&p).Error; err == nil {
				db.Model(&entity.Anggota{}).Where("jabatan = ?", acc.Jabatan).Update("id_pengguna", p.IdPengguna)
			}
		}
	}

	log.Println("✅ [Seeder] Seluruh Akun Slot Jabatan HMIK berhasil di-seed (Password: 'qwerqwer')")
	return nil
}
