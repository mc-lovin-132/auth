package service

import (
	"context"

	"github.com/mc-lovin-132/auth/internal/domain"
)

func (s *Service) Refresh(ctx context.Context, token string) (*domain.AccessToken, *domain.RefreshToken, error) {
	// получаем информацию о рефреш токене
	refreshTokenObject, err := s.refreshTokenRepo.Get(ctx, token)
	if err != nil {
		return nil, nil, err
	}
	// проверяем не был ли отозван
	if refreshTokenObject.Revoked {
		// если кто то пытается отозывнный токен,
		// скорее всего присутствует угроза безопастности
		// для безопастности отзываем все токены

		// TODO: подумать что делать с ошибкой возникающей при отзыве после обнаружения угрозы
		// отзыв всех аксес токенов
		err = s.RevokeByUser(ctx, refreshTokenObject.UserID)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, domain.ErrRefreshTokenRevoked
	}
	// проверяем не был ли ранее использован
	if refreshTokenObject.Used {
		err = s.RevokeByUser(ctx, refreshTokenObject.UserID)
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, domain.ErrRefreshTokenAlreadyUsed
	}

	// проверки безопасности пройдены

	// опредделяем с какого устройства производится запрос
	deviceID, err := s.deviceFingerprinter.Detect(ctx)
	if err != nil {
		return nil, nil, err
	}

	// озываем текущие токены
	// аксес отзываем
	err = s.accessTokenDeviceRepo.UpdateVersion(ctx, refreshTokenObject.UserID, deviceID)
	if err != nil {
		return nil, nil, err
	}
	// рефреш помечаем как использованный
	err = s.refreshTokenRepo.MarkAsUsed(ctx, refreshTokenObject.Value)
	if err != nil {
		return nil, nil, err
	}

	// создаем новую сессию
	sessionID := s.sessioner.Gen(ctx)

	// получаем актуальные версии аксес токена
	actualDeviceVersion, err := s.accessTokenDeviceRepo.GetVersion(ctx, refreshTokenObject.UserID, deviceID)
	if err != nil {
		return nil, nil, err
	}
	actualGlobalVersion, err := s.accessTokenGlobalRepo.GetVersion(ctx, refreshTokenObject.UserID)
	if err != nil {
		return nil, nil, err
	}

	// создаем новую пару
	// аксес
	payload := &domain.AccessToken{
		UserID:        refreshTokenObject.UserID,
		GlobalVersion: actualGlobalVersion,
		DeviceVersion: actualDeviceVersion,
		DeviceID:      deviceID,
		SessionID:     sessionID,
	}
	accessTokenObject, err := s.accessManager.Gen(ctx, payload)
	if err != nil {
		return nil, nil, err
	}
	// рефреш
	newRefreshTokenObject := s.refreshManager.Gen(
		ctx,
		accessTokenObject.UserID,
		deviceID,
		sessionID,
	)

	// сохраняем рефреш
	err = s.refreshTokenRepo.Create(ctx, newRefreshTokenObject)
	if err != nil {
		return nil, nil, err
	}

	return accessTokenObject, newRefreshTokenObject, nil
}
