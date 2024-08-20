package celeritas

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/tsawler/celeritas/filesystems"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func (c *Celeritas) UploadFile(r *http.Request, destination, field string, fs filesystems.FS) error {
	fileName, err := c.getFileToUpload(r, field)
	if err != nil {
		c.ErrorLog.Println(err)
		return err
	}

	if fs != nil {
		// Uploading to a remote filesystem
		err = fs.Put(fileName, destination)
		if err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	} else {
		// Uploading to a (server) local filesystem
		err = os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			c.ErrorLog.Println(err)
			return err
		}
	}

	defer func() {
		_ = os.Remove(fileName)
	}()

	return nil
}

func (c *Celeritas) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(c.config.uploads.maxUploadSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// Looks at the first ~500 bytes to try to figure out the correct mimetype
	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}
	// Return to start of file after mimetype check
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	if !inSlice(c.config.uploads.allowedMimeTypes, mimeType.String()) {
		return "", fmt.Errorf("invalid file mime type: %s", mimeType.String())
	}

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

func inSlice(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
