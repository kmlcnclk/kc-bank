package query

import (
	"context"
	"errors"
	"kc-bank/app/repository"
	"kc-bank/domain"
)

type IUserQueryService interface {
	GetUser(ctx context.Context, Id string) (*domain.User, error)
}

type userQueryService struct {
	userRepository repository.IUserRepository
}

func NewUserQueryService(userRepository repository.IUserRepository) IUserQueryService {
	return &userQueryService{
		userRepository: userRepository,
	}
}

func (u *userQueryService) GetUser(ctx context.Context, Id string) (*domain.User, error) {
	user, err := u.userRepository.GetUser(ctx, Id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
