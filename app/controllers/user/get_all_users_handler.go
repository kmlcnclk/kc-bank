package user

import (
	"context"
	"kc-bank/app/controllers/user/response"
	"kc-bank/app/services/user/query"
)

type GetUserAllRequest struct{}

type GetUserAllResponse struct {
	Users []response.UserResponse `json:"users"`
}

type GetUserAllHandler struct {
	queryService query.IUserQueryService
}

func NewGetUserAllHandler(queryService query.IUserQueryService) *GetUserAllHandler {
	return &GetUserAllHandler{
		queryService: queryService,
	}
}

func (h *GetUserAllHandler) Handle(ctx context.Context, req *GetUserAllRequest) (*GetUserAllResponse, error) {
	users, err := h.queryService.GetAllUsers(ctx)

	if err != nil {
		return nil, err
	}

	return &GetUserAllResponse{Users: response.ToUserResponseList(users)}, nil
}
