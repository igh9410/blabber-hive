package docs

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

func byteHandler(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Write(b)
	}
}

func HandleSpec(swagger *openapi3.T) http.HandlerFunc {
	b, err := swagger.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return byteHandler(b)
}

// UseSwagger adapted for Gin
func UseSwagger(r *gin.Engine, swagger *openapi3.T) {
	// Serve the Swagger JSON
	r.GET("/swagger.json", gin.WrapH(HandleSpec(swagger)))

	// Serve the Swagger UI files
	r.Static("/docs", "docs/swagger-ui") // Adjust the path as per your folder structure
}
