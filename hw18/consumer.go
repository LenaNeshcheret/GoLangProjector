package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   orangesTopic,
	})

	var small, medium, large int

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Oranges: small=%d, medium=%d, large=%d\n", small, medium, large)
				small, medium, large = 0, 0, 0
			}
		}
	}()

	for {
		message, err := kafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read message")
		}

		var orangeSize int32
		if err := json.Unmarshal(message.Value, &orangeSize); err != nil {
			log.Printf("Failed to decode message: ", err)
			continue
		}

		switch {
		case orangeSize < 250:
			small++
		case orangeSize <= 400:
			medium++
		default:
			large++
		}
	}
}
