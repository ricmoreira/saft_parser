package mresponse

type (
	// ErrorResponse represents the error model for App error response messages
	SAFTUpload struct {
		FilesUploaded []string `json:"files_uploaded"`
	}

	FileToKafka struct {
		ProductsCount int      `json:"products_count"`
		ProductsCodes []string `json:"products_codes"`
	}
)
