package utils

import (
	"os"
	"sync"
)

type Config struct {
	KafkaHost string
	KafkaPort string
	Topics    []string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

func loadConfig() *Config {
	// Load Kafka host and port
	kafkaHost := os.Getenv("KAFKA_HOST")
	if kafkaHost == "" {
		kafkaHost = "localhost"
	}

	kafkaPort := os.Getenv("KAFKA_PORT")
	if kafkaPort == "" {
		kafkaPort = "9092"
	}

	// Load Kafka Topics from environment variables or use default values
	topics := []string{
		loadKafkaTopic(UserRegistrationEventTopic),
		loadKafkaTopic(AccountDeletionEventTopic),
		loadKafkaTopic(AccountDeletionEmailEventTopic),
		loadKafkaTopic(PasswordChangeEmailEventTopic),
	}

	return &Config{
		KafkaHost: kafkaHost,
		KafkaPort: kafkaPort,
		Topics:    topics,
	}
}

// loadKafkaTopic returns the environment variable value or the default value
func loadKafkaTopic(topic KafkaTopic) string {
	if value := os.Getenv(topic.envVarName); value != "" {
		topic.DefaultValue = value
	}

	return topic.DefaultValue
}
