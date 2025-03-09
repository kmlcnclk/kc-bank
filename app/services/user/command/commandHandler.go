package command

import (
	"context"
	"fmt"
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
	userRepository  repository.IUserRepository
	passwordService services.IPasswordService
}

func NewCommandHandler(userRepository repository.IUserRepository, passwordService services.IPasswordService) ICommandHandler {
	return &commandHandler{
		userRepository:  userRepository,
		passwordService: passwordService,
	}
}

func (c *commandHandler) Save(ctx context.Context, command Command) error {
	// TODO: add national id check

	hashedPassword, err := c.passwordService.HashPassword(command.Password)

	if err != nil {
		return fmt.Errorf("password could not hash: %s", err.Error())
	}

	newUser := c.BuildEntity(command, hashedPassword)

	err = c.userRepository.CreateUser(ctx, newUser)

	if err != nil {
		return err
	}

	return nil
}

func (c *commandHandler) BuildEntity(command Command, hashedPassword string) *domain.User {
	return &domain.User{
		Id:        uuid.New().String(),
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Email:     command.Email,
		Password:  hashedPassword,
		Age:       command.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
