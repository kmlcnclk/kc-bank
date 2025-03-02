package user

import (
	"context"
	"kc-bank/app/controllers/user/response"
	"kc-bank/app/services/query"
)

type GetUserRequest struct {
	Id string `json:"id" param:"id"`
}

type GetUserResponse struct {
	User response.UserResponse `json:"user"`
}

type GetUserHandler struct {
	queryService query.IUserQueryService
}

func NewGetUserHandler(queryService query.IUserQueryService) *GetUserHandler {
	return &GetUserHandler{
		queryService: queryService,
	}
}

func (h *GetUserHandler) Handle(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, err := h.queryService.GetUser(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &GetUserResponse{User: response.ToUserResponse(user)}, nil
}
