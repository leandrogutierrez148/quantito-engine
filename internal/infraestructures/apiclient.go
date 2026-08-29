package infraestructures

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const apiTimeout = time.Duration(1 * time.Minute)

const authorizationHeader = "Authorization"
const contentTypeHeader = "Content-Type"

const applicationJson = "application/json"
const formUrlEncoded = "application/x-www-form-urlencoded"

type apiClient struct {
	baseUrl     string
	headers     map[string]string
	bearerToken string
	// Login function used when making a request that requires authentication
	loginFn func() (string, error)
	// Mutex to prevent concurrent access when updating the token
	tokenMutex sync.Mutex
	client     interfaces.IHTTPClient
}

type httpErrorResponse struct {
	Name    string `json:"name,omitempty"`
	Code    int    `json:"code,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewAPIClient(baseUrl string, defaultHeaders map[string]string) interfaces.IAPIClient {
	return &apiClient{
		baseUrl: baseUrl,
		headers: defaultHeaders,
		client: &http.Client{
			Timeout: apiTimeout,
		},
	}
}

// RequestQuery adds query parameters to the request
func RequestQuery(query url.Values) models.APIClientOptionsFunc {
	return func(o *models.APIRequestOptions) error {
		o.Query = query
		return nil
	}
}

// RequestHeaders adds headers to the request
func RequestHeaders(headers map[string]string) models.APIClientOptionsFunc {
	return func(o *models.APIRequestOptions) error {
		o.Headers = headers
		return nil
	}
}

// RequestJsonBody adds a JSON payload to the request
func RequestJsonBody(payload any) models.APIClientOptionsFunc {
	return func(o *models.APIRequestOptions) error {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[contentTypeHeader] = applicationJson

		o.Body = data
		return nil
	}
}

func RequestUrlencodedBody(data url.Values) models.APIClientOptionsFunc {
	return func(o *models.APIRequestOptions) error {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[contentTypeHeader] = formUrlEncoded
		o.Body = []byte(data.Encode())
		return nil
	}
}

// RequestWithAuth configures the request to use authentication. Requires the login function to be set via SetLoginFn
func RequestWithAuth() models.APIClientOptionsFunc {
	return func(o *models.APIRequestOptions) error {
		o.Auth = true
		return nil
	}
}

// RequestWithBasicAuth configures the request to include basicAuth
func RequestWithBasicAuth(user, password string) models.APIClientOptionsFunc {
	return func(o *models.APIRequestOptions) error {
		auth := user + ":" + password
		header := base64.StdEncoding.EncodeToString([]byte(auth))
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[authorizationHeader] = "Basic " + header
		return nil
	}
}

// Get performs a GET request to the API, decoding the JSON response into result.
// Usage example:
//
//	apiClient.Get("/client", &result, infraestructures.RequestQuery(url.Values{"id": []string{"42"}}))
func (c *apiClient) Get(path string, result any, options ...models.APIClientOptionsFunc) error {
	return c.makeAndLogRequest(http.MethodGet, path, options, result)
}

// Post performs a POST request to the API, decoding the JSON response into result.
// Usage example:
//
//	apiClient.Post("/task", &result, infraestructures.RequestJsonBody(requestData), infraestructures.RequestWithAuth())
func (c *apiClient) Post(path string, result any, options ...models.APIClientOptionsFunc) error {
	return c.makeAndLogRequest(http.MethodPost, path, options, result)
}

// Put performs a PUT request to the API, decoding the JSON response into result.
// Usage example:
//
//	apiClient.Put("/task", &result, infraestructures.RequestJsonBody(requestData), infraestructures.RequestWithAuth())
func (c *apiClient) Put(path string, result any, options ...models.APIClientOptionsFunc) error {
	return c.makeAndLogRequest(http.MethodPut, path, options, result)
}

// Patch performs a PATCH request to the API, decoding the JSON response into result.
// Usage example:
//
//	apiClient.Patch("/task", &result, infraestructures.RequestJsonBody(requestData), infraestructures.RequestWithAuth())
func (c *apiClient) Patch(path string, result any, options ...models.APIClientOptionsFunc) error {
	return c.makeAndLogRequest(http.MethodPatch, path, options, result)
}

// Delete performs a DELETE request to the API, decoding the JSON response into result.
// Usage example:
//
//	apiClient.Delete("/client", &result, infraestructures.RequestQuery(url.Values{"id": []string{"42"}}), infraestructures.RequestWithAuth())
func (c *apiClient) Delete(path string, result any, options ...models.APIClientOptionsFunc) error {
	return c.makeAndLogRequest(http.MethodDelete, path, options, result)
}

// SetLoginFn sets the function to be used for login or token refresh.
func (c *apiClient) SetLoginFn(loginFn func() (string, error)) {
	c.loginFn = loginFn
}

func (c *apiClient) makeAndLogRequest(method string, path string, options []models.APIClientOptionsFunc, result any) error {
	apiUrl, err := c.buildUrl(path)
	if err != nil {
		return err
	}
	start := time.Now()

	reqOptions, err := c.buildOptions(options)
	if err != nil {
		return err
	}

	err = c.makeRequest(method, apiUrl, reqOptions, result)
	duration := time.Since(start)

	if err != nil {
		log.Error().Err(err).Str("duration", duration.String()).Str("url", apiUrl).Str("method", method).Str("query", reqOptions.Query.Encode()).Send()
		return err
	}

	log.Debug().Str("duration", duration.String()).Str("url", apiUrl).Str("method", method).Str("query", reqOptions.Query.Encode()).Send()
	return nil
}

func (c *apiClient) buildUrl(path string) (string, error) {
	destination, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	if destination.IsAbs() {
		return path, nil
	}
	return url.JoinPath(c.baseUrl, path)
}

func (c *apiClient) makeRequest(method string, apiUrl string, options *models.APIRequestOptions, result any) error {
	req, err := c.newRequest(method, apiUrl, options)
	if err != nil {
		return err
	}

	if len(options.Query) > 0 {
		req.URL.RawQuery = options.Query.Encode()
	}

	apiResponse, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer apiResponse.Body.Close()

	// Check if response is Unauthorized and retry once after fetching a new token
	if apiResponse.StatusCode == http.StatusUnauthorized && options.Auth {
		log.Debug().Str("url", apiUrl).Msg("Received Unauthorized, retrying after fetching token")
		if err := c.resetToken(c.tokenFromRequest(req)); err != nil {
			return err
		}
		return c.makeRequest(method, apiUrl, options, result)
	}

	if apiResponse.StatusCode != http.StatusOK && apiResponse.StatusCode != http.StatusCreated {
		return c.decodeError(apiResponse)
	}

	return json.NewDecoder(apiResponse.Body).Decode(result)
}

func (c *apiClient) buildOptions(options []models.APIClientOptionsFunc) (*models.APIRequestOptions, error) {
	reqOptions := &models.APIRequestOptions{}
	for _, option := range options {
		err := option(reqOptions)
		if err != nil {
			return nil, err
		}
	}
	return reqOptions, nil
}

func (c *apiClient) decodeError(resp *http.Response) error {
	var e httpErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return errors.New(resp.Status)
	}

	if e.Error != "" {
		return errors.New(e.Error)
	}

	if e.Message != "" {
		return errors.New(e.Message)
	}

	return errors.New(resp.Status)
}

func (c *apiClient) newRequest(method string, url string, options *models.APIRequestOptions) (*http.Request, error) {
	buffer, err := c.requestBuffer(options.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, err
	}

	req.Close = true

	for headerName, headerValue := range c.headers {
		req.Header.Set(headerName, headerValue)
	}

	for headerName, headerValue := range options.Headers {
		req.Header.Set(headerName, headerValue)
	}

	if options.Auth {
		token, err := c.getToken()
		if err != nil {
			return nil, err
		}
		req.Header.Set(authorizationHeader, "Bearer "+token)
	}

	return req, nil
}

func (c *apiClient) resetToken(invalidToken string) error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// If the invalid token is no longer in the client, it is already being updated or was updated
	if c.bearerToken != invalidToken {
		return nil
	}

	c.bearerToken = ""
	_, err := c.requestToken()
	return err
}

func (c *apiClient) getToken() (string, error) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	return c.requestToken()
}

func (c *apiClient) requestToken() (string, error) {
	if c.loginFn == nil {
		return "", errors.New("login function is not defined in apiClient")
	}

	if c.bearerToken != "" {
		return c.bearerToken, nil
	}

	token, err := c.loginFn()
	if err != nil {
		return "", err
	}

	c.bearerToken = token
	return token, nil
}

func (c *apiClient) tokenFromRequest(req *http.Request) string {
	header := req.Header.Get(authorizationHeader)
	return strings.TrimPrefix(header, "Bearer ")
}

func (c *apiClient) requestBuffer(payload []byte) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}
	return bytes.NewReader(payload), nil
}
