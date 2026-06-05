package config

import (
	"time"

	configx "github.com/iamKienb/go-core/config"
)

type TemporalConfig struct {
	Address         string        `env:"TEMPORAL_ADDRESS"`
	Namespace       string        `env:"TEMPORAL_NAMESPACE"`
	OrderTaskQueue  string        `env:"TEMPORAL_ORDER_TASK_QUEUE"`
	ActivityTimeout time.Duration `env:"TEMPORAL_ACTIVITY_TIMEOUT"`
	RollbackTimeout time.Duration `env:"TEMPORAL_ROLLBACK_TIMEOUT"`
}

type ApiConfig struct {
	APIGatewayAddr string `env:"API_GATEWAY_SERVICE_ADDR"`
}

type OrderCommandConfig struct {
	Postgres configx.PostgresConfig `envPrefix:"ORDER_COMMAND_SERVICE"`
	Server   configx.Server         `envPrefix:"ORDER_COMMAND_SERVICE"`
	Temporal TemporalConfig         `envPrefix:"ORDER_COMMAND_SERVICE"`
	Gateway  ApiConfig              `envPrefix:"ORDER_COMMAND_SERVICE"`
}
