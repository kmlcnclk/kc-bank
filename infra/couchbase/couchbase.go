package couchbase

import (
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"
)

type Couchbase struct {
	cluster *gocb.Cluster
}

func NewCouchbase(url, username, password string) (*Couchbase, error) {
	cluster, err := ConnectCouchbase(url, username, password)

	if err != nil {
		zap.L().Error("Failed to initialize Couchbase:", zap.Error(err))
		return nil, err // Return nil instead of an invalid instance.
	}

	return &Couchbase{cluster: cluster}, nil
}

func ConnectCouchbase(url, username, password string) (*gocb.Cluster, error) {
	cluster, err := gocb.Connect(url, gocb.ClusterOptions{
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

func (c *Couchbase) NewBucket(bucketName string) (*gocb.Bucket, error) {
	bucket := c.cluster.Bucket(bucketName)

	err := bucket.WaitUntilReady(3*time.Second, &gocb.WaitUntilReadyOptions{})

	if err != nil {
		zap.L().Fatal("Failed to open bucket", zap.Error(err))
		return nil, err
	}

	return bucket, nil
}

func (c *Couchbase) InitializeBucket(bucketName string) *gocb.Bucket {
	bucket, err := c.NewBucket(bucketName)

	if err != nil {
		zap.L().Fatal("failed to create new bucket", zap.Error(err))
	}

	return bucket
}
