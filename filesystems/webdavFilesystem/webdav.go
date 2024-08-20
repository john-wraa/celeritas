package webdavFilesystem

import (
	"fmt"
	"github.com/john-wraa/celeritas/filesystems"
	"github.com/studio-b12/gowebdav"
	"io"
	"os"
	"path"
	"strings"
)

type WebDAV struct {
	Host string
	User string
	Pass string
}

// getCredentials generates sftp client using the credentials stored in the SFTP type
func (w *WebDAV) getCredentials() *gowebdav.Client {
	c := gowebdav.NewClient(w.Host, w.User, w.Pass)
	return c
}

// Put transfers a file to the remote file system
func (w *WebDAV) Put(fileName, folder string) error {
	client := w.getCredentials()

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = client.WriteStream(fmt.Sprintf("%s/%s", folder, path.Base(fileName)), file, 0664)
	if err != nil {
		return err
	}

	return nil
}

// List returns a listing of all files in the remote bucket with the given prefix, except for files named with a leading .
func (w *WebDAV) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing
	client := w.getCredentials()
	files, err := client.ReadDir(prefix)
	if err != nil {
		return listing, err
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			b := float64(file.Size())
			kb := b / 1024
			mb := kb / 1024
			current := filesystems.Listing{
				LastModified: file.ModTime(),
				Key:          file.Name(),
				Size:         mb,
				IsDir:        file.IsDir(),
			}
			listing = append(listing, current)
		}
	}

	return listing, nil
}

// Delete removes one or more files from the remote filesystem
func (w *WebDAV) Delete(itemsToDelete []string) bool {
	client := w.getCredentials()

	for _, item := range itemsToDelete {
		err := client.Remove(item)
		if err != nil {
			return false
		}
	}

	return true
}

// Get pulls a file from the remote file system and saves it somewhere on our server
func (w *WebDAV) Get(destination string, items ...string) error {
	client := w.getCredentials()

	for _, item := range items {
		err := func() error {
			webdavFilepath := item
			localFilepath := fmt.Sprintf("%s/%s", destination, path.Base(item))

			reader, err := client.ReadStream(webdavFilepath)
			if err != nil {
				return err
			}
			defer func(reader io.ReadCloser) {
				_ = reader.Close()
			}(reader)

			file, err := os.Create(localFilepath)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
