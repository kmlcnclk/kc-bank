package query

import (
	"context"
	"errors"
	"kc-bank/app/repository"
	"kc-bank/domain"
)

type IAccountQueryService interface {
	GetAccount(ctx context.Context, Id string) (*domain.Account, error)
	GetAllAccounts(ctx context.Context) ([]*domain.Account, error)
}

type accountQueryService struct {
	accountRepository repository.IAccountRepository
}

func NewAccountQueryService(accountRepository repository.IAccountRepository) IAccountQueryService {
	return &accountQueryService{
		accountRepository: accountRepository,
	}
}

func (u *accountQueryService) GetAccount(ctx context.Context, Id string) (*domain.Account, error) {
	account, err := u.accountRepository.GetAccount(ctx, Id)

	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}

	return account, nil
}

func (u *accountQueryService) GetAllAccounts(ctx context.Context) ([]*domain.Account, error) {
	accounts, err := u.accountRepository.GetAllAccounts(ctx)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}
