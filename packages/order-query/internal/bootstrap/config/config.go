package config

import configx "github.com/iamKienb/go-core/config"

type InventoryQueryConfig struct {
	ES     configx.ElasticSearchConfig `envPrefix:"INVENTORY_QUERY_SERVICE"`
	Server configx.Server              `envPrefix:"INVENTORY_QUERY_SERVICE"`
}
