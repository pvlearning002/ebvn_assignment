package model

type ShortenLinkInput struct {
	ExpirationTime uint64 `json:"expiration_time" default:"3600" example:"3600"  binding:"required,gte=500"`
	Url            string `json:"url" default:"https://www.google.com/" example:"https://www.google.com/"  binding:"required,url"`
}
