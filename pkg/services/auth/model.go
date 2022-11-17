package auth

type (
	SignUpRequest struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	OTPInfo struct {
		OTP   string `json:"otp"`
		Email string `json:"email"`
	}

	SignInResponse struct {
		AccessToken  *string `json:"access_token"`
		ExpiresIn    int32   `json:"expires_in"`
		IdToken      *string `json:"id_token"`
		RefreshToken *string `json:"refresh_token"`
		TokenType    *string `json:"token_type"`
	}

	ForgotPasswordObject struct {
		Email string `json:"email"`
	}

	ConfirmForgotPasswordObject struct {
		ConfirmationCode string `json:"confirmation_code"`
		Username         string `json:"username"`
		Password         string `json:"password"`
	}
)
