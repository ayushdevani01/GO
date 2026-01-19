package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {

	// 1. Configure the Writer (Producer)
	// We point it to ONE broker. It will discover the rest automatically!
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "critical-data",
		Balancer: &kafka.LeastBytes{}, // Distributes messages evenly across partitions
	}

	defer writer.Close()

	fmt.Println("Starting Producer...")

	// 2. Send 10 messages
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("Message-%d", i)

		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fmt.Sprintf("Key-%d", i)), // Keys determine which partition gets the data
				Value: []byte(msg),
			},
		)

		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

		fmt.Printf("Sent: %s\n", msg)
		time.Sleep(500 * time.Millisecond) // Slow down so we can see it happening
	}
}
