package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"tender-service/internal/repository/repoerrs"
)

type AuthRepository struct {
	db *Postgres
}

func NewAuthRepository(db *Postgres) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserId(ctx context.Context, username string) (string, error) {
	query, args, err := r.db.Builder.
		Select("id").
		From("employee").
		Where("username = ?", username).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("AuthRepository.GetUserId  - r.db.Builder.ToSql: %v", err)
	}

	var id string
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrs.ErrNotFound
		}
		return "", fmt.Errorf("AuthRepository.GetUserId  - r.db.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *AuthRepository) GetUserOrganizationId(ctx context.Context, userId string) (string, error) {
	query, args, err := r.db.Builder.
		Select("organization_id").
		From("organization_responsible").
		Where("user_id = ?", userId).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("AuthRepository.GetUserOrganizationId  - r.db.Builder.ToSql: %v", err)
	}

	var id string
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repoerrs.ErrNotFound
		}
		return "", fmt.Errorf("AuthRepository.GetUserOrganizationId  - r.db.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *AuthRepository) OrganizationIsExist(ctx context.Context, organizationId string) (bool, error) {
	query, args, err := r.db.Builder.
		Select("id").
		From("organization").
		Where("id = ?", organizationId).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("AuthRepository.OrganizationIsExist  - r.db.Builder.ToSql: %v", err)
	}

	var id string
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, repoerrs.ErrNotFound
		}
		return false, fmt.Errorf("AuthRepository.OrganizationIsExist  - r.db.Pool.QueryRow: %v", err)
	}

	return true, nil
}

func (r *AuthRepository) UserIsExist(ctx context.Context, userId string) (bool, error) {
	query, args, err := r.db.Builder.
		Select("id").
		From("employee").
		Where("id = ?", userId).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("AuthRepository.UserIsExist  - r.db.Builder.ToSql: %v", err)
	}

	var id string
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, repoerrs.ErrNotFound
		}
		return false, fmt.Errorf("AuthRepository.UserIsExist  - r.db.Pool.QueryRow: %v", err)
	}

	return true, nil
}
