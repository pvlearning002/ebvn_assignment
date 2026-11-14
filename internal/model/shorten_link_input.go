package model

type ShortenLinkInput struct {
	ExpirationTime int    `json:"expiration_time" default:"3600" example:"3600"`
	Url            string `json:"url" default:"https://www.google.com/" example:"https://www.google.com/"`
}
