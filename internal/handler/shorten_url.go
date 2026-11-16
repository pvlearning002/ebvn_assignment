package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
	"github.com/redis/go-redis/v9"
)

type ShortenURL interface {
	CreateShortenLink(ctx *gin.Context)
}

// shortenURLHandler handles the HTTP requests for shortening URLs.
type shortenURLHandler struct {
	shortenURLService service.ShortenURL
	cfg               config.Config
	cfgError          config.ConfigError
	client            *redis.Client
}

// NewShortenURLHandler creates a new instance of the ShortenURLHandler with the provided service.
func NewShortenURLHandler(svc service.ShortenURL, cfg config.Config, cfgError config.ConfigError, client *redis.Client) *shortenURLHandler {
	return &shortenURLHandler{
		shortenURLService: svc,
		cfg:               cfg,
		cfgError:          cfgError,
		client:            client,
	}
}

// CreateShortenLink handles the creation of a shortened URL.
// @Summary Create Shortened URL
// @Description CreateShortenLink handles the HTTP request for creating a shortened URL.
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param input body model.ShortenLinkInput true "Input required"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (self *shortenURLHandler) CreateShortenLink(ctx *gin.Context) {
	// Implementation for creating a shortened link
	if err := self.client.Ping(ctx).Err(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorRedisConnection})
		return
	}
	input := &model.ShortenLinkInput{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorInvalidInput})
		return
	}
	shortenLinkOutput, err := self.shortenURLService.CreateShortenLink(ctx, input.Url)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorInternalServer})
		return
	}
	ctx.JSON(http.StatusOK, shortenLinkOutput)
}
