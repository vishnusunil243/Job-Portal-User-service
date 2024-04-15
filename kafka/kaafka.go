package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type FixedPartitioner struct{}

func (p *FixedPartitioner) Partition(msg *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	return 0, nil
}
func (p *FixedPartitioner) RequiresConsistency() bool {
	return true
}
func NewFixedPartitioner() sarama.PartitionerConstructor {
	return func(topic string) sarama.Partitioner {
		return &FixedPartitioner{}
	}
}

func ProduceShortlistUserMessage(email, company, designation, jobId string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewFixedPartitioner()

	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}
	defer producer.Close()

	topic := "ShortlistUser"
	msg := fmt.Sprintf(`{"Email": "%s", "Designation": "%s", "JobID": "%s", "Company": "%s"}`, email, designation, jobId, company)
	message := []byte(msg)
	p, o, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
		Partition: 0,
	})
	fmt.Println("partition ", p, "offset ", o)
	fmt.Println("message sent", msg)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil
}
func InterviewScheduledMessage(email, company, date, designation, roomId string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewFixedPartitioner()
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}
	defer producer.Close()

	topic := "InterviewSchedule"
	msg := fmt.Sprintf(`{"Email": "%s", "Designation": "%s", "Date": "%s", "Company": "%s","RoomId":"%s"}`, email, designation, date, company, roomId)
	message := []byte(msg)
	p, o, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
		Partition: 0,
	})
	fmt.Println("partition ", p, "offset ", o)
	fmt.Println("message sent", msg)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil
}
func HiredMessage(email, company, designation, interviewDate string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewFixedPartitioner()
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}
	defer producer.Close()
	topic := "Hired"
	msg := fmt.Sprintf(`{"Email":"%s","Designation":"%s","Date":"%s","Company":"%s"}`, email, designation, interviewDate, company)
	message := []byte(msg)
	p, o, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
		Partition: 0,
	})
	fmt.Println("partition ", p, "offset ", o)
	fmt.Println("message sent", msg)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil
}
func WarningEmail(userName, company, designation, interviewDate, companyEmail, userEmail, roomId string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewFixedPartitioner()
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}
	defer producer.Close()
	topic := "Warning"
	msg := fmt.Sprintf(`{"UserName":"%s","Designation":"%s","Date":"%s","Company":"%s","CompanyEmail":"%s","UserEmail":"%s","RoomId":"%s"}`, userName, designation, interviewDate, company, companyEmail, userEmail, roomId)
	message := []byte(msg)
	p, o, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
		Partition: 0,
	})
	fmt.Println("partition ", p, "offset ", o)
	fmt.Println("message sent", msg)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
		return err
	}

	return nil
}
