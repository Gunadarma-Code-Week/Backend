package helper

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type MutationResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

type ConflictResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateConflictResponse(message string) ConflictResponse {
	return ConflictResponse{
		Message: message,
		Data:    "nil",
	}
}

func CreateNotFoundResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Data:    "nil",
		Errors:  "nil",
	}
}

var NotFoundResponse ErrorResponse = CreateNotFoundResponse("Data yang dicari tidak ditemukan")

func CreateSuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func CreateMutationResponse(message string, data interface{}) MutationResponse {
	return MutationResponse{
		Message: message,
		Data:    data,
	}
}

func CreateErrorResponse(message string, errors interface{}) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Data:    "nil",
		Errors:  errors,
	}
}
