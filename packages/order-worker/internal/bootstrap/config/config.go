package config

import configx "github.com/iamKienb/go-core/config"

type InventoryWorkerConfig struct {
	ES       configx.ElasticSearchConfig `envPrefix:"INVENTORY_WORKER_SERVICE"`
	Kafka    configx.KafkaConfig         `envPrefix:"INVENTORY_WORKER_SERVICE"`
	Consumer configx.ConsumerConfig      `envPrefix:"INVENTORY_WORKER_SERVICE"`
}
