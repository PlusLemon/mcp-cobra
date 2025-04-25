package mcp

import (
	"log"
	"os"
)

var globalLogger *FileLogger

func init() {
	globalLogger, _ = NewFileLogger("/tmp/mcp-cobra.log")
	if globalLogger == nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Panicf("Failed to get user home directory: %v", err)
			return
		}
		globalLogger, _ = NewFileLogger(homeDir + "/mcp-cobra.log")
		if globalLogger == nil {
			log.Panicf("Failed to create log file")
		}
	}
}

type FileLogger struct {
	file   *os.File
	logger *log.Logger
}

func NewFileLogger(path string) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, err
	}
	l := log.New(f, "", log.LstdFlags)
	return &FileLogger{file: f, logger: l}, nil
}

func (fl *FileLogger) WriteLog(msg string) {
	fl.logger.Println(msg)
}

func (fl *FileLogger) Close() error {
	return fl.file.Close()
}

func CloseGlobalLogger() {
	globalLogger.Close()
}

func LogInfo(msg string) {
	globalLogger.WriteLog(msg)
}
