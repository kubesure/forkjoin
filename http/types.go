package http

//METHOD http methods supported by http dispatcher
type METHOD string

//METHOD http methods supported by http dispatcher
const (
	GET   METHOD = "GET"
	POST         = "POST"
	PUT          = "PUT"
	PATCH        = "PATCH"
)

type AUTHENTICATION string

const (
	NONE  AUTHENTICATION = "NONE"
	BASIC AUTHENTICATION = "BASIC"
)

type BasicAuthCredentials struct {
	UserName          string
	Password          string
	ServerCertificate string
}

//HTTPRequest URL and method to be dispatched too
type HTTPRequest struct {
	ID      string
	Message HTTPMessage
}

//HTTPResponse URL and method to be dispatched too
type HTTPResponse struct {
	Message HTTPMessage
}

//HTTPMessage URL and method to be dispatched too
type HTTPMessage struct {
	ID                   string
	URL                  string
	Method               METHOD
	Payload              string
	Headers              map[string]string
	ActiveDeadLine       uint32
	StatusCode           int
	Authentication       AUTHENTICATION
	BasicAtuhCredentials BasicAuthCredentials
}

//Add adds headers to messsage
func (hm *HTTPMessage) Add(key, value string) {
	if hm.Headers == nil {
		hm.Headers = make(map[string]string)
	}
	hm.Headers[key] = value
}
