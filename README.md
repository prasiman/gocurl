# cURL for GitHub Action

You can use this action to perform REST API requests, using the net/http package.

# Usage

```yaml
name: Example of cURL action

on: [push]
jobs:
  test-curl-action:
    name: "Perform REST API request"
    runs-on: ubuntu-latest
    steps:
      - name: "Call API"
        uses: prasiman/gocurl@v1
        with:
          # The target URL
          # Required: true if custom-config is not set
          url: https://reqres.in/api/users

          # The request method, basically it's one of GET|POST|PUT|PATCH
          # Default is GET
          method: "POST"

          # List of response status codes to be accepted, else it will set the job to be failed
          # If more than one value is needed, you can use comma (,) as separator
          # In this case if the response status code is not one of 200, 201 and 204, the job will be failed
          # Default is 200,201,204
          accept: 200,201,204

          # Headers can be passed through with an escaped json object string
          # Default: "{}"
          headers: '{ "custom-header": "value" }'

          # Params can be passed through with an escaped json object string
          # Default: "{}"
          params: '{ "param1": "value", "param2": "value2" }'

          # Body request passed through with an escaped json object string
          # Apply only to POST|PUT request
          # Default: "{}"
          body: '{ "name": "breeze",  "job": "devops" }'

          # Send raw body instead of JSON-parsed or YAML-parsed body
          # Default: false
          send-raw-body: false

          # Request timeout (milliseconds)
          # Default: 1000
          timeout: 1000

          # Basic authentication using username and password
          # This will override the Authorization header, for example Authorization: Basic QWxhZGRpbjpPcGVuU2VzYW1l
          # Format => username:password as base 64
          # Default: ""
          basic-auth-token: ${{ secrets.curl_auth_token }}

          # The authentication using token
          # This will override the Authorization header, for example Authorization: Bearer QWxhZGRpbjpPcGVuU2VzYW1l
          # Default: ""
          bearer-token: ${{ secrets.bearer_token }}

          # If you want to use proxy with the request, you can use proxy-url
          # Format => host:port
          # Default: ""
          proxy-url: https://proxy-url:3000

          # Proxy authentication username
          # Default: ""
          proxy-auth-username: ${{ secrets.proxy_auth_username }}

          # Proxy authentication password
          # Default: ""
          proxy-auth-password: ${{ secrets.proxy_auth_password }}

          # If it is set to true, it will show the response log in the GitHub UI
          # Default: false
          log-response: false

          # The number of attempts before giving up
          # Default: 1
          retries: 3
```

# Response object

```javascript
{
  // `data` is the response that was provided by the server
  "data": {},

  // `status` is the HTTP status code from the server response
  "status": 200,

  // `headers` the HTTP headers that the server responded with
  // All header names are lower cased and can be accessed using the bracket notation.
  // Example: `response.headers['content-type']`
  "headers": {},

}

```

# Use Response

```yaml
name: Example of cURL action

on: [push]
jobs:
  test-curl-action:
    name: "Perform REST API"
    runs-on: ubuntu-latest
    steps:
      - name: "Call API 1"
        uses: prasiman/gocurl@v1
        id: api
        with:
          url: https://reqres.in/api/users
          method: "POST"
          accept: 201
          body: '{ "name": "breeze", "job": "devops" }'
          log-response: true
      - name: "Call API 2"
        uses: prasiman/gocurl@v1
        id: api2
        with:
          url: https://reqres.in/api/users
          method: "POST"
          accept: 200
          # You can send raw body instead of JSON-parsed or YAML-parsed body
          body: '"string"'
          log-response: true
      - name: "Use response"
        run: echo ${{ steps.api.outputs.response }}
```
