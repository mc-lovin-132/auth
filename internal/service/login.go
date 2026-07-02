package service

import (
	"context"

	"github.com/mc-lovin-132/auth/internal/domain"
)

func (s *Service) Login(ctx context.Context, email, password string) (*domain.AccessToken, *domain.RefreshToken, error) {
	// получаем пользователя и проверяем пароль
	user, err := s.userClient.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user.Password != password {
		return nil, nil, domain.ErrInvalidPassword
	}

	// определяем с какого устройства был сделан запрос
	// TODO: должны передаваться данные устройства
	deviceID, err := s.deviceFingerprinter.Detect(ctx)
	if err != nil {
		return nil, nil, err
	}

	// отзываем все аксес токены пользователя на устройстве
	// примечание:
	// есть ли в этом действии необходимость?
	// может ли пользователь использовать логин до момента истечения сессии?
	// с точки зрения клиентской логики - не может, но вообще сервер этого не запрещает, поэтому отзыв необходим
	err = s.accessTokenDeviceRepo.UpdateVersion(ctx, user.ID, deviceID)
	if err != nil {
		return nil, nil, err
	}

	// отзываем все ревреш токены пользователя на устройстве
	err = s.refreshTokenRepo.MarkAsRevokedByDevice(ctx, user.ID, deviceID)
	if err != nil {
		return nil, nil, err
	}

	// получаем актуальные версии аксес токена
	actualDeviceVersion, err := s.accessTokenDeviceRepo.GetVersion(ctx, user.ID, deviceID)
	if err != nil {
		return nil, nil, err
	}
	actualGlobalVersion, err := s.accessTokenGlobalRepo.GetVersion(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	// генерируем сессию
	sessionID := s.sessioner.Gen(ctx)

	// создаем аксес токен со всем пейлоадом
	payload := &domain.AccessToken{
		UserID:        user.ID,
		GlobalVersion: actualGlobalVersion,
		DeviceVersion: actualDeviceVersion,
		DeviceID:      deviceID,
		SessionID:     sessionID,
	}
	accessTokenObject, err := s.accessManager.Gen(ctx, payload)
	if err != nil {
		return nil, nil, err
	}

	// создаем рефреш токен со всем пейлоадом
	refreshTokenObject := s.refreshManager.Gen(ctx, user.ID, deviceID, sessionID)

	// сохраняем рефреш в бд
	err = s.refreshTokenRepo.Create(ctx, refreshTokenObject)
	if err != nil {
		return nil, nil, err
	}

	return accessTokenObject, refreshTokenObject, nil
}
