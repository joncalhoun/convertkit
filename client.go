package convertkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Default values
const (
	DefaultBaseURL = "https://api.convertkit.com/v3/"
)

// Client is used to make API calls to the Convert Kit API
type Client struct {
	Secret     string
	BaseURL    string
	HTTPClient interface {
		Do(*http.Request) (*http.Response, error)
	}
}

// Do will perform any API query by:
//
// 1. Adding the API Secret wherever it is needed.
//
// 2. Encoding the params.
//
// 3. Performing the HTTP request for the API call with data from (1) and (2).
//
// 4. Decoding the response body to the provided response variable.
//
// 5. Handling errors from 400+ status codes and parsing the body into an
// ErrorResponse error.
//
// You generally should NOT be using this directly unless you need access an API
// endpoint that isn't supported, you are adding a new API method to this
// package, or you are extending this package in some way (eg creating an
// integration version).
//
// It also happens to simplify testing a bit, since we can define endpoints with
// various errors in our tests and call Do(...) with the relevant method/path.
func (c *Client) Do(method, path string, params, response interface{}) error {
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	req, err := c.request(method, path, params)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 400 {
		if resp.StatusCode == 404 {
			return ErrorResponse{
				StatusCode: 404,
				Type:       "not_found_error",
				Message:    fmt.Sprintf("resource not found or path invalid: %v %v", method, path),
			}
		}
		return c.decodeError(resp)
	}
	err = c.decode(resp.Body, response)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) request(method, path string, params interface{}) (*http.Request, error) {
	var reqURL string
	var reqBody io.Reader
	reqHeader := make(http.Header)

	// All requests need an API Secret
	paramsMap, err := jsonMap(params)
	if err != nil {
		return nil, fmt.Errorf("constructing request: %w", err)
	}
	paramsMap["api_secret"] = c.Secret

	switch method {
	case http.MethodGet, http.MethodDelete:
		query := make(url.Values)
		for k, v := range paramsMap {
			query.Set(k, fmt.Sprintf("%v", v))
		}
		reqURL = c.url(path) + "?" + query.Encode()
	default:
		var buffer bytes.Buffer
		err = json.NewEncoder(&buffer).Encode(paramsMap)
		if err != nil {
			return nil, fmt.Errorf("constructing request: %w", err)
		}
		reqURL = c.url(path)
		reqBody = &buffer
		reqHeader.Set("Content-Type", "application/json")
	}

	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("constructing request: %w", err)
	}
	req.Header = reqHeader

	return req, nil
}

func jsonMap(params interface{}) (map[string]interface{}, error) {
	if params == nil {
		return make(map[string]interface{}), nil
	}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("json map: %w", err)
	}
	var m map[string]interface{}
	err = json.NewDecoder(&buffer).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("json map: %w", err)
	}
	if m == nil {
		return make(map[string]interface{}), nil
	}
	return m, nil
}

func (c *Client) url(path string) string {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
	return fmt.Sprintf("%s%s", baseURL, path)
}

func (c *Client) decodeError(r *http.Response) error {
	var errResp ErrorResponse
	var b bytes.Buffer
	_, err := io.Copy(&b, r.Body)
	if err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}
	errResp.RawBody = b.String()
	err = c.decode(&b, &errResp)
	if err != nil {
		return err
	}
	errResp.StatusCode = r.StatusCode
	return errResp
}

func (c *Client) decode(r io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("decoding: %w", err)
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("decoding: %w", err)
	}
	return nil
}

// ErrorResponse is returned when the server returns a JSON error.
type ErrorResponse struct {
	StatusCode int
	Type       string `json:"error"`
	Message    string `json:"message"`
	RawBody    string
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%d - %v: %v", e.StatusCode, e.Type, e.Message)
}
