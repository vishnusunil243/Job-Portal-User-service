package kafka

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ProduceShortlistUserMessage(email, company, userId, jobId string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}
	defer p.Close()

	topic := "ShortlistUser"
	message := fmt.Sprintf(`{"Email": "%s", "UserID": "%s", "JobID": "%s", "Company": "%s"}`, email, userId, jobId, company)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	p.Flush(15 * 1000)

	return nil
}
