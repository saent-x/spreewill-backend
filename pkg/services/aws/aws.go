package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"os"
)

type CognitoClient struct {
	AppClientID string
	UserPoolID  string
	*cip.Client
}

func Init() *CognitoClient {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		panic(err)
	}

	return &CognitoClient{
		os.Getenv("COGNITO_APP_CLIENT_ID"),
		os.Getenv("COGNITO_USER_POOL_ID"),
		cip.NewFromConfig(cfg),
	}
}
