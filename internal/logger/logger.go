package logger

import (
	"log"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *log.Logger
)

func InitializeLogger() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	})
	return logger
}

func getLoggerInstance() *log.Logger {
	return InitializeLogger()
}

func Info(message string) {
	getLoggerInstance().Println(message)
}

func Error(message string) {
	getLoggerInstance().Println(message)
}

func Fatal(message string) {
	getLoggerInstance().Println(message)
	os.Exit(1)
}

func Panic(message string) {
	getLoggerInstance().Println(message)
	panic(message)
}
