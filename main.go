package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/config"
	redisservice "github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/repository"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/rental"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/service/tracker"
	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/internal/transfer/rest/api"
)

const configPath = "internal/config/default.env"

func main() {
	cfg, err := config.NewConfig(context.Background(), configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("config retrieval failed: %w", err))
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting Scootin Aboot")

	validate := validator.New()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	initializeRedis(redisClient)

	redisService := redisservice.NewRedisService(redisClient)
	trackerService := tracker.NewTrackingService(logger, redisService)
	rentalService := rental.NewRentalService(redisService)

	router := mux.NewRouter()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP),
		Handler: router,
	}

	users := cfg.GetUsersMap()

	server := api.NewServer(logger, validate, httpServer, router, redisService, rentalService, trackerService, users)

	server.Run()

	//clients := []*client.Client{
	//	{
	//		ClientUUID: uuid.New(),
	//		Longitude:  73.4,
	//		Latitude:   45.4,
	//		City:       "Ottawa",
	//		Height:     50000.0,
	//		Width:      55000.0,
	//	},
	//	{
	//		ClientUUID: uuid.New(),
	//		Longitude:  73.5,
	//		Latitude:   45.5,
	//		City:       "Ottawa",
	//		Height:     50000.0,
	//		Width:      55000.0,
	//	},
	//	{
	//		ClientUUID: uuid.New(),
	//		Longitude:  73.6,
	//		Latitude:   45.6,
	//		City:       "Ottawa",
	//		Height:     50000.0,
	//		Width:      55000.0,
	//	},
	//}
}

func initializeRedis(redisClient *redis.Client) {
	if _, err := redisClient.GeoAdd(context.Background(), "Ottawa", &redis.GeoLocation{
		Name:      "0dae4f8c-dbbf-4bac-90f2-b80f07255ba5",
		Longitude: 73.5673,
		Latitude:  45.5017,
	}).Result(); err != nil {
		panic(err)
	}

	if err := redisClient.Set(context.Background(), "0dae4f8c-dbbf-4bac-90f2-b80f07255ba5", true, 0).Err(); err != nil {
		panic(err)
	}

	if _, err := redisClient.GeoAdd(context.Background(), "Ottawa", &redis.GeoLocation{
		Name:      "61637887-385e-47bd-ad8c-5ace4fbd2877",
		Longitude: 73.5548,
		Latitude:  45.5088,
	}).Result(); err != nil {
		panic(err)
	}

	if err := redisClient.Set(context.Background(), "61637887-385e-47bd-ad8c-5ace4fbd2877", true, 0).Err(); err != nil {
		panic(err)
	}

	if _, err := redisClient.GeoAdd(context.Background(), "Ottawa", &redis.GeoLocation{
		Name:      "4117b009-5e61-4b3a-aac5-c9d6a75483cb",
		Longitude: 73.5637,
		Latitude:  45.4724,
	}).Result(); err != nil {
		panic(err)
	}

	if err := redisClient.Set(context.Background(), "4117b009-5e61-4b3a-aac5-c9d6a75483cb", true, 0).Err(); err != nil {
		panic(err)
	}

	if _, err := redisClient.GeoAdd(context.Background(), "Montreal", &redis.GeoLocation{
		Name:      "bad9f260-e3f5-4375-a4b3-3f6e258eb21f",
		Longitude: 65.5637,
		Latitude:  30.5234,
	}).Result(); err != nil {
		panic(err)
	}

	if err := redisClient.Set(context.Background(), "bad9f260-e3f5-4375-a4b3-3f6e258eb21f", true, 0).Err(); err != nil {
		panic(err)
	}

	if _, err := redisClient.GeoAdd(context.Background(), "Montreal", &redis.GeoLocation{
		Name:      "32341255-c86a-4106-94e0-28dd9b3f88f2",
		Longitude: 65.1207,
		Latitude:  30.2827,
	}).Result(); err != nil {
		panic(err)
	}

	if err := redisClient.Set(context.Background(), "32341255-c86a-4106-94e0-28dd9b3f88f2", true, 0).Err(); err != nil {
		panic(err)
	}

	if _, err := redisClient.GeoAdd(context.Background(), "Montreal", &redis.GeoLocation{
		Name:      "b55fcd8c-383c-4169-9e4a-1c1bf15fdb76",
		Longitude: 65.5537,
		Latitude:  30.5234,
	}).Result(); err != nil {
		panic(err)
	}

	if err := redisClient.Set(context.Background(), "b55fcd8c-383c-4169-9e4a-1c1bf15fdb76", true, 0).Err(); err != nil {
		panic(err)
	}
}
