// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package oapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// NewSite defines model for NewSite.
type NewSite struct {
	// Site Site URL
	Site string `json:"site"`

	// Slogin login for site
	Slogin string `json:"slogin"`

	// Spw passwor for site
	Spw string `json:"spw"`
}

// NewUser defines model for NewUser.
type NewUser struct {
	// Login User login from registration
	Login string `json:"login"`

	// Password User pass from registartion
	Password string `json:"password"`
}

// Site defines model for Site.
type Site struct {
	// Site Site URL
	Site string `json:"site"`

	// SiteID site id (credintial_id)
	SiteID string `json:"siteID"`

	// Slogin login for site
	Slogin string `json:"slogin"`

	// Spw passwor for site
	Spw string `json:"spw"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = NewUser

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = NewUser

// AddSiteJSONRequestBody defines body for AddSite for application/json ContentType.
type AddSiteJSONRequestBody = NewSite

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// User login
	// (POST /api/auth/login)
	Login(w http.ResponseWriter, r *http.Request)
	// User registration
	// (POST /api/auth/register)
	CreateUser(w http.ResponseWriter, r *http.Request)
	// Add new site
	// (POST /api/user/site/add)
	AddSite(w http.ResponseWriter, r *http.Request)
	// get all users sites data
	// (GET /api/user/site/list)
	ListSite(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// User login
// (POST /api/auth/login)
func (_ Unimplemented) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// User registration
// (POST /api/auth/register)
func (_ Unimplemented) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Add new site
// (POST /api/user/site/add)
func (_ Unimplemented) AddSite(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// get all users sites data
// (GET /api/user/site/list)
func (_ Unimplemented) ListSite(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateUser operation middleware
func (siw *ServerInterfaceWrapper) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateUser(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddSite operation middleware
func (siw *ServerInterfaceWrapper) AddSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddSite(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ListSite operation middleware
func (siw *ServerInterfaceWrapper) ListSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListSite(w, r)
	}))

	for i := len(siw.HandlerMiddlewares) - 1; i >= 0; i-- {
		handler = siw.HandlerMiddlewares[i](handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
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
		r.Post(options.BaseURL+"/api/auth/login", wrapper.Login)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/auth/register", wrapper.CreateUser)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/user/site/add", wrapper.AddSite)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/user/site/list", wrapper.ListSite)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8yWX2/bNhDAvwrBDWgGqJaTdg/T07Ik6IwF6TAv2EMRDIR4lrhKJHuk6nqBv/twR9mx",
	"I7kx9q95Mk3eHe/Pj3e6l6VrvbNgY5DFvQxlDa3i5RWiQ1p4dB4wGuDt0mmgXw2hROOjcVYWSVjwWSYX",
	"DlsVZSGNja/OZCbjykP6CxWgXGeyhRBUddDQ5nirGiIaW8n1OpMIHzqDoGXxTvYXbsTv1pm8geXcRBg6",
	"Hvrd/ftIVtz+cj28KpOhcZWxQx3eFguHgm2OafrlUM2rEJYOP6P4KLpeqHcjWe1DvA0wUpsD/pKw6J1G",
	"1wqEyoSIio9HvO8d1Qcs0fGuIYXjhh5Fs4lia51C2ZRKNc3bhSze3cuvERaykF/lD2DmPZX5prbrbKy4",
	"s8uhw7QvjBYnJYI2NhrV/G70N7uQ9u4eU4zZpbxb363pyNiFS8/BRlVGWkKrTCMLqbyJoNrvw1JVFeDE",
	"OJlJq1oyPk974vznmfgVVCsz2SEp1TH6UOT5jtI6exTMG+frnwA8oDBBqAAfOjRxJTYZFQHwoylBnNCC",
	"K4WRYm1MCTZwons/zr0qaxBnk+nAg+VyOVF8PHFY5b1uyK9nF1c386uXZ5PppI5tQ/5FwDa8XczTvaNh",
	"5CyTU35NbEhmVjl8EcQb955jkZn8CBhSiKeT6WRKlp0Hq7yRhXzFW8RNrLnWufImV12s8y3v3oU4LP6P",
	"yuoGULTqPQj4ZEI0thIdMUzq4oT1KUGEEj+HmZaFvO5JpfJDiD84vdqUGizfo7xvTMka+R/B2YfWSasn",
	"EOa3yxDtu8vXihtYCpbYxS9iB8xj8I6KQbecTacj7awrSwhh0TXpwVMqX48Jpj6rgTqo4BBY8nQo+Rs6",
	"W/Xtw+EWtoTnQnVN/NeSk0bOSGo6C588lBG0gF4mk6FrW4WrvQbHBw+ApAbVt8qnGbGwTHjgbn/ch+MC",
	"QUXYqdD/SAjHqbQW0YlYgwjRIRwByumBTl5yLPoJSFiKXk6tQp0a/y4Er6ffDTUTyyZQOlnq2zH7F8q+",
	"iKKCKG5vZ5dEF2ef7nl2eO0NzS1l5G9OoyFXWh+m7FxrhovHEY0i4FEUhHjJe1n/vpTV29xOBuidaz1P",
	"nwT/EXdpuA7zwx9Jf5e7o50zEdrwlJf9/N8Ma4WoVgefynF8H9ME2Zx1kQeHQ/PnZ6ie2QhoVZNYSjVX",
	"WvcfV88F6l0kx3huTAK5ghGe6fAfwnxtQtzSPBxrX4Yacgo0RxGeATVN+mR5TthQs1ZNw406pEQJraIi",
	"U+u/AgAA//9a44pJTg4AAA==",
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
