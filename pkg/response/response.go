package response

// Pagination struct represents pagination metadata.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PaginationResponse struct is a struct to define pagination response.
type PaginationResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Pagination *Pagination `json:"pagination"`
	Data       any         `json:"data"`
}

// DataResponse struct is a struct to define data response.
type DataResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// MessageResponse struct is a struct to define message response.
type MessageResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
