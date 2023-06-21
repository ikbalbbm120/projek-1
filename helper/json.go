package helper

import (
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder	:= json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func WriteToResponseBody(writer *http.Request, result interface{}) {
	writer.Header.Add("content-type", "application/json")
	encoder := json.NewDecoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err)
}