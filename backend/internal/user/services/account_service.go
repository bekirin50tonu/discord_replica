package services

import "backend/internal/user/repositories"

// Account servisi yapısı tanımıdır.
// It is a Account Service struct.
type AccountService struct {
	repository_account repositories.AccountRepository
}

// Yeni Account Servisi yaratır. Hata varsa sonuç olarak yansıtabilir.
// You can create Account Service. If there are any errors, can return its.
func NewAccountService(repository repositories.AccountRepository) (*AccountService, error) {

	// Eğer eklenecek başka altyapılar varsa buraya eklenebilir.
	// If you want to add another infrastructure, you can add here.
	return &AccountService{
		repository_account: repository,
	}, nil
}

func ConnectWithUser() {

}
