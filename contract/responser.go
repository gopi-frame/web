package contract

import "net/http"

// Responser the response interface
type Responser interface {
	// SetStatusCode sets the response http status code
	SetStatusCode(statusCode int)
	// StatusCode gets the response http status code
	StatusCode() int
	// SetContent sets the response body content
	SetContent(content any)
	// Content gets the response body content
	Content() any
	// SetHeader sets the response header, if replace is true, it will replace the existing header,
	// and if replace is false, it appends the new value into existing header
	SetHeader(key, header string, replace ...bool)
	// SetHeaders sets headers map to the response
	SetHeaders(headers map[string]string)
	// HasHeader returns if the specific key is exist
	HasHeader(key string) bool
	// Header returns header value of specific header
	Header(key string) string
	// Headers returns all headers as [http.Header]
	Headers() http.Header
	// SetCookie sets cookie to response
	SetCookie(cookie *http.Cookie)
	// Cookies returns all response cookies
	Cookies() []*http.Cookie
	// Send sends the response
	Send(w http.ResponseWriter, r *http.Request)
}
