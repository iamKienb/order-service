package config

import configx "github.com/iamKienb/go-core/config"

type OrderQueryConfig struct {
	ES     configx.ElasticSearchConfig `envPrefix:"ORDER_QUERY_SERVICE"`
	Server configx.Server              `envPrefix:"ORDER_QUERY_SERVICE"`
}
