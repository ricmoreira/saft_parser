package services

import (
	"encoding/json"
	"fmt"
	"saft_parser/config"
	"saft_parser/models/response"
	"saft_parser/models/saft-pt-4"
	"saft_parser/util"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer allows sending messages to Kafka broker
type KafkaProducer struct {
	prod   *kafka.Producer
	config *config.Config
}

// NewKafkaProducer is the KafkaProducer constructor
func NewKafkaProducer(config *config.Config) *KafkaProducer {
	return &KafkaProducer{
		config: config,
	}
}

// Connect connects instance to Kafka server enabling instance to send messages
func (kp *KafkaProducer) Connect() error {
	config := &kafka.ConfigMap{
		"bootstrap.servers":    kp.config.BootstrapServers,
		"client.id":            "campaignservice",
		"default.topic.config": kafka.ConfigMap{"acks": "all"},
		"message.max.bytes":    kp.config.MessageMaxBytes,
	}

	p, err := kafka.NewProducer(config)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return err
	}

	kp.prod = p
	return nil
}

// SendProductsToTopic sends an array of *msaft.Product to provided topic name
func (kp *KafkaProducer) SendProductsToTopic(topic string, request []*msaft.Product) (*mresponse.FileToKafka, *mresponse.ErrorResponse) {

	jptBytes, err := json.Marshal(request)
	if err != nil {
		return nil, util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, err.Error())
	}

	// don't allow to process more than max size allowed of data
	if len(jptBytes) > kp.config.MessageMaxBytes {
		s := strconv.Itoa(kp.config.MessageMaxBytes)
		return nil, util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, "To many products to processs. Max allowed: " + s + " bytes of data.")
	}

	e := kp.sendMessage(topic, jptBytes)

	if e != nil {
		return nil, e
	}

	resp := &mresponse.FileToKafka{
		ProductsCount: len(request),
		ProductsCodes: reduceToProductsCodes(request),
	}

	return resp, nil
}

func reduceToProductsCodes(products []*msaft.Product) []string {
	res := make([]string, len(products))

	for i, product := range products {
		res[i] = product.ProductCode
	}

	return res
}

// SendMessage sends a message to Kafka to provided topic
// errors.HandleErrorResponse(500, errors.SERVICE_UNAVAILABLE, errors.DefaultErrorsMessages[errors.SERVICE_UNAVAILABLE], nil)
func (kp *KafkaProducer) sendMessage(topic string, message []byte) *mresponse.ErrorResponse {

	deliveryChan := make(chan kafka.Event, 10000)

	err := kp.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message},
		deliveryChan,
	)

	if err != nil {
		return util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, err.Error())
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		close(deliveryChan)
		return util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, m.TopicPartition.Error.Error())
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		close(deliveryChan)
		return nil
	}
}