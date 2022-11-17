package util

import (
	"encoding/json"
	"net/http"
	"spreewill-core/pkg/models"
	"strings"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetHeaderToken(w http.ResponseWriter, r *http.Request) string {
	authHeader := strings.Split(r.Header.Get("Authorization"), " ")

	if len(authHeader) != 2 {
		return ""
	}

	return authHeader[1]
}

func GetAllInCursor[T models.Entity](cur *mongo.Cursor, w http.ResponseWriter) {
	var arr []T

	err := cur.All(mgm.Ctx(), &arr)
	if err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	SendSuccess(w, http.StatusOK, arr)
}
