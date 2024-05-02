package gpt3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultAPIVersion     = "2023-03-15-preview"
	defaultUserAgent      = "kubectl-openai"
	defaultTimeoutSeconds = 30
)

// Client  is an api client to communicate with the openAI gtp3 apis
type Client interface {
	// ChatCompletion creates a completion with the Chat completion endpoint which
	// is what powers the ChatGPT experience.
	ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error)

	// Completion creates a completion with the default engine. This is the main endpoint of the API
	// which auto-completes based on the given prompt.
	Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
}

type client struct {
	endpoint       string
	apiKey         string
	deploymentName string
	apiVersion     string
	userAgent      string
	httpClient     *http.Client
}

// NewClient create a new gpt-3 client with the specified params
//
//	@param endpoint
//	@param apiKey
//	@param deploymentName
//	@param options
//	@return Client
//	@return error
func NewClient(endpoint string, apiKey string, deploymentName string, options ...ClientOption) (Client, error) {
	// Create a new HTTP client with a default timeout.
	httpClient := &http.Client{
		Timeout: defaultTimeoutSeconds * time.Second,
	}

	// Create a new client instance with the provided parameters.
	c := &client{
		endpoint:       endpoint,
		apiKey:         apiKey,
		deploymentName: deploymentName,
		apiVersion:     defaultAPIVersion,
		userAgent:      defaultUserAgent,
		httpClient:     httpClient,
	}

	// Apply any additional client options provided.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Completion sends a completion request to the OpenAI API and returns the completion response.
//
//	@receiver c
//	@param ctx
//	@param request
//	@return *CompletionResponse
//	@return error
func (c *client) Completion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	// Set the Stream field of the request to false
	request.Stream = false

	// Create a new request using the context, HTTP method, and endpoint URL
	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/openai/deployments/%s/completions", c.deploymentName), request)
	if err != nil {
		return nil, err
	}

	// Perform the request and get the response
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	// Create a new CompletionResponse object to store the response data
	output := new(CompletionResponse)

	// Parse the response and populate the output object
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}

	// Return the output object and nil error
	return output, nil
}

// ChatCompletion sends a chat completion request to the OpenAI API and returns the response.
//
//	@receiver c
//	@param ctx
//	@param request
//	@return *ChatCompletionResponse
//	@return error
func (c *client) ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	request.Stream = false

	req, err := c.newRequest(ctx, "POST", fmt.Sprintf("/openai/deployments/%s/chat/completions", c.deploymentName), request)
	if err != nil {
		return nil, err
	}

	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ChatCompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

// performRequest Getting called in completion, chatCompletion and multiple other functions above in this file
//
//	@receiver c
//	@param req
//	@return *http.Response
//	@return error
func (c *client) performRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkForSuccess(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// checkForSuccess Getting called in the performRequest function above
// returns an error if this response includes an error.
//
//	@param resp
//	@return error
func checkForSuccess(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	var result APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := APIError{
			StatusCode: resp.StatusCode,
			Type:       "Unexpected",
			Message:    string(data),
		}
		return apiError
	}
	result.Error.StatusCode = resp.StatusCode
	return result.Error
}

// Getting called in this file above, in the completion and chat completion functions
func getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
}

// jsonBodyReader Getting called in the newRequest function below
// jsonBodyReader is a helper function that converts the given body interface{} into a JSON-encoded io.Reader.
//
//	@param body
//	@return io.Reader
//	@return error
func jsonBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}

	raw, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed encoding json: %w", err)
	}

	return bytes.NewBuffer(raw), nil
}

// newRequest Getting called in completion, chatCompletion and other functions above
// newRequest creates a new HTTP request with the specified method, path, and payload.
//
//	@receiver c
//	@param ctx
//	@param method
//	@param path
//	@param payload
//	@return *http.Request
//	@return error
func (c *client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	// Create a JSON body reader from the payload
	bodyReader, err := jsonBodyReader(payload)
	if err != nil {
		return nil, err
	}

	// Construct the request URL with the endpoint, path, and API version
	reqURL := fmt.Sprintf("%s%s?api-version=%s", c.endpoint, path, c.apiVersion)

	// Create a new HTTP request with the specified method, URL, and body reader
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set the Content-type and api-key headers
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	return req, nil
}
