package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"math/rand"
)

const orangesTopic = "oranges"

func main() {
	ctx := context.Background()

	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", orangesTopic, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to dial kafka")
	}
	defer conn.Close()

	for {
		orangeSize := rand.Intn(401) + 100
		message := fmt.Sprintf("%d", orangeSize)
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(message)},
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write kafka messages")
		}
	}
}
