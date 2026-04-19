package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwaggerHandler struct {
	specFile []byte
}

func NewSwaggerHandler(specFile []byte) *SwaggerHandler {
	return &SwaggerHandler{
		specFile: specFile,
	}
}

func (h *SwaggerHandler) ServeSpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml", h.specFile)
}

func (h *SwaggerHandler) ServeUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Source Score API Documentation</title>
    <link
        rel="stylesheet"
        type="text/css"
        href="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui.css"
        integrity="sha384-rcbEi6xgdPk0iWkAQzT2F3FeBJXdG+ydrawGlfHAFIZG7wU6aKbQaRewysYpmrlW"
        crossorigin="anonymous"
    >
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; padding:0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script
        src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-bundle.js"
        integrity="sha384-NXtFPpN61oWCuN4D42K6Zd5Rt2+uxeIT36R7kpXBuY9tLnZorzrJ4ykpqwJfgjpZ"
        crossorigin="anonymous"
    ></script>
    <script
        src="https://unpkg.com/swagger-ui-dist@5.18.2/swagger-ui-standalone-preset.js"
        integrity="sha384-qr68CD0cvHa88PmVu7e1a58Ego4qvKtcvcLdS2a8Mo5zILI01gyIV9jVwJk7X2NU"
        crossorigin="anonymous"
    ></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/swagger/spec",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                onComplete: function () {
                    ui.preauthorizeApiKey("ApiKeyAuth", "demo-api-key");
                }
            });
        };
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
