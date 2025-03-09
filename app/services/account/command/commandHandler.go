package command

import (
	"context"
	"errors"
	"kc-bank/app/repository"
	"kc-bank/domain"
	"kc-bank/pkg/services"
	"time"

	"github.com/google/uuid"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) error
	TransferMoney(ctx context.Context, command TransferMoneyCommand) error
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

func (c *commandHandler) TransferMoney(ctx context.Context, command TransferMoneyCommand) error {
	// TODO: check user id for existence

	fromIbanId, err := c.accountRepository.FindByIban(ctx, command.FromIBAN)

	if err != nil {
		return err
	}

	if len(fromIbanId) == 0 {
		return errors.New("from iban does not exist")
	}

	toIbanId, err := c.accountRepository.FindByIban(ctx, command.ToIBAN)

	if err != nil {
		return err
	}

	if len(toIbanId) == 0 {
		return errors.New("to iban does not exist")
	}

	isBalanceEnough, err := c.accountRepository.CheckAmountForFromIban(ctx, command.FromIBAN, command.Amount)

	if err != nil {
		return err
	}

	if !isBalanceEnough {
		return errors.New("balance is not enough")
	}

	err = c.accountRepository.TransferMoney(ctx, fromIbanId, toIbanId, command.Amount)

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
