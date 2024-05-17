package response

import (
	"net/http"
	"os"
)

// FileResponse is used to send a file response
type FileResponse struct {
	*ReaderResponse
	filename string
}

// SetFile sets the filename
func (fileResponse *FileResponse) SetFile(filename string) *FileResponse {
	fileResponse.filename = filename
	return fileResponse
}

// Send reads the file content and sends it
func (fileResponse *FileResponse) Send(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(fileResponse.filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileResponse.SetReader(f)
	fileResponse.ReaderResponse.Send(w, r)
}
