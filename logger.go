package main

import (
	"log"
)

type Logger interface {
	Log(message string)
}

type StdLogger struct{}

func (s *StdLogger) Log(message string) {
	log.Print(message)

}
