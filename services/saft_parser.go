package services

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"saft_parser/config"
	"saft_parser/models/saft/go_SaftT104"
)

type (
	// SAFTParser represents the service for parsing SAFT files
	SAFTParser struct {
		config *config.Config
	}
)

// SAFTParser is the constructor of SAFTParser
func NewSAFTParser(config *config.Config) *SAFTParser {
	return &SAFTParser{
		config: config,
	}
}

// ParseFile receives a xml file name located in SAFT_FILES_FOLDER and marshalls file data to an *msaft.AuditFile model
func (sp *SAFTParser) ParseFile(fileName string) (*go_SaftT104.TxsdAuditFile, error) {
	xmlFile, err := os.Open(config.MustGetEnv(config.SAFT_FILES_FOLDER) + string(os.PathSeparator) + fileName)

	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	v := &go_SaftT104.TxsdAuditFile{}
	if err := xml.Unmarshal(byteValue, v); err != nil {
		return nil, err
	}

	return v, nil
}
