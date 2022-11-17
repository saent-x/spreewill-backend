package auth

import (
	"encoding/json"
	"net/http"
	"spreewill-core/pkg/util"

	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go/aws"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cognitoClient, ok := r.Context().Value("CognitoClient").(*CognitoClient)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "could not retrieve cognitoClient from context")
		return
	}

	awsReq := &cip.AdminInitiateAuthInput{
		AuthFlow:   "ADMIN_USER_PASSWORD_AUTH",
		ClientId:   aws.String(cognitoClient.AppClientID),
		UserPoolId: aws.String(cognitoClient.UserPoolID),
		AuthParameters: map[string]string{
			"USERNAME": req.Email,
			"PASSWORD": req.Password,
		},
	}

	output, err := cognitoClient.AdminInitiateAuth(r.Context(), awsReq)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	res := &SignInResponse{
		AccessToken:  output.AuthenticationResult.AccessToken,
		ExpiresIn:    output.AuthenticationResult.ExpiresIn,
		IdToken:      output.AuthenticationResult.IdToken,
		RefreshToken: output.AuthenticationResult.RefreshToken,
		TokenType:    output.AuthenticationResult.TokenType,
	}

	util.SendSuccess(w, http.StatusOK, res)
	return
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cognitoClient, ok := r.Context().Value("CognitoClient").(*CognitoClient)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "could not retrieve cognitoClient from context")
		return
	}

	awsReq := &cip.SignUpInput{
		ClientId: aws.String(cognitoClient.AppClientID),
		Username: aws.String(req.Email),
		Password: aws.String(req.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(req.PhoneNumber),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(req.Email),
			},
		},
	}

	_, err = cognitoClient.SignUp(r.Context(), awsReq)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, "success")
	return
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	var req OTPInfo

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cognitoClient, ok := r.Context().Value("CognitoClient").(*CognitoClient)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "could not retrieve cognitoClient from context")
		return
	}

	confirmSignUp := &cip.ConfirmSignUpInput{
		ClientId:         aws.String(cognitoClient.AppClientID),
		Username:         aws.String(req.Email),
		ConfirmationCode: aws.String(req.OTP),
	}

	_, err = cognitoClient.ConfirmSignUp(r.Context(), confirmSignUp)

	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, "verified!")
	return
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordObject

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cognitoClient, ok := r.Context().Value("CognitoClient").(*CognitoClient)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "could not retrieve cognitoClient from context")
		return
	}

	cognitoUser := &cip.ForgotPasswordInput{
		ClientId: &cognitoClient.AppClientID,
		Username: &req.Email,
	}

	forgotPasswordOutput, err := cognitoClient.ForgotPassword(r.Context(), cognitoUser)

	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, forgotPasswordOutput)
	return
}

func ConfirmForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ConfirmForgotPasswordObject

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cognitoClient, ok := r.Context().Value("CognitoClient").(*CognitoClient)
	if !ok {
		util.SendError(w, http.StatusInternalServerError, "could not retrieve cognitoClient from context")
		return
	}

	cognitoUser := &cip.ConfirmForgotPasswordInput{
		ClientId:         &cognitoClient.AppClientID,
		Username:         &req.Username,
		ConfirmationCode: &req.ConfirmationCode,
		Password:         &req.Password,
	}

	resp, err := cognitoClient.ConfirmForgotPassword(r.Context(), cognitoUser)

	if err != nil {
		util.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.SendSuccess(w, http.StatusOK, resp)
	return
}

func Protected(w http.ResponseWriter, r *http.Request) {
	util.SendSuccess(w, http.StatusOK, "protected!!")
	return
}
