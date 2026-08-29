package infraestructures

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"io"
// 	"math"
// 	"quantito-engine/internal/interfaces/mocks"
// 	"quantito-engine/internal/models"
// 	"net/http"
// 	"net/url"
// 	"testing"

// 	"github.com/jarcoal/httpmock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// var successfullResponseContent = map[string]string{
// 	"response": "OK",
// }

// var errorResponse = map[string]string{
// 	"message": ":-(",
// }

// var requestData = map[string]string{
// 	"Id": "42",
// }

// const urlWithQuery = "http://BASEURL/endpoint?a=b"
// const urlWithoutQuery = "http://BASEURL/endpoint"

// func TestRequestQuery(t *testing.T) {
// 	options := &models.ApiRequestOptions{}
// 	query := url.Values{"key": []string{"value"}}

// 	err := RequestQuery(query)(options)

// 	assert.NoError(t, err)
// 	assert.Equal(t, query, options.Query)
// }

// func TestRequestHeaders(t *testing.T) {
// 	options := &models.ApiRequestOptions{}
// 	headers := map[string]string{"key": "value"}

// 	err := RequestHeaders(headers)(options)

// 	assert.NoError(t, err)
// 	assert.Equal(t, headers, options.Headers)
// }

// func TestRequestJsonBody(t *testing.T) {
// 	t.Run("Happy path", func(t *testing.T) {
// 		options := &models.ApiRequestOptions{}
// 		body := map[string]string{"key": "value"}

// 		err := RequestJsonBody(body)(options)
// 		bodyBytes, _ := json.Marshal(body)

// 		assert.NoError(t, err)
// 		assert.Equal(t, bodyBytes, options.Body)
// 		assert.Equal(t, applicationJson, options.Headers[contentTypeHeader])
// 	})

// 	t.Run("Error when cannot be encoded as JSON", func(t *testing.T) {
// 		options := &models.ApiRequestOptions{}
// 		body := map[string]any{"key": math.Inf(1)}

// 		err := RequestJsonBody(body)(options)

// 		assert.Error(t, err)
// 	})
// }

// func TestRequestUrlencodedBody(t *testing.T) {
// 	options := &models.ApiRequestOptions{}
// 	data := url.Values{
// 		"key": []string{"value"},
// 	}

// 	err := RequestUrlencodedBody(data)(options)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "key=value", string(options.Body))
// 	assert.Equal(t, formUrlEncoded, options.Headers[contentTypeHeader])
// }

// func TestRequestWithAuth(t *testing.T) {
// 	options := &models.ApiRequestOptions{}

// 	err := RequestWithAuth()(options)

// 	assert.NoError(t, err)
// 	assert.True(t, options.Auth)
// }

// func TestRequestWithBasicAuth(t *testing.T) {
// 	options := &models.ApiRequestOptions{}

// 	err := RequestWithBasicAuth("username", "password")(options)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "Basic dXNlcm5hbWU6cGFzc3dvcmQ=", options.Headers[authorizationHeader])
// }

// func TestGet(t *testing.T) {
// 	t.Run("Happy path", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithQuery, nil)
// 		expectedReq.Header.Add("key2", "value2")
// 		httpResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)

// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestQuery(url.Values{"a": []string{"b"}}), RequestHeaders(map[string]string{"key2": "value2"}))

// 		assert.NoError(t, err)
// 		assert.Equal(t, successfullResponseContent, response)
// 	})

// 	t.Run("Error when native HTTP client fails", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		mockNativeClient.EXPECT().Do(mock.Anything).Return(nil, errors.New("CLIENT_ERROR"))

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestQuery(url.Values{"a": []string{"b"}}))

// 		assert.Error(t, err)
// 		assert.Nil(t, response)
// 	})

// 	t.Run("Error when server fails with message field", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithQuery, nil)

// 		httpResponse, _ := httpmock.NewJsonResponse(500, errorResponse)

// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestQuery(url.Values{"a": []string{"b"}}))

// 		assert.Error(t, err)
// 		assert.Equal(t, errors.New(":-("), err)
// 		assert.Nil(t, response)
// 	})

// 	t.Run("Error when server fails with error field", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithQuery, nil)

// 		expectedResponse := map[string]string{"error": ":-("}
// 		httpResponse, _ := httpmock.NewJsonResponse(500, expectedResponse)

// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestQuery(url.Values{"a": []string{"b"}}))

// 		assert.Error(t, err)
// 		assert.Equal(t, errors.New(":-("), err)
// 		assert.Nil(t, response)
// 	})

// 	t.Run("Error when server fails without message", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithQuery, nil)

// 		httpResponse, _ := httpmock.NewJsonResponse(404, map[string]string{})
// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestQuery(url.Values{"a": []string{"b"}}))

// 		assert.Error(t, err)
// 		assert.Equal(t, errors.New("404"), err)
// 		assert.Nil(t, response)
// 	})

// 	t.Run("Error if options processing fails", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		mockOptions := func(*models.ApiRequestOptions) error {
// 			return errors.New("ERROR")
// 		}

// 		err := client.Get("endpoint", struct{}{}, mockOptions)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "ERROR")
// 	})
// }

// func TestPost(t *testing.T) {
// 	t.Run("Happy path", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodPost, urlWithoutQuery, requestData)

// 		httpResponse, _ := httpmock.NewJsonResponse(201, successfullResponseContent)
// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Post("endpoint", &response, RequestJsonBody(requestData))

// 		assert.NoError(t, err)
// 		assert.Equal(t, successfullResponseContent, response)
// 	})

// 	t.Run("Error when native HTTP client fails", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		mockNativeClient.EXPECT().Do(mock.Anything).Return(nil, errors.New("CLIENT_ERROR"))

// 		var response map[string]string
// 		err := client.Post("endpoint", &response)

// 		assert.Error(t, err)
// 		assert.Nil(t, response)
// 	})

// 	t.Run("Error when server fails with message", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodPost, urlWithoutQuery, requestData)

// 		httpResponse, _ := httpmock.NewJsonResponse(500, errorResponse)
// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Post("endpoint", &response, RequestJsonBody(requestData))

// 		assert.Error(t, err)
// 		assert.Equal(t, errors.New(":-("), err)
// 		assert.Nil(t, response)
// 	})

// 	t.Run("Error when server fails without message", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodPost, urlWithoutQuery, requestData)

// 		httpResponse, _ := httpmock.NewJsonResponse(404, map[string]string{})
// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Post("endpoint", &response, RequestJsonBody(requestData))

// 		assert.Error(t, err)
// 		assert.Equal(t, errors.New("404"), err)
// 		assert.Nil(t, response)
// 	})
// }

// func TestPut(t *testing.T) {
// 	mockNativeClient := mocks.NewMockIHttpClient(t)
// 	client := newTestApiClient(mockNativeClient)

// 	expectedReq := newMockRequest(http.MethodPut, urlWithoutQuery, requestData)

// 	httpResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)
// 	mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 	var response map[string]string
// 	err := client.Put("endpoint", &response, RequestJsonBody(requestData))

// 	assert.NoError(t, err)
// 	assert.Equal(t, successfullResponseContent, response)
// }

// func TestPatch(t *testing.T) {
// 	mockNativeClient := mocks.NewMockIHttpClient(t)
// 	client := newTestApiClient(mockNativeClient)

// 	expectedReq := newMockRequest(http.MethodPatch, urlWithoutQuery, requestData)

// 	httpResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)
// 	mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 	var response map[string]string
// 	err := client.Patch("endpoint", &response, RequestJsonBody(requestData))

// 	assert.NoError(t, err)
// 	assert.Equal(t, successfullResponseContent, response)
// }

// func TestDelete(t *testing.T) {
// 	mockNativeClient := mocks.NewMockIHttpClient(t)
// 	client := newTestApiClient(mockNativeClient)

// 	expectedReq := newMockRequest(http.MethodDelete, urlWithQuery, nil)

// 	httpResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)
// 	mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 	var response map[string]string
// 	err := client.Delete("endpoint", &response, RequestQuery(url.Values{"a": []string{"b"}}))

// 	assert.NoError(t, err)
// 	assert.Equal(t, successfullResponseContent, response)
// }

// type mockLoginStruct struct {
// 	mock.Mock
// }

// func (l *mockLoginStruct) Login() (string, error) {
// 	args := l.Called()
// 	return args.String(0), args.Error(1)
// }

// func TestLogin(t *testing.T) {
// 	t.Run("Happy path", func(t *testing.T) {
// 		mockLogin := &mockLoginStruct{}
// 		mockLogin.On("Login").Return("token", nil).Once()

// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		client.SetLoginFn(mockLogin.Login)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithoutQuery, nil)
// 		expectedReq.Header.Add("Authorization", "Bearer token")

// 		httpResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)

// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestWithAuth())

// 		assert.NoError(t, err)
// 		assert.Equal(t, successfullResponseContent, response)
// 		mockLogin.AssertExpectations(t)
// 	})

// 	t.Run("Error if login function is not defined and request made with auth", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithoutQuery, nil)
// 		expectedReq.Header.Add("Authorization", "Bearer token")

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestWithAuth())

// 		assert.Error(t, err)
// 	})

// 	t.Run("Error if login fails", func(t *testing.T) {
// 		mockLogin := &mockLoginStruct{}
// 		mockLogin.On("Login").Return("", errors.New("ERROR")).Once()

// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		client.SetLoginFn(mockLogin.Login)

// 		expectedReq := newMockRequest(http.MethodGet, urlWithoutQuery, nil)
// 		expectedReq.Header.Add("Authorization", "Bearer token")

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestWithAuth())

// 		assert.Error(t, err)
// 		mockLogin.AssertExpectations(t)
// 	})

// 	t.Run("Resets token on 401 and retries request", func(t *testing.T) {
// 		mockLogin := &mockLoginStruct{}
// 		mockLogin.On("Login").Return("token1", nil).Once()
// 		mockLogin.On("Login").Return("token2", nil).Once()

// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		client.SetLoginFn(mockLogin.Login)

// 		failedReq := newMockRequest(http.MethodGet, urlWithoutQuery, nil)
// 		failedReq.Header.Add("Authorization", "Bearer token1")

// 		successfulReq := newMockRequest(http.MethodGet, urlWithoutQuery, nil)
// 		successfulReq.Header.Add("Authorization", "Bearer token2")

// 		failedResponse, _ := httpmock.NewJsonResponse(401, errorResponse)
// 		successResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)

// 		mockNativeClient.EXPECT().Do(mock.Anything).Return(failedResponse, nil).Once()
// 		mockNativeClient.EXPECT().Do(mock.Anything).Return(successResponse, nil).Once()

// 		var response map[string]string
// 		err := client.Get("endpoint", &response, RequestWithAuth())

// 		assert.NoError(t, err)
// 		assert.Equal(t, successfullResponseContent, response)
// 		mockLogin.AssertExpectations(t)
// 		mockNativeClient.AssertExpectations(t)
// 	})

// 	t.Run("Does not reset token if request does not need authentication", func(t *testing.T) {
// 		mockLogin := &mockLoginStruct{}

// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		client.SetLoginFn(mockLogin.Login)

// 		failedReq := newMockRequest(http.MethodGet, urlWithoutQuery, nil)
// 		failedResponse, _ := httpmock.NewJsonResponse(401, errorResponse)

// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, failedReq))).Return(failedResponse, nil).Once()

// 		var response map[string]string
// 		err := client.Get("endpoint", &response)

// 		assert.Error(t, err)
// 		mockLogin.AssertNotCalled(t, "Login")
// 		mockNativeClient.AssertExpectations(t)
// 	})
// }

// func TestNewApiClient(t *testing.T) {
// 	baseUrl := "baseUrl"
// 	defaultHeaders := map[string]string{"Header1": "val1"}

// 	client := NewAPIClient(baseUrl, defaultHeaders)

// 	assert.Equal(t, baseUrl, client.(*apiClient).baseUrl)
// 	assert.Equal(t, defaultHeaders, client.(*apiClient).headers)
// 	assert.IsType(t, &http.Client{}, client.(*apiClient).client)
// }

// func TestApiClient(t *testing.T) {
// 	t.Run("Ignores baseUrl if url is absolute", func(t *testing.T) {
// 		mockNativeClient := mocks.NewMockIHttpClient(t)
// 		client := newTestApiClient(mockNativeClient)

// 		requestPath := "http://otherdomain.com/endpoint"

// 		expectedReq := newMockRequest(http.MethodGet, requestPath, nil)
// 		httpResponse, _ := httpmock.NewJsonResponse(200, successfullResponseContent)

// 		mockNativeClient.EXPECT().Do(mock.MatchedBy(matchRequest(t, expectedReq))).Return(httpResponse, nil)

// 		var response map[string]string
// 		err := client.Get(requestPath, &response)

// 		assert.NoError(t, err)
// 		assert.Equal(t, successfullResponseContent, response)
// 	})
// }

// func newTestApiClient(mockNativeClient *mocks.MockIHttpClient) apiClient {
// 	return apiClient{
// 		baseUrl: "http://BASEURL",
// 		headers: map[string]string{
// 			"key": "value",
// 		},
// 		client: mockNativeClient,
// 	}
// }

// func newMockRequest(method string, url string, body any) *http.Request {
// 	var request *http.Request
// 	if body != nil {
// 		request, _ = http.NewRequest(method, url, bodyReader(body))
// 	} else {
// 		request, _ = http.NewRequest(method, url, nil)
// 	}
// 	request.Header.Set("key", "value")

// 	if body != nil {
// 		request.Header.Set("Content-Type", "application/json")
// 	}

// 	return request
// }

// func matchRequest(t *testing.T, expectedRequest *http.Request) func(req *http.Request) bool {
// 	return func(req *http.Request) bool {
// 		return assert.Equal(t, expectedRequest.Host, req.Host) &&
// 			assert.Equal(t, expectedRequest.Header, req.Header) &&
// 			assert.Equal(t, expectedRequest.Method, req.Method) &&
// 			assert.Equal(t, decodeBody(t, expectedRequest.Body), decodeBody(t, req.Body)) &&
// 			assert.Equal(t, expectedRequest.URL, req.URL)
// 	}
// }

// func decodeBody(t *testing.T, body io.ReadCloser) any {
// 	if body == nil {
// 		return nil
// 	}
// 	data, err := io.ReadAll(body)
// 	assert.NoError(t, err)
// 	return data
// }

// func bodyReader(body any) *bytes.Reader {
// 	jsonBody, _ := json.Marshal(body)
// 	return bytes.NewReader(jsonBody)
// }
