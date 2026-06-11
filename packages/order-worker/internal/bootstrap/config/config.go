package config

import configx "github.com/iamKienb/go-core/config"

type OrderWorkerConfig struct {
	ES       configx.ElasticSearchConfig `envPrefix:"ORDER_WORKER_SERVICE"`
	Redis    configx.RedisConfig         `envPrefix:"ORDER_WORKER_SERVICE"`
	Kafka    configx.KafkaConfig         `envPrefix:"ORDER_WORKER_SERVICE"`
	Consumer configx.ConsumerConfig      `envPrefix:"ORDER_WORKER_SERVICE"`
}
