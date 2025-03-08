package repository

import (
	"context"
	"errors"
	"kc-bank/domain"
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}

type userRepository struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
}

func NewUserRepository(cluster *gocb.Cluster, bucket *gocb.Bucket) IUserRepository {
	return &userRepository{
		cluster: cluster,
		bucket:  bucket,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.bucket.DefaultCollection().Insert(user.Id, user, &gocb.InsertOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})

	if err != nil {
		zap.L().Error("Failed to create user", zap.Error(err))
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {

	_, err := r.bucket.DefaultCollection().Replace(user.Id, user, &gocb.ReplaceOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})

	if err != nil {
		zap.L().Error("Failed to update user", zap.Error(err))
		return err
	}

	return nil
}

func (r *userRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	data, err := r.bucket.DefaultCollection().Get(id, &gocb.GetOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("user not found")
		}

		zap.L().Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	var user domain.User
	if err := data.Content(&user); err != nil {
		zap.L().Error("Failed to unmarshal user", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := "SELECT META(u).id, u.* FROM `users` u ORDER BY u.CreatedAt DESC"

	rows, err := r.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
	})

	if err != nil {
		zap.L().Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Row(&user); err != nil {
			zap.L().Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("Error iterating rows", zap.Error(err))
		return nil, err
	}

	return users, nil
}
