package doubles

import (
	"io"
	"log"
)

func NewLoggerStub() *log.Logger {
	return log.New(io.Discard, "", 0)
}
