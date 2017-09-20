package goproxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"context"

	"github.com/andrexus/goproxmox/pveauth"
	"github.com/hashicorp/logutils"
)

const (
	libraryVersion  = "0.1.2"
	logLevelEnvName = "GOPROXMOX_LOGLEVEL"
	apiBasePath     = "/api2/json/"
	mediaType       = "application/json"
)

func init() {
	logLevel := os.Getenv(logLevelEnvName)
	if logLevel == "" {
		logLevel = "INFO"
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

// Client manages communication with proxmox API.
type Client struct {
	// HTTP client used to communicate with the proxmox API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// Control panel username
	Username string

	// Control panel password
	Password string

	// Services used for communicating with the API
	Nodes NodesService
	VMs   QemuService

	// Optional function called after every successful request made to the proxmox API
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Errors map
	Errors map[string]string `json:"errors"`

	// ResponseCode returned from the API
	ResponseCode int
}

// NewClient returns a new proxmox API client.
func NewClient(host, username, password string) *Client {
	config := pveauth.Config{
		Username:  username,
		Password:  password,
		TicketURL: fmt.Sprintf("%s%saccess/ticket", host, apiBasePath),
	}
	ticket, err := config.PasswordCredentialsTicket(context.Background())
	if err != nil {
		panic(err)
	}
	httpClient := config.Client(context.Background(), ticket)

	apiServerBaseUrl := fmt.Sprintf("%s%s", host, apiBasePath)
	baseURL, _ := url.Parse(apiServerBaseUrl)

	c := &Client{client: httpClient, BaseURL: baseURL}

	log.Printf("[DEBUG] Base URL: %s\n", baseURL)

	c.Nodes = &NodesServiceOp{client: c}
	c.VMs = &QemuServiceOp{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, body map[string]string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	urlValues := url.Values{}
	if body != nil {
		for k, v := range body {
			urlValues.Add(k, v)
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBufferString(urlValues.Encode()))
	if err != nil {
		return nil, err
	}

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Add("Accept", mediaType)

	return req, nil
}

// OnRequestCompleted sets the API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()
	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("[DEBUG] Response: %s\n", string(body))
			err = json.NewDecoder(bytes.NewReader(body)).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}
	return resp, err
}

func (r *ErrorResponse) Error() string {
	errors := ""
	for  key, value := range r.Errors {
		errors += fmt.Sprintf("\t%s: %s\n", key, value)
	}
	return fmt.Sprintf("Respose code: %d\nErrors:\n%s", r.ResponseCode, errors)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r, ResponseCode: r.StatusCode}
	data, err := ioutil.ReadAll(r.Body)
	log.Printf("[DEBUG] Check response: %s\n", string(data))
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}
