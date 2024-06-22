package infrastructure

import (
	"io"
	"log"
	"metrix/internal/usecases"
	"os"
)

const (
	FilePerm600 os.FileMode = 0o600
)

type Logger struct{}

func NewLogger() usecases.Logger {
	return &Logger{}
}

func (l *Logger) log(filePath string, format string, v ...any) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, FilePerm600)
	if err != nil {
		log.Printf("%s", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close log error file")
		}
	}()

	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	log.Printf(format, v...)
}

func (l *Logger) LogError(format string, v ...any) {
	l.log("./log/error.log", format, v...)
}

func (l *Logger) LogAccess(format string, v ...any) {
	l.log("./log/access.log", format, v...)
}
