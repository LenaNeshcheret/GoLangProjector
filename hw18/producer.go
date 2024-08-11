package main

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"hw18/common"
	"math/rand"
)

func main() {
	ctx := context.Background()

	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", common.OrangesTopic, 0)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to dial kafka")
	}
	defer func(conn *kafka.Conn) {
		err := conn.Close()
		if err != nil {
			
		}
	}(conn)

	for {
		orange := common.Orange{
			Size: int32(rand.Intn(401) + 100),
		}
		message, err := json.Marshal(orange)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to marshal orange to JSON")
		}
		_, err = conn.WriteMessages(
			kafka.Message{Value: message},
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write kafka messages")
		}
	}
}
