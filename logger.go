package crawler

import (
	"log"
	"os"
)

var (
	Debug *log.Logger
	Log *log.Logger
	Error *log.Logger
)

func LogInit() {
	format := log.Ldate | log.Ltime | log.Lshortfile
	Debug = log.New(os.Stdout, "[DEBUG]: ", format)
	Log = log.New(os.Stdout, "[INFO]: ", format)
	Error = log.New(os.Stderr, "[ERROR]: ", format|log.Llongfile)
}

