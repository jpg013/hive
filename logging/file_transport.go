package logging

import (
	"fmt"
	"log"
	"os"
)

// FileOut represents ioWriter for file transport.
type FileOut struct {
	handle *os.File
}

// Write method takes byte array and writes to log file
func (out *FileOut) Write(p []byte) (n int, err error) {
	log.SetOutput(out.handle)
	log.Println(string(p))

	return 0, nil
}

func createIfNotExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0755)
	}
}

// NewFileTransport configures and returns a file transport
func NewFileTransport(conf *TransportConfig) Transport {
	filePath := conf.FilePath
	fileName := conf.FileName

	if filePath == "" {
		wd, _ := os.Getwd()
		filePath = fmt.Sprintf("%v/logs", wd)
	}

	if fileName == "" {
		fileName = "exceptions.log"
	}

	createIfNotExists(filePath)

	// Open log file
	f, err := os.OpenFile(fmt.Sprintf("%v/%v", filePath, fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	fileOut := &FileOut{
		handle: f,
	}

	return Transport{
		Level: conf.Level,
		Out:   fileOut,
	}
}
