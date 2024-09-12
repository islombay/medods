package model

type LoginRequest struct {
	UserId string `json:"user_id" binding:"required"`
	IP     string `json:"-"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Register struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	IP        string `json:"-"`
}
