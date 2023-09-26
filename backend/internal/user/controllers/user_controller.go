package controllers

import "backend/internal/user/services"

type UserController struct {
	service_user    *services.UserService
	service_session *services.SessionService
	service_account *services.AccountService
}

func NewUserController(user *services.UserService, session *services.SessionService, account *services.AccountService) (*UserController, error) {
	return &UserController{
		service_user:    user,
		service_session: session,
		service_account: account,
	}, nil
}
