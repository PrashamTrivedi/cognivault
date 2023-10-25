package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title My Go Project API
// @version 1.0
// @description This is a sample API for My Go Project.
// @host localhost:8080
// @BasePath /api/v1
func SwaggerHandler() http.HandlerFunc {
	return httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition"
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	)
}
