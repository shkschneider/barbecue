package main

import(
	"fmt"
	"io"
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

func (l Log) log(lvl LogLevel, any interface{}) {
	t := time.Now().Format(fmt.Sprintf("%s %s.000", time.DateOnly, time.TimeOnly))
	if l.colors {
		l.writer.Write([]byte(fmt.Sprintf("\033[%dm%s [%3s] %+v\033[0m\n", lvl.Color, t, lvl.Tag, any)))
	} else {
		l.writer.Write([]byte(fmt.Sprintf("%s [%3s] %+v\n", t, lvl.Tag, any)))
	}
}

func (l Log) Debug(many ...interface{}) {
	for _, any := range many {
		l.log(LogLevelDebug, any)
	}
}

func (l Log) Info(many ...interface{}) {
	for _, any := range many {
		l.log(LogLevelInfo, any)
	}
}

func (l Log) Warning(many ...interface{}) {
	for _, any := range many {
		l.log(LogLevelWarning, any)
	}
}

func (l Log) Error(many ...interface{}) {
	for _, any := range many {
		l.log(LogLevelError, any)
	}
}

func (l Log) Panic(many ...interface{}) {
	if len(many) == 0 {
		panic("!!!")
	}
	for i := 0 ; i < len(many) - 1 ; i++ {
		l.log(LogLevelPanic, many[i])
	}
	panic(many[len(many) - 1])
}
