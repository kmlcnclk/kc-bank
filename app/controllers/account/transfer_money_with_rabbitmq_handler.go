package account

import (
	"context"
	"kc-bank/app/services/account/command"
)

type TransferMoneyWithRabbitMQRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	FromIBAN string  `json:"fromIBAN" validate:"required"`
	ToIBAN   string  `json:"toIBAN" validate:"required"`
}

func (req *TransferMoneyWithRabbitMQRequest) ToCommand() command.TransferMoneyCommand {
	return command.TransferMoneyCommand{
		Amount:   req.Amount,
		FromIBAN: req.FromIBAN,
		ToIBAN:   req.ToIBAN,
	}
}

type TransferMoneyWithRabbitMQResponse struct {
	Message string `json:"message"`
}

type TransferMoneyWithRabbitMQHandler struct {
	command command.ICommandHandler
}

func NewTransferMoneyWithRabbitMQHandler(command command.ICommandHandler) *TransferMoneyWithRabbitMQHandler {
	return &TransferMoneyWithRabbitMQHandler{
		command: command,
	}
}

func (h *TransferMoneyWithRabbitMQHandler) Handle(ctx context.Context, req *TransferMoneyWithRabbitMQRequest) (*TransferMoneyWithRabbitMQResponse, error) {
	err := h.command.TransferMoneyWithRabbitMQPublisher(ctx, req.ToCommand())

	if err != nil {
		return nil, err
	}

	return &TransferMoneyWithRabbitMQResponse{
		Message: "Amount successfully transferred!",
	}, nil
}
