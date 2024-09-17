package service

import (
	"context"
	"tender-service/internal/entity"
	"tender-service/internal/repository"
)

type Auth interface {
	GetUserId(ctx context.Context, username string) (string, error)
	CheckResponsibility(ctx context.Context, userId string, organizationId string) (bool, error)
}

type Tender interface {
	GetTenders(ctx context.Context, limit int, offset int, serviceType []string) ([]entity.Tender, error)
	CreateTender(ctx context.Context, userId string, input entity.CreateTenderInput) (entity.Tender, error)
	GetTendersByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId string, userId *string) (string, error)
	UpdateTenderStatus(ctx context.Context, tenderId, status, userId string) (entity.Tender, error)
	EditTender(ctx context.Context, tenderId string, userId string, input entity.EditTenderInput) (entity.Tender, error)
	RollbackTender(ctx context.Context, tenderId string, version int, userId string) (entity.Tender, error)
}

type Bid interface {
	CreateBid(ctx context.Context, input entity.CreateBidInput) (entity.Bid, error)
	GetBidsByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Bid, error)
	GetBidsForTender(ctx context.Context, tenderId string, userId string, limit int, offset int) ([]entity.Bid, error)
	GetBidStatus(ctx context.Context, bidId string, userId string) (string, error)
	UpdateBidStatus(ctx context.Context, bidId, status, userId string) (entity.Bid, error)
	EditBid(ctx context.Context, bidId string, userId string, input entity.EditBidInput) (entity.Bid, error)
	SubmitBidDecision(ctx context.Context, bidId, decision, userId string) (entity.Bid, error)
	RollbackBid(ctx context.Context, bidId string, version int, userId string) (entity.Bid, error)
}

type Service struct {
	Auth   Auth
	Tender Tender
	Bid    Bid
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:   NewAuthService(repos),
		Tender: NewTenderService(repos),
		Bid:    NewBidService(repos),
	}
}
