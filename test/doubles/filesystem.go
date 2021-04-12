package doubles

import (
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
