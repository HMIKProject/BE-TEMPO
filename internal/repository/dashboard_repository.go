package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDashboardStats(ctx context.Context) (*entity.DashboardStatsResponse, error)
	GetTotalAnggotaAktif(ctx context.Context) (*entity.CardDetail, error)
	GetTotalModul(ctx context.Context) (*entity.CardDetail, error)
	GetTotalProgramKerja(ctx context.Context) (*entity.CardDetail, error)
	GetVotingAktif(ctx context.Context) (*entity.CardDetail, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// getAnggotaCountFromJSON membaca cmd/seeder/data_anggota.json sebagai fallback jika database belum terisi
func getAnggotaCountFromJSON() int64 {
	paths := []string{
		"cmd/seeder/data_anggota.json",
		"../cmd/seeder/data_anggota.json",
		"../../cmd/seeder/data_anggota.json",
	}

	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			var anggotas []map[string]interface{}
			if err := json.Unmarshal(data, &anggotas); err == nil && len(anggotas) > 0 {
				return int64(len(anggotas))
			}
		}
	}
	return 50
}

func (r *dashboardRepository) GetTotalAnggotaAktif(ctx context.Context) (*entity.CardDetail, error) {
	var total, bulanIni int64
	if r.db != nil {
		_ = r.db.WithContext(ctx).Table("anggota").
			Where("status IN ?", []string{"Pengurus", "Magang"}).
			Count(&total).Error

		_ = r.db.WithContext(ctx).Table("anggota").
			Where("status IN ? AND dibuat_pada >= DATE_TRUNC('month', CURRENT_DATE)", []string{"Pengurus", "Magang"}).
			Count(&bulanIni).Error
	}

	return &entity.CardDetail{
		CardKey:             "total_anggota_aktif",
		Judul:               "Total Anggota Aktif",
		Total:               total,
		PertumbuhanBulanIni: bulanIni,
		Keterangan:          fmt.Sprintf("+ %d Mahasiswa Bulan ini", bulanIni),
		WarnaTema:           "blue",
		TautanAksi:          "/admin/buku-induk-anggota",
	}, nil
}

func (r *dashboardRepository) GetTotalModul(ctx context.Context) (*entity.CardDetail, error) {
	var total, bulanIni int64
	if r.db != nil {
		_ = r.db.WithContext(ctx).Table("modul_pembelajaran").Count(&total).Error
		_ = r.db.WithContext(ctx).Table("modul_pembelajaran").
			Where("dibuat_pada >= DATE_TRUNC('month', CURRENT_DATE)").
			Count(&bulanIni).Error
	}

	return &entity.CardDetail{
		CardKey:             "total_modul",
		Judul:               "Total Modul",
		Total:               total,
		PertumbuhanBulanIni: bulanIni,
		Keterangan:          fmt.Sprintf("+ %d Modul Bulan ini", bulanIni),
		WarnaTema:           "yellow",
		TautanAksi:          "/admin/modul",
	}, nil
}

func (r *dashboardRepository) GetTotalProgramKerja(ctx context.Context) (*entity.CardDetail, error) {
	var total, bulanIni int64
	if r.db != nil {
		_ = r.db.WithContext(ctx).Table("program_kerja").Count(&total).Error
		_ = r.db.WithContext(ctx).Table("program_kerja").
			Where("dibuat_pada >= DATE_TRUNC('month', CURRENT_DATE)").
			Count(&bulanIni).Error
	}

	return &entity.CardDetail{
		CardKey:             "total_program_kerja",
		Judul:               "Total Program Kerja",
		Total:               total,
		PertumbuhanBulanIni: bulanIni,
		Keterangan:          fmt.Sprintf("+ %d Program Kerja Bulan ini", bulanIni),
		WarnaTema:           "green",
		TautanAksi:          "/admin/program-kerja",
	}, nil
}

func (r *dashboardRepository) GetVotingAktif(ctx context.Context) (*entity.CardDetail, error) {
	var total, bulanIni int64
	if r.db != nil {
		_ = r.db.WithContext(ctx).Table("pemira").Where("status = ?", "BERJALAN").Count(&total).Error
		_ = r.db.WithContext(ctx).Table("suara").Count(&bulanIni).Error
	}

	ket := "Belum ada Voting Aktif"
	if total > 0 {
		ket = fmt.Sprintf("%d Suara Masuk", bulanIni)
	}

	return &entity.CardDetail{
		CardKey:             "voting_aktif",
		Judul:               "Voting Aktif",
		Total:               total,
		PertumbuhanBulanIni: bulanIni,
		Keterangan:          ket,
		WarnaTema:           "orange",
		TautanAksi:          "/admin/pemira",
	}, nil
}

func (r *dashboardRepository) GetDashboardStats(ctx context.Context) (*entity.DashboardStatsResponse, error) {
	cardAnggota, _ := r.GetTotalAnggotaAktif(ctx)
	cardModul, _ := r.GetTotalModul(ctx)
	cardProker, _ := r.GetTotalProgramKerja(ctx)
	cardVoting, _ := r.GetVotingAktif(ctx)

	return &entity.DashboardStatsResponse{
		TotalMahasiswaAktif: cardAnggota.Total,
		TotalModul:          cardModul.Total,
		TotalProgramKerja:   cardProker.Total,
		VotingAktif:         cardVoting.Total,
		MahasiswaAktif: entity.StatItem{
			Total:               cardAnggota.Total,
			PertumbuhanBulanIni: cardAnggota.PertumbuhanBulanIni,
			Keterangan:          cardAnggota.Keterangan,
		},
		Modul: entity.StatItem{
			Total:               cardModul.Total,
			PertumbuhanBulanIni: cardModul.PertumbuhanBulanIni,
			Keterangan:          cardModul.Keterangan,
		},
		ProgramKerja: entity.StatItem{
			Total:               cardProker.Total,
			PertumbuhanBulanIni: cardProker.PertumbuhanBulanIni,
			Keterangan:          cardProker.Keterangan,
		},
		Voting: entity.VotingStatItem{
			Total:          cardVoting.Total,
			SedangBerjalan: cardVoting.PertumbuhanBulanIni,
			Keterangan:     cardVoting.Keterangan,
		},
		CardTotalAnggotaAktif: *cardAnggota,
		CardTotalModul:        *cardModul,
		CardTotalProgramKerja: *cardProker,
		CardVotingAktif:       *cardVoting,
	}, nil
}