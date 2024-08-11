package main

import (
	"context"
	"encoding/json"
	"hw18/common"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   common.OrangesTopic,
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

		var orange common.Orange
		if err := json.Unmarshal(message.Value, &orange); err != nil {
			log.Warn().Err(err).Msg("Failed to decode message")
			continue
		}

		switch {
		case orange.Size < 250:
			small++
		case orange.Size <= 400:
			medium++
		default:
			large++
		}
	}
}
