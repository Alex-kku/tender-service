package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"tender-service/internal/entity"
	"tender-service/internal/repository"
	"tender-service/internal/repository/repoerrs"
)

type TenderService struct {
	repo *repository.Repository
}

func NewTenderService(repo *repository.Repository) *TenderService {
	return &TenderService{repo: repo}
}

func (s *TenderService) GetTenders(ctx context.Context, limit int, offset int, serviceType []string) ([]entity.Tender, error) {
	return s.repo.Tender.GetTenders(ctx, limit, offset, serviceType)
}
func (s *TenderService) CreateTender(ctx context.Context, userId string, input entity.CreateTenderInput) (entity.Tender, error) {
	tender := entity.Tender{
		Name:           input.Name,
		Description:    input.Description,
		ServiceType:    input.ServiceType,
		Status:         "Created",
		OrganizationId: input.OrganizationId,
		Version:        1,
	}
	return s.repo.Tender.CreateTender(ctx, userId, tender)
}
func (s *TenderService) GetTendersByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Tender, error) {
	return s.repo.Tender.GetTendersByUserId(ctx, limit, offset, userId)
}
func (s *TenderService) GetTenderStatus(ctx context.Context, tenderId string, userId *string) (string, error) {
	tender, err := s.repo.Tender.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrTenderNotFound
		}
		logrus.Errorf("TenderService.GetTenderStatus: cannot get tender: %v", err)
		return "", ErrCannotGetTender
	}
	if tender.Status != "Published" {
		if userId == nil {
			return "", ErrNotEnoughPermissions
		}
		userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, *userId)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return "", ErrNotEnoughPermissions
			}
			logrus.Errorf("TenderService.GetTenderStatus: cannot get user organization id: %v", err)
			return "", ErrNotEnoughPermissions
		}
		if userOrganizationId != tender.OrganizationId {
			return "", ErrNotEnoughPermissions
		}
	}
	return tender.Status, nil
}
func (s *TenderService) UpdateTenderStatus(ctx context.Context, tenderId, status, userId string) (entity.Tender, error) {
	tender, err := s.repo.Tender.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get tender: %v", err)
		return entity.Tender{}, ErrCannotGetTender
	}

	if tender.Status == "Closed" {
		return entity.Tender{}, fmt.Errorf("cannot edit closed tender")
	} else if tender.Status == status {
		return entity.Tender{}, fmt.Errorf("impossible change the status to the same")
	} else if tender.Status == "Published" && status == "Created" {
		return entity.Tender{}, fmt.Errorf("impossible change the status to the previous one")
	}

	userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrNotEnoughPermissions
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get user organization id: %v", err)
		return entity.Tender{}, ErrNotEnoughPermissions
	}
	if userOrganizationId != tender.OrganizationId {
		return entity.Tender{}, ErrNotEnoughPermissions
	}

	tender, err = s.repo.Tender.UpdateTenderStatus(ctx, tenderId, status)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.GetTenderStatus: cannot update tender: %v", err)
		return entity.Tender{}, ErrCannotUpdateTender
	}

	return tender, nil
}
func (s *TenderService) EditTender(ctx context.Context, tenderId string, userId string, input entity.EditTenderInput) (entity.Tender, error) {
	tender, err := s.repo.Tender.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get tender: %v", err)
		return entity.Tender{}, ErrCannotGetTender
	}
	if tender.Status == "Closed" {
		return entity.Tender{}, fmt.Errorf("cannot edit closed tender")
	}

	userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrNotEnoughPermissions
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get user organization id: %v", err)
		return entity.Tender{}, ErrNotEnoughPermissions
	}
	if userOrganizationId != tender.OrganizationId {
		return entity.Tender{}, ErrNotEnoughPermissions
	}

	tender, err = s.repo.Tender.UpdateTender(ctx, tenderId, input)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.GetTenderStatus: cannot update tender: %v", err)
		return entity.Tender{}, ErrCannotUpdateTender
	}

	return tender, nil
}
func (s *TenderService) RollbackTender(ctx context.Context, tenderId string, version int, userId string) (entity.Tender, error) {
	tender, err := s.repo.Tender.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get tender: %v", err)
		return entity.Tender{}, ErrCannotGetTender
	}
	if tender.Version <= version {
		return entity.Tender{}, fmt.Errorf("cannot roll back to the current and higher version")
	}
	if tender.Status == "Closed" {
		return entity.Tender{}, fmt.Errorf("cannot edit closed tender")
	}

	userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrNotEnoughPermissions
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get user organization id: %v", err)
		return entity.Tender{}, ErrNotEnoughPermissions
	}
	if userOrganizationId != tender.OrganizationId {
		return entity.Tender{}, ErrNotEnoughPermissions
	}

	tender, err = s.repo.Tender.RollbackTender(ctx, tenderId, version)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Tender{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.GetTenderStatus: cannot update tender: %v", err)
		return entity.Tender{}, ErrCannotUpdateTender
	}

	return tender, nil
}
