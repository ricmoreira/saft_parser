package services

import (
	"encoding/json"
	"saft_parser/config"
	"saft_parser/models/response"
	"saft_parser/models/saft/go_SaftT104"
	"saft_parser/util"
	"strconv"
	"log"

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
		log.Printf("Failed to create producer: %s\n", err)
		return err
	}

	kp.prod = p
	return nil
}

// SendProductsToTopic sends an array of *go_SaftT104.TxsdProduct to provided topic name
func (kp *KafkaProducer) SendProductsToTopic(topic string, request []*go_SaftT104.TxsdProduct) *mresponse.FileToKafkaProducts {

	res := mresponse.FileToKafkaProducts {}

	jptBytes, err := json.Marshal(request)
	if err != nil {
		res.Error = util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, err.Error())
		return &res
	}

	// don't allow to process more than max size allowed of data
	if len(jptBytes) > kp.config.MessageMaxBytes {
		s := strconv.Itoa(kp.config.MessageMaxBytes)
		res.Error = util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, "To many products to processs. Max allowed: " + s + " bytes of data.")
		return &res
	}

	e := kp.sendMessage(topic, jptBytes)

	if e != nil {
		res.Error = e
		return &res
	}

	res.ProductsCount = len(request)
	res.ProductsCodes = reduceToProductsCodes(request)

	return &res
}


// SendInvoicesToTopic sends an array of *go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice to provided topic name
func (kp *KafkaProducer) SendInvoicesToTopic(topic string, request []*go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice) *mresponse.FileToKafkaInvoices {

	resp := mresponse.FileToKafkaInvoices{}

	jptBytes, err := json.Marshal(request)
	if err != nil {
		resp.Error = util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, err.Error())
		return &resp 
	}

	// don't allow to process more than max size allowed of data
	if len(jptBytes) > kp.config.MessageMaxBytes {
		s := strconv.Itoa(kp.config.MessageMaxBytes)
		resp.Error = util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, "To many invoices to processs. Max allowed: " + s + " bytes of data.")
		return &resp
	}

	e := kp.sendMessage(topic, jptBytes)

	if e != nil {
		resp.Error = e
		return &resp
	}

	resp.InvoicesCount = len(request)
	resp.InvoicesCodes = reduceToInvoicesCodes(request)

	return &resp
}

// SendHeaderToTopic sends the Header of SAFT-T file, containing company info to provided topic name
func (kp *KafkaProducer) SendHeaderToTopic(topic string, request *go_SaftT104.TxsdHeader) *mresponse.FileToKafkaHeader {

	resp := mresponse.FileToKafkaHeader{}

	jptBytes, err := json.Marshal(request)
	if err != nil {
		resp.Error = util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, err.Error())
		return &resp 
	}

	// don't allow to process more than max size allowed of data
	if len(jptBytes) > kp.config.MessageMaxBytes {
		s := strconv.Itoa(kp.config.MessageMaxBytes)
		resp.Error = util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, "To many invoices to processs. Max allowed: " + s + " bytes of data.")
		return &resp
	}

	e := kp.sendMessage(topic, jptBytes)

	if e != nil {
		resp.Error = e
		return &resp
	}

	return &resp
}

func reduceToProductsCodes(products []*go_SaftT104.TxsdProduct) []string {
	res := make([]string, len(products))

	for i, product := range products {
		res[i] = string(product.ProductCode)
	}

	return res
}

func reduceToInvoicesCodes(invoices []*go_SaftT104.TxsdSourceDocumentsSequenceSalesInvoicesSequenceInvoice) []string {
	res := make([]string, len(invoices))

	for i, invoice := range invoices {
		res[i] = string(invoice.InvoiceNo)
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
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		close(deliveryChan)
		return util.HandleErrorResponse(util.SERVICE_UNAVAILABLE, nil, m.TopicPartition.Error.Error())
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		close(deliveryChan)
		return nil
	}
}
