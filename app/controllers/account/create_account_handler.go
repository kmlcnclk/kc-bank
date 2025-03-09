package account

import (
	"context"
	"kc-bank/app/services/account/command"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" validate:"required"`
	UserId   string `json:"userId" validate:"required"`
}

func (req *CreateAccountRequest) ToCommand() command.Command {
	return command.Command{
		Currency: req.Currency,
		UserId:   req.UserId,
	}
}

type CreateAccountResponse struct {
	Message string `json:"message"`
}

type CreateAccountHandler struct {
	command command.ICommandHandler
}

func NewCreateAccountHandler(command command.ICommandHandler) *CreateAccountHandler {
	return &CreateAccountHandler{
		command: command,
	}
}

func (h *CreateAccountHandler) Handle(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	err := h.command.Save(ctx, req.ToCommand())

	if err != nil {
		return nil, err
	}

	return &CreateAccountResponse{
		Message: "Account Created Successfully",
	}, nil
}
