package contract

import "net/http"

// Resolver resolver interface
type Resolver interface {
	Resolve(*http.Request, any) error
}
