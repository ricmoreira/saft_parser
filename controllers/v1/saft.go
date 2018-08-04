package controllers

import (
	"saft_parser/util"

	"saft_parser/services"

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

	// send products to Kafka
	resp, e := this.kafkaClient.SendProductsToTopic("products", auditFile.Products)
	if e != nil {
		c.JSON(e.HttpCode, e)
		return
	}

	c.JSON(200, resp)
}
