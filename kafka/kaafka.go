package kafka

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func ProduceShortlistUserMessage(email, company, userId, jobId string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}
	defer producer.Close()

	topic := "ShortlistUser"
	message := []byte(fmt.Sprintf(`{"Email": "%s", "UserID": "%s", "JobID": "%s", "Company": "%s"}`, email, userId, jobId, company))
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Partition: 0,
		Value:     sarama.StringEncoder(message),
	})
	fmt.Println("hiiii")
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil
}
