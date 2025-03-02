package couchbase

import (
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"
)

type Couchbase struct {
	cluster *gocb.Cluster
}

func (c *Couchbase) NewBucket(name string) (*gocb.Bucket, error) {
	bucket := c.cluster.Bucket(name)

	err := bucket.WaitUntilReady(3*time.Second, &gocb.WaitUntilReadyOptions{})

	if err != nil {
		zap.L().Fatal("Failed to open bucket", zap.Error(err))
		return nil, err
	}

	return bucket, nil
}

func NewCouchbase(username, password string) (*Couchbase, error) {
	cluster, err := ConnectCouchbase(username, password)

	if err != nil {
		zap.L().Error("Failed to initialize Couchbase:", zap.Error(err))
		return nil, err // Return nil instead of an invalid instance.
	}

	return &Couchbase{cluster: cluster}, nil
}

func ConnectCouchbase(username, password string) (*gocb.Cluster, error) {
	cluster, err := gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout: 3 * time.Second,
			KVTimeout:      3 * time.Second,
			QueryTimeout:   3 * time.Second,
		},
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})

	if err != nil {
		zap.L().Fatal("Failed to connect to couchbase", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Connected to couchbase")

	return cluster, nil
}

func InitializeCouchbase(username, password string) *gocb.Bucket {
	cb, err := NewCouchbase(username, password)

	if err != nil {
		zap.L().Fatal("Failed to initialize Couchbase instance", zap.Error(err))
	}

	bucket, err := cb.NewBucket("users")

	if err != nil {
		zap.L().Fatal("failed to create new bucket", zap.Error(err))
	}

	return bucket
}
