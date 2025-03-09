package command

import (
	"context"
	"encoding/json"
	"errors"
	"kc-bank/app/repository"
	"kc-bank/domain"
	"kc-bank/infra/rabbitmq"
	"kc-bank/pkg/services"
	"log"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ICommandHandler interface {
	Save(ctx context.Context, command Command) error
	TransferMoney(ctx context.Context, command TransferMoneyCommand) error
	TransferMoneyWithRabbitMQPublisher(ctx context.Context, command TransferMoneyCommand) error
	TransferMoneyWithRabbitMQConsumer()
}

type commandHandler struct {
	accountRepository repository.IAccountRepository
	ibanService       services.IIbanService
	rmqService        rabbitmq.IRabbitMQService
	exchangeName      string
}

func NewCommandHandler(
	accountRepository repository.IAccountRepository,
	ibanService services.IIbanService,
	rmqService rabbitmq.IRabbitMQService,
	exchangeName string,
) ICommandHandler {
	return &commandHandler{
		accountRepository: accountRepository,
		ibanService:       ibanService,
		rmqService:        rmqService,
		exchangeName:      exchangeName,
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

func (c *commandHandler) TransferMoneyWithRabbitMQPublisher(ctx context.Context, command TransferMoneyCommand) error {

	serializedData, err := json.Marshal(command)

	if err != nil {
		zap.L().Error("Failed to serialize data", zap.Error(err))
		return err
	}

	err = c.rmqService.Publish(c.exchangeName, "", serializedData)

	if err != nil {
		zap.L().Error("Failed to publish message", zap.Error(err))
		return err
	}

	zap.L().Info("Message published successfully")

	return nil
}

func (c *commandHandler) TransferMoneyWithRabbitMQConsumer() {

	msgs, err := c.rmqService.Consume()

	if err != nil {
		zap.L().Error("Failed to consume message", zap.Error(err))
		log.Fatal("Failed to start consuming messages:", err)
	}

	for msg := range msgs {
		var transferReq TransferMoneyCommand
		if err := json.Unmarshal(msg.Body, &transferReq); err != nil {
			log.Println("Failed to parse transfer request:", err)
			continue
		}

		err := c.TransferMoney(context.Background(), transferReq)

		if err != nil {
			zap.L().Error("Failed to transfer money", zap.Error(err))
			log.Println("Failed to transfer money:", err)
		}

		zap.L().Info("Message consumed successfully")
	}
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
