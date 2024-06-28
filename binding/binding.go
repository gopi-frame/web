package binding

import "net/http"

// Binding is the alias of [Parser]
type Binding func(r *http.Request, dest any) error

// Resolve resolve
func (b Binding) Resolve(r *http.Request, dest any) error {
	return b(r, dest)
}
