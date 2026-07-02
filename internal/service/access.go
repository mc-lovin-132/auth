package service

import (
	"context"

	"github.com/mc-lovin-132/auth/internal/domain"
)

func (s *Service) Access(ctx context.Context, token string) (*domain.AccessToken, error) {
	// парс проверяет не истек ли и валидна ли подпись
	accessTokenObject, err := s.accessManager.Parse(ctx, token)
	if err != nil {
		return nil, err
	}
	// получаем актульную версию токена на девайсе
	actualDeviceVersion, err := s.accessTokenDeviceRepo.GetVersion(ctx, accessTokenObject.UserID, accessTokenObject.DeviceID)
	if err != nil {
		return nil, err
	}
	// проверяем совпадает ли актуальная версия с версией токена
	if accessTokenObject.DeviceVersion != actualDeviceVersion {
		err = s.RevokeByUser(ctx, accessTokenObject.UserID)
		if err != nil {
			return nil, err
		}
		return nil, domain.ErrRefreshTokenRevoked
	}
	// получаем актуальную глобальную версию
	actualGlobalVersion, err := s.accessTokenGlobalRepo.GetVersion(ctx, accessTokenObject.UserID)
	if err != nil {
		return nil, err
	}
	// проверяем совпадение версий
	if accessTokenObject.GlobalVersion != actualGlobalVersion {
		err = s.RevokeByUser(ctx, accessTokenObject.UserID)
		if err != nil {
			return nil, err
		}
		return nil, domain.ErrRefreshTokenRevoked
	}
	// проверяем не отозван ли конкретно этот токен
	exists, err := s.accessTokenBlackList.Exists(ctx, accessTokenObject.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		err = s.RevokeByUser(ctx, accessTokenObject.UserID)
		if err != nil {
			return nil, err
		}
		return nil, domain.ErrAccessTokenRevoked
	}

	// проверяем не использовался ли раньше токен
	exists, err = s.accessTokenUsedRepo.Exists(ctx, accessTokenObject.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		err = s.RevokeByUser(ctx, accessTokenObject.UserID)
		if err != nil {
			return nil, err
		}
		return nil, domain.ErrAccessTokenAlreadyUsed
	}

	// помечаем токен как использованный
	err = s.accessTokenUsedRepo.Add(ctx, accessTokenObject.ID, accessTokenObject.Value)
	if err != nil {
		return nil, err
	}

	return accessTokenObject, nil
}
