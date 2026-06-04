package config

import configx "github.com/iamKienb/go-core/config"

type InventoryCommandConfig struct {
	Postgres configx.PostgresConfig `envPrefix:"INVENTORY_COMMAND_SERVICE"`
	Server   configx.Server         `envPrefix:"INVENTORY_COMMAND_SERVICE"`
}
