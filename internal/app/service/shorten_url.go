package service

import (
	"context"

	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/app/repository"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_log"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_string"

	"errors"

	"github.com/redis/go-redis/v9"
)

// ShortenURL defines the interface for shortening URLs.
type ShortenURL interface {
	CreateShortenLink(ctx context.Context, url string) (*model.ShortenLinkOutput, error)
	GetUrlByCode(ctx context.Context, keyCode string) (string, error)
}

// shortenURLService is the implementation of the ShortenURLService interface.
type shortenURLService struct {
	urlStorage    repository.URLStorage
	cfg           config.Config
	cfgError      config.ConfigError
	logger        utils_log.CustomLogger
	custom_string utils_string.UtilsString
}

// NewShortenURLService creates a new instance of the shortenURLService with the provided URL storage and configuration.
func NewShortenURLService(urlStorage repository.URLStorage, cfg config.Config, cfgError config.ConfigError, logger utils_log.CustomLogger, custom_string utils_string.UtilsString) ShortenURL {
	return &shortenURLService{
		urlStorage:    urlStorage,
		cfg:           cfg,
		cfgError:      cfgError,
		logger:        logger,
		custom_string: custom_string,
	}
}

func (self *shortenURLService) generateKey(ctx context.Context) (string, error) {
	key, err := self.custom_string.GenerateRandomString(self.cfg.KeycodeLengthDefault)
	if err != nil {
		self.logger.WriteLogError(err.Error())
		return self.cfg.StringEmpty, err
	}
	url_check, _ := self.urlStorage.GetURL(ctx, key)
	for url_check != self.cfg.StringEmpty {
		key, err := self.custom_string.GenerateRandomString(self.cfg.KeycodeLengthDefault)
		if err != nil {
			self.logger.WriteLogError(err.Error())
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
		self.logger.WriteLogError(err.Error())
		return &model.ShortenLinkOutput{KeyCode: self.cfg.StringEmpty, Message: err.Error()}, err
	}

	err = self.urlStorage.StoreURL(ctx, key, url)
	if err != nil {
		self.logger.WriteLogError(err.Error())
		return &model.ShortenLinkOutput{KeyCode: self.cfg.StringEmpty, Message: err.Error()}, err
	}
	return &model.ShortenLinkOutput{KeyCode: key, Message: "Shorten URL generated successfully"}, nil
}

// ShortenURL shortens the provided original URL and returns the shortened URL.

func (self *shortenURLService) GetUrlByCode(
	ctx context.Context,
	keyCode string,
) (string, error) {

	url, err := self.urlStorage.GetURL(ctx, keyCode)

	if errors.Is(err, redis.Nil) {
		return self.cfg.StringEmpty, nil
	}

	if err != nil {
		self.logger.WriteLogError(err.Error())
		return self.cfg.StringEmpty, err
	}

	return url, nil
}
