// pkg/utils/logger.go

package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// TODO: Panic and Fatal loggers should write to a file
	FatalLogger = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

}
