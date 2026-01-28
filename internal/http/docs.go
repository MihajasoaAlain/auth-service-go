package httpapi

import (
	_ "embed"
	"net/http"
)

//go:embed docs/openapi.yaml
var openapiSpec []byte

//go:embed docs/swagger.html
var swaggerHTML []byte

func (a AuthAPI) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(swaggerHTML)
}

func (a AuthAPI) OpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openapiSpec)
}
