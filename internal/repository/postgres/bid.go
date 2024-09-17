package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"tender-service/internal/entity"
	"tender-service/internal/repository/repoerrs"
)

type BidRepository struct {
	db *Postgres
}

func NewBidRepository(db *Postgres) *BidRepository {
	return &BidRepository{db: db}
}

func (r *BidRepository) CreateBid(ctx context.Context, bid entity.Bid) (entity.Bid, error) {
	query, args, err := r.db.Builder.
		Insert("bids").
		Columns("name", "description", "status", "tender_id", "author_type", "author_id", "version").
		Values(
			bid.Name,
			bid.Description,
			bid.Status,
			bid.TenderId,
			bid.AuthorType,
			bid.AuthorId,
			bid.Version,
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.CreateBid  - r.db.Builder.ToSql: %v", err)
	}

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&bid.Id,
		&bid.CreatedAt,
	)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.CreateBid  - r.db.Pool.QueryRow: %v", err)
	}

	return bid, nil
}

func (r *BidRepository) GetBidsByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Bid, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bids").
		Where("author_id = ?", userId).
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("BidRepository.GetBidsByUserId  - r.db.Builder.ToSql: %v", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("BidRepository.GetBidsByUserId - r.db.Pool.Query: %v", err)
	}
	defer rows.Close()

	var bids []entity.Bid
	for rows.Next() {
		var bid entity.Bid
		err := rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderId,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
			&bid.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("BidRepository.GetBidsByUserId - rows.Scan: %v", err)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *BidRepository) GetBidsForTender(ctx context.Context, tenderId string, limit int, offset int) ([]entity.Bid, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bids").
		Where("tender_id = ?", tenderId).
		Where("status = ANY(?)", []string{"Published", "Canceled"}).
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("BidRepository.GetBidsForTender  - r.db.Builder.ToSql: %v", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("BidRepository.GetBidsForTender - r.db.Pool.Query: %v", err)
	}
	defer rows.Close()

	var bids []entity.Bid
	for rows.Next() {
		var bid entity.Bid
		err := rows.Scan(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderId,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
			&bid.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("BidRepository.GetBidsForTender - rows.Scan: %v", err)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *BidRepository) GetBidById(ctx context.Context, id string) (entity.Bid, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bids").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.GetBidById  - r.db.Builder.ToSql: %v", err)
	}

	var bid entity.Bid
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("BidRepository.GetBidById  - r.db.Pool.QueryRow: %v", err)
	}

	return bid, nil
}

func (r *BidRepository) UpdateBidStatus(ctx context.Context, bidId string, status string) (entity.Bid, error) {
	query, args, err := r.db.Builder.
		Update("bids").
		Set("status", status).
		Where("id = ?", bidId).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBidStatus  - r.db.Builder.ToSql: %v", err)
	}

	var bid entity.Bid
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBidStatus  - r.db.Pool.QueryRow: %v", err)
	}

	return bid, nil
}

func (r *BidRepository) UpdateBid(ctx context.Context, bidId string, input entity.EditBidInput) (entity.Bid, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	getBid, args, err := r.db.Builder.
		Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bids").
		Where("id = ?", bidId).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid  - r.db.Builder.ToSql: %v", err)
	}

	var bid entity.Bid
	err = tx.QueryRow(ctx, getBid, args...).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid  - tx.QueryRow: %v", err)
	}

	createOldVersionBid, args, err := r.db.Builder.
		Insert("bids_old_version").
		Columns("bid_id", "name", "description", "status", "tender_id", "author_type", "author_id", "version").
		Values(
			&bid.Id,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderId,
			&bid.AuthorType,
			&bid.AuthorId,
			&bid.Version,
		).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, createOldVersionBid, args...)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid - tx.Exec: %v", err)
	}

	updateBidQuery := r.db.Builder.
		Update("bids").
		Where("id = ?", bidId)
	bid.Version += 1
	updateBidQuery = updateBidQuery.Set("version", bid.Version)
	if input.Name != nil {
		bid.Name = *input.Name
		updateBidQuery = updateBidQuery.Set("name", bid.Name)
	}
	if input.Description != nil {
		bid.Description = *input.Description
		updateBidQuery = updateBidQuery.Set("description", bid.Description)
	}
	updateBid, args, err := updateBidQuery.ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, updateBid, args...)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.UpdateBid - tx.Commit: %v", err)
	}

	return bid, nil
}

func (r *BidRepository) RollbackBid(ctx context.Context, bidId string, version int) (entity.Bid, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	getCurrentBid, args, err := r.db.Builder.
		Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version").
		From("bids").
		Where("id = ?", bidId).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid  - r.db.Builder.ToSql: %v", err)
	}

	var currentBid entity.Bid
	err = tx.QueryRow(ctx, getCurrentBid, args...).Scan(
		&currentBid.Id,
		&currentBid.Name,
		&currentBid.Description,
		&currentBid.Status,
		&currentBid.TenderId,
		&currentBid.AuthorType,
		&currentBid.AuthorId,
		&currentBid, version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid  - tx.QueryRow: %v", err)
	}

	createOldVersionBid, args, err := r.db.Builder.
		Insert("bids_old_version").
		Columns("bid_id", "name", "description", "status", "tender_id", "author_type", "author_id", "version").
		Values(
			currentBid.Id,
			currentBid.Name,
			currentBid.Description,
			currentBid.Status,
			currentBid.TenderId,
			currentBid.AuthorType,
			currentBid.AuthorId,
			currentBid.Version,
		).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, createOldVersionBid, args...)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid - tx.Exec: %v", err)
	}

	getBidReqVersion, args, err := r.db.Builder.
		Select("bid_id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bids_old_version").
		Where("bid_id = ?", bidId).
		Where("version = ?", version).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid  - r.db.Builder.ToSql: %v", err)
	}

	var bidReqVersion entity.Bid
	err = tx.QueryRow(ctx, getBidReqVersion, args...).Scan(
		&bidReqVersion.Id,
		&bidReqVersion.Name,
		&bidReqVersion.Description,
		&bidReqVersion.Status,
		&bidReqVersion.TenderId,
		&bidReqVersion.AuthorType,
		&bidReqVersion.AuthorId,
		&bidReqVersion.Version,
		&bidReqVersion.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackBid  - tx.QueryRow: %v", err)
	}

	currentBid.Version += 1
	currentBid.Name = bidReqVersion.Name
	currentBid.Description = bidReqVersion.Description

	rollbackBidReqVersion, args, err := r.db.Builder.
		Update("bids").
		Where("id = ?", bidId).
		Set("version", currentBid.Version).
		Set("name", currentBid.Name).
		Set("description", currentBid.Description).
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackTender  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, rollbackBidReqVersion, args...)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackTender - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.RollbackTender - tx.Commit: %v", err)
	}

	return currentBid, nil
}

func (r *BidRepository) SubmitBidDecision(ctx context.Context, bidId string, decision string) (entity.Bid, error) {
	query, args, err := r.db.Builder.
		Update("bids").
		Where("id = ?", bidId).
		Set("decision", decision).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	if err != nil {
		return entity.Bid{}, fmt.Errorf("BidRepository.SubmitBidDecision  - r.db.Builder.ToSql: %v", err)
	}

	var bid entity.Bid
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&bid.Id,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderId,
		&bid.AuthorType,
		&bid.AuthorId,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bid{}, repoerrs.ErrNotFound
		}
		return entity.Bid{}, fmt.Errorf("BidRepository.SubmitBidDecision  - r.db.Pool.QueryRow: %v", err)
	}

	return bid, nil
}
