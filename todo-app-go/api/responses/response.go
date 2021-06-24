package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorType struct {
	Detail string `json:"detail"`
	Status string `json:"status"`
}

func JsonResponse(writer http.ResponseWriter, statusCode int, responseData interface{}) {
	writer.WriteHeader(statusCode)
	encodingError := json.NewEncoder(writer).Encode(responseData)
	if encodingError != nil {
		fmt.Fprintf(writer, "%s", encodingError.Error())
	}

}

func ErrorResponse(writer http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JsonResponse(writer, statusCode, ErrorType{Detail: err.Error(), Status: "failed"})
	} else {
		JsonResponse(writer, statusCode, nil)
	}

}
