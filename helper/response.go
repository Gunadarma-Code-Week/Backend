package helper

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func CreateConflictResponse(message string) Response {
	return Response{
		Message: message,
		Data:    "nil",
		Errors:  "nil",
	}
}

func CreateNotFoundResponse(message string) Response {
	return Response{
		Message: message,
		Data:    "nil",
		Errors:  "nil",
	}
}

var NotFoundResponse Response = CreateNotFoundResponse("Data yang dicari tidak ditemukan")

func standardizeData(d interface{}) interface{} {
	if d == nil {
		return "nil"
	}
	return d
}

// CreateSuccessResponse returns a 201 response with data and nil errors.
func CreateSuccessResponse(message string, data interface{}) Response {
	return Response{
		Message: message,
		Data:    standardizeData(data),
		Errors:  "nil",
	}
}

// CreateMutationResponse returns a 201 response with data and nil errors.
func CreateMutationResponse(message string, data interface{}) Response {
	return Response{
		Message: message,
		Data:    standardizeData(data),
		Errors:  "nil",
	}
}

// CreateErrorResponse returns a 400 response.
// errors must be a map[string][]string (dictionary per AGENTS.md).
func CreateErrorResponse(message string, errors interface{}) Response {
	e := errors
	if e == nil {
		e = "nil"
	}
	return Response{
		Message: message,
		Data:    "nil",
		Errors:  e,
	}
}

// CreateInternalErrorResponse returns a 500 response with nil data and nil errors.
func CreateInternalErrorResponse(message string) Response {
	return Response{
		Message: message,
		Data:    "nil",
		Errors:  "nil",
	}
}

