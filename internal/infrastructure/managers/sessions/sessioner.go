package service

import (
	"context"

	"github.com/google/uuid"
)

type SessionManager struct {
}

func New() *SessionManager {
	return &SessionManager{}
}

func (s *SessionManager) Gen(ctx context.Context) string {
	return uuid.New().String()
}
