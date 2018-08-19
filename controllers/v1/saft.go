package controllers

import (
	"saft_parser/models/response"
	"saft_parser/services"
	"saft_parser/util"
	"sync"

	"github.com/gin-gonic/gin"
)

type (
	// SAFTController represents the controller for operating on SAFT files
	SAFTController struct {
		saftFileTransfer *services.SAFTFileTransfer
		saftParser       *services.SAFTParser
		kafkaClient      *services.KafkaProducer
	}
)

// NewSAFTController is the constructor of SAFTController
func NewSAFTController(saftFileTransfer *services.SAFTFileTransfer,
	saftParser *services.SAFTParser,
	kafkaClient *services.KafkaProducer,
) (*SAFTController, error) {

	e := kafkaClient.Connect()

	if e != nil {
		return nil, e
	}

	inst := &SAFTController{
		saftFileTransfer: saftFileTransfer,
		saftParser:       saftParser,
		kafkaClient:      kafkaClient,
	}

	return inst, nil
}

// CreateAction creates a new role
func (this SAFTController) FileToKafkaAction(c *gin.Context) {

	fileHeader, err := c.FormFile("SAFT-XML")
	if err != nil {
		e := util.HandleErrorResponse(util.INVALID_REQUEST, nil, err.Error())
		c.JSON(e.HttpCode, e)
		return
	}

	// save file in disk
	res, e := this.saftFileTransfer.UploadSAFTFile(fileHeader)
	if e != nil {
		c.JSON(e.HttpCode, e)
		return
	}

	// for now, only allows one file upload per request
	fileName := res.FilesUploaded[0]

	// parse saft file sent
	auditFile, err := this.saftParser.ParseFile(fileName)
	if err != nil {
		e := util.HandleErrorResponse(util.INVALID_REQUEST, nil, err.Error())
		c.JSON(e.HttpCode, e)
		return
	}

	// call Kafka producer to send messages
	var wg sync.WaitGroup
	wg.Add(3)

	cProducts := make(chan *mresponse.FileToKafkaProducts, 1)
	cInvoices := make(chan *mresponse.FileToKafkaInvoices, 1)
	cHeader := make(chan *mresponse.FileToKafkaHeader, 1)

	var productsResult *mresponse.FileToKafkaProducts
	var invoicesResult *mresponse.FileToKafkaInvoices
	var headerResult *mresponse.FileToKafkaHeader

	// send products to Kafka
	go func() {
		defer wg.Done()
		cProducts <- this.kafkaClient.SendProductsToTopic("products", auditFile.MasterFiles.Products)
		productsResult = <-cProducts
	}()

	// send invoices to Kafka
	go func() {
		defer wg.Done()
		cInvoices <- this.kafkaClient.SendInvoicesToTopic("invoices", auditFile.SourceDocuments.SalesInvoices.Invoices)
		invoicesResult = <-cInvoices
	}()

	// send header to Kafka
	go func() {
		defer wg.Done()
		cHeader <- this.kafkaClient.SendHeaderToTopic("header", auditFile.Header)
		headerResult = <-cHeader
	}()

	// wait for all results
	wg.Wait()

	// send combined result
	resp := mresponse.FileToKafka{
		Products: productsResult,
		Invoices: invoicesResult,
		Header: headerResult,
	}

	c.JSON(200, resp)
}
