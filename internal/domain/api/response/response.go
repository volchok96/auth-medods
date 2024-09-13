package response

type UserResponse struct {
	AccessToken  string `json:"access_token"`
	GetRefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
	GUID         string `json:"guid"`
}