package main

import(
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// LogLevel

type LogLevel struct {
	Tag string
	Color uint8
}

var (
	LogLevelDebug = LogLevel { "DBG", 34 }
	LogLevelInfo = LogLevel { "INF", 32 }
	LogLevelWarning = LogLevel { "WRN", 33 }
	LogLevelError = LogLevel { "ERR", 31 }
	LogLevelPanic = LogLevel { "WTF", 35 }
)

// Log

type Log struct {
	writer	io.Writer
	level	LogLevel
	colors	bool
}

var	log Log

func NewLog(writer io.Writer, level LogLevel, colors bool) {
	log = Log { writer, level, colors }
}

func (l Log) log(lvl LogLevel, any string) {
	var buffer bytes.Buffer
	if l.colors {
		buffer.WriteString(fmt.Sprintf("\033[%dm", lvl.Color))
	}
	buffer.WriteString(fmt.Sprintf("%s [%3s] %s",
		time.Now().Format(fmt.Sprintf("%s %s.000", time.DateOnly, time.TimeOnly)), lvl.Tag, any))
	if l.colors {
		buffer.WriteString(fmt.Sprintf("\033[0m"))
	}
	buffer.WriteString("\n")
	l.writer.Write(buffer.Bytes())
}

func any(many interface{}) string {
	any := fmt.Sprintf("%+v", many)
	any = strings.TrimPrefix(any, "[")
	any = strings.TrimSuffix(any, "]")
	return any
}

func (l Log) Debug(many ...interface{}) {
	l.log(LogLevelDebug, any(many))
}

func (l Log) Info(many ...interface{}) {
	l.log(LogLevelInfo, any(many))
}

func (l Log) Warning(many ...interface{}) {
	l.log(LogLevelWarning, any(many))
}

func (l Log) Error(many ...interface{}) {
	l.log(LogLevelError, any(many))
}

func (l Log) Panic(many ...interface{}) {
	l.log(LogLevelPanic, any(many))
	if len(many) == 0 {
		panic("!!!")
	} else {
		panic(any(many))
	}
}
