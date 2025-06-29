package main

import (
	"bufio"
	"flag"
	"log"
	"orderedmap/internal/queue"
	"os"
	"strings"
)

type Command struct {
	Action string `json:"action"`
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
}

func main() {
	filePath := flag.String("file", "", "Path to commands file")
	flag.Parse()

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

	if *filePath != "" {
		processFile(*filePath, rmq)
	} else {
		processStdin(rmq)
	}
}

func processFile(filePath string, rmq *queue.RabbitMQ) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processLine(scanner.Text(), rmq)
	}
}

func processStdin(rmq *queue.RabbitMQ) {
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Enter commands (addItem key value, deleteItem key, getItem key, getAllItems):")
	for scanner.Scan() {
		processLine(scanner.Text(), rmq)
	}
}

func processLine(line string, rmq *queue.RabbitMQ) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 1 {
		return
	}

	cmd := Command{Action: parts[0]}
	switch cmd.Action {
	case "addItem":
		if len(parts) != 3 {
			log.Println("Invalid addItem command. Usage: addItem key value")
			return
		}
		cmd.Key = parts[1]
		cmd.Value = parts[2]
	case "deleteItem", "getItem":
		if len(parts) != 2 {
			log.Printf("Invalid %s command. Usage: %s key", cmd.Action, cmd.Action)
			return
		}
		cmd.Key = parts[1]
	case "getAllItems":
		// No arguments needed
	default:
		log.Printf("Unknown command: %s", cmd.Action)
		return
	}

	if err := rmq.Publish(cmd); err != nil {
		log.Printf("Failed to publish command: %v", err)
	} else {
		log.Printf("Sent command: %s", cmd.Action)
	}
}
