package user

import (
	"context"
	"kc-bank/app/services/user/command"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required,min=8,max=16"`
	Age       int32  `json:"age" validate:"required"`
}

func (req *CreateUserRequest) ToCommand() command.Command {
	return command.Command{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Age:       req.Age,
		Id:        "",
	}
}

type CreateUserResponse struct {
	Message string `json:"message"`
}

type CreateUserHandler struct {
	command command.ICommandHandler
}

func NewCreateUserHandler(command command.ICommandHandler) *CreateUserHandler {
	return &CreateUserHandler{
		command: command,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	err := h.command.Save(ctx, req.ToCommand())

	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		Message: "User Created Successfully",
	}, nil
}
