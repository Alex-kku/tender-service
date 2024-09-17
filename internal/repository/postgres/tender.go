package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"tender-service/internal/entity"
	"tender-service/internal/repository/repoerrs"
)

type TenderRepository struct {
	db *Postgres
}

func NewTenderRepository(db *Postgres) *TenderRepository {
	return &TenderRepository{db: db}
}

func (r *TenderRepository) GetTenders(ctx context.Context, limit int, offset int, serviceType []string) ([]entity.Tender, error) {
	getTenders := r.db.Builder.
		Select("id", "name", "description", "service_type", "status", "organization_id", "version", "created_at").
		From("tenders").
		Where("status = ?", "Published").
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset))
	if len(serviceType) != 0 {
		getTenders = getTenders.Where("service_type = ANY(?)", serviceType)
	}

	query, args, err := getTenders.ToSql()
	if err != nil {
		return nil, fmt.Errorf("TenderRepository.GetTenders  - getTenders.ToSql: %v", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("TenderRepository.GetTenders - r.db.Pool.Query: %v", err)
	}
	defer rows.Close()

	var tenders []entity.Tender
	for rows.Next() {
		var tender entity.Tender
		err := rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationId,
			&tender.Version,
			&tender.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("TenderRepository.GetTenders - rows.Scan: %v", err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (r *TenderRepository) CreateTender(ctx context.Context, userId string, tender entity.Tender) (entity.Tender, error) {
	query, args, err := r.db.Builder.
		Insert("tenders").
		Columns("name", "description", "service_type", "status", "organization_id", "version", "creator_id").
		Values(
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.Status,
			tender.OrganizationId,
			tender.Version,
			userId,
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.CreateTender  - r.db.Builder.ToSql: %v", err)
	}

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&tender.Id,
		&tender.CreatedAt,
	)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.CreateTender  - r.db.Pool.QueryRow: %v", err)
	}

	return tender, nil
}

func (r *TenderRepository) GetTendersByUserId(ctx context.Context, limit int, offset int, userId string) ([]entity.Tender, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "service_type", "status", "organization_id", "version", "created_at").
		From("tenders").
		Where("creator_id = ?", userId).
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("TenderRepository.GetTendersByUserId  - r.db.Builder.ToSql: %v", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("TenderRepository.GetTendersByUserId - r.db.Pool.Query: %v", err)
	}
	defer rows.Close()

	var tenders []entity.Tender
	for rows.Next() {
		var tender entity.Tender
		err := rows.Scan(
			&tender.Id,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationId,
			&tender.Version,
			&tender.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("TenderRepository.GetTendersByUserId - rows.Scan: %v", err)
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (r *TenderRepository) GetTenderById(ctx context.Context, id string) (entity.Tender, error) {
	query, args, err := r.db.Builder.
		Select("id", "name", "description", "service_type", "status", "organization_id", "version", "created_at").
		From("tenders").
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.GetTenderById  - r.db.Builder.ToSql: %v", err)
	}

	var tender entity.Tender
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		return entity.Tender{}, fmt.Errorf("TenderRepository.GetTenderById  - r.db.Pool.QueryRow: %v", err)
	}

	return tender, nil
}

func (r *TenderRepository) UpdateTenderStatus(ctx context.Context, tenderId string, status string) (entity.Tender, error) {
	query, args, err := r.db.Builder.
		Update("tenders").
		Set("status", status).
		Where("id = ?", tenderId).
		Suffix("RETURNING id, name, description,  service_type, status, organization_id, version, created_at").
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTenderStatus  - r.db.Builder.ToSql: %v", err)
	}

	var tender entity.Tender
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTenderStatus  - r.db.Pool.QueryRow: %v", err)
	}

	return tender, nil
}

func (r *TenderRepository) UpdateTender(ctx context.Context, tenderId string, input entity.EditTenderInput) (entity.Tender, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	getTender, args, err := r.db.Builder.
		Select("id", "name", "description", "service_type", "status", "organization_id", "version", "created_at", "creator_id").
		From("tenders").
		Where("id = ?", tenderId).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender  - r.db.Builder.ToSql: %v", err)
	}

	var tender entity.Tender
	var creatorId string
	err = tx.QueryRow(ctx, getTender, args...).Scan(
		&tender.Id,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.Version,
		&tender.CreatedAt,
		&creatorId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender  - tx.QueryRow: %v", err)
	}

	createOldVersionTender, args, err := r.db.Builder.
		Insert("tenders_old_version").
		Columns("tender_id", "name", "description", "service_type", "status", "organization_id", "version", "creator_id").
		Values(
			tender.Id,
			tender.Name,
			tender.Description,
			tender.ServiceType,
			tender.Status,
			tender.OrganizationId,
			tender.Version,
			creatorId,
		).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, createOldVersionTender, args...)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender - tx.Exec: %v", err)
	}

	updateTenderQuery := r.db.Builder.
		Update("tenders").
		Where("id = ?", tenderId)
	tender.Version += 1
	updateTenderQuery = updateTenderQuery.Set("version", tender.Version)
	if input.Name != nil {
		tender.Name = *input.Name
		updateTenderQuery = updateTenderQuery.Set("name", tender.Name)
	}
	if input.Description != nil {
		tender.Description = *input.Description
		updateTenderQuery = updateTenderQuery.Set("description", tender.Description)
	}
	if input.ServiceType != nil {
		tender.ServiceType = *input.ServiceType
		updateTenderQuery = updateTenderQuery.Set("service_type", tender.ServiceType)
	}
	updateTender, args, err := updateTenderQuery.ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTenderStatus  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, updateTender, args...)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.UpdateTender - tx.Commit: %v", err)
	}

	return tender, nil
}

func (r *TenderRepository) RollbackTender(ctx context.Context, tenderId string, version int) (entity.Tender, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	getCurrentTender, args, err := r.db.Builder.
		Select("id", "name", "description", "service_type", "status", "organization_id", "version", "creator_id").
		From("tenders").
		Where("id = ?", tenderId).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender  - r.db.Builder.ToSql: %v", err)
	}

	var currentTender entity.Tender
	var creatorId string
	err = tx.QueryRow(ctx, getCurrentTender, args...).Scan(
		&currentTender.Id,
		&currentTender.Name,
		&currentTender.Description,
		&currentTender.ServiceType,
		&currentTender.Status,
		&currentTender.OrganizationId,
		&currentTender.Version,
		&creatorId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender  - tx.QueryRow: %v", err)
	}

	createOldVersionTender, args, err := r.db.Builder.
		Insert("tenders_old_version").
		Columns("tender_id", "name", "description", "service_type", "status", "organization_id", "version", "creator_id").
		Values(
			currentTender.Id,
			currentTender.Name,
			currentTender.Description,
			currentTender.ServiceType,
			currentTender.Status,
			currentTender.OrganizationId,
			currentTender.Version,
			creatorId,
		).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, createOldVersionTender, args...)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender - tx.Exec: %v", err)
	}

	getTenderReqVersion, args, err := r.db.Builder.
		Select("id", "name", "description", "service_type", "status", "organization_id", "version", "created_at", "creator_id").
		From("tenders_old_version").
		Where("tender_id = ?", tenderId).
		Where("version = ?", version).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender  - r.db.Builder.ToSql: %v", err)
	}

	var tenderReqVersion entity.Tender
	err = tx.QueryRow(ctx, getTenderReqVersion, args...).Scan(
		&tenderReqVersion.Id,
		&tenderReqVersion.Name,
		&tenderReqVersion.Description,
		&tenderReqVersion.ServiceType,
		&tenderReqVersion.Status,
		&tenderReqVersion.OrganizationId,
		&tenderReqVersion.Version,
		&tenderReqVersion.CreatedAt,
		&creatorId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Tender{}, repoerrs.ErrNotFound
		}
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender  - tx.QueryRow: %v", err)
	}

	currentTender.Version += 1
	currentTender.Name = tenderReqVersion.Name
	currentTender.Description = tenderReqVersion.Description
	currentTender.ServiceType = tenderReqVersion.ServiceType

	rollbackTenderReqVersion, args, err := r.db.Builder.
		Update("tenders").
		Where("id = ?", tenderId).
		Set("version", currentTender.Version).
		Set("name", currentTender.Name).
		Set("description", currentTender.Description).
		Set("service_type", currentTender.ServiceType).
		ToSql()

	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender  - r.db.Builder.ToSql: %v", err)
	}

	_, err = tx.Exec(ctx, rollbackTenderReqVersion, args...)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Tender{}, fmt.Errorf("TenderRepository.RollbackTender - tx.Commit: %v", err)
	}

	return currentTender, nil
}
