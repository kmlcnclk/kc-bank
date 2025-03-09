package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port                              string `yaml:"port" mapstructure:"port"`
	RabbitMQURL                       string `yaml:"rabbitmq_url" mapstructure:"rabbitmq_url"`
	RabbitMQTransferMoneyQueueName    string `yaml:"rabbitmq_transfer_money_queue_name" mapstructure:"rabbitmq_transfer_money_queue_name"`
	RabbitMQTransferMoneyExchangeName string `yaml:"rabbitmq_transfer_money_exchange_name" mapstructure:"rabbitmq_transfer_money_exchange_name"`
	RabbitMQTransferMoneyExchangeType string `yaml:"rabbitmq_transfer_money_exchange_type" mapstructure:"rabbitmq_transfer_money_exchange_type"`
	CouchbaseUrl                      string `yaml:"couchbase_url" mapstructure:"couchbase_url"`
	CouchbaseUsername                 string `yaml:"couchbase_username" mapstructure:"couchbase_username"`
	CouchbasePassword                 string `yaml:"couchbase_password" mapstructure:"couchbase_password"`
}

func Read() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &appConfig
}
