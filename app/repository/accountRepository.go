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
	query := "SELECT META(u).id, u.* FROM `accounts` u ORDER BY u.CreatedAt DESC"

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
