package main

import (
	consulServiceManager "github.com/ayushs-2k4/go-consul-service-manager"
	"kontest-email-service/consumer"
	"kontest-email-service/service"
	"kontest-email-service/utils"
	"log"
	"os"
	"strconv"
)

var (
	applicationHost = "localhost"             // Default value for local development
	applicationPort = "5151"                  // Default value for local development
	serviceName     = "KONTEST-EMAIL-SERVICE" // Service name for Service Registry
	consulHost      = "localhost"             // Default value for local development
	consulPort      = 5157                    // Port as a constant (can be constant if it won't change)
)

func initializeVariables() {
	// Attempt to read the KONTEST_API_SERVER_HOST environment variable
	if host := os.Getenv("KONTEST_EMAIL_SERVICE_HOST"); host != "" {
		applicationHost = host // Override with the environment variable if set
	}

	// Attempt to read the KONTEST_API_SERVER_PORT environment variable
	if port := os.Getenv("KONTEST_EMAIL_SERVICE_PORT"); port != "" {
		applicationPort = port // Override with the environment variable if set
	}

	// Attempt to read the CONSUL_ADDRESS environment variable
	if host := os.Getenv("CONSUL_HOST"); host != "" {
		consulHost = host // Override with the environment variable if set
	}

	// Attempt to read the CONSUL_PORT environment variable
	if port := os.Getenv("CONSUL_PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil {
			consulPort = portInt // Override with the environment variable if set and valid
		}
	}
}

func main() {
	initializeVariables()

	portInt, err := strconv.Atoi(applicationPort)
	if err != nil {
		log.Fatalf("Failed to convert applicationPort to integer: %v", err)
	}

	// Load config
	config := utils.GetConfig()

	consulService := consulServiceManager.NewConsulService(consulHost, consulPort)
	consulService.Start(applicationHost, portInt, serviceName)

	// Load Kafka configuration
	broker := config.KafkaHost + ":" + config.KafkaPort

	// Initialize the EmailService and KafkaMessageListener
	emailService := &service.MailJetEmailService{}
	listener := consumer.NewKafkaMessageListener(emailService)

	// Start consuming messages
	listener.ConsumeMessages(config.Topics, "kontest-email-service", broker)
}
