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

type BidService struct {
	repo *repository.Repository
}

func NewBidService(repo *repository.Repository) *BidService {
	return &BidService{repo: repo}
}

func (s *BidService) CreateBid(ctx context.Context, input entity.CreateBidInput) (entity.Bid, error) {
	tender, err := s.repo.Tender.GetTenderById(ctx, input.TenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrTenderNotFound
		}
		logrus.Errorf("TenderService.UpdateTenderStatus: cannot get tender: %v", err)
		return entity.Bid{}, ErrCannotGetTender
	}
	if tender.Status == "Closed" {
		return entity.Bid{}, fmt.Errorf("cannot edit closed tender")
	}

	ok, err := s.repo.Auth.UserIsExist(ctx, input.AuthorId)
	if !ok {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrUserNotFound
		}
		logrus.Errorf("BidService.CreateBid: cannot check existence of the user: %v", err)
		return entity.Bid{}, ErrUserNotFound
	}

	if input.AuthorType == "Organization" {
		_, err := s.repo.Auth.GetUserOrganizationId(ctx, input.AuthorId)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return entity.Bid{}, ErrNotEnoughPermissions
			}
			logrus.Errorf("BidService.CreateBid: cannot get user organization id: %v", err)
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	}

	bid := entity.Bid{
		Name:        input.Name,
		Description: input.Description,
		Status:      "Created",
		TenderId:    input.TenderId,
		AuthorType:  input.AuthorType,
		AuthorId:    input.AuthorId,
		Version:     1,
	}

	return s.repo.Bid.CreateBid(ctx, bid)
}
func (s *BidService) GetBidsByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Bid, error) {
	return s.repo.Bid.GetBidsByUserId(ctx, limit, offset, userId)
}
func (s *BidService) GetBidsForTender(ctx context.Context, tenderId string, userId string, limit int, offset int) ([]entity.Bid, error) {
	tender, err := s.repo.Tender.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrTenderNotFound
		}
		logrus.Errorf("BidService.GetBidsForTender: cannot get tender: %v", err)
		return nil, ErrCannotGetTender
	}

	userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrNotEnoughPermissions
		}
		logrus.Errorf("BidService.GetBidsForTender: cannot get user organization id: %v", err)
		return nil, ErrNotEnoughPermissions
	}
	if userOrganizationId != tender.OrganizationId {
		return nil, ErrNotEnoughPermissions
	}

	bids, err := s.repo.Bid.GetBidsForTender(ctx, tenderId, limit, offset)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		logrus.Errorf("BidService.GetBidsForTender: cannot get bid: %v", err)
		return nil, ErrCannotGetBid
	}

	return bids, nil
}
func (s *BidService) GetBidStatus(ctx context.Context, bidId string, userId string) (string, error) {
	bid, err := s.repo.Bid.GetBidById(ctx, bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return "", ErrBidNotFound
		}
		logrus.Errorf("BidService.GetBidStatus: cannot get bid: %v", err)
		return "", ErrCannotGetBid
	}

	if bid.AuthorType == "User" {
		if bid.AuthorId != userId {
			return "", ErrNotEnoughPermissions
		}
	} else {
		userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return "", ErrNotEnoughPermissions
			}
			logrus.Errorf("BidService.GetBidStatus: cannot get user organization id: %v", err)
			return "", ErrNotEnoughPermissions
		}
		tenderOrganizationId, _ := s.repo.Auth.GetUserOrganizationId(ctx, bid.AuthorId)
		if userOrganizationId != tenderOrganizationId {
			return "", ErrNotEnoughPermissions
		}
	}

	return bid.Status, nil
}

func (s *BidService) UpdateBidStatus(ctx context.Context, bidId, status, userId string) (entity.Bid, error) {
	bid, err := s.repo.Bid.GetBidById(ctx, bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.UpdateBidStatus: cannot get bid: %v", err)
		return entity.Bid{}, ErrCannotGetBid
	}

	if bid.Status == "Canceled" {
		return entity.Bid{}, fmt.Errorf("cannot edit canceled bid")
	} else if bid.Status == status {
		return entity.Bid{}, fmt.Errorf("impossible change the status to the same")
	} else if bid.Status == "Published" && status == "Created" {
		return entity.Bid{}, fmt.Errorf("impossible change the status to the previous one")
	}

	if bid.AuthorType == "User" {
		if bid.AuthorId != userId {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	} else {
		userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return entity.Bid{}, ErrNotEnoughPermissions
			}
			logrus.Errorf("BidService.UpdateBidStatus: cannot get user organization id: %v", err)
			return entity.Bid{}, ErrNotEnoughPermissions
		}

		tenderOrganizationId, _ := s.repo.Auth.GetUserOrganizationId(ctx, bid.AuthorId)

		if userOrganizationId != tenderOrganizationId {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	}

	bid, err = s.repo.Bid.UpdateBidStatus(ctx, bidId, status)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.UpdateBidStatus: cannot update bid: %v", err)
		return entity.Bid{}, ErrCannotUpdateBid
	}

	return bid, nil
}
func (s *BidService) EditBid(ctx context.Context, bidId string, userId string, input entity.EditBidInput) (entity.Bid, error) {
	bid, err := s.repo.Bid.GetBidById(ctx, bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.EditBid: cannot get bid: %v", err)
		return entity.Bid{}, ErrCannotGetBid
	}

	if bid.Status == "Canceled" {
		return entity.Bid{}, fmt.Errorf("cannot edit canceled bid")
	}

	if bid.AuthorType == "User" {
		if bid.AuthorId != userId {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	} else {
		userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return entity.Bid{}, ErrNotEnoughPermissions
			}
			logrus.Errorf("BidService.EditBid: cannot get user organization id: %v", err)
			return entity.Bid{}, ErrNotEnoughPermissions
		}

		tenderOrganizationId, _ := s.repo.Auth.GetUserOrganizationId(ctx, bid.AuthorId)

		if userOrganizationId != tenderOrganizationId {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	}

	bid, err = s.repo.Bid.UpdateBid(ctx, bidId, input)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.EditBid: cannot update bid: %v", err)
		return entity.Bid{}, ErrCannotUpdateBid
	}

	return bid, nil
}
func (s *BidService) SubmitBidDecision(ctx context.Context, bidId, decision, userId string) (entity.Bid, error) {
	bid, err := s.repo.Bid.GetBidById(ctx, bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.SubmitBidDecision: cannot get bid: %v", err)
		return entity.Bid{}, ErrCannotGetBid
	}
	if bid.Status == "Canceled" {
		return entity.Bid{}, fmt.Errorf("cannot edit canceled bid")
	}
	if bid.Status == "Created" {
		return entity.Bid{}, ErrNotEnoughPermissions
	}

	userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
		logrus.Errorf("BidService.SubmitBidDecision: cannot get user organization id: %v", err)
		return entity.Bid{}, ErrNotEnoughPermissions
	}
	tender, _ := s.repo.Tender.GetTenderById(ctx, bid.TenderId)
	if userOrganizationId != tender.OrganizationId {
		return entity.Bid{}, ErrNotEnoughPermissions
	}

	if decision == "Approved" && tender.Status != "Closed" {
		_, err := s.repo.Tender.UpdateTenderStatus(ctx, tender.Id, "Closed")
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return entity.Bid{}, ErrTenderNotFound
			}
			logrus.Errorf("TenderService.GetTenderStatus: cannot update tender: %v", err)
			return entity.Bid{}, ErrCannotUpdateTender
		}
	}

	bid, err = s.repo.Bid.SubmitBidDecision(ctx, bidId, decision)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.SubmitBidDecision: cannot update bid: %v", err)
		return entity.Bid{}, ErrCannotUpdateBid
	}

	return bid, nil
}

func (s *BidService) RollbackBid(ctx context.Context, bidId string, version int, userId string) (entity.Bid, error) {
	bid, err := s.repo.Bid.GetBidById(ctx, bidId)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.EditBid: cannot get bid: %v", err)
		return entity.Bid{}, ErrCannotGetBid
	}
	if bid.Version <= version {
		return entity.Bid{}, fmt.Errorf("cannot roll back to the current and higher version")
	}
	if bid.Status == "Canceled" {
		return entity.Bid{}, fmt.Errorf("cannot edit canceled bid")
	}

	if bid.AuthorType == "User" {
		if bid.AuthorId != userId {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	} else {
		userOrganizationId, err := s.repo.Auth.GetUserOrganizationId(ctx, userId)
		if err != nil {
			if errors.Is(err, repoerrs.ErrNotFound) {
				return entity.Bid{}, ErrNotEnoughPermissions
			}
			logrus.Errorf("BidService.RollbackBid: cannot get user organization id: %v", err)
			return entity.Bid{}, ErrNotEnoughPermissions
		}

		tenderOrganizationId, _ := s.repo.Auth.GetUserOrganizationId(ctx, bid.AuthorId)

		if userOrganizationId != tenderOrganizationId {
			return entity.Bid{}, ErrNotEnoughPermissions
		}
	}

	bid, err = s.repo.Bid.RollbackBid(ctx, bidId, version)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return entity.Bid{}, ErrBidNotFound
		}
		logrus.Errorf("BidService.RollbackBid: cannot update bid: %v", err)
		return entity.Bid{}, ErrCannotUpdateBid
	}

	return bid, nil
}
