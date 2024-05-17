package response

import (
	"io"
	"net/http"
)

// StreamedResponse used to send a streamed response
type StreamedResponse struct {
	*Response
	step func(w io.Writer) bool
}

// SetStep sets the step func
func (streamed *StreamedResponse) SetStep(step func(w io.Writer) bool) *StreamedResponse {
	streamed.step = step
	return streamed
}

// Send sends the response
func (streamed *StreamedResponse) Send(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
	default:
		for {
			if !streamed.step(w) {
				break
			}
		}
		w.(http.Flusher).Flush()
	}
}
