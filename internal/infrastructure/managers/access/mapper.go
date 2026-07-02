package managers

import (
	"errors"
	"fmt"

	"github.com/mc-lovin-132/auth/internal/domain"

	"github.com/golang-jwt/jwt/v4"
)

func accessManagerErrorMapper(err error) error {
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return fmt.Errorf("%w: %w", domain.ErrInvalidAccessToken, jwt.ErrTokenMalformed)
	} else if errors.Is(err, jwt.ErrTokenUnverifiable) {
		return fmt.Errorf("%w: %w", domain.ErrInvalidAccessToken, jwt.ErrTokenUnverifiable)
	} else if errors.Is(err, jwt.ErrSignatureInvalid) {
		return fmt.Errorf("%w: %w", domain.ErrInvalidSignature, jwt.ErrSignatureInvalid)
	} else if errors.Is(err, domain.ErrAccessTokenExpired) {
		return domain.ErrAccessTokenExpired
	}
	return fmt.Errorf("%w: %w", domain.ErrInternal, err)
}
