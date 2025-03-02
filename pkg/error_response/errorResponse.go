package errorresponse

type ErrorResponse struct {
	Status      int32         `json:"status"`
	ErrorDetail []ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Description string `json:"description"`
	FieldName   string `json:"fieldName"`
}

type CustomError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
