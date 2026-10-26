package usecase

import (
	"context"

	"github.com/HMIKProject/hmik-corex-backend/internal/entity"
	"github.com/HMIKProject/hmik-corex-backend/internal/repository"
)

type DashboardUsecase interface {
	GetDashboardStats(ctx context.Context) (*entity.DashboardStatsResponse, error)
	GetTotalAnggotaAktif(ctx context.Context) (*entity.CardDetail, error)
	GetTotalModul(ctx context.Context) (*entity.CardDetail, error)
	GetTotalProgramKerja(ctx context.Context) (*entity.CardDetail, error)
	GetVotingAktif(ctx context.Context) (*entity.CardDetail, error)
}

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

func (u *dashboardUsecase) GetDashboardStats(ctx context.Context) (*entity.DashboardStatsResponse, error) {
	return u.repo.GetDashboardStats(ctx)
}

func (u *dashboardUsecase) GetTotalAnggotaAktif(ctx context.Context) (*entity.CardDetail, error) {
	return u.repo.GetTotalAnggotaAktif(ctx)
}

func (u *dashboardUsecase) GetTotalModul(ctx context.Context) (*entity.CardDetail, error) {
	return u.repo.GetTotalModul(ctx)
}

func (u *dashboardUsecase) GetTotalProgramKerja(ctx context.Context) (*entity.CardDetail, error) {
	return u.repo.GetTotalProgramKerja(ctx)
}

func (u *dashboardUsecase) GetVotingAktif(ctx context.Context) (*entity.CardDetail, error) {
	return u.repo.GetVotingAktif(ctx)
}