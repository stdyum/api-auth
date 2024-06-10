package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stdyum/api-auth/internal/app/entities"
)

func (r *repository) CreateUser(ctx context.Context, user entities.User) error {
	_, err := r.database.ExecContext(ctx, `
	INSERT INTO users (id, email, verified_email, login, password, picture) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.Email, user.VerifiedEmail, user.Login, user.Password, user.Picture)
	return err
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (entities.User, error) {
	row := r.database.QueryRowContext(ctx, `
	SELECT id, email, verified_email, login, password, picture FROM users
		WHERE id = $1 
	`, id)

	return r.rowToUser(row)
}

func (r *repository) GetUserByLogin(ctx context.Context, login string) (entities.User, error) {
	row := r.database.QueryRowContext(ctx, `
	SELECT id, email, verified_email, login, password, picture FROM users
		WHERE login = $1 
	`, login)

	return r.rowToUser(row)
}

func (r *repository) GetUserByLoginAndEmail(ctx context.Context, login string, email string) (entities.User, error) {
	row := r.database.QueryRowContext(ctx, `
	SELECT id, email, verified_email, login, password, picture FROM users
		WHERE login = $1 AND email = $2 
	`, login, email)

	return r.rowToUser(row)
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	row := r.database.QueryRowContext(ctx, `
	SELECT id, email, verified_email, login, password, picture FROM users
		WHERE email = $1 
	`, email)

	return r.rowToUser(row)
}

func (r *repository) SetEmailConfirmed(ctx context.Context, userId uuid.UUID) error {
	res, err := r.database.ExecContext(ctx, "UPDATE users SET verified_email = true WHERE id = $1", userId)
	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) SetPassword(ctx context.Context, userId uuid.UUID, password string) error {
	_, err := r.database.ExecContext(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, userId)
	return err
}

func (r *repository) rowToUser(row *sql.Row) (user entities.User, err error) {
	err = row.Scan(&user.ID, &user.Email, &user.VerifiedEmail, &user.Login, &user.Password, &user.Picture)
	return
}
