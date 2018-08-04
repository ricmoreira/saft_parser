package services

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"saft_parser/config"
	"saft_parser/models/response"
	"saft_parser/util"
)

type (
	// SAFTService represents the service for operating on SAFT files
	SAFTService struct {
		config *config.Config
		saftParser *SAFTParser
	}
)

// NewSAFTService is the constructor of SAFTService
func NewSAFTService(config *config.Config, saftParser *SAFTParser) *SAFTService {
	return &SAFTService{
		config: config,
		saftParser: saftParser,
	}
}

// UploadAction stores a file in disk
// TODO: a better file storage solution should be implemented
func (this SAFTService) UploadAction(fileHeader *multipart.FileHeader) (*mresponse.SAFTUpload, *mresponse.ErrorResponse) {

	file, err := fileHeader.Open()
	if err != nil {
		return nil, util.HandleErrorResponse(util.INVALID_REQUEST, nil, err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, util.HandleErrorResponse(util.INVALID_REQUEST, nil, err.Error())
	}

	err = ioutil.WriteFile(this.config.SaftFilesFolder+string(os.PathSeparator)+fileHeader.Filename, data, 0666)
	if err != nil {
		return nil, util.HandleErrorResponse(util.INVALID_REQUEST, nil, err.Error())
	}

	this.saftParser.ParseFile(fileHeader.Filename)
	
	r := &mresponse.SAFTUpload{
		FilesUploaded: []string{fileHeader.Filename},
	}

	return r, nil
}
