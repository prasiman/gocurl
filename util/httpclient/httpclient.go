package httpclient

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var AcceptedCode = []string{"200", "201", "204"}
var RetryCount = 1

type retryableTransport struct {
	transport http.RoundTripper
}

func NewRetryableClient(timeout time.Duration, accept []string, retry *int, proxyUrl *string, proxyUsername *string, proxyPassword *string) *http.Client {
	var transport = &retryableTransport{
		transport: &http.Transport{},
	}

	if *proxyUrl != "" {
		proxy, _ := url.Parse(*proxyUrl)
		if *proxyUsername != "" {
			proxy.User = url.UserPassword(*proxyUsername, *proxyPassword)
		}
		transport = &retryableTransport{
			transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	RetryCount = *retry
	AcceptedCode = accept

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	// Send the request
	resp, err := t.transport.RoundTrip(req)
	// Retry logic
	retries := 0
	for shouldRetry(err, resp) && retries < RetryCount {
		// Wait for the specified backoff period
		time.Sleep(backoff(retries))
		// We're going to retry, consume any response to reuse the connection.
		drainBody(resp)
		// Clone the request body again
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		// Retry the request
		resp, err = t.transport.RoundTrip(req)
		retries++
	}
	// Return the response
	return resp, err
}

func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	for _, value := range AcceptedCode {
		statusCodeInt, _ := strconv.Atoi(value)

		if resp.StatusCode == statusCodeInt {
			return false
		}
	}

	return true
}

func drainBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
