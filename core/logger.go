package core

import(
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

type Level struct {
	Id int
	Tag string
	Color uint8
}

var (
	LevelDebug 	 = Level { 1, "DBG", 34 }
	LevelInfo 	 = Level { 2, "INF", 32 }
	LevelWarning = Level { 3, "WRN", 33 }
	LevelError 	 = Level { 4, "ERR", 31 }
	LevelPanic 	 = Level { 5, "WTF", 35 }
)

type Logger struct {
	writer		io.Writer
	level		Level
	colors		bool
	datetime	bool
}

func NewLogger(writer io.Writer, level Level) *Logger {
	return &Logger { writer, level, true, false }
}

func (l *Logger) log(lvl Level, any string) {
	if lvl.Id < l.level.Id { return }
	var buffer bytes.Buffer
	if l.colors {
		buffer.WriteString(fmt.Sprintf("\033[%dm", lvl.Color))
	}
	if l.datetime {
		buffer.WriteString(fmt.Sprintf("%s [%3s] %s",
			time.Now().Format(fmt.Sprintf("%s %s.000", time.DateOnly, time.TimeOnly)), lvl.Tag, any))
	} else {
		buffer.WriteString(fmt.Sprintf("[%3s] %s", lvl.Tag, any))
	}
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

func (l *Logger) Debug(many ...interface{}) {
	l.log(LevelDebug, any(many))
}

func (l *Logger) Info(many ...interface{}) {
	l.log(LevelInfo, any(many))
}

func (l *Logger) Warning(many ...interface{}) {
	l.log(LevelWarning, any(many))
}

func (l *Logger) Error(many ...interface{}) {
	l.log(LevelError, any(many))
}

func (l *Logger) Panic(many ...interface{}) {
	l.log(LevelPanic, any(many))
	if len(many) == 0 {
		panic("!!!")
	} else {
		panic(any(many))
	}
}
