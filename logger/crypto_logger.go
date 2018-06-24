package logger

import (
	"time"
	"log"
)

type CryptoLogger struct {
}

type Event struct {
	Id         string
	EventType  string
	MarketName string
	EventTime  time.Time
}

type SearchEvent struct {
}

func NewCryptoLogger() *CryptoLogger {
	logger := &CryptoLogger{}

	return logger
}

func (l *CryptoLogger) PushEvent(event Event) error {

	log.Println(event)

	return nil
}

func (l *CryptoLogger) SearchEvent(event SearchEvent) error {

	return nil
}
