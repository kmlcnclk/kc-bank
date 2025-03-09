package command

import (
	"context"
	"kc-bank/app/repository"
	"kc-bank/domain"
	"kc-bank/pkg/services"
	"time"

	"github.com/google/uuid"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) error
}

type commandHandler struct {
	accountRepository repository.IAccountRepository
	ibanService       services.IIbanService
}

func NewCommandHandler(accountRepository repository.IAccountRepository, ibanService services.IIbanService) ICommandHandler {
	return &commandHandler{
		accountRepository: accountRepository,
		ibanService:       ibanService,
	}
}

func (c *commandHandler) Save(ctx context.Context, command Command) error {
	// TODO: check user id for existence

	iban := c.ibanService.GenerateIBAN("TR", 5, 16)

	newAccount := c.BuildEntity(command, iban)

	err := c.accountRepository.CreateAccount(ctx, newAccount)

	if err != nil {
		return err
	}

	return nil
}

func (c *commandHandler) BuildEntity(command Command, iban string) *domain.Account {
	return &domain.Account{
		Id:        uuid.New().String(),
		Currency:  command.Currency,
		Iban:      iban,
		Balance:   0.0,
		UserId:    command.UserId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
