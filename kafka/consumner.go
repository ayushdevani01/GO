package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	// 1. Configure the Reader (Consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "critical-data",
		GroupID:     "learning-group-1", // Identifies this specific consumer
		StartOffset: kafka.FirstOffset,  // Ask to read from the very beginning
	})

	defer reader.Close()

	fmt.Println("Starting Consumer...")

	// 2. Read messages forever
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("failed to read message:", err)
			break
		}

		// Print the message AND where it came from (Partition/Offset)
		fmt.Printf("Received: %s | Partition: %d | Offset: %d\n",
			string(m.Value), m.Partition, m.Offset)
	}
}
