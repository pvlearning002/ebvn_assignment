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

func (self *shortenURLService) generateKey(ctx context.Context) (string, error) {
	key, err := self.genPass.GeneratePassword()
	if err != nil {
		return self.cfg.StringEmpty, err
	}
	url_check, _ := self.urlStorage.GetURL(ctx, key)
	for url_check != self.cfg.StringEmpty {
		key, err = self.genPass.GeneratePassword()
		if err != nil {
			return self.cfg.StringEmpty, err
		}
		url_check, _ = self.urlStorage.GetURL(ctx, key)
	}

	return key, nil
}

// ShortenURL shortens the provided original URL and returns the shortened URL.
func (self *shortenURLService) CreateShortenLink(ctx context.Context, url string) (string, error) {
	key, err := self.generateKey(ctx)
	if err != nil {
		return self.cfg.StringEmpty, err
	}

	err = self.urlStorage.StoreURL(ctx, key, url)
	if err != nil {
		return self.cfg.StringEmpty, err
	}
	return key, nil
}
