package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SmoothWay/booking/internal/config"
	"github.com/SmoothWay/booking/internal/infrastructure/broker/kafka"
	"github.com/SmoothWay/booking/internal/infrastructure/cache"
	"github.com/SmoothWay/booking/internal/infrastructure/delivery/http"
	"github.com/SmoothWay/booking/internal/infrastructure/repository/postgres"
	"github.com/SmoothWay/booking/internal/usecase"
	confluentkafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	redisCache := cache.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	defer redisCache.Close()

	bookingRepo := postgres.NewBookingRepository(db)
	userRepo := postgres.NewUserRepository(db)
	unitRepo := postgres.NewUnitRepository(db, redisCache)

	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	unitUsecase := usecase.NewUnitUsecase(unitRepo)

	// TODO: wrap usecases into domain.Usecase interface maybe?
	kafkaConsumer := kafka.InitKafkaConsumer(cfg, bookingUsecase, userUsecase, unitUsecase)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start Kafka consumer in a goroutine
	go func() {
		if err := kafkaConsumer.ConsumeLoop(ctx, func(msg *confluentkafka.Message) error {
			// This will be handled by the specific handlers
			return nil
		}); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	server := http.NewServer(cfg.HTTP, bookingUsecase, userUsecase, unitUsecase)
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down gracefully...")
	cancel()
}
