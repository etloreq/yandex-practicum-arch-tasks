// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
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
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// CurrentState defines model for CurrentState.
type CurrentState struct {
	// DeviceId ID девайса
	DeviceId int `json:"device_id"`

	// Temperature Температура
	Temperature int `json:"temperature"`

	// Timestamp Таймстамп замера
	Timestamp int `json:"timestamp"`
}

// Error defines model for Error.
type Error struct {
	// Code Код ошибки
	Code string `json:"code"`

	// Message Сообщение об ошибке
	Message string `json:"message"`
}

// PostTelemetryJSONRequestBody defines body for PostTelemetry for application/json ContentType.
type PostTelemetryJSONRequestBody = CurrentState

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /telemetry)
	PostTelemetry(w http.ResponseWriter, r *http.Request)

	// (GET /telemetry/{device_id}/latest)
	GetTelemetryDeviceIdLatest(w http.ResponseWriter, r *http.Request, deviceId int)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (POST /telemetry)
func (_ Unimplemented) PostTelemetry(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /telemetry/{device_id}/latest)
func (_ Unimplemented) GetTelemetryDeviceIdLatest(w http.ResponseWriter, r *http.Request, deviceId int) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostTelemetry operation middleware
func (siw *ServerInterfaceWrapper) PostTelemetry(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostTelemetry(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetTelemetryDeviceIdLatest operation middleware
func (siw *ServerInterfaceWrapper) GetTelemetryDeviceIdLatest(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "device_id" -------------
	var deviceId int

	err = runtime.BindStyledParameterWithOptions("simple", "device_id", chi.URLParam(r, "device_id"), &deviceId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "device_id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTelemetryDeviceIdLatest(w, r, deviceId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/telemetry", wrapper.PostTelemetry)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/telemetry/{device_id}/latest", wrapper.GetTelemetryDeviceIdLatest)
	})

	return r
}

type PostTelemetryRequestObject struct {
	Body *PostTelemetryJSONRequestBody
}

type PostTelemetryResponseObject interface {
	VisitPostTelemetryResponse(w http.ResponseWriter) error
}

type PostTelemetry201Response struct {
}

func (response PostTelemetry201Response) VisitPostTelemetryResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type PostTelemetry400JSONResponse Error

func (response PostTelemetry400JSONResponse) VisitPostTelemetryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostTelemetry500JSONResponse Error

func (response PostTelemetry500JSONResponse) VisitPostTelemetryResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetTelemetryDeviceIdLatestRequestObject struct {
	DeviceId int `json:"device_id"`
}

type GetTelemetryDeviceIdLatestResponseObject interface {
	VisitGetTelemetryDeviceIdLatestResponse(w http.ResponseWriter) error
}

type GetTelemetryDeviceIdLatest200JSONResponse CurrentState

func (response GetTelemetryDeviceIdLatest200JSONResponse) VisitGetTelemetryDeviceIdLatestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetTelemetryDeviceIdLatest400JSONResponse Error

func (response GetTelemetryDeviceIdLatest400JSONResponse) VisitGetTelemetryDeviceIdLatestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetTelemetryDeviceIdLatest404JSONResponse Error

func (response GetTelemetryDeviceIdLatest404JSONResponse) VisitGetTelemetryDeviceIdLatestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetTelemetryDeviceIdLatest500JSONResponse Error

func (response GetTelemetryDeviceIdLatest500JSONResponse) VisitGetTelemetryDeviceIdLatestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /telemetry)
	PostTelemetry(ctx context.Context, request PostTelemetryRequestObject) (PostTelemetryResponseObject, error)

	// (GET /telemetry/{device_id}/latest)
	GetTelemetryDeviceIdLatest(ctx context.Context, request GetTelemetryDeviceIdLatestRequestObject) (GetTelemetryDeviceIdLatestResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// PostTelemetry operation middleware
func (sh *strictHandler) PostTelemetry(w http.ResponseWriter, r *http.Request) {
	var request PostTelemetryRequestObject

	var body PostTelemetryJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostTelemetry(ctx, request.(PostTelemetryRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostTelemetry")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostTelemetryResponseObject); ok {
		if err := validResponse.VisitPostTelemetryResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetTelemetryDeviceIdLatest operation middleware
func (sh *strictHandler) GetTelemetryDeviceIdLatest(w http.ResponseWriter, r *http.Request, deviceId int) {
	var request GetTelemetryDeviceIdLatestRequestObject

	request.DeviceId = deviceId

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetTelemetryDeviceIdLatest(ctx, request.(GetTelemetryDeviceIdLatestRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetTelemetryDeviceIdLatest")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetTelemetryDeviceIdLatestResponseObject); ok {
		if err := validResponse.VisitGetTelemetryDeviceIdLatestResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9SVXWsTTxTGv8py/v/LNdlqvdlLrUhRUNC7UmRJTtOR7M46MymUsNCkiIpCLwVBRT/B",
	"tnZx+5LNVzjzjeTMxrw08U6LXmWH3TnnOc/8nkkfWjJOZYKJ0RD2Qbd2MY7c492eUpiYJyYyyOtUyRSV",
	"EejetnFPtPCZaNcL3VIiNUImEMLmhkenVNAJ5XRmB5SDD2Y/RQhBJAY7qCDzwWCcoopMT+FyCfpKBV3S",
	"mAp7QLkd2kP+XV1HxKhNFKcrq+R0Rpd2YIeUcz2PvruH4hflMh8UvugJhW0It+aGXNQ733U78+GeUlIt",
	"e9SS7VWzfaCKTj2q7Gsq6ZjOqZwp0UaJpMNzxah11Fm1/wtVVNGxfUMFjaikwuPlfMFiueCVyZy0WZPt",
	"jD8QyY7kfkaYLu812MUYjdq/EaUCfNhDpWsNQWOtEbBKmWLCL0O41QgaAfiQRmbXTd+cbnfOSG1WzPKe",
	"chpTaQf2nUcjqujEvrUvmZ+cRjRyCzukgi4cEIUd2gMqnWPSHYeQyWYbQngstXk6bVgPi9rcke39+igS",
	"g4lTEKVpV7TczuZzzTJ+Ys9P/yvcgRD+a85y0ZyEormQiCyrPdWpTHR93jeDteURHz1gn9aD4LfJqGlz",
	"7a+4+XGSugsq6ZTdo7Oa+LE9oMoOWMnta1Hyacpi7rnupcdHS2Oq6IJGE3LLBXWUT0ydkdPsTyOYNbuR",
	"wRqiDq5i6TPXtof21TQXjpxze+ii8o0qzw64kR1SZY/4G3u0fFUtcnUfZ1htOC2b7Ye1EGZdRTEaVBrC",
	"rT4IlsH8gw9JFHOE5m+QWf6M6qE/5/HSNbS9RFbwxyj2/3Jg14P161HCFxD/XzAQowUu/oncZNmPAAAA",
	"///to224zQcAAA==",
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
