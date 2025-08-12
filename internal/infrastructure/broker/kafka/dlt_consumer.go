package kafka

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/SmoothWay/booking/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// DLTConsumer handles messages from Dead Letter Topics
type DLTConsumer struct {
	consumer    *kafka.Consumer
	retryConfig RetryConfig
	handlers    map[string]MessageHandlerWithAck
}

// NewDLTConsumer creates a new DLT consumer
func NewDLTConsumer(cfg *config.Config, retryConfig RetryConfig) *DLTConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(cfg.Kafka.Brokers, ","),
		"group.id":           cfg.Kafka.GroupID + "-dlt",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Fatalf("Failed to create DLT consumer: %v", err)
	}

	return &DLTConsumer{
		consumer:    consumer,
		retryConfig: retryConfig,
		handlers:    make(map[string]MessageHandlerWithAck),
	}
}

// SubscribeToDLT subscribes to a Dead Letter Topic
func (d *DLTConsumer) SubscribeToDLT(originalTopic string, handler MessageHandlerWithAck) error {
	dltTopic := originalTopic + ".dlt"
	log.Printf("Subscribing to DLT topic: %s", dltTopic)

	d.handlers[dltTopic] = handler
	return d.consumer.Subscribe(dltTopic, nil)
}

// ConsumeDLTLoop processes messages from DLT topics
func (d *DLTConsumer) ConsumeDLTLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			d.consumer.Close()
			return ctx.Err()
		default:
			msg, err := d.consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Received DLT message on topic %s: %s", *msg.TopicPartition.Topic, string(msg.Value))

				if handler, exists := d.handlers[*msg.TopicPartition.Topic]; exists {
					if err := d.processDLTMessage(*msg.TopicPartition.Topic, msg, handler); err != nil {
						log.Printf("Error processing DLT message: %v", err)
						continue
					}
				} else {
					log.Printf("No handler found for DLT topic: %s", *msg.TopicPartition.Topic)
				}
			} else {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
					continue
				}
				log.Printf("Error reading DLT message: %v", err)
			}
		}
	}
}

// processDLTMessage handles DLT message processing with additional retries
func (d *DLTConsumer) processDLTMessage(topic string, msg *kafka.Message, handler MessageHandlerWithAck) error {
	// Extract DLT retry count
	dltRetryCount := d.getDLTRetryCount(msg)

	// Check if we've exceeded DLT max retries
	if dltRetryCount >= d.retryConfig.DLTMaxRetries {
		log.Printf("DLT max retries (%d) exceeded for topic %s, message permanently failed",
			d.retryConfig.DLTMaxRetries, topic)

		// Log the permanent failure and acknowledge to remove from DLT
		log.Printf("Permanently failed message: %s", string(msg.Value))
		_, err := d.consumer.CommitMessage(msg)
		return err
	}

	// Process the message
	shouldAck, err := handler(msg)
	if err != nil {
		log.Printf("DLT handler error for topic %s (attempt %d/%d): %v",
			topic, dltRetryCount+1, d.retryConfig.DLTMaxRetries, err)

		// Schedule DLT retry
		return d.scheduleDLTRetry(topic, msg, dltRetryCount)
	}

	if shouldAck {
		// Successfully processed - acknowledge
		_, err := d.consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Failed to commit DLT message: %v", err)
			return err
		}
		log.Printf("DLT message successfully processed and acknowledged for topic %s", topic)
	} else {
		log.Printf("DLT message NOT acknowledged for topic %s", topic)
	}

	return nil
}

// getDLTRetryCount extracts DLT retry count from message headers
func (d *DLTConsumer) getDLTRetryCount(msg *kafka.Message) int {
	for _, header := range msg.Headers {
		if string(header.Key) == "dlt_retry_count" {
			if count, err := strconv.Atoi(string(header.Value)); err == nil {
				return count
			}
		}
	}
	return 0
}

// scheduleDLTRetry schedules a DLT message for retry
func (d *DLTConsumer) scheduleDLTRetry(topic string, msg *kafka.Message, dltRetryCount int) error {
	// Calculate DLT retry delay with exponential backoff
	delay := d.retryConfig.DLTRetryDelay * time.Duration(1<<dltRetryCount)

	log.Printf("Scheduling DLT retry for topic %s in %v (attempt %d/%d)",
		topic, delay, dltRetryCount+1, d.retryConfig.DLTMaxRetries)

	// Add DLT retry count to message headers
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   "dlt_retry_count",
		Value: []byte(strconv.Itoa(dltRetryCount + 1)),
	})

	// For DLT, we'll just log and wait - in a real implementation, you might use a scheduler
	log.Printf("DLT retry scheduled for topic %s (would wait %v in production)", topic, delay)

	// Don't acknowledge - will be retried
	return nil
}

// InitDLTConsumer initializes the DLT consumer with handlers
func InitDLTConsumer(cfg *config.Config, userUsecase domain.UserUsecase) *DLTConsumer {
	// Configure DLT retry behavior
	retryConfig := RetryConfig{
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		DLTMaxRetries: 2,                // Additional retries in DLT
		DLTRetryDelay: 30 * time.Second, // Longer delay in DLT
	}

	log.Printf("Initializing DLT consumer with DLT max retries: %d, DLT delay: %v",
		retryConfig.DLTMaxRetries, retryConfig.DLTRetryDelay)

	dltConsumer := NewDLTConsumer(cfg, retryConfig)

	// Set up DLT handlers
	dltHandlers := map[string]MessageHandlerWithAck{
		"user.created": UserCreateHandlerWithAck(userUsecase),
	}

	for topic, handler := range dltHandlers {
		if err := dltConsumer.SubscribeToDLT(topic, handler); err != nil {
			log.Printf("Failed to subscribe to DLT topic %s: %v", topic, err)
		} else {
			log.Printf("Successfully subscribed to DLT topic: %s.dlt", topic)
		}
	}

	log.Printf("DLT consumer initialized with %d topic handlers", len(dltHandlers))
	return dltConsumer
}
