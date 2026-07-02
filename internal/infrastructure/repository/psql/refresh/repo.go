package repo

import (
	"context"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var data RefreshTokenModel
	err := r.db.GetContext(ctx, &data, getQuery, token)
	if err != nil {
		return nil, sqlErrorMapper(err)
	}
	return toDomain(&data), nil
}
func (r *Repository) GetBySessionID(ctx context.Context, sessionID string) (*domain.RefreshToken, error) {
	var data RefreshTokenModel
	err := r.db.GetContext(ctx, &data, getQuery, sessionID)
	if err != nil {
		return nil, sqlErrorMapper(err)
	}
	return toDomain(&data), nil
}
func (r *Repository) Create(ctx context.Context, token *domain.RefreshToken) error {
	data := fromDomain(token)
	_, err := r.db.NamedExecContext(ctx, createQuery, data)
	if err != nil {
		return sqlErrorMapper(err)
	}
	return nil
}

// value токена по совместительству является его id
func (r *Repository) MarkAsUsed(ctx context.Context, value string) error {
	rows, err := r.db.ExecContext(ctx, markAsUsedQuery, value)
	if err != nil {
		return sqlErrorMapper(err)
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return sqlErrorMapper(err)
	}
	if rowsAffected == 0 {
		return sqlErrorMapper(err)
	}
	return nil
}
func (r *Repository) MarkAsRevokedByDevice(ctx context.Context, userID int, deviceID string) error {
	_, err := r.db.ExecContext(ctx, markAsRevokedByDeviceQuery, userID, deviceID)
	if err != nil {
		return sqlErrorMapper(err)
	}
	// ПРИМЕЧАНИЕ: ситуации когда ни одной строки не будет изменено могут возникать
	// у пользователя может просто не быть ни одного токена
	// rowsAffected, err := rows.RowsAffected()
	// if err != nil {
	// 	return sqlErrorMapper(err)
	// }
	// if rowsAffected == 0 {
	// 	return sqlErrorMapper(err)
	// }
	return nil
}
func (r *Repository) MarkAsRevokedByUser(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, markAsRevokedByUserQuery, userID)
	if err != nil {
		return sqlErrorMapper(err)
	}
	// ПРИМЕЧАНИЕ: ситуации когда ни одной строки не будет изменено могут возникать
	// rowsAffected, err := rows.RowsAffected()
	// if err != nil {
	// 	return sqlErrorMapper(err)
	// }
	// if rowsAffected == 0 {
	// 	return sqlErrorMapper(err)
	// }
	return nil
}
func (r *Repository) MarkAsRevokedByConcrete(ctx context.Context, value string) error {
	rows, err := r.db.ExecContext(ctx, markAsRevokedByConcreteQuery, value)
	if err != nil {
		return sqlErrorMapper(err)
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return sqlErrorMapper(err)
	}
	if rowsAffected == 0 {
		return sqlErrorMapper(err)
	}
	return nil
}
