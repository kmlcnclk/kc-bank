package account

import (
	"context"
	"kc-bank/app/services/account/command"
)

type TransferMoneyRequest struct {
	Amount   float64 `json:"amount" validate:"required"`
	FromIBAN string  `json:"fromIBAN" validate:"required"`
	ToIBAN   string  `json:"toIBAN" validate:"required"`
}

func (req *TransferMoneyRequest) ToCommand() command.TransferMoneyCommand {
	return command.TransferMoneyCommand{
		Amount:   req.Amount,
		FromIBAN: req.FromIBAN,
		ToIBAN:   req.ToIBAN,
	}
}

type TransferMoneyResponse struct {
	Message string `json:"message"`
}

type TransferMoneyHandler struct {
	command command.ICommandHandler
}

func NewTransferMoneyHandler(command command.ICommandHandler) *TransferMoneyHandler {
	return &TransferMoneyHandler{
		command: command,
	}
}

func (h *TransferMoneyHandler) Handle(ctx context.Context, req *TransferMoneyRequest) (*TransferMoneyResponse, error) {
	err := h.command.TransferMoney(ctx, req.ToCommand())

	if err != nil {
		return nil, err
	}

	return &TransferMoneyResponse{
		Message: "Amount successfully transferred!",
	}, nil
}
