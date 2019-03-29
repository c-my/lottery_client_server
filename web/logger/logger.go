package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "Info: ", log.LstdFlags)
	Warning = log.New(os.Stdout, "Warning: ", log.LstdFlags|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, os.Stdout), "Error: ", log.LstdFlags|log.Lshortfile)
}
