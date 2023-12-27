package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prasiman/gocurl/util/httpclient"
)

var url = flag.String("url", "", "Request URL")
var method = flag.String("method", "GET", "Request method")
var accept = flag.String("accept", "200,201,204", "List of accepted response code, separated by comma")
var headers = flag.String("headers", `{}`, "Request headers")
var params = flag.String("params", `{}`, "Request query params")
var body = flag.String("body", `{}`, "Request body")
var basicAuthToken = flag.String("basic-auth-token", "", "Basic Auth token")
var bearerToken = flag.String("bearer-token", "", "Bearer token")
var logResponse = flag.Bool("log-response", false, "Log response to outputs")
var retries = flag.Int("retries", 1, "Number of retries")
var timeout = flag.Int("timeout", 1000, "Request timeout in ms")
var proxyUrl = flag.String("proxy-url", "", "Proxy URL")
var proxyAuthUsername = flag.String("proxy-auth-username", "", "Proxy auth username")
var proxyAuthPassword = flag.String("proxy-auth-password", "", "Proxy auth username")

func main() {
	flag.Parse()

	time, _ := time.ParseDuration(strconv.Itoa(*timeout) + "ms")
	client := httpclient.NewRetryableClient(time, retries, proxyUrl, proxyAuthUsername, proxyAuthPassword)

	var requestHeaders map[string]any
	var requestParams map[string]any
	var requestBody io.Reader

	json.Unmarshal([]byte(*headers), &requestHeaders)
	json.Unmarshal([]byte(*params), &requestParams)
	requestBody = strings.NewReader(*body)

	// Set request
	request, err := http.NewRequest(*method, *url, requestBody)
	if err != nil {
		log.Fatalln("Error: " + err.Error())
	}

	request.Header.Set("Content-Type", "application/json")

	if *basicAuthToken != "" {
		request.Header.Set("Authorization", "Basic "+*basicAuthToken)
	}

	if *bearerToken != "" {
		request.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Set query params
	q := request.URL.Query()
	for key, value := range requestParams {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	request.URL.RawQuery = q.Encode()

	// Set headers
	for key, value := range requestHeaders {
		request.Header.Set(key, fmt.Sprintf("%v", value))
	}

	// Get response
	response, _ := client.Do(request)

	acceptedCode := strings.Split(*accept, ",")
	accepted := false

	for _, value := range acceptedCode {
		statusCodeInt, _ := strconv.Atoi(value)

		if response.StatusCode == statusCodeInt {
			accepted = true
		}
	}

	if !accepted {
		log.Fatalln("Error: Request failed with status code " + strconv.Itoa(response.StatusCode))
	}

	body, _ := io.ReadAll(response.Body)

	var respBody map[string]any

	json.Unmarshal(body, &respBody)

	// Set response object
	responseArray := map[string]any{
		"data":        respBody,
		"status_code": response.StatusCode,
		"headers":     response.Header,
	}

	responseOutput, _ := json.Marshal(responseArray)

	// Set response to output
	msg := []byte("response=" + string(responseOutput) + "\r\n")
	f, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Unable to write command to the environment file: %s", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Error occurred. Reason: " + err.Error())
		}
	}()

	if _, err := f.Write(msg); err != nil {
		log.Fatalf("Unable to write command to the environment file: %s", err)
	}

	if *logResponse {
		fmt.Println("Response: " + string(responseOutput))
	}
}
