package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Logger struct {
	out      io.Writer  // 日志输出的目标
	minLevel Level      // 日志级别
	pretty   bool       // 是否以 JSON 格式输出日志
	mu       sync.Mutex // 互斥锁，用于保证线程安全
}

func New(out io.Writer, minLevel Level, pretty bool) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
		pretty:   pretty,
	}
}

// 打印日志
func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

// 打印错误日志
func (l *Logger) Error(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

// 打印错误日志，并退出程序
func (l *Logger) Fatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte
	var err error
	if l.pretty {
		line, err = json.MarshalIndent(aux, "", "  ") // 以 JSON 格式输出日志，两个空格缩进
		if err != nil {

		}
		if err != nil {
			line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
		}
		s := string(line)
		s = strings.ReplaceAll(s, "\\n", "\n    ") // 将 \n 替换为换行符和四个空格，便于观看
		s = strings.ReplaceAll(s, "\\t", "  ")     // 将 \t 替换为两个空格，便于观看
		line = []byte(s)
	} else {
		line, err = json.Marshal(aux)
		if err != nil {
			line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
