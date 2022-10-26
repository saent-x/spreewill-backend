package util

import (
	"encoding/json"
	"net/http"
	"spreewill-core/pkg/models"
)

func SendError[T any](w http.ResponseWriter, status int, data T) {
	responseObject := models.ResponseObject{
		Message: "error",
		Status:  status,
		Data:    data,
	}

	sendJSON(w, status, &responseObject)
}

func SendSuccess[T any](w http.ResponseWriter, status int, data T) {
	responseObject := models.ResponseObject{
		Message: "success",
		Status:  status,
		Data:    data,
	}

	sendJSON(w, status, &responseObject)
}

func sendJSON[T any](w http.ResponseWriter, status int, resObj *T) {
	response, _ := json.Marshal(resObj)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
