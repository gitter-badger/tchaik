// Copyright 2015, David Howden
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// FileSystem is an extension of the http.FileSystem interface
// which includes the RemoteOpen method which returns a *File
// from which the file data can be fetched.
type FileSystem interface {
	http.FileSystem
	RemoteOpen(string) (*File, error)
}

// fileSystem implements FileSystem
type fileSystem struct {
	client Client
}

// NewFileSystem creates a new file system using the given Client to handle
// file requests.
func NewFileSystem(c Client) *fileSystem {
	return &fileSystem{
		client: c,
	}
}

// file is a basic representation of a remote file such that all operations (i.e.
// seeking) will work correctly.
type file struct {
	io.ReadSeeker
	stat *fileInfo
}

// RemoteOpen returns a *File which represents the remote file, and implements
// io.ReadCloser which reads the file contents from the remote system.
func (fs *fileSystem) RemoteOpen(path string) (*File, error) {
	rf, err := fs.client.Get(path)
	if err != nil {
		return nil, err
	}
	return rf, nil
}

// Open the given file and return an http.File implementation representing it.  This method
// will block until the file has been completely fetched (http.File implements io.Seeker
// which means that for a trivial implementation we need all the underlying data).
func (fs *fileSystem) Open(path string) (http.File, error) {
	rf, err := fs.RemoteOpen(path)
	if err != nil {
		return nil, err
	}
	defer rf.Close()

	buf, err := ioutil.ReadAll(rf)
	if err != nil {
		return nil, err
	}

	return &file{
		ReadSeeker: bytes.NewReader(buf),
		stat: &fileInfo{
			name:    rf.Name,
			size:    rf.Size,
			modTime: rf.ModTime,
		},
	}, nil
}

// Close is a nop as we have already closed the original file.
func (f *file) Close() error {
	return nil
}

// Implements http.File.
func (f *file) Readdir(int) ([]os.FileInfo, error) {
	return nil, nil
}

// FileInfo is a simple implementation of os.FileInfo.
type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) Size() int64        { return f.size }
func (f *fileInfo) Mode() os.FileMode  { return os.FileMode(0777) }
func (f *fileInfo) ModTime() time.Time { return f.modTime }
func (f *fileInfo) IsDir() bool        { return false }
func (f *fileInfo) Sys() interface{}   { return nil }

func (f *file) Stat() (os.FileInfo, error) {
	return f.stat, nil
}
