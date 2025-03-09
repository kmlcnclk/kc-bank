package account

import (
	"context"
	"kc-bank/app/controllers/account/response"
	"kc-bank/app/services/account/query"
)

type GetAccountAllRequest struct{}

type GetAccountAllResponse struct {
	Accounts []response.AccountResponse `json:"accounts"`
}

type GetAccountAllHandler struct {
	queryService query.IAccountQueryService
}

func NewGetAccountAllHandler(queryService query.IAccountQueryService) *GetAccountAllHandler {
	return &GetAccountAllHandler{
		queryService: queryService,
	}
}

func (h *GetAccountAllHandler) Handle(ctx context.Context, req *GetAccountAllRequest) (*GetAccountAllResponse, error) {
	// TODO: user just should be able to see his/her accounts
	accounts, err := h.queryService.GetAllAccounts(ctx)

	if err != nil {
		return nil, err
	}

	return &GetAccountAllResponse{Accounts: response.ToAccountResponseList(accounts)}, nil
}
