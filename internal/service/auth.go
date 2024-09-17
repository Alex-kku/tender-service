package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"tender-service/internal/repository"
	"tender-service/internal/repository/repoerrs"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GetUserId(ctx context.Context, username string) (string, error) {
	id, err := s.repo.Auth.GetUserId(ctx, username)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrUserNotFound
		}
		logrus.Errorf("AuthService.GetUserId: cannot get user id: %v", err)
		return "", ErrCannotGetUserId
	}

	return id, nil
}
func (s *AuthService) CheckResponsibility(ctx context.Context, userId string, organizationId string) (bool, error) {
	ok, err := s.repo.Auth.OrganizationIsExist(ctx, organizationId)
	if !ok {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return false, ErrOrganizationNotFound
		}
		logrus.Errorf("AuthService.CheckResponsibility: cannot check existence of the organization: %v", err)
		return false, ErrOrganizationNotFound
	}

	userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return false, ErrUserOrganizationNotFound
		}
		logrus.Errorf("AuthService.GetUserId: cannot get user organization id: %v", err)
		return false, ErrUserOrganizationNotFound
	}

	if userOrganizationId != organizationId {
		return false, ErrNotEnoughPermissions
	}

	return true, nil
}
