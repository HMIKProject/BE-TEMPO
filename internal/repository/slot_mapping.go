package repository

import (
	"strings"
)

// SlotOrgMapping mendefinisikan relasi struktural sebuah slot akun dengan struktur organisasi HMIK
type SlotOrgMapping struct {
	Jabatan        string `json:"jabatan"`
	IdDepartemen   *int64 `json:"id_departemen,omitempty"`
	NamaDepartemen string `json:"nama_departemen,omitempty"`
	IdDivisi       *int64 `json:"id_divisi,omitempty"`
	NamaDivisi     string `json:"nama_divisi,omitempty"`
}

func int64Ptr(v int64) *int64 {
	return &v
}

// GetSlotOrgMapping memetakan slot akun secara deterministik ke Jabatan, Departemen, dan Divisi resmi
func GetSlotOrgMapping(username string, fullName string) SlotOrgMapping {
	u := strings.ToLower(strings.TrimSpace(username))

	switch {
	// BPH (Badan Pengurus Harian)
	case strings.HasPrefix(u, "hmik.kahim"):
		return SlotOrgMapping{Jabatan: "Ketua Himpunan"}
	case strings.HasPrefix(u, "hmik.wakahim"):
		return SlotOrgMapping{Jabatan: "Wakil Ketua Himpunan"}
	case u == "hmik.sekretaris.01":
		return SlotOrgMapping{Jabatan: "Sekretaris 1"}
	case u == "hmik.sekretaris.02":
		return SlotOrgMapping{Jabatan: "Sekretaris 2"}
	case u == "hmik.bendahara.01":
		return SlotOrgMapping{Jabatan: "Bendahara 1"}
	case u == "hmik.bendahara.02":
		return SlotOrgMapping{Jabatan: "Bendahara 2"}

	// Departemen 1: Riset dan Teknologi
	case u == "hmik.ristek.01":
		return SlotOrgMapping{Jabatan: "Ketua Departemen Riset dan Teknologi", IdDepartemen: int64Ptr(1), NamaDepartemen: "Riset dan Teknologi"}
	case u == "hmik.ristek.02":
		return SlotOrgMapping{Jabatan: "Sekretaris Departemen Riset dan Teknologi", IdDepartemen: int64Ptr(1), NamaDepartemen: "Riset dan Teknologi"}
	case u == "hmik.ristek.inovasi.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Riset & Inovasi", IdDepartemen: int64Ptr(1), NamaDepartemen: "Riset dan Teknologi", IdDivisi: int64Ptr(1), NamaDivisi: "Riset & Inovasi"}
	case u == "hmik.ristek.inovasi.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Riset & Inovasi", IdDepartemen: int64Ptr(1), NamaDepartemen: "Riset dan Teknologi", IdDivisi: int64Ptr(1), NamaDivisi: "Riset & Inovasi"}
	case u == "hmik.ristek.produk.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Pengembangan Produk", IdDepartemen: int64Ptr(1), NamaDepartemen: "Riset dan Teknologi", IdDivisi: int64Ptr(2), NamaDivisi: "Pengembangan Produk"}
	case u == "hmik.ristek.produk.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Pengembangan Produk", IdDepartemen: int64Ptr(1), NamaDepartemen: "Riset dan Teknologi", IdDivisi: int64Ptr(2), NamaDivisi: "Pengembangan Produk"}

	// Departemen 2: Media Informasi
	case u == "hmik.medinfo.01":
		return SlotOrgMapping{Jabatan: "Ketua Departemen Media Informasi", IdDepartemen: int64Ptr(2), NamaDepartemen: "Media Informasi"}
	case u == "hmik.medinfo.02":
		return SlotOrgMapping{Jabatan: "Sekretaris Departemen Media Informasi", IdDepartemen: int64Ptr(2), NamaDepartemen: "Media Informasi"}
	case u == "hmik.medinfo.kreatif.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Media Kreatif", IdDepartemen: int64Ptr(2), NamaDepartemen: "Media Informasi", IdDivisi: int64Ptr(3), NamaDivisi: "Media Kreatif"}
	case u == "hmik.medinfo.kreatif.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Media Kreatif", IdDepartemen: int64Ptr(2), NamaDepartemen: "Media Informasi", IdDivisi: int64Ptr(3), NamaDivisi: "Media Kreatif"}
	case u == "hmik.medinfo.pubdok.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Publikasi & Dokumentasi", IdDepartemen: int64Ptr(2), NamaDepartemen: "Media Informasi", IdDivisi: int64Ptr(4), NamaDivisi: "Publikasi & Dokumentasi"}
	case u == "hmik.medinfo.pubdok.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Publikasi & Dokumentasi", IdDepartemen: int64Ptr(2), NamaDepartemen: "Media Informasi", IdDivisi: int64Ptr(4), NamaDivisi: "Publikasi & Dokumentasi"}

	// Departemen 3: Hubungan External
	case u == "hmik.humas.01":
		return SlotOrgMapping{Jabatan: "Ketua Departemen Hubungan External", IdDepartemen: int64Ptr(3), NamaDepartemen: "Hubungan External"}
	case u == "hmik.humas.02":
		return SlotOrgMapping{Jabatan: "Sekretaris Departemen Hubungan External", IdDepartemen: int64Ptr(3), NamaDepartemen: "Hubungan External"}
	case u == "hmik.humas.kemitraan.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Kemitraan Strategis", IdDepartemen: int64Ptr(3), NamaDepartemen: "Hubungan External", IdDivisi: int64Ptr(5), NamaDivisi: "Kemitraan Strategis"}
	case u == "hmik.humas.kemitraan.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Kemitraan Strategis", IdDepartemen: int64Ptr(3), NamaDepartemen: "Hubungan External", IdDivisi: int64Ptr(5), NamaDivisi: "Kemitraan Strategis"}
	case u == "hmik.humas.pengmas.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Pengabdian Masyarakat", IdDepartemen: int64Ptr(3), NamaDepartemen: "Hubungan External", IdDivisi: int64Ptr(6), NamaDivisi: "Pengabdian Masyarakat"}
	case u == "hmik.humas.pengmas.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Pengabdian Masyarakat", IdDepartemen: int64Ptr(3), NamaDepartemen: "Hubungan External", IdDivisi: int64Ptr(6), NamaDivisi: "Pengabdian Masyarakat"}

	// Departemen 4: Manajemen Internal
	case u == "hmik.internal.01":
		return SlotOrgMapping{Jabatan: "Ketua Departemen Manajemen Internal", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal"}
	case u == "hmik.internal.02":
		return SlotOrgMapping{Jabatan: "Sekretaris Departemen Manajemen Internal", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal"}
	case u == "hmik.internal.msda.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi MSDA", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal", IdDivisi: int64Ptr(7), NamaDivisi: "MSDA"}
	case u == "hmik.internal.msda.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi MSDA", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal", IdDivisi: int64Ptr(7), NamaDivisi: "MSDA"}
	case u == "hmik.internal.aset.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Aset & Perlengkapan", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal", IdDivisi: int64Ptr(8), NamaDivisi: "Aset & Perlengkapan"}
	case u == "hmik.internal.aset.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Aset & Perlengkapan", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal", IdDivisi: int64Ptr(8), NamaDivisi: "Aset & Perlengkapan"}
	case u == "hmik.internal.minatbakat.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Minat & Bakat Internal", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal", IdDivisi: int64Ptr(9), NamaDivisi: "Minat & Bakat"}
	case u == "hmik.internal.minatbakat.02", u == "hmik.human.01":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Minat & Bakat Internal", IdDepartemen: int64Ptr(4), NamaDepartemen: "Manajemen Internal", IdDivisi: int64Ptr(9), NamaDivisi: "Minat & Bakat"}

	// Departemen 5: Academy UP
	case u == "hmik.academy.01":
		return SlotOrgMapping{Jabatan: "Ketua Departemen Academy UP", IdDepartemen: int64Ptr(5), NamaDepartemen: "Academy UP"}
	case u == "hmik.academy.02":
		return SlotOrgMapping{Jabatan: "Sekretaris Departemen Academy UP", IdDepartemen: int64Ptr(5), NamaDepartemen: "Academy UP"}
	case u == "hmik.academy.akademik.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Akademik", IdDepartemen: int64Ptr(5), NamaDepartemen: "Academy UP", IdDivisi: int64Ptr(10), NamaDivisi: "Akademik"}
	case u == "hmik.academy.akademik.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Akademik", IdDepartemen: int64Ptr(5), NamaDepartemen: "Academy UP", IdDivisi: int64Ptr(10), NamaDivisi: "Akademik"}
	case u == "hmik.academy.keprofesian.01":
		return SlotOrgMapping{Jabatan: "Ketua Divisi Keprofesian", IdDepartemen: int64Ptr(5), NamaDepartemen: "Academy UP", IdDivisi: int64Ptr(11), NamaDivisi: "Keprofesian"}
	case u == "hmik.academy.keprofesian.02":
		return SlotOrgMapping{Jabatan: "Anggota Divisi Keprofesian", IdDepartemen: int64Ptr(5), NamaDepartemen: "Academy UP", IdDivisi: int64Ptr(11), NamaDivisi: "Keprofesian"}

	default:
		j := fullName
		if j == "" {
			j = username
		}
		return SlotOrgMapping{Jabatan: j}
	}
}
