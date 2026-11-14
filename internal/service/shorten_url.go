package service

import (
	"context"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
	"github.com/pvlearning002/ebvn_assignment/internal/repository"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

// ShortenURL defines the interface for shortening URLs.
type ShortenURL interface {
	CreateShortenLink(ctx context.Context, url string) (*model.ShortenLinkOutput, error)
	GetUrlByCode(ctx context.Context, keyCode string) (string, error)
}

// shortenURLService is the implementation of the ShortenURLService interface.
type shortenURLService struct {
	urlStorage repository.URLStorage
	genPass    GenPass
	cfg        config.Config
	cfgError   config.ConfigError
}

// NewShortenURLService creates a new instance of the shortenURLService with the provided URL storage and configuration.
func NewShortenURLService(urlStorage repository.URLStorage, genPass GenPass, cfg config.Config, cfgError config.ConfigError) ShortenURL {
	return &shortenURLService{
		urlStorage: urlStorage,
		genPass:    genPass,
		cfg:        cfg,
		cfgError:   cfgError,
	}
}

func (self *shortenURLService) generateKey(ctx context.Context) (string, error) {
	key, err := self.genPass.GenerateRandomString(self.cfg.KeycodeLengthDefault)
	if err != nil {
		utils.WriteLogError(err.Error())
		return self.cfg.StringEmpty, err
	}
	url_check, _ := self.urlStorage.GetURL(ctx, key)
	for url_check != self.cfg.StringEmpty {
		key, err = self.genPass.GenerateRandomString(self.cfg.KeycodeLengthDefault)
		if err != nil {
			utils.WriteLogError(err.Error())
			return self.cfg.StringEmpty, err
		}
		url_check, _ = self.urlStorage.GetURL(ctx, key)
	}

	return key, nil
}

// ShortenURL shortens the provided original URL and returns the shortened URL.
func (self *shortenURLService) CreateShortenLink(ctx context.Context, url string) (*model.ShortenLinkOutput, error) {
	key, err := self.generateKey(ctx)
	if err != nil {
		utils.WriteLogError(err.Error())
		return &model.ShortenLinkOutput{KeyCode: self.cfg.StringEmpty, Message: err.Error()}, err
	}

	err = self.urlStorage.StoreURL(ctx, key, url)
	if err != nil {
		utils.WriteLogError(err.Error())
		return &model.ShortenLinkOutput{KeyCode: self.cfg.StringEmpty, Message: err.Error()}, err
	}
	return &model.ShortenLinkOutput{KeyCode: key, Message: "Shorten URL generated successfully"}, nil
}

// ShortenURL shortens the provided original URL and returns the shortened URL.
func (self *shortenURLService) GetUrlByCode(ctx context.Context, keyCode string) (string, error) {
	url, err := self.urlStorage.GetURL(ctx, keyCode)
	if err != nil {
		utils.WriteLogError(err.Error())
		return self.cfg.StringEmpty, err
	}
	return url, nil
}
