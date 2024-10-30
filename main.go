package main

import (
	consulServiceManager "github.com/ayushs-2k4/go-consul-service-manager/consulservicemanager"
	"kontest-email-service/consumer"
	"kontest-email-service/service"
	"kontest-email-service/utils/kafka_utils"
	"log"
	"os"
	"strconv"
)

var (
	applicationHost = "localhost"             // Default value for local development
	applicationPort = 5151                    // Default value for local development
	serviceName     = "KONTEST-EMAIL-SERVICE" // Service name for Service Registry
	consulHost      = "localhost"             // Default value for local development
	consulPort      = 5150                    // Port as a constant (can be constant if it won't change)
)

func initializeVariables() {
	// Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error fetching hostname: %v", err)
	}

	// Attempt to read the KONTEST_API_SERVER_HOST environment variable
	if host := os.Getenv("KONTEST_EMAIL_SERVICE_HOST"); host != "" {
		applicationHost = host // Override with the environment variable if set
	} else {
		applicationHost = hostname // Use the machine's hostname if the env var is not set
	}

	// Attempt to read the KONTEST_API_SERVER_PORT environment variable
	if port := os.Getenv("KONTEST_EMAIL_SERVICE_PORT"); port != "" {
		parsedPort, err := strconv.Atoi(port)
		if err != nil {
			log.Fatalf("Invalid port value: %v", err)
		}
		applicationPort = parsedPort // Override with the environment variable if set
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

	// Load kafkaConfig
	kafkaConfig := kafka_utils.GetKafkaConfig()

	consulService := consulServiceManager.NewConsulService(consulHost, consulPort)
	consulService.Start(applicationHost, applicationPort, serviceName, []string{})

	// Load kafka_utils configuration
	broker := kafkaConfig.KafkaHost + ":" + kafkaConfig.KafkaPort

	// Initialize the EmailService and KafkaMessageListener
	emailService := service.NewMailJetEmailService()
	listener := consumer.NewKafkaMessageListener(emailService)

	// Start consuming messages
	listener.ConsumeMessages(kafkaConfig.Topics, "kontest-email-service", broker)
}
