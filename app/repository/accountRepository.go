package repository

import (
	"context"
	"errors"
	"kc-bank/domain"
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	GetAllAccounts(ctx context.Context) ([]*domain.Account, error)
	FindByIban(ctx context.Context, iban string) (string, error)
	CheckAmountForFromIban(ctx context.Context, iban string, amount float64) (bool, error)
	TransferMoney(ctx context.Context, fromIbanId, toIbanId string, amount float64) error
}

type accountRepository struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
}

func NewAccountRepository(cluster *gocb.Cluster, bucket *gocb.Bucket) IAccountRepository {
	return &accountRepository{
		cluster: cluster,
		bucket:  bucket,
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	_, err := r.bucket.DefaultCollection().Insert(account.Id, account, &gocb.InsertOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})

	if err != nil {
		zap.L().Error("Failed to create account", zap.Error(err))
		return err
	}

	return nil
}

func (r *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {

	_, err := r.bucket.DefaultCollection().Replace(account.Id, account, &gocb.ReplaceOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})

	if err != nil {
		zap.L().Error("Failed to update account", zap.Error(err))
		return err
	}

	return nil
}

func (r *accountRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	data, err := r.bucket.DefaultCollection().Get(id, &gocb.GetOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("account not found")
		}

		zap.L().Error("Failed to get account", zap.Error(err))
		return nil, err
	}

	var account domain.Account
	if err := data.Content(&account); err != nil {
		zap.L().Error("Failed to unmarshal account", zap.Error(err))
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) GetAllAccounts(ctx context.Context) ([]*domain.Account, error) {
	query := "SELECT META(a).id, a.* FROM `accounts` a ORDER BY a.CreatedAt DESC"

	rows, err := r.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
	})

	if err != nil {
		zap.L().Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Row(&account); err != nil {
			zap.L().Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("Error iterating rows", zap.Error(err))
		return nil, err
	}

	return accounts, nil
}

func (r *accountRepository) FindByIban(ctx context.Context, iban string) (string, error) {
	query := "SELECT META(a).id, a.id FROM `accounts` a WHERE a.Iban = $iban LIMIT 1"

	rows, err := r.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: map[string]interface{}{"iban": iban},
		Adhoc:           true,
	})

	if err != nil {
		zap.L().Error("Failed to execute query", zap.Error(err))
		return "", err
	}

	defer rows.Close()

	var account struct {
		Id string
	}

	if rows.Next() {
		if err := rows.Row(&account); err != nil {
			zap.L().Error("Failed to scan row", zap.Error(err))
			return "", err
		}

		return account.Id, nil
	}

	return "", nil
}

func (r *accountRepository) CheckAmountForFromIban(ctx context.Context, iban string, amount float64) (bool, error) {
	query := "SELECT Balance FROM `accounts` WHERE Iban = $iban LIMIT 1"

	rows, err := r.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: map[string]interface{}{"iban": iban},
		Adhoc:           true,
	})
	if err != nil {
		zap.L().Error("Failed to execute query", zap.Error(err))
		return false, err
	}
	defer rows.Close()

	var rawBalance struct {
		Balance float64
	}

	// Check if there is a row
	if rows.Next() {
		if err := rows.Row(&rawBalance); err != nil {
			zap.L().Error("Failed to scan row", zap.Error(err))
			return false, err
		}

		return rawBalance.Balance >= amount, nil
	}

	// No matching account found
	return false, nil
}

func (r *accountRepository) TransferMoney(ctx context.Context, fromIbanId, toIbanId string, amount float64) error {
	// r.bucket.DefaultCollection().MutateIn(fromIbanId, []gocb.MutateInSpec{
	// 	gocb.DecrementSpec("Balance", int64(amount), nil),
	// }, &gocb.MutateInOptions{
	// 	Context: ctx,
	// })

	// // Perform the transfer within a transaction
	// _, err := r.cluster.Query(`
	// 	UPDATE accounts SET Balance = Balance - $amount WHERE Iban = $fromIban;
	// 	UPDATE accounts SET Balance = Balance + $amount WHERE Iban = $toIban;
	// `, &gocb.QueryOptions{
	// 	Context: ctx,
	// 	NamedParameters: map[string]interface{}{
	// 		"amount":   amount,
	// 		"fromIban": fromIban,
	// 		"toIban":   toIban,
	// 	},
	// 	Adhoc: true,
	// })
	// if err != nil {
	// 	zap.L().Error("Failed to execute transfer query", zap.Error(err))
	// 	return err
	// }

	return nil
}
