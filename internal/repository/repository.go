package repository

import (
	"context"
	"tender-service/internal/entity"
	"tender-service/internal/repository/postgres"
)

type Auth interface {
	GetUserId(ctx context.Context, username string) (string, error)
	GetUserOrganizationId(ctx context.Context, userId string) (string, error)
	OrganizationIsExist(ctx context.Context, organizationId string) (bool, error)
	UserIsExist(ctx context.Context, userId string) (bool, error)
}

type Tender interface {
	GetTenders(ctx context.Context, limit int, offset int, serviceType []string) ([]entity.Tender, error)
	CreateTender(ctx context.Context, userId string, tender entity.Tender) (entity.Tender, error)
	GetTendersByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Tender, error)
	GetTenderById(ctx context.Context, id string) (entity.Tender, error)
	UpdateTenderStatus(ctx context.Context, tenderId string, status string) (entity.Tender, error)
	UpdateTender(ctx context.Context, tenderId string, input entity.EditTenderInput) (entity.Tender, error)
	RollbackTender(ctx context.Context, tenderId string, version int) (entity.Tender, error)
}

type Bid interface {
	CreateBid(ctx context.Context, bid entity.Bid) (entity.Bid, error)
	GetBidsByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Bid, error)
	GetBidsForTender(ctx context.Context, tenderId string, limit int, offset int) ([]entity.Bid, error)
	GetBidById(ctx context.Context, id string) (entity.Bid, error)
	UpdateBidStatus(ctx context.Context, bidId string, status string) (entity.Bid, error)
	UpdateBid(ctx context.Context, bidId string, input entity.EditBidInput) (entity.Bid, error)
	RollbackBid(ctx context.Context, bidId string, version int) (entity.Bid, error)
	SubmitBidDecision(ctx context.Context, bidId string, decision string) (entity.Bid, error)
}

type Repository struct {
	Auth   Auth
	Tender Tender
	Bid    Bid
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		Auth:   postgres.NewAuthRepository(db),
		Tender: postgres.NewTenderRepository(db),
		Bid:    postgres.NewBidRepository(db),
	}
}
