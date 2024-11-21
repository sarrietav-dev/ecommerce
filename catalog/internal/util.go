package internal

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(statusCode)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
		Title:  "Error",
		Detail: err.Error(),
		Status: fmt.Sprint(statusCode),
	}})
}
