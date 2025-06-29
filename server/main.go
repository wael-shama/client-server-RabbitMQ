package main

import (
	"encoding/json"
	"log"
	"orderedmap/internal/logger"
	"orderedmap/internal/orderedmap"
	"orderedmap/internal/queue"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

type Command struct {
	Action string `json:"action"`
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
}

func main() {
	// Get RabbitMQ URL from environment
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/" // Default for local dev
	}

	rmq, err := queue.NewRabbitMQ(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	om := orderedmap.NewOrderedMap()

	// Update logger paths
	getItemLogger := logger.NewSafeLogger("/logs/getItem.log")
	getAllItemsLogger := logger.NewSafeLogger("/logs/getAllItems.log")

	msgs, err := rmq.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Server started. Waiting for commands...")

	for {
		select {
		case d := <-msgs:
			go func(delivery amqp.Delivery) {
				var cmd Command
				if err := json.Unmarshal(delivery.Body, &cmd); err != nil {
					log.Printf("Error decoding command: %v", err)
					delivery.Nack(false, false)
					return
				}

				processCommand(om, &cmd, getItemLogger, getAllItemsLogger)
				delivery.Ack(false)
			}(d)
		case <-sigChan:
			log.Println("Shutting down server...")
			return
		}
	}
}

func processCommand(
	om *orderedmap.OrderedMap,
	cmd *Command,
	getItemLogger *logger.SafeLogger,
	getAllItemsLogger *logger.SafeLogger,
) {
	switch cmd.Action {
	case "addItem":
		om.Add(cmd.Key, cmd.Value)
		log.Printf("Added item: %s=%s", cmd.Key, cmd.Value)
	case "deleteItem":
		om.Delete(cmd.Key)
		log.Printf("Deleted item: %s", cmd.Key)
	case "getItem":
		if value, exists := om.Get(cmd.Key); exists {
			logEntry := map[string]string{"key": cmd.Key, "value": value}
			if jsonData, err := json.Marshal(logEntry); err == nil {
				getItemLogger.WriteString(string(jsonData))
			}
		}
	case "getAllItems":
		items := om.GetAll()
		if jsonData, err := json.Marshal(items); err == nil {
			getAllItemsLogger.WriteString(string(jsonData))
		}
	default:
		log.Printf("Unknown command: %s", cmd.Action)
	}
}
