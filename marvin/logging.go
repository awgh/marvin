package marvin

import (
	"log"

	irc "github.com/fluffle/goirc/client"
)

type sLogger struct{}

// LinePrinter prints a line for debugging
func LinePrinter(line *irc.Line) {
	log.Println("Public:", line.Public())
	log.Println("Target:", line.Target())
	log.Println("Text:", line.Text())
	log.Println("Args:", line.Args)
	log.Println("Cmd:", line.Cmd)
	log.Println("Host:", line.Host)
	log.Println("Ident:", line.Ident)
	log.Println("Nick:", line.Nick)
	log.Println("Raw:", line.Raw)
	log.Println("Src:", line.Src)
	log.Println("Tags:", line.Tags)
	log.Println("Time:", line.Time)
}

func (s sLogger) Debug(f string, a ...interface{}) {
	log.Printf(f, a...)
}
func (s sLogger) Info(f string, a ...interface{}) {
	log.Printf(f, a...)
}
func (s sLogger) Warn(f string, a ...interface{}) {
	log.Printf(f, a...)
}
func (s sLogger) Error(f string, a ...interface{}) {
	log.Printf(f, a...)
}
