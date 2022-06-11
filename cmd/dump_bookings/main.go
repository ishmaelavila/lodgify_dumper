package main

import (
	"encoding/json"
	"os"

	"github.com/ishmaelavila/lodgify_dumper/internal/lodgify"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	if err != nil {
		sugar.Fatal("Error loading .env file")
	}

	lodgifyAPIKey := os.Getenv("LODGIFY_API_KEY")

	if lodgifyAPIKey == "" {
		sugar.Fatal("LODGIFY_API_KEY must not be empty")
	}
	lodgifyBaseURL := os.Getenv("LODGIFY_BASE_URL")

	if lodgifyBaseURL == "" {
		sugar.Fatal("LODGIFY_BASE_URL must not be empty")
	}

	args := lodgify.LodgifyClientArgs{
		BaseURL: lodgifyBaseURL,
		APIKey:  lodgifyAPIKey,
		Logger:  sugar,
	}

	sugar.Info("creating lodgify client")
	client, err := lodgify.NewClient(args)

	if err != nil {
		sugar.Fatalf("error initializing lodgify client %w", err)
	}

	bookings, err := client.GetBookings()

	if err != nil {
		sugar.Fatalf("error retreiving bookings from lodgify: %v", err)
	}

	j, _ := json.Marshal(bookings)

	sugar.Infof("%s", string(j))

}
