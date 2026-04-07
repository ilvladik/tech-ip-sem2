package usecases

type LoginInput struct {
	Login    string
	Password string
}

type LoginOutput struct {
	AccessToken string
	TokenType   string
}

type VerifyInput struct {
	Token string
}

type VerifyOutput struct {
	Subject string
}
