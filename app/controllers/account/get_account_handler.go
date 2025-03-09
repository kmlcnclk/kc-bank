package account

import (
	"context"
	"kc-bank/app/controllers/account/response"
	"kc-bank/app/services/account/query"
)

type GetAccountRequest struct {
	Id string `json:"id" param:"id"`
}

type GetAccountResponse struct {
	Account response.AccountResponse `json:"account"`
}

type GetAccountHandler struct {
	queryService query.IAccountQueryService
}

func NewGetAccountHandler(queryService query.IAccountQueryService) *GetAccountHandler {
	return &GetAccountHandler{
		queryService: queryService,
	}
}

func (h *GetAccountHandler) Handle(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	// TODO: user just should be able to see his/her accounts
	account, err := h.queryService.GetAccount(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &GetAccountResponse{Account: response.ToAccountResponse(account)}, nil
}
