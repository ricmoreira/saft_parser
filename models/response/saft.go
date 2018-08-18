package mresponse

type (
	// ErrorResponse represents the error model for App error response messages
	SAFTUpload struct {
		FilesUploaded []string `json:"files_uploaded"`
	}

	FileToKafkaProducts struct {
		ProductsCount int            `json:"products_count,omitempty"`
		ProductsCodes []string       `json:"products_codes,omitempty"`
		Error         *ErrorResponse `json:"error,omitempty"`
	}

	FileToKafkaInvoices struct {
		InvoicesCount int            `json:"invoices_count,omitempty"`
		InvoicesCodes []string       `json:"invoices_codes,omitempty"`
		Error         *ErrorResponse `json:"error,omitempty"`
	}

	FileToKafka struct {
		Products *FileToKafkaProducts `json:"products"`
		Invoices *FileToKafkaInvoices `json:"invoices"`
	}
)
