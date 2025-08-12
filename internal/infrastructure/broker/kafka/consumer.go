package kafka

import (
	"context"
	"log"
	"strings"
	"time"

	"strconv"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/SmoothWay/booking/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type MessageHandler func(message *kafka.Message) error

// MessageHandlerWithAck provides acknowledgment capability
type MessageHandlerWithAck func(message *kafka.Message) (bool, error)

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxRetries    int           // Maximum retry attempts before DLT
	RetryDelay    time.Duration // Delay between retries
	DLTMaxRetries int           // Maximum retry attempts in DLT
	DLTRetryDelay time.Duration // Delay between DLT retries
}

type Consumer struct {
	consumer      *kafka.Consumer
	topicHandlers map[string]MessageHandler
	retryConfig   RetryConfig
	dlqProducer   *Producer
}

func newConsumer(config *KafkaConfig, retryConfig RetryConfig, dlqProducer *Producer) *Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(config.Brokers, ","),
		"group.id":           config.GroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false, // Disable auto-commit for manual control
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	return &Consumer{
		consumer:      consumer,
		topicHandlers: make(map[string]MessageHandler),
		retryConfig:   retryConfig,
		dlqProducer:   dlqProducer,
	}
}

func (c *Consumer) Subscribe(topic string, handler MessageHandler) error {
	log.Printf("Subscribing to topic: %s", topic)
	c.topicHandlers[topic] = handler
	return c.consumer.Subscribe(topic, nil)
}

// SubscribeWithAck allows handlers to control acknowledgment with retry and DLT support
func (c *Consumer) SubscribeWithAck(topic string, handler MessageHandlerWithAck) error {
	log.Printf("Subscribing to topic with ack and retry: %s", topic)

	// Convert MessageHandlerWithAck to MessageHandler with retry logic
	wrappedHandler := func(msg *kafka.Message) error {
		return c.processMessageWithRetry(topic, msg, handler)
	}

	c.topicHandlers[topic] = wrappedHandler
	return c.consumer.Subscribe(topic, nil)
}

// processMessageWithRetry handles message processing with retry logic and DLT
func (c *Consumer) processMessageWithRetry(topic string, msg *kafka.Message, handler MessageHandlerWithAck) error {
	// Extract retry count from message headers
	retryCount := c.getRetryCount(msg)

	// Check if we've exceeded max retries
	if retryCount >= c.retryConfig.MaxRetries {
		log.Printf("Max retries (%d) exceeded for topic %s, sending to DLT", c.retryConfig.MaxRetries, topic)
		return c.sendToDLT(topic, msg, "max_retries_exceeded")
	}

	// Process the message
	shouldAck, err := handler(msg)
	if err != nil {
		log.Printf("Handler error for topic %s (attempt %d/%d): %v", topic, retryCount+1, c.retryConfig.MaxRetries, err)

		// Don't acknowledge - will be retried
		return c.scheduleRetry(topic, msg, retryCount)
	}

	if shouldAck {
		// Manually commit the offset
		_, err := c.consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Failed to commit message for topic %s: %v", topic, err)
			return err
		}
		log.Printf("Message acknowledged for topic %s at offset %v", topic, msg.TopicPartition.Offset)
	} else {
		log.Printf("Message NOT acknowledged for topic %s at offset %v", topic, msg.TopicPartition.Offset)
	}

	return nil
}

// getRetryCount extracts retry count from message headers
func (c *Consumer) getRetryCount(msg *kafka.Message) int {
	for _, header := range msg.Headers {
		if string(header.Key) == "retry_count" {
			if count, err := strconv.Atoi(string(header.Value)); err == nil {
				return count
			}
		}
	}
	return 0
}

// scheduleRetry schedules a message for retry with exponential backoff
func (c *Consumer) scheduleRetry(topic string, msg *kafka.Message, retryCount int) error {
	// Calculate delay with exponential backoff
	delay := c.retryConfig.RetryDelay * time.Duration(1<<retryCount)

	log.Printf("Scheduling retry for topic %s in %v (attempt %d/%d)", topic, delay, retryCount+1, c.retryConfig.MaxRetries)

	// Add retry count to message headers
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   "retry_count",
		Value: []byte(strconv.Itoa(retryCount + 1)),
	})

	// Send to retry topic with delay
	retryTopic := topic + ".retry"
	err := c.dlqProducer.SendMessage(context.Background(), retryTopic, msg.Value)
	if err != nil {
		log.Printf("Failed to send message to retry topic %s: %v", retryTopic, err)
		return err
	}

	// Don't acknowledge original message - it will be retried
	return nil
}

// sendToDLT sends a message to the Dead Letter Topic
func (c *Consumer) sendToDLT(topic string, msg *kafka.Message, reason string) error {
	dltTopic := topic + ".dlt"

	// Add DLT metadata to message headers
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   "dlt_reason",
		Value: []byte(reason),
	})

	log.Printf("Sending message to DLT topic %s: %s", dltTopic, reason)
	err := c.dlqProducer.SendMessage(context.Background(), dltTopic, msg.Value)
	if err != nil {
		log.Printf("Failed to send message to DLT topic %s: %v", dltTopic, err)
		return err
	}

	// Acknowledge the original message to prevent infinite retries
	_, err = c.consumer.CommitMessage(msg)
	if err != nil {
		log.Printf("Failed to commit message after sending to DLT: %v", err)
		return err
	}

	log.Printf("Message sent to DLT and original message acknowledged")
	return nil
}

func (c *Consumer) ConsumeLoop(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			c.consumer.Close()
			return ctx.Err()
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Received message on topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))

				// Route to specific handler if available
				if topicHandler, exists := c.topicHandlers[*msg.TopicPartition.Topic]; exists {
					if err := topicHandler(msg); err != nil {
						log.Printf("Error processing message from topic %s: %v", *msg.TopicPartition.Topic, err)

						continue
					}
				} else if handler != nil {
					log.Printf("Received message on topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
					if err := handler(msg); err != nil {
						log.Printf("Error processing message with generic handler: %v", err)
						continue
					}
				} else {
					log.Printf("No handler found for topic: %s", *msg.TopicPartition.Topic)
				}
			} else {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
					continue
				}
				log.Printf("Error reading message: %v", err)
			}
		}
	}
}

func InitKafkaConsumer(cfg *config.Config, bookingUsecase domain.BookingUsecase, userUsecase domain.UserUsecase, unitUsecase domain.UnitUsecase) *Consumer {
	// Initialize Kafka configuration
	kafkaConfig := &KafkaConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
	}

	// Configure retry behavior
	retryConfig := RetryConfig{
		MaxRetries:    3,                // Retry 3 times before DLT
		RetryDelay:    5 * time.Second,  // Start with 5 second delay
		DLTMaxRetries: 2,                // Retry 2 more times in DLT
		DLTRetryDelay: 30 * time.Second, // Longer delay in DLT
	}

	// Create DLT producer
	dlqProducer := NewProducer(cfg.Kafka.Brokers, cfg.Kafka.GroupID+"-dlq")

	log.Printf("Initializing Kafka consumer with brokers: %v, group: %s", cfg.Kafka.Brokers, cfg.Kafka.GroupID)
	log.Printf("Retry config: max=%d, delay=%v, dlt_max=%d, dlt_delay=%v",
		retryConfig.MaxRetries, retryConfig.RetryDelay, retryConfig.DLTMaxRetries, retryConfig.DLTRetryDelay)

	consumer := newConsumer(kafkaConfig, retryConfig, dlqProducer)

	// Set up topic subscriptions with handlers
	topics := map[string]func(msg *kafka.Message) error{
		"booking.created":   BookingCreateHandler(bookingUsecase),
		"booking.updated":   BookingUpdateHandler(bookingUsecase),
		"booking.cancelled": BookingCancelHandler(bookingUsecase),
		"user.created":      UserCreateHandler(userUsecase),
		"user.updated":      UserUpdateHandler(userUsecase),
		"unit.created":      UnitCreateHandler(unitUsecase),
	}

	// Subscribe to all topics
	for topic, handler := range topics {
		log.Printf("Setting up handler for topic: %s", topic)
		if err := consumer.Subscribe(topic, handler); err != nil {
			log.Printf("Failed to subscribe to topic %s: %v", topic, err)
		} else {
			log.Printf("Successfully subscribed to topic: %s", topic)
		}
	}

	topicsWithAck := map[string]func(msg *kafka.Message) (bool, error){
		"user.created": UserCreateHandlerWithAck(userUsecase),
	}

	for topic, handler := range topicsWithAck {
		log.Printf("Setting up handler with ack and retry for topic: %s", topic)
		if err := consumer.SubscribeWithAck(topic, handler); err != nil {
			log.Printf("Failed to subscribe to topic %s: %v", topic, err)
		} else {
			log.Printf("Successfully subscribed to topic: %s (with ack and retry)", topic)
		}
	}

	log.Printf("Kafka consumer initialized with %d topic handlers", len(topics)+len(topicsWithAck))
	return consumer
}
