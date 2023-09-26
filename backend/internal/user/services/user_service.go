package services

import "backend/internal/user/repositories"

// Account servisi yapısı tanımıdır.
// It is a Account Service struct.
type UserService struct {
	repository_user repositories.UserRepository
}

// Yeni Account Servisi yaratır. Hata varsa sonuç olarak yansıtabilir.
// You can create Account Service. If there are any errors, can return its.
func NewUserService(repository repositories.UserRepository) (*UserService, error) {

	// Eğer eklenecek başka altyapılar varsa buraya eklenebilir.
	// If you want to add another infrastructure, you can add here.
	return &UserService{
		repository_user: repository,
	}, nil
}

func (u *UserService) GetUserWithToken(token string) {

}
