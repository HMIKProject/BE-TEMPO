package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PemiraRepository mendefinisikan seluruh kontrak akses data untuk sistem Pemira (E-Voting)
type PemiraRepository interface {
	// 1. Pemira
	GetActivePemira(ctx context.Context) (*entity.Pemira, error)
	GetPemiraByID(ctx context.Context, idPemira int64) (*entity.Pemira, error)
	GetLatestPemira(ctx context.Context) (*entity.Pemira, error)
	UpdatePemiraStatus(ctx context.Context, idPemira int64, status string) error
	UpdatePemiraSchedule(ctx context.Context, idPemira int64, tahunPeriode string, tMulai, tSelesai time.Time) error
	CreatePemira(ctx context.Context, pemira *entity.Pemira) error

	// 2. Paslon (Kandidat)
	GetPaslonByPemiraID(ctx context.Context, idPemira int64) ([]entity.Paslon, error)
	GetPaslonByID(ctx context.Context, idPaslon int64) (*entity.Paslon, error)
	CreatePaslon(ctx context.Context, paslon *entity.Paslon) error
	UpdatePaslon(ctx context.Context, paslon *entity.Paslon) error
	DeletePaslon(ctx context.Context, idPaslon int64) error
	CountSuaraByPaslon(ctx context.Context, idPaslon int64) (int64, error)

	// 3. Pemilih (DPT & Autentikasi OTP)
	CreatePemilih(ctx context.Context, pemilih *entity.Pemilih) error
	FindPemilihByNimAndPemira(ctx context.Context, nim string, idPemira int64) (*entity.Pemilih, error)
	UpdateOTP(ctx context.Context, idPemilih int64, otpCode string, expiredAt time.Time) error
	ClearOTP(ctx context.Context, idPemilih int64) error
	GetPemilihList(ctx context.Context, idPemira int64, search string, status *bool, limit, offset int) ([]entity.Pemilih, int64, error)
	ResetPemilihStatus(ctx context.Context, idPemilih int64) error

	// 4. Suara & Transaksi Atomik (Anti Double-Voting & LUBER JURDIL)
	CastVote(ctx context.Context, idPemira int64, nim string, idPaslon int64) error
	CountTotalSuara(ctx context.Context, idPemira int64) (int64, error)
	CountSuaraPerPaslon(ctx context.Context, idPemira int64) (map[int64]int64, error)
	CountTotalDPT(ctx context.Context, idPemira int64) (int64, error)
	ResetKotakSuara(ctx context.Context, idPemira int64) error

	// 5. Pelantikan MPIK & Transisi Manajemen Himpunan
	GetWinningPaslon(ctx context.Context, idPemira int64) (*entity.Paslon, int64, error)
	PromoteKabinet(ctx context.Context, kabinet *entity.Kabinet) error
	UpsertAnggotaBPH(ctx context.Context, nim, nama, jabatan string, angkatan int, foto string) (*entity.Anggota, error)
	DemisionerkanPengurusLama(ctx context.Context, exceptNIMs []string) error
	PurgeVotersAndVotes(ctx context.Context, idPemira int64) error
}

type pemiraRepository struct {
	db *gorm.DB
}

// NewPemiraRepository membuat instance baru PemiraRepository
func NewPemiraRepository(db *gorm.DB) PemiraRepository {
	return &pemiraRepository{db: db}
}

// ==========================================
// 1. PEMIRA QUERIES
// ==========================================

func (r *pemiraRepository) GetActivePemira(ctx context.Context) (*entity.Pemira, error) {
	var pemira entity.Pemira
	err := r.db.WithContext(ctx).
		Where("status = ?", "BERJALAN").
		Order("tanggal_mulai DESC").
		First(&pemira).Error
	if err != nil {
		return nil, err
	}
	return &pemira, nil
}

func (r *pemiraRepository) GetPemiraByID(ctx context.Context, idPemira int64) (*entity.Pemira, error) {
	var pemira entity.Pemira
	err := r.db.WithContext(ctx).
		Where("id_pemira = ?", idPemira).
		First(&pemira).Error
	if err != nil {
		return nil, err
	}
	return &pemira, nil
}

// ==========================================
// 2. PASLON QUERIES
// ==========================================

func (r *pemiraRepository) GetPaslonByPemiraID(ctx context.Context, idPemira int64) ([]entity.Paslon, error) {
	var paslons []entity.Paslon
	err := r.db.WithContext(ctx).
		Where("id_pemira = ?", idPemira).
		Order("nomor_urut ASC").
		Find(&paslons).Error
	return paslons, err
}

func (r *pemiraRepository) GetPaslonByID(ctx context.Context, idPaslon int64) (*entity.Paslon, error) {
	var paslon entity.Paslon
	err := r.db.WithContext(ctx).
		Where("id_paslon = ?", idPaslon).
		First(&paslon).Error
	if err != nil {
		return nil, err
	}
	return &paslon, nil
}

// ==========================================
// 3. PEMILIH (DPT & OTP) QUERIES
// ==========================================

func (r *pemiraRepository) CreatePemilih(ctx context.Context, pemilih *entity.Pemilih) error {
	return r.db.WithContext(ctx).Create(pemilih).Error
}

func (r *pemiraRepository) FindPemilihByNimAndPemira(ctx context.Context, nim string, idPemira int64) (*entity.Pemilih, error) {
	var pemilih entity.Pemilih
	err := r.db.WithContext(ctx).
		Where("nim = ? AND id_pemira = ?", nim, idPemira).
		First(&pemilih).Error
	if err != nil {
		return nil, err
	}
	return &pemilih, nil
}

func (r *pemiraRepository) UpdateOTP(ctx context.Context, idPemilih int64, otpCode string, expiredAt time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.Pemilih{}).
		Where("id_pemilih = ?", idPemilih).
		Updates(map[string]interface{}{
			"kode_otp":    otpCode,
			"otp_expired": expiredAt,
		}).Error
}

func (r *pemiraRepository) ClearOTP(ctx context.Context, idPemilih int64) error {
	return r.db.WithContext(ctx).Model(&entity.Pemilih{}).
		Where("id_pemilih = ?", idPemilih).
		Updates(map[string]interface{}{
			"kode_otp":    gorm.Expr("NULL"),
			"otp_expired": gorm.Expr("NULL"),
		}).Error
}

func (r *pemiraRepository) GetPemilihList(ctx context.Context, idPemira int64, search string, status *bool, limit, offset int) ([]entity.Pemilih, int64, error) {
	var pemilihs []entity.Pemilih
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Pemilih{}).Where("id_pemira = ?", idPemira)

	if search != "" {
		query = query.Where("nim ILIKE ? OR email_kampus ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if status != nil {
		query = query.Where("status_memilih = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("nim ASC").Limit(limit).Offset(offset).Find(&pemilihs).Error
	return pemilihs, total, err
}

func (r *pemiraRepository) ResetPemilihStatus(ctx context.Context, idPemilih int64) error {
	return r.db.WithContext(ctx).Model(&entity.Pemilih{}).
		Where("id_pemilih = ?", idPemilih).
		Update("status_memilih", false).Error
}

// ==========================================
// 4. SUARA & TRANSAKSI ATOMIK (LUBER JURDIL & ANTI DOUBLE-VOTING)
// ==========================================

// CastVote mengeksekusi transaksi atomik pencatatan suara:
// 1. Mengunci baris data pemilih dengan row-level lock (SELECT FOR UPDATE)
// 2. Memastikan pemilih belum pernah memilih
// 3. Mengubah status_memilih = true
// 4. Menambahkan surat suara anonim murni ke tabel suara tanpa menyimpan relasi NIM pemilih
func (r *pemiraRepository) CastVote(ctx context.Context, idPemira int64, nim string, idPaslon int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var pemilih entity.Pemilih

		// 1. Kunci baris pemilih untuk mencegah race condition pencoblosan ganda (Anti Double-Voting)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("nim = ? AND id_pemira = ?", nim, idPemira).
			First(&pemilih).Error; err != nil {
			return fmt.Errorf("pemilih dengan NIM %s tidak ditemukan dalam DPT: %w", nim, err)
		}

		// 2. Validasi status hak suara
		if pemilih.StatusMemilih {
			return errors.New("hak suara Anda sudah digunakan sebelumnya")
		}

		// 3. Ubah status hak suara pemilih menjadi telah memilih
		if err := tx.Model(&pemilih).Update("status_memilih", true).Error; err != nil {
			return fmt.Errorf("gagal memperbarui status hak suara: %w", err)
		}

		// 4. Masukkan surat suara murni anonim (TANPA kolom nim atau id_pemilih demi asas Luber Jurdil)
		suara := entity.Suara{
			IdPemira:     idPemira,
			IdPaslon:     idPaslon,
			WaktuMemilih: time.Now(),
		}
		if err := tx.Create(&suara).Error; err != nil {
			return fmt.Errorf("gagal mencatat surat suara ke kotak suara: %w", err)
		}

		return nil
	})
}

func (r *pemiraRepository) CountTotalSuara(ctx context.Context, idPemira int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("suara").
		Where("id_pemira = ?", idPemira).
		Count(&count).Error
	return count, err
}

type paslonVoteAggregate struct {
	IdPaslon int64 `gorm:"column:id_paslon"`
	Total    int64 `gorm:"column:total"`
}

func (r *pemiraRepository) CountSuaraPerPaslon(ctx context.Context, idPemira int64) (map[int64]int64, error) {
	var results []paslonVoteAggregate
	err := r.db.WithContext(ctx).Table("suara").
		Select("id_paslon, count(*) as total").
		Where("id_pemira = ?", idPemira).
		Group("id_paslon").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[int64]int64)
	for _, res := range results {
		counts[res.IdPaslon] = res.Total
	}
	return counts, nil
}

func (r *pemiraRepository) CountTotalDPT(ctx context.Context, idPemira int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("pemilih").
		Where("id_pemira = ?", idPemira).
		Count(&count).Error
	return count, err
}

func (r *pemiraRepository) ResetKotakSuara(ctx context.Context, idPemira int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Hapus seluruh surat suara di kotak suara untuk pemira ini
		if err := tx.Where("id_pemira = ?", idPemira).Delete(&entity.Suara{}).Error; err != nil {
			return fmt.Errorf("gagal mengosongkan kotak suara: %w", err)
		}

		// 2. Reset status hak suara seluruh pemilih pada pemira ini ke false dan bersihkan OTP
		if err := tx.Model(&entity.Pemilih{}).Where("id_pemira = ?", idPemira).
			Updates(map[string]interface{}{
				"status_memilih": false,
				"kode_otp":       gorm.Expr("NULL"),
				"otp_expired":    gorm.Expr("NULL"),
			}).Error; err != nil {
			return fmt.Errorf("gagal mereset status pemilih: %w", err)
		}

		return nil
	})
}

func (r *pemiraRepository) GetLatestPemira(ctx context.Context) (*entity.Pemira, error) {
	var pemira entity.Pemira
	err := r.db.WithContext(ctx).Order("id_pemira DESC").First(&pemira).Error
	if err != nil {
		return nil, err
	}
	return &pemira, nil
}

func (r *pemiraRepository) UpdatePemiraStatus(ctx context.Context, idPemira int64, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Pemira{}).
		Where("id_pemira = ?", idPemira).
		Update("status", status).Error
}

func (r *pemiraRepository) UpdatePemiraSchedule(ctx context.Context, idPemira int64, tahunPeriode string, tMulai, tSelesai time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.Pemira{}).
		Where("id_pemira = ?", idPemira).
		Updates(map[string]interface{}{
			"tahun_periode":   tahunPeriode,
			"tanggal_mulai":   tMulai,
			"tanggal_selesai": tSelesai,
			"diperbarui_pada": time.Now(),
		}).Error
}

func (r *pemiraRepository) CreatePemira(ctx context.Context, pemira *entity.Pemira) error {
	return r.db.WithContext(ctx).Create(pemira).Error
}

func (r *pemiraRepository) CreatePaslon(ctx context.Context, paslon *entity.Paslon) error {
	return r.db.WithContext(ctx).Create(paslon).Error
}

func (r *pemiraRepository) UpdatePaslon(ctx context.Context, paslon *entity.Paslon) error {
	return r.db.WithContext(ctx).Save(paslon).Error
}

func (r *pemiraRepository) DeletePaslon(ctx context.Context, idPaslon int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Paslon{}, idPaslon).Error
}

func (r *pemiraRepository) CountSuaraByPaslon(ctx context.Context, idPaslon int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("suara").
		Where("id_paslon = ?", idPaslon).
		Count(&count).Error
	return count, err
}

func (r *pemiraRepository) GetWinningPaslon(ctx context.Context, idPemira int64) (*entity.Paslon, int64, error) {
	type voteResult struct {
		IdPaslon int64 `gorm:"column:id_paslon"`
		Total    int64 `gorm:"column:total"`
	}
	var topVote voteResult
	err := r.db.WithContext(ctx).Table("suara").
		Select("id_paslon, count(*) as total").
		Where("id_pemira = ?", idPemira).
		Group("id_paslon").
		Order("total DESC").
		Limit(1).
		Scan(&topVote).Error
	if err != nil {
		return nil, 0, err
	}

	if topVote.IdPaslon == 0 {
		return nil, 0, errors.New("belum ada surat suara yang masuk untuk menentukan pemenang")
	}

	var paslon entity.Paslon
	if err := r.db.WithContext(ctx).First(&paslon, topVote.IdPaslon).Error; err != nil {
		return nil, 0, err
	}

	return &paslon, topVote.Total, nil
}

func (r *pemiraRepository) PromoteKabinet(ctx context.Context, kabinet *entity.Kabinet) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Nonaktifkan kabinet lama
		if err := tx.Model(&entity.Kabinet{}).Where("is_aktif = ?", true).Update("is_aktif", false).Error; err != nil {
			return err
		}

		// 2. Tambahkan kabinet baru
		kabinet.IsAktif = true
		if err := tx.Create(kabinet).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *pemiraRepository) UpsertAnggotaBPH(ctx context.Context, nim, nama, jabatan string, angkatan int, foto string) (*entity.Anggota, error) {
	var existing entity.Anggota
	err := r.db.WithContext(ctx).Where("nim = ?", nim).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAnggota := entity.Anggota{
				Nim:          nim,
				NamaLengkap:  nama,
				ProgramStudi: "Ilmu komputer",
				Angkatan:     angkatan,
				Jabatan:      jabatan,
				Status:       "Pengurus",
				FotoProfil:   foto,
			}
			if err := r.db.WithContext(ctx).Create(&newAnggota).Error; err != nil {
				return nil, err
			}
			return &newAnggota, nil
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"nama_lengkap": nama,
		"jabatan":      jabatan,
		"status":       "Pengurus",
		"diubah_pada":  time.Now(),
	}
	if foto != "" {
		updates["foto_profil"] = foto
	}
	if err := r.db.WithContext(ctx).Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}

	r.db.WithContext(ctx).First(&existing, existing.IdAnggota)
	return &existing, nil
}

func (r *pemiraRepository) DemisionerkanPengurusLama(ctx context.Context, exceptNIMs []string) error {
	query := r.db.WithContext(ctx).Model(&entity.Anggota{}).
		Where("jabatan IN ?", []string{"Ketua Himpunan", "Wakil Ketua Himpunan"}).
		Where("status = ?", "Pengurus")
	if len(exceptNIMs) > 0 {
		query = query.Where("nim NOT IN ?", exceptNIMs)
	}
	return query.Update("status", "Demisioner").Error
}

func (r *pemiraRepository) PurgeVotersAndVotes(ctx context.Context, idPemira int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Hapus seluruh surat suara di kotak suara
		if err := tx.Where("id_pemira = ?", idPemira).Delete(&entity.Suara{}).Error; err != nil {
			return fmt.Errorf("gagal membersihkan surat suara pasca-pelantikan: %w", err)
		}

		// 2. Hapus seluruh data pemilih (NIM, email kampus, dll) demi menjaga privasi
		if err := tx.Where("id_pemira = ?", idPemira).Delete(&entity.Pemilih{}).Error; err != nil {
			return fmt.Errorf("gagal membersihkan data pemilih pasca-pelantikan: %w", err)
		}

		return nil
	})
}
