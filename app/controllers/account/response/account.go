package response

import (
	"kc-bank/domain"
	"time"
)

type AccountResponse struct {
	Id        string    `json:"id"`
	Currency  string    `json:"currency"`
	Iban      string    `json:"iban"`
	Balance   float64   `json:"balance"`
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToAccountResponse(account *domain.Account) AccountResponse {
	return AccountResponse{
		Id:        account.Id,
		Currency:  account.Currency,
		Iban:      account.Iban,
		Balance:   account.Balance,
		UserId:    account.UserId,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

func ToAccountResponseList(accounts []*domain.Account) []AccountResponse {
	var response = make([]AccountResponse, 0)

	for _, account := range accounts {
		response = append(response, ToAccountResponse(account))
	}

	return response
}
