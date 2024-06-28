package request

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/support/maps"
	"github.com/gopi-frame/types"
	validationcontract "github.com/gopi-frame/validation/contract"
	"github.com/gopi-frame/web/binding"
	"github.com/gopi-frame/web/contract"
	"github.com/gopi-frame/web/mimetype"
	"github.com/julienschmidt/httprouter"
)

// NewRequest creates a new [Request] instance with http.Request and httprouter.Params
func NewRequest(r *http.Request, p httprouter.Params) *Request {
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if r.Header.Get("content-type") == mimetype.FormData {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				panic(err)
			}
		} else {
			if err := r.ParseForm(); err != nil {
				panic(err)
			}
		}
	}

	req := &Request{
		Request: r,
		Params:  p,
		Values:  maps.NewMap[string, any](),
	}
	return req
}

// Request http request
type Request struct {
	Request *http.Request
	Params  httprouter.Params
	Values  *maps.Map[string, any]
	form    validationcontract.Form
	locale  string
}

// Clone clones a new [Request] instance from current one
func (r *Request) Clone() *Request {
	return &Request{
		Request: r.Request,
		Params:  r.Params,
		Values:  r.Values,
	}
}

// Set sets a value with specific key to current request
func (r *Request) Set(key string, value any) {
	r.Values.Lock()
	defer r.Values.Unlock()
	r.Values.Set(key, value)
}

// Get gets a value from current request by specific key
//
// if the specific key is not exist, it returns nil and false
func (r *Request) Get(key string) (any, bool) {
	r.Values.RLock()
	defer r.Values.RUnlock()
	if ok := r.Values.ContainsKey(key); ok {
		return r.Values.Get(key)
	}
	return nil, false
}

// MustGet gets a value from current request by specific key
//
// if the specific key is not exist, it will panic
func (r *Request) MustGet(key string) any {
	if value, exists := r.Get(key); !exists {
		panic(fmt.Errorf("key \"%s\" does not exists in context Values", key))
	} else {
		return value
	}
}

// Method returns the http request method
func (r *Request) Method() string {
	return r.Request.Method
}

// IsGet returns if the http request method is GET
func (r *Request) IsGet() bool {
	return r.Method() == http.MethodGet
}

// IsPost returns if the http request method is POST
func (r *Request) IsPost() bool {
	return r.Method() == http.MethodPost
}

// IsPut returns if the http request method is PUT
func (r *Request) IsPut() bool {
	return r.Method() == http.MethodPut
}

// IsPatch returns if the http request method is PATCH
func (r *Request) IsPatch() bool {
	return r.Method() == http.MethodPatch
}

// IsDelete returns if the http request method is DELETE
func (r *Request) IsDelete() bool {
	return r.Method() == http.MethodDelete
}

// IsHead returns if the http request method is HEAD
func (r *Request) IsHead() bool {
	return r.Method() == http.MethodHead
}

// IsConnect returns if the http request method is CONNECT
func (r *Request) IsConnect() bool {
	return r.Method() == http.MethodConnect
}

// IsOptions returns if the http request method is OPTIONS
func (r *Request) IsOptions() bool {
	return r.Method() == http.MethodOptions
}

// IsTrace returns if the http request method is GET
func (r *Request) IsTrace() bool {
	return r.Method() == http.MethodTrace
}

// Host returns the request's host
func (r *Request) Host() string {
	return r.Request.Host
}

// RequestURI returns the request's uri
func (r *Request) RequestURI() string {
	return r.Request.RequestURI
}

// Path returns the request's url path
func (r *Request) Path() string {
	return r.Request.URL.Path
}

// Query returns the query value by specific key
//
// if the specific key is not exist, default value will be returned
//
// if default value is not provided, nil will be returned
func (r *Request) Query(key string, defaultValue ...string) types.String {
	if values, exists := r.QueryArray(key); exists {
		return values.Get(0)
	} else if len(defaultValue) > 0 {
		return types.String(defaultValue[0])
	} else {
		return ""
	}
}

// QueryArray returns the query values by specific key
//
// if the specific key is not exist, the second return value will be false
func (r *Request) QueryArray(key string) (*lists.List[types.String], bool) {
	values := lists.NewList[types.String]()
	if r.Request.URL.Query().Has(key) {
		for _, value := range r.Request.URL.Query()[key] {
			values.Push(types.String(value))
		}
		return values, true
	}
	return values, false
}

// QueryMap returns the query values by specific key as a map
//
// if the specific key is not exist, the second return value will be false
func (r *Request) QueryMap(key string) (*maps.Map[string, types.String], bool) {
	queries := r.Request.URL.Query()
	result := maps.NewMap[string, types.String]()
	exists := false
	for k, values := range queries {
		if i := strings.IndexByte(k, '['); i > 0 && k[:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j > 0 {
				exists = true
				result.Set(k[i+1:][:j], types.String(values[0]))
			}
		}
	}
	return result, exists
}

// PostForm returns the form value by specific key
//
// if the specific key is not exist, default value will be returned
//
// if default value is not provided, nil will be returned
func (r *Request) PostForm(key string, defaultValue ...string) types.String {
	if values, exists := r.PostFormArray(key); exists {
		return values.Get(0)
	} else if len(defaultValue) > 0 {
		return types.String(defaultValue[0])
	} else {
		return ""
	}
}

// PostFormArray returns the form values by specific key
//
// if the specific key is not exist, the second return value will be false
func (r *Request) PostFormArray(key string) (*lists.List[types.String], bool) {
	values := lists.NewList[types.String]()
	if r.Request.PostForm.Has(key) {
		for _, value := range r.Request.PostForm[key] {
			values.Push(types.String(value))
		}
		return values, true
	}
	return values, false
}

// PostFormMap returns the form values by specific key as a map
//
// if the specific key is not exist, the second return value will be false
func (r *Request) PostFormMap(key string) (*maps.Map[string, types.String], bool) {
	posts := r.Request.PostForm
	result := maps.NewMap[string, types.String]()
	exists := false
	for k, values := range posts {
		if i := strings.IndexByte(k, '['); i > 0 && k[:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j > 0 {
				exists = true
				result.Set(k[i+1:][:j], types.String(values[0]))
			}
		}
	}
	return result, exists
}

// Param returns path param value of the specific key
func (r *Request) Param(key string) types.String {
	return types.String(r.Params.ByName(key))
}

// File returns an instance of [formdata.UploadedFile] of the specific name
func (r *Request) File(name string) (*binding.UploadedFile, error) {
	file, fileHeader, err := r.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	return binding.NewUploadedFile(file, fileHeader)
}

// Files returns an instance of [binding.UploadedFiles] of the specific name
func (r *Request) Files(name string) *binding.UploadedFiles {
	if err := r.Request.ParseMultipartForm(32 << 20); err != nil {
		panic(err)
	}
	fileHeaders := r.Request.MultipartForm.File[name]
	return binding.NewUploadedFiles(fileHeaders)
}

// Header returns the header value by specific key
//
// if the specific key is not exist, default value will be returned
//
// if default value is not provided, nil will be returned
func (r *Request) Header(key string, defaultValue ...string) types.String {
	if headers, exists := r.HeaderArray(key); exists {
		return headers.Get(0)
	} else if len(defaultValue) > 0 {
		return types.String(defaultValue[0])
	} else {
		return ""
	}
}

// HeaderArray returns the header values by specific key
//
// if the specific key is not exist, the second return value will be false
func (r *Request) HeaderArray(key string) (*lists.List[types.String], bool) {
	headers := r.Request.Header.Values(key)
	values := lists.NewList[types.String]()
	for _, header := range headers {
		values.Push(types.String(header))
	}
	return values, len(headers) > 0
}

// ClientIP returns the ip address of client
func (r *Request) ClientIP() string {
	xForwardedFor := r.Header("X-Forwarded-For", "").String()
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header("X-Real-IP", "").String())
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(r.Request.RemoteAddr); err == nil {
		return ip
	}
	return ""
}

// Bind parses request into an instance of [validationcontract.Form]
//
// if bindings is not provided, it uses [binding.Form] and
// binding implements according to Content-Type header
func (r *Request) Bind(form validationcontract.Form, bindings ...contract.Resolver) error {
	if len(bindings) == 0 {
		bindings = append(bindings, binding.Form)
		if h := r.Header("Content-Type"); h != "" {
			contentType := h.String()
			if contentType == mimetype.JSON {
				bindings = append(bindings, binding.JSON)
			} else if contentType == mimetype.XML {
				bindings = append(bindings, binding.XML)
			} else if contentType == mimetype.FormData || contentType == mimetype.FormURLEncode {
				bindings = append(bindings, binding.Form)
			}
		}
	}
	for _, binding := range bindings {
		if err := binding.Resolve(r.Request, form); err != nil {
			return err
		}
	}
	r.form = form
	return nil
}

// Form returns the validated instance of [validationcontract.Form]
func (r *Request) Form() validationcontract.Form {
	return r.form
}

// SetLocale set locale
func (r *Request) SetLocale(locale string) {
	r.locale = locale
}

// Locale get locale
func (r *Request) Locale() string {
	return r.locale
}
