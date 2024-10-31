package loggers

import (
	"log"
	"os"
)

var LogErr = log.New(os.Stderr, "ERRORR\t", log.Lshortfile|log.Ltime|log.Ldate)
var LogInfo = log.New(os.Stdout, "INFOO\t", log.Lshortfile|log.Ltime|log.Ldate)
