package binding

import (
	"mime/multipart"
	"path/filepath"

	"github.com/gopi-frame/support/lists"
)

// UploadedFiles is a slice of [UploadedFile] instance
type UploadedFiles struct {
	*lists.List[*UploadedFile]
}

// NewUploadedFiles creates an instance of [UploadedFiles]
func NewUploadedFiles(fileHeaders []*multipart.FileHeader) *UploadedFiles {
	files := lists.NewList[*UploadedFile]()
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			panic(err)
		}
		if file, err := NewUploadedFile(file, fileHeader); err != nil {
			panic(err)
		} else {
			files.Push(file)
		}
	}
	return &UploadedFiles{files}
}

// Save saves all uploaded files under the given directory
func (uploadedFiles *UploadedFiles) Save(dirpath string) {
	uploadedFiles.Each(func(index int, uploadedFile *UploadedFile) bool {
		if err := uploadedFile.SaveAs(filepath.Join(dirpath, uploadedFile.Name())); err != nil {
			panic(err)
		}
		return true
	})
}

// Close closes all uploaded files
func (uploadedFiles *UploadedFiles) Close() (err error) {
	uploadedFiles.Each(func(index int, uploadedFile *UploadedFile) bool {
		if err = uploadedFile.Close(); err != nil {
			return false
		}
		return true
	})
	return
}
