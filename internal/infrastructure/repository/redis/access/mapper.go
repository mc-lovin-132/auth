package repo

import (
	"fmt"

	"github.com/mc-lovin-132/auth/internal/domain"
)

const intZeroValue = 0
const stringZeroValue = ""

func redisErrorMappeer(err error) error {
	return fmt.Errorf("%w: %w", domain.ErrInternal, err)
}
