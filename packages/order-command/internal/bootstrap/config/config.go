package config

import (
	"time"

	configx "github.com/iamKienb/go-core/config"
)

type TemporalConfig struct {
	Address         string        `env:"_TEMPORAL_ADDRESS"`
	Namespace       string        `env:"_TEMPORAL_NAMESPACE"`
	OrderTaskQueue  string        `env:"_TEMPORAL_ORDER_TASK_QUEUE"`
	ActivityTimeout time.Duration `env:"_TEMPORAL_ACTIVITY_TIMEOUT"`
	RollbackTimeout time.Duration `env:"_TEMPORAL_ROLLBACK_TIMEOUT"`
}

type UpstreamConfig struct {
	UserCommandURL      string `env:"_USER_COMMAND_URL" envDefault:"http://localhost:8001"`
	ProductQueryURL     string `env:"_PRODUCT_QUERY_URL" envDefault:"http://localhost:8103"`
	InventoryCommandURL string `env:"_INVENTORY_COMMAND_URL" envDefault:"http://localhost:8004"`
}

type OrderCommandConfig struct {
	Postgres configx.PostgresConfig `envPrefix:"ORDER_COMMAND_SERVICE"`
	Server   configx.Server         `envPrefix:"ORDER_COMMAND_SERVICE"`
	Temporal TemporalConfig         `envPrefix:"ORDER_COMMAND_SERVICE"`
	Upstream UpstreamConfig         `envPrefix:"ORDER_COMMAND_SERVICE"`
}
