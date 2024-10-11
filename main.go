package main

import (
	"kontest-email-service/consumer"
	"kontest-email-service/service"
	"kontest-email-service/utils"
)

func main() {
	println("Hello, World!")

	// Load config
	config := utils.GetConfig()

	// Load Kafka configuration
	broker := config.KafkaHost + ":" + config.KafkaPort

	// Initialize the EmailService and KafkaMessageListener
	emailService := &service.MailJetEmailService{}
	listener := consumer.NewKafkaMessageListener(emailService)

	// Start consuming messages
	listener.ConsumeMessages(config.Topics, "kontest-email-service", broker)
}
