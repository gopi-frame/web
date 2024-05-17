package binding

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

// NewUploadedFile creates an instance of [UploadedFile]
func NewUploadedFile(file multipart.File, fileHeader *multipart.FileHeader) (*UploadedFile, error) {
	mime, err := mimetype.DetectReader(file)
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}
	uploadedFile := &UploadedFile{
		fileHeader: fileHeader,
		file:       file,
		mime:       mime,
	}
	return uploadedFile, err
}

// UploadedFile is an object contains the information of file which is uploaded through multipart-form
// it alse provides some functions to operate the file uploaded
type UploadedFile struct {
	fileHeader *multipart.FileHeader
	file       multipart.File
	mime       *mimetype.MIME
	content    *[]byte
}

// Name returns the name of the uploaded file
func (uploadedFile *UploadedFile) Name() string {
	return uploadedFile.fileHeader.Filename
}

// ClientExtension returns the ext from the name of the uploaded file
func (uploadedFile *UploadedFile) ClientExtension() string {
	return filepath.Ext(uploadedFile.fileHeader.Filename)
}

// ClientMimeType returns the mimetype of the uploaded file from "Content-Type" header
func (uploadedFile *UploadedFile) ClientMimeType() string {
	return uploadedFile.fileHeader.Header.Get("Content-type")
}

// MimeType returns the mimetype detected from file content
func (uploadedFile *UploadedFile) MimeType() string {
	return uploadedFile.mime.String()
}

// Extension returns the ext of mimetype detected from file content
func (uploadedFile *UploadedFile) Extension() string {
	return uploadedFile.mime.Extension()
}

// Size returns the size of the uploaded file
func (uploadedFile *UploadedFile) Size() int64 {
	return uploadedFile.fileHeader.Size
}

// File get multipart.File
func (uploadedFile *UploadedFile) File() multipart.File {
	return uploadedFile.file
}

// Content returns the content of the uploaded file
func (uploadedFile *UploadedFile) Content() ([]byte, error) {
	if uploadedFile.content != nil {
		return *uploadedFile.content, nil
	}
	content := make([]byte, 0, uploadedFile.fileHeader.Size)
	scanner := bufio.NewScanner(uploadedFile.file)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return content, err
	}
	return content, nil
}

// SaveAs saves the uploaded file as the given filename
func (uploadedFile *UploadedFile) SaveAs(filename string) error {
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(dst, uploadedFile.file); err != nil {
		return err
	}
	return nil
}

// Close closes the file resource
func (uploadedFile *UploadedFile) Close() error {
	return uploadedFile.file.Close()
}
