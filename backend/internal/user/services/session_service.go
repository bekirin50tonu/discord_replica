package services

import "backend/internal/user/repositories"

type SessionService struct {
	repository_session *repositories.SessionRepository
}

func NewSessionService(repo *repositories.SessionRepository) (*SessionService, error) {

	return &SessionService{
		repository_session: repo,
	}, nil
}
