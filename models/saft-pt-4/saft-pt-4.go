package msaft

import (
	"encoding/xml"
)

type Product struct {
	XMLName            xml.Name        `xml:"Product"`
	ProductType        string          `json:"ProductType" xml:"ProductType"`
	ProductCode        string          `json:"ProductCode" xml:"ProductCode"`
	ProductGroup       string          `json:"ProductGroup" xml:"ProductGroup"`
	ProductDescription string          `json:"ProductDescription" xml:"ProductDescription"`
	ProductNumberCode  string          `json:"ProductNumberCode" xml:"ProductNumberCode"`
	CustomsDetails     *CustomsDetails `json:"CustomsDetails" xml:"CustomsDetails"`
}

type CustomsDetails struct {
	XMLName  xml.Name `xml:"CustomsDetails"`
	CNCode   string   `json:"CNCode" xml:"CNCode"`
	UNNumber string   `json:"UNNumber" xml:"UNNumber"`
}

type AuditFile struct {
	XMLName  xml.Name  `xml:"AuditFile"`
	Products []*Product `json:"Products" xml:"MasterFiles>Product"`
}
