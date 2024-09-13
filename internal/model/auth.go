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
	Email     string `json:"email" binding:"required"`
	IP        string `json:"-"`
}

type RefreshRequest struct {
	AccessToken  string `form:"access_token" binding:"required" header:"access_token"`
	RefreshToken string `form:"refresh_token" binding:"required" header:"refresh_token"`
	IP           string `json:"-"`
}

type DeviceInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}
