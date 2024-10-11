package consumer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"kontest-email-service/service"
	"kontest-email-service/utils"
	"log"
)

type KafkaMessageListener struct {
	emailService service.EmailService
	config       *utils.Config
}

// NewKafkaMessageListener creates a new KafkaMessageListener with the provided email service
func NewKafkaMessageListener(emailService service.EmailService) *KafkaMessageListener {
	return &KafkaMessageListener{
		emailService: emailService,
	}
}

// ConsumeMessages listens for messages on the specified Kafka topics and dispatches to appropriate handlers
func (l *KafkaMessageListener) ConsumeMessages(topics []string, groupId string, broker string) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	// Subscribe to the topics
	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %s", err)
	}

	log.Printf("Subscribed to topics: %v", topics)

	// Poll for messages and process them
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Printf("Message received on topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
			l.processMessage(string(msg.Value), *msg.TopicPartition.Topic)
		} else {
			log.Printf("Error while consuming message: %v", err)
		}
	}
}

//// processMessage determines the topic and calls the appropriate handler function
//func (l *KafkaMessageListener) processMessage(message string, topic string) {
//	switch utils.KafkaTopic(topic) {
//	case utils.UserRegistrationEventTopic.:
//		l.handleUserRegistrationMessage(message)
//
//	case utils.AccountDeletionEventTopic:
//		log.Printf("Received message for Account Deletion Event Topic, but no action is defined.")
//
//	case utils.AccountDeletionEmailEventTopic:
//		l.handleAccountDeletionMessage(message)
//
//	case utils.PasswordChangeEmailEventTopic:
//		l.handlePasswordChangeMessage(message)
//
//	default:
//		log.Printf("Unknown topic: %s", topic)
//	}
//}

// processMessage determines the topic and calls the appropriate handler function
func (l *KafkaMessageListener) processMessage(message string, topic string) {
	switch topic {
	case utils.UserRegistrationEventTopic.DefaultValue:
		l.handleUserRegistrationMessage(message)

	case utils.AccountDeletionEventTopic.DefaultValue:
		log.Printf("Received message for Account Deletion Event Topic, but no action is defined.")

	case utils.AccountDeletionEmailEventTopic.DefaultValue:
		l.handleAccountDeletionMessage(message)

	case utils.PasswordChangeEmailEventTopic.DefaultValue:
		l.handlePasswordChangeMessage(message)
	}
}

// handleUserRegistrationMessage processes user registration events
func (l *KafkaMessageListener) handleUserRegistrationMessage(message string) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(message), &jsonData)
	if err != nil {
		log.Printf("Failed to parse user registration message: %v", err)
		return
	}

	email, ok := jsonData["email"].(string)
	if !ok {
		log.Printf("Invalid email in user registration message: %s", message)
		return
	}

	registrationDate, _ := jsonData["registrationDate"].(string)
	l.emailService.SendEmail(email, "Welcome to Kontest", "Thank you for registering with us!")
	log.Printf("Processed user registration message for email: %s, registration date: %s", email, registrationDate)
}

// handleAccountDeletionMessage processes account deletion events
func (l *KafkaMessageListener) handleAccountDeletionMessage(message string) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(message), &jsonData)
	if err != nil {
		log.Printf("Failed to parse account deletion message: %v", err)
		return
	}

	email, ok := jsonData["email"].(string)
	if !ok {
		log.Printf("Invalid email in account deletion message: %s", message)
		return
	}

	deletionDate, _ := jsonData["deletionDate"].(string)
	l.emailService.SendEmail(email, "Account Deletion", "Your account has been deleted!")
	log.Printf("Processed account deletion message for email: %s, deletion date: %s", email, deletionDate)
}

// handlePasswordChangeMessage processes password change events
func (l *KafkaMessageListener) handlePasswordChangeMessage(message string) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(message), &jsonData)
	if err != nil {
		log.Printf("Failed to parse password change message: %v", err)
		return
	}

	email, ok := jsonData["email"].(string)
	if !ok {
		log.Printf("Invalid email in password change message: %s", message)
		return
	}

	updateDate, _ := jsonData["updateDate"].(string)
	l.emailService.SendEmail(email, "Password Changed", "Your account password has been changed! If you did not make this change, please contact us immediately.")
	log.Printf("Processed password change message for email: %s, update date: %s", email, updateDate)
}
