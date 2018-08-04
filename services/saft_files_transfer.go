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
	// SAFTFileTransfer represents the service for transfering SAFT files
	SAFTFileTransfer struct {
		config *config.Config
	}
)

// NewSAFTService is the constructor of SAFTService
func NewSAFTFileTransfer (config *config.Config) *SAFTFileTransfer {
	return &SAFTFileTransfer{
		config: config,
	}
}

// UploadSAFTFile stores a file in disk
// TODO: a better file storage solution should be implemented
func (this SAFTFileTransfer) UploadSAFTFile(fileHeader *multipart.FileHeader) (*mresponse.SAFTUpload, *mresponse.ErrorResponse) {

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

	r := &mresponse.SAFTUpload{
		FilesUploaded: []string{fileHeader.Filename},
	}

	return r, nil
}
