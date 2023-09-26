package integrations

import "backend/internal/user/services"

type UserIntegrations struct {
	service_user *services.UserService
}

func NewIntegrations(user *services.UserService) (*UserIntegrations, error) {

	return &UserIntegrations{
		service_user: user,
	}, nil
}
