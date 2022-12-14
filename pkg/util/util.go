package util

import (
	"encoding/json"
	"errors"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"net/http"
	"spreewill-core/pkg/models"
	"spreewill-core/pkg/services/aws"
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

func GetUserIdForFromAccessToken(w http.ResponseWriter, r *http.Request) (*string, error) {
	accessToken := GetHeaderToken(w, r)

	var emptyString *string

	if accessToken == "" {
		return emptyString, errors.New("invalid authorization header")
	}

	// TODO: verify that the userID exists in cognito
	cognitoClient, ok := r.Context().Value("CognitoClient").(*aws.CognitoClient)
	if !ok {
		return emptyString, errors.New("could not retrieve cognitoClient from context")
	}

	getUserInput := &cip.GetUserInput{
		AccessToken: &accessToken,
	}

	output, err := cognitoClient.GetUser(r.Context(), getUserInput)
	if err != nil {
		return emptyString, err
	}

	return output.Username, nil
}
