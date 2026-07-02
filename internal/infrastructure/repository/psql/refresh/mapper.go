package repo

import (
	"database/sql"
	"fmt"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/lib/pq"
)

func sqlErrorMapper(err error) error {
	fmt.Println(err.Error())
	if pqErr, ok := err.(*pq.Error); ok {
		// нарушение уникальности
		if pqErr.Code == "23505" {
			return fmt.Errorf("%w: %w", domain.ErrRefreshTokenNotUnique, err)
		}
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("%w: %w", domain.ErrRefreshTokenNotFound, err)
	}
	return fmt.Errorf("%w: %w", domain.ErrInternal, err)
}
