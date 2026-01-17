package postgres

import (
	"context"

	"github.com/ansarctica/domashka4/internal/models"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(ctx, query, user.Email, user.Password).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var u models.User

	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE id = $1`

	var u models.User

	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
