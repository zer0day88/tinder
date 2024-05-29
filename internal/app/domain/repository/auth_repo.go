package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zer0day88/tinder/internal/app/domain/entities"
	"time"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Insert(ctx context.Context, auth *entities.Auth) error {

	q := `INSERT INTO auths (id, email, enc_password) VALUES ($1,$2,$3)`

	_, err := r.db.Exec(ctx, q, auth.ID, auth.Email, auth.EncPassword)
	if err != nil {
		return fmt.Errorf("unable to insert row: %v", err)
	}

	return nil
}

func (r *AuthRepository) FindOneByEmail(ctx context.Context, email string) (*entities.Auth, error) {

	q := `select id, email, enc_password from auths where email = $1`

	rows := r.db.QueryRow(ctx, q, email)

	auth := entities.Auth{}
	err := rows.Scan(&auth.ID, &auth.Email, &auth.EncPassword)

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *AuthRepository) UpdateLastSigning(ctx context.Context, id string) error {

	q := `update auths set last_sign_in_at = $1 where id = $2`

	_, err := r.db.Exec(ctx, q, time.Now(), id)
	if err != nil {
		return fmt.Errorf("unable to update row: %v", err)
	}

	return nil
}
