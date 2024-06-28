package controller

import (
	"io"
	"net/http"

	"github.com/gopi-frame/web/request"
	"github.com/gopi-frame/web/response"
)

// Controller basic controller implemention
type Controller struct {
	*request.Request
}

// Init inits controller
func (controller *Controller) Init(request *request.Request) {
	controller.Request = request
}

// Response returns a basic response
func (controller *Controller) Response(statusCode int, content ...any) *response.Response {
	return response.New(statusCode, content...)
}

// JSON returns a json response
func (controller *Controller) JSON(statusCode int, content ...any) *response.JSONResponse {
	return controller.Response(statusCode, content...).JSON()
}

// XML returns a xml response
func (controller *Controller) XML(statusCode int, content ...any) *response.XMLResponse {
	return controller.Response(statusCode, content...).XML()
}

// Reader returns a reader response
func (controller *Controller) Reader(statusCode int, r io.Reader) *response.ReaderResponse {
	return controller.Response(statusCode).Reader(r)
}

// Redirect returns a redirect response
func (controller *Controller) Redirect(location string, statusCode ...int) *response.RedirectResponse {
	code := http.StatusFound
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	return controller.Response(code).Redirect(location)
}

// File returns a file response
func (controller *Controller) File(statusCode int, file string) *response.FileResponse {
	return controller.Response(statusCode).File(file)
}

// Stream returns a streamed response
func (controller *Controller) Stream(step func(io.Writer) bool) *response.StreamedResponse {
	return (&response.StreamedResponse{Response: &response.Response{}}).SetStep(step)
}
