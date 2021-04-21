package doubles

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

// NewFileSystemStub returns a new fs.FS stub with provided files.
func NewFileSystemStub(files map[string][]byte) *filesystemStub {
	return &filesystemStub{
		files: files,
	}
}

// implements type FS interface
type filesystemStub struct {
	files map[string][]byte
}

func (f *filesystemStub) Open(name string) (fs.File, error) {
	if name == "badfile.txt" {
		return &badFile{}, nil
	}

	tmpfile, err := ioutil.TempFile("", "fsdemo")
	if err != nil {
		log.Fatal(err)
	}

	contents, ok := f.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	tmpfile.Write([]byte(contents))
	tmpfile.Seek(0, 0)

	return tmpfile, nil
}

type badFile struct{}

func (f *badFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("bad file")
}

func (f *badFile) Read(_ []byte) (int, error) {
	return 0, errors.New("bad file")
}

func (f *badFile) Close() error {
	return errors.New("bad file")
}
