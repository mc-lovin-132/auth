package service

import "context"

func (s *Service) RevokeByUser(ctx context.Context, userID int) error {
	err := s.accessTokenGlobalRepo.UpdateVersion(ctx, userID)
	if err != nil {
		return err
	}
	err = s.refreshTokenRepo.MarkAsRevokedByUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) RevokeByDevice(ctx context.Context, userID int, deviceID string) error {
	err := s.accessTokenDeviceRepo.UpdateVersion(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	err = s.refreshTokenRepo.MarkAsRevokedByDevice(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) RevokeConcrete(ctx context.Context, token string) error {
	accessTokenObject, err := s.accessManager.Parse(ctx, token)
	if err != nil {
		return err
	}
	// примечание: нет смысла делать лишние запросы в редис
	// для проверки отозван или использован токен или нет
	// в любом случае мы помечааем его как отозванный

	// получаем по сессии рефреш
	sessionID := accessTokenObject.SessionID
	refreshTokenObject, err := s.refreshTokenRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	// отзываем аксес
	// TODO: есть ли необходимость в хранении значения токена или его ID достаточно?
	err = s.accessTokenBlackList.Add(ctx, accessTokenObject.ID, accessTokenObject.Value)

	// отзываем рефреш
	err = s.refreshTokenRepo.MarkAsRevokedByConcrete(ctx, refreshTokenObject.Value)
	if err != nil {
		return err
	}

	return nil
}
