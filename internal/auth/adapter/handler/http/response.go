package http

// authResponse represents an authentication Response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// newAuthResponse is a helper function to create a Response body for handling authentication data
func newAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}
