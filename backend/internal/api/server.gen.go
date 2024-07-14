// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /api/chats)
	ChatServiceCreateChatRoom(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// ChatServiceCreateChatRoom operation middleware
func (siw *ServerInterfaceWrapper) ChatServiceCreateChatRoom(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ChatServiceCreateChatRoom(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/api/chats", wrapper.ChatServiceCreateChatRoom)
}

type ChatServiceCreateChatRoomRequestObject struct {
	Body *ChatServiceCreateChatRoomJSONRequestBody
}

type ChatServiceCreateChatRoomResponseObject interface {
	VisitChatServiceCreateChatRoomResponse(w http.ResponseWriter) error
}

type ChatServiceCreateChatRoom200JSONResponse CreateChatRoomResponse

func (response ChatServiceCreateChatRoom200JSONResponse) VisitChatServiceCreateChatRoomResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ChatServiceCreateChatRoomdefaultJSONResponse struct {
	Body       Status
	StatusCode int
}

func (response ChatServiceCreateChatRoomdefaultJSONResponse) VisitChatServiceCreateChatRoomResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /api/chats)
	ChatServiceCreateChatRoom(ctx context.Context, request ChatServiceCreateChatRoomRequestObject) (ChatServiceCreateChatRoomResponseObject, error)
}

type StrictHandlerFunc = strictgin.StrictGinHandlerFunc
type StrictMiddlewareFunc = strictgin.StrictGinMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// ChatServiceCreateChatRoom operation middleware
func (sh *strictHandler) ChatServiceCreateChatRoom(ctx *gin.Context) {
	var request ChatServiceCreateChatRoomRequestObject

	var body ChatServiceCreateChatRoomJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.ChatServiceCreateChatRoom(ctx, request.(ChatServiceCreateChatRoomRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ChatServiceCreateChatRoom")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(ChatServiceCreateChatRoomResponseObject); ok {
		if err := validResponse.VisitChatServiceCreateChatRoomResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RWTW/jNhD9KwO2hy7gldPdm0/reoNF0EONJJciCLBjciTNluKo5CiBG+S/F6Tk2I6N",
	"boC2N/Fr3ps3jxw9GStdL4GCJrN4Msm21GH5XLWo1yJd/u6j9BSVqazYSKjklpoHtcQO1SyMQ6X3yh2Z",
	"mdFtT2ZhkkYOjXmeGXZ578l0wI7OLDy/RJDNN7Kat64K6I7UNf05UNJTbv8iYuolJDqT7oEQP0aqzcL8",
	"MN/LNp80m78Idhbti0jjaR1FZTPUy7DN4dA5VpaAfn0AqXGgmXGUbOQ+L5uFWUlQ5JAAA2DcsEaMW0gU",
	"GT3/RQ46SgkbAvQSGnhkbQHhU6YB2qLCGG5DCbQlKPNSl+/TIJWZvdLg05jP0ytWt28K9YZa3CjqkM4D",
	"fB0Xv45IjmoOlADBS8MWPVCMEqETR35MlROkgRU3nqCWCI7rmiIFhT5KE7HrODRA4YGjhC7XcAYcrB9c",
	"nr++vLmF5foqS+3ger0qgwquSuAhkYPNFu6a6/Xq/qdWtU+L+bxhbYdNZaWbN7G37yq4RNvuqe+qY3dl",
	"1DYSQc9kKWX1HCouplSsOJrt0hoPzgqZccqRIvtUwe8ygMUANQcHMih0Eglwkz+15XSkTD7fyiOowKPE",
	"P0aHsAKHUrm75foKPlPiJsCXgR3tc7NeBlc1xb8lQ+w5zV3ZOy8Q6d2pY3IS5+uZiiZTlo8t2xZSK4N3",
	"sKFsbwpDBw/oh2Kruwk49rZaiaP7k4mM/fIKcdCPH/ae46DUUMwWm2Q75bQEz0kz1iR2Gm1kMcZtEeeV",
	"7nDbUqTsBgQrXScBEh0GKE5NxXvFSCrZN5knK3Xpew/J6VOxvzIYI5bxBHUuHUcP5HMt3tdoi9WPrfRa",
	"cw5wGRrPqa1gGbaZ6/mjB2e82OmiZ2OlfLl2Tjqoz2j/alLu/h+WoGbybgYSD0JvRvmtZwr6ppckT3Go",
	"pbQAVp/X8rt8Q/GBLeVymJl5oJhGsS6qi+rnLKf0FLBnszAfq4vqY7YzaltKlf0+z02gjHoZ207WF7Pk",
	"V+4Y47itmJmJY6/6Rdx2vBhBKZQY2PeebYky/5YyoV0H/m6vOdsNS/oZjiO5sY+UibGzFfofLi7+NxJT",
	"Ay0sjj3526/jDaxx8Pqf4U9N4wze5xFpcm/cM5sZxSaZxd1hxcz988HK7h/iaMfz/fPfAQAA//82aSHG",
	"KwkAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
