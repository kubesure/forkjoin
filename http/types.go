package http

// METHOD represents HTTP methods supported by the HTTP dispatcher.
type METHOD string

// Predefined HTTP methods supported by the HTTP dispatcher.
// These constants represent the most commonly used HTTP methods.
const (
	GET   METHOD = "GET"   // HTTP GET method, typically used to retrieve data from a server.
	POST  METHOD = "POST"  // HTTP POST method, typically used to send data to a server.
	PUT   METHOD = "PUT"   // HTTP PUT method, typically used to update or replace data on a server.
	PATCH METHOD = "PATCH" // HTTP PATCH method, typically used to partially update data on a server.
)

// AUTHENTICATION represents the type of authentication used for HTTP requests.
type AUTHENTICATION string

// Predefined authentication types supported by the HTTP dispatcher.
// These constants define the available authentication mechanisms.
const (
	NONE   AUTHENTICATION = "NONE"   // No authentication required.
	BASIC  AUTHENTICATION = "BASIC"  // Basic authentication using username and password.
	MUTUAL AUTHENTICATION = "MUTUAL" // Mutual authentication using certificates.
)

// BasicAuthCredentials represents the credentials required for basic authentication.
type BasicAuthCredentials struct {
	UserName          string // The username for basic authentication.
	Password          string // The password for basic authentication.
	ServerCertificate string // The server certificate for secure communication.
}

// MutualAuthCredentials represents the credentials required for mutual authentication.
type MutualAuthCredentials struct {
	ClientCertificate string // The client certificate for mutual authentication.
	ClientKey         string // The client private key for mutual authentication.
	CACertificate     string // The CA certificate for verifying the server.
}

// HTTPRequest represents an HTTP request to be dispatched.
// It contains an ID and the HTTP message to be sent.
type HTTPRequest struct {
	ID      string      // Unique identifier for the HTTP request.
	Message HTTPMessage // The HTTP message containing details of the request.
}

// HTTPResponse represents an HTTP response received from the server.
// It contains the HTTP message with the response details.
type HTTPResponse struct {
	Message HTTPMessage // The HTTP message containing details of the response.
}

// HTTPMessage represents the details of an HTTP request or response.
// It includes the URL, method, payload, headers, and authentication details.
type HTTPMessage struct {
	ID                    string                // Unique identifier for the HTTP message.
	URL                   string                // The URL to which the HTTP request is sent.
	Method                METHOD                // The HTTP method (e.g., GET, POST).
	Payload               string                // The payload or body of the HTTP request.
	Headers               map[string]string     // The headers to be included in the HTTP request.
	ActiveDeadLine        uint32                // The maximum time (in seconds) the request can remain active.
	StatusCode            int                   // The HTTP status code received in the response.
	Authentication        AUTHENTICATION        // The type of authentication used for the request.
	BasicAtuhCredentials  BasicAuthCredentials  // Credentials for basic authentication.
	MutualAuthCredentials MutualAuthCredentials // Credentials for mutual authentication.
}

// Add adds a header to the HTTP message.
// If the Headers map is nil, it initializes the map before adding the header.
//
// Parameters:
//   - key: The header name.
//   - value: The header value.
func (hm *HTTPMessage) Add(key, value string) {
	if hm.Headers == nil {
		hm.Headers = make(map[string]string)
	}
	hm.Headers[key] = value
}
