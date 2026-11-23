package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

type ShortenURL interface {
	CreateShortenLink(ctx *gin.Context)
	RedirectToOriginalURL(ctx *gin.Context)
}

// shortenURLHandler handles the HTTP requests for shortening URLs.
type shortenURLHandler struct {
	shortenURLService service.ShortenURL
	cfg               config.Config
	cfgError          config.ConfigError
	response          utils.ResponseUtils
	cfgParameter      config.ConfigParameter
	logger            utils.LoggerUtils
}

// NewShortenURLHandler creates a new instance of the ShortenURLHandler with the provided service.
func NewShortenURLHandler(svc service.ShortenURL, cfg config.Config, cfgError config.ConfigError, response utils.ResponseUtils, cfgParameter config.ConfigParameter, logger utils.LoggerUtils) *shortenURLHandler {
	return &shortenURLHandler{
		shortenURLService: svc,
		cfg:               cfg,
		cfgError:          cfgError,
		response:          response,
		cfgParameter:      cfgParameter,
		logger:            logger,
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
	input := &model.ShortenLinkInput{}
	if err := ctx.ShouldBindJSON(input); err != nil {
		self.logger.WriteLogError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, self.response.GetInputFieldError(err))
		return
	}
	shortenLinkOutput, err := self.shortenURLService.CreateShortenLink(ctx, input.Url)
	if err != nil {
		self.logger.WriteLogError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, self.response.GetErrorInternalServer())
		return
	}
	ctx.JSON(http.StatusOK, shortenLinkOutput)
}

// RedirectToOriginalURL handles the redirection to the original URL based on the provided key code.
// @Summary Redirect to Original URL
// @Description RedirectToOriginalURL handles the HTTP request for redirecting to the original URL based on the provided key code.
// @Tags link
// @Produce application/json
// @Param code path string true "Key Code"
// @Success 302 {string} string "Redirects to the original URL"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/redirect/{code} [get]
func (self *shortenURLHandler) RedirectToOriginalURL(ctx *gin.Context) {
	keyCode := ctx.Param(self.cfgParameter.Code)
	if keyCode == self.cfg.StringEmpty {
		self.logger.WriteLogError(self.cfgError.ErrorInvalidCode)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, self.response.GetErrorInvalidInput())
		return
	}
	url, err := self.shortenURLService.GetUrlByCode(ctx, keyCode)
	if err != nil {
		self.logger.WriteLogError(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, self.response.GetErrorInternalServer())
		return
	}
	if url == self.cfg.StringEmpty {
		ctx.AbortWithStatusJSON(http.StatusNotFound, self.response.GetErrorUrlNotFound())
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, url)
}
