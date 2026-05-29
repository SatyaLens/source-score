package handlers

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwaggerHandler struct {
	specFile []byte
}

var (
	//go:embed "swagger-ui.html"
	html string
)

func NewSwaggerHandler(specFile []byte) *SwaggerHandler {
	return &SwaggerHandler{
		specFile: specFile,
	}
}

func (h *SwaggerHandler) ServeSpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml", h.specFile)
}

func (h *SwaggerHandler) ServeUI(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
