package responsedata

import ()

type AuthResponseData struct {
	AccessToken string
	TokenType   string
}

type CreateSessionResponseData struct {
	SessionId string
}
