package log

import (
	"log"
	"os"
)

var Debug = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
var Info = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var Warning = log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime|log.Lshortfile)
var Error = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var Critical = log.New(os.Stderr, "CRITICAL\t", log.Ldate|log.Ltime|log.Lshortfile)
