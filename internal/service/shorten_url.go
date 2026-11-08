package service

import (
	"context"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/repository"
)

// ShortenURL defines the interface for shortening URLs.
type ShortenURL interface {
	CreateShortenLink(ctx context.Context, url string) (string, error)
}

// shortenURLService is the implementation of the ShortenURLService interface.
type shortenURLService struct {
	urlStorage repository.URLStorage
	genPass    GenPass
	cfg        config.Config
}

// NewShortenURLService creates a new instance of the shortenURLService with the provided URL storage and configuration.
func NewShortenURLService(urlStorage repository.URLStorage, genPass GenPass, cfg config.Config) ShortenURL {
	return &shortenURLService{
		urlStorage: urlStorage,
		genPass:    genPass,
		cfg:        cfg,
	}
}

// ShortenURL shortens the provided original URL and returns the shortened URL.
func (self *shortenURLService) CreateShortenLink(ctx context.Context, url string) (string, error) {
	code, err := self.genPass.GeneratePassword()
	if err != nil {
		return self.cfg.StringEmpty, err
	}
	err = self.urlStorage.StoreURL(ctx, code, url)
	if err != nil {
		return self.cfg.StringEmpty, err
	}
	return code, nil
}
