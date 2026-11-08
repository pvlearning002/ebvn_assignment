package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
)

type ShortenURL interface {
	CreateShortenLink(ctx *gin.Context)
}

// shortenURLHandler handles the HTTP requests for shortening URLs.
type shortenURLHandler struct {
	shortenURLService service.ShortenURL
	cfg               config.Config
	cfgError          config.ConfigError
}

// NewShortenURLHandler creates a new instance of the ShortenURLHandler with the provided service.
func NewShortenURLHandler(svc service.ShortenURL, cfg config.Config, cfgError config.ConfigError) *shortenURLHandler {
	return &shortenURLHandler{
		shortenURLService: svc,
		cfg:               cfg,
		cfgError:          cfgError,
	}
}

type shortenLinkInput struct {
	Url string `json:"url"`
}

// CreateShortenLink handles the creation of a shortened URL.
// @Summary Create Shortened URL
// @Description CreateShortenLink handles the HTTP request for creating a shortened URL.
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param input body shortenLinkInput true "Input required"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (self *shortenURLHandler) CreateShortenLink(ctx *gin.Context) {
	// Implementation for creating a shortened link
	input := &shortenLinkInput{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorInvalidInput})
		return
	}
	key, err := self.shortenURLService.CreateShortenLink(ctx, input.Url)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorInternalServer})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{self.cfg.KeyTitle: key})
}
