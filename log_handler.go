// The nslog package provides non-structured logging, which is one of the implementation of "log/slog".
package nslog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const DEFAULT_LEVEL = slog.LevelInfo
const DEFAULT_TIME_LAYOUT = "2006/01/02 15:04:05"
const DEFAULT_SOURCE_LEVEL = slog.LevelWarn

type LogHandler struct {
	options LogHandlerOptions
	attrs   []slog.Attr
	groups  []string
	mutex   *sync.Mutex
	writer  io.Writer
}

// An option to customize output of log message.
type LogHandlerOptions struct {
	Level          slog.Leveler // Set level to output log message. By default, Error, Warn, and Info logs are output.
	AddColor       bool         // Add console color for level if it is true. (default: false)
	TimeLayout     string       // Set own time layout for [Time.Format]. (default: "2006/01/02 15:04:05")
	AddPID         bool         // Add PID as hex string if it is true. (default: false)
	AddGoroutineID bool         // Add Goroutine ID as hex string if it is true. (default: false)
	AddSourceLevel slog.Leveler // Set level to output log source, which is the file and line number that called the function. By default, Error and Warn logs are output with source.
	SourceFilePath bool         // Use filepath for source if it is true. Use filename for source if it is false.
}

// Create a new [slog.Logger] object that implements [nslog.LogHandler].
func NewLogger(writer io.Writer, options *LogHandlerOptions) *slog.Logger {
	handler := NewLogHandler(writer, options)
	return slog.New(handler)
}

// Create a new [nslog.LogHandler] object.
func NewLogHandler(writer io.Writer, options *LogHandlerOptions) *LogHandler {
	// set default parameters
	if options == nil {
		options = &LogHandlerOptions{}
	}
	if options.Level == nil {
		options.Level = DEFAULT_LEVEL
	}
	if options.TimeLayout == "" {
		options.TimeLayout = DEFAULT_TIME_LAYOUT
	}
	if options.AddSourceLevel == nil {
		options.AddSourceLevel = DEFAULT_SOURCE_LEVEL
	}

	// override parameters by environment variables
	switch os.Getenv("GO_NSLOG_LEVEL") {
	case "ERROR":
		options.Level = slog.LevelError
	case "WARN":
		options.Level = slog.LevelWarn
	case "INFO":
		options.Level = slog.LevelInfo
	case "DEBUG":
		options.Level = slog.LevelDebug
	default:
		// do not use environment variable for Level
	}
	nslogAddColor := os.Getenv("GO_NSLOG_ADD_COLOR")
	if strings.EqualFold(nslogAddColor, "false") || nslogAddColor == "0" {
		options.AddColor = false
	} else if strings.EqualFold(nslogAddColor, "true") || nslogAddColor == "1" {
		options.AddColor = true
	} else {
		// do not use environment variable for AddColor flag
	}
	nslogTimeLayout := os.Getenv("GO_NSLOG_TIME_LAYOUT")
	if nslogTimeLayout != "" {
		options.TimeLayout = nslogTimeLayout
	}
	nslogAddPID := os.Getenv("GO_NSLOG_ADD_PID")
	if strings.EqualFold(nslogAddPID, "false") || nslogAddPID == "0" {
		options.AddPID = false
	} else if strings.EqualFold(nslogAddPID, "true") || nslogAddPID == "1" {
		options.AddPID = true
	} else {
		// do not use environment variable for AddPID flag
	}
	nslogAddGoroutineID := os.Getenv("GO_NSLOG_ADD_GOROUTINEID")
	if strings.EqualFold(nslogAddGoroutineID, "false") || nslogAddGoroutineID == "0" {
		options.AddGoroutineID = false
	} else if strings.EqualFold(nslogAddGoroutineID, "true") || nslogAddGoroutineID == "1" {
		options.AddGoroutineID = true
	} else {
		// do not use environment variable for AddGoroutineID flag
	}
	switch os.Getenv("GO_NSLOG_ADD_SOURCE_LEVEL") {
	case "ERROR":
		options.AddSourceLevel = slog.LevelError
	case "WARN":
		options.AddSourceLevel = slog.LevelWarn
	case "INFO":
		options.AddSourceLevel = slog.LevelInfo
	case "DEBUG":
		options.AddSourceLevel = slog.LevelDebug
	default:
		// do not use environment variable for Source level
	}
	nslogSourceFilePath := os.Getenv("GO_NSLOG_SOURCE_FILE_PATH")
	if strings.EqualFold(nslogSourceFilePath, "false") || nslogSourceFilePath == "0" {
		options.SourceFilePath = false
	} else if strings.EqualFold(nslogSourceFilePath, "true") || nslogSourceFilePath == "1" {
		options.SourceFilePath = true
	} else {
		// do not use environment variable for SourceFilePath flag
	}

	return &LogHandler{
		options: *options,
		mutex:   &sync.Mutex{},
		writer:  writer,
	}
}

func (handler *LogHandler) clone() *LogHandler {
	return &LogHandler{
		options: handler.options,
		attrs:   slices.Clip(handler.attrs),
		groups:  slices.Clip(handler.groups),
		writer:  handler.writer,
		mutex:   handler.mutex,
	}
}

func (handler *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= handler.options.Level.Level()
}

func (handler *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	new_handler := handler.clone()
	new_handler.attrs = append(new_handler.attrs, attrs...)
	return new_handler
}

func (handler *LogHandler) WithGroup(name string) slog.Handler {
	new_handler := handler.clone()
	new_handler.groups = append(new_handler.groups, name)
	return new_handler
}

func (handler *LogHandler) Handle(_ context.Context, record slog.Record) error {
	// time
	time := record.Time.Format(handler.options.TimeLayout)

	// pid
	var pid int = 0
	if handler.options.AddPID {
		pid = os.Getpid()
	}

	// goroutineid
	var goroutineID uint64 = 0
	if handler.options.AddGoroutineID {
		b := make([]byte, 64)
		b = b[:runtime.Stack(b, false)]
		b = bytes.TrimPrefix(b, []byte("goroutine "))
		idField := b[:bytes.IndexByte(b, ' ')]
		goroutineID, _ = strconv.ParseUint(string(idField), 10, 64)
	}

	// level
	var level string
	switch record.Level {
	case slog.LevelError:
		if handler.options.AddColor {
			level = color.HiRedString("ERROR")
		} else {
			level = "ERROR"
		}
	case slog.LevelWarn:
		if handler.options.AddColor {
			level = color.HiYellowString("WARN.")
		} else {
			level = "WARN."
		}
	case slog.LevelInfo:
		if handler.options.AddColor {
			level = color.HiGreenString("INFO.")
		} else {
			level = "INFO."
		}
	case slog.LevelDebug:
		if handler.options.AddColor {
			level = color.HiCyanString("DEBUG")
		} else {
			level = "DEBUG"
		}
	default:
		level = "UNSET"
	}

	// withGroup
	withGroup := ""
	if len(handler.groups) > 0 {
		withGroup = strings.Join(handler.groups, ".")
	}

	// withAttributes
	var withAttributes []string
	for _, attribute := range handler.attrs {
		withAttributes = append(withAttributes, attribute.Key+"="+attribute.Value.String())
	}

	// with
	with := withGroup
	if len(withAttributes) > 0 {
		with += "[" + strings.Join(withAttributes, " ") + "]"
	}
	if with != "" {
		with += ":"
	}

	// message
	message := record.Message

	// attributes
	var attributes []string
	record.Attrs(func(attribute slog.Attr) bool {
		attributes = append(attributes, attribute.Key+"="+attribute.Value.String())
		return true
	})

	// source
	var source string
	if record.Level >= handler.options.AddSourceLevel.Level() {
		if record.PC != 0 {
			frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
			if handler.options.SourceFilePath {
				source = "(" + frame.File + ":" + strconv.Itoa(frame.Line) + ")"
			} else {
				source = "(" + filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line) + ")"
			}
		}
	}

	log_strings := []string{time}
	if pid > 0 {
		log_strings = append(log_strings, fmt.Sprintf("%04X", pid))
	}
	if goroutineID > 0 {
		log_strings = append(log_strings, fmt.Sprintf("%08X", goroutineID))
	}
	log_strings = append(log_strings, level)
	if with != "" {
		log_strings = append(log_strings, with)
	}
	log_strings = append(log_strings, message)
	if len(attributes) > 0 {
		log_strings = append(log_strings, strings.Join(attributes, " "))
	}
	if source != "" {
		log_strings = append(log_strings, source)
	}
	log_bytes := []byte(strings.Join(log_strings, " ") + "\n")

	handler.mutex.Lock()
	defer handler.mutex.Unlock()
	_, err := handler.writer.Write(log_bytes)
	return err
}
