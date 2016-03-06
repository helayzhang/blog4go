// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os"
	"time"
)

// ConsoleWriter is a console logger
type ConsoleWriter struct {
	blog *BLog

	closed bool

	colored bool

	// log hook
	hook      Hook
	hookLevel Level
}

// NewConsoleWriter initialize a console writer, singlton
func NewConsoleWriter() (consoleWriter *ConsoleWriter, err error) {
	singltonLock.Lock()
	defer singltonLock.Unlock()
	if nil != blog {
		consoleWriter, _ := blog.(*ConsoleWriter)
		return consoleWriter, ErrAlreadyInit
	}

	consoleWriter, err = newConsoleWriter()
	if nil != err {
		return nil, err
	}

	blog = consoleWriter
	go consoleWriter.daemon()
	return consoleWriter, nil
}

// newConsoleWriter initialize a console writer, not singlton
func newConsoleWriter() (consoleWriter *ConsoleWriter, err error) {
	consoleWriter = new(ConsoleWriter)
	consoleWriter.blog = NewBLog(os.Stdout)

	consoleWriter.closed = false

	consoleWriter.colored = false

	// log hook
	consoleWriter.hook = nil
	consoleWriter.hookLevel = DEBUG

	go consoleWriter.daemon()

	blog = consoleWriter
	return consoleWriter, nil
}

func (writer *ConsoleWriter) daemon() {
	f := time.Tick(10 * time.Second)

DaemonLoop:
	for {
		select {
		case <-f:
			if writer.closed {
				break DaemonLoop
			}

			writer.flush()
		}
	}
}

func (writer *ConsoleWriter) write(level Level, format string) {
	defer func() {
		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string) {
				writer.hook.Fire(level, format)
			}(level, format)
		}
	}()

	if writer.closed {
		return
	}

	writer.blog.write(level, format)
}

func (writer *ConsoleWriter) writef(level Level, format string, args ...interface{}) {
	defer func() {

		if nil != writer.hook && !(level < writer.hookLevel) {
			go func(level Level, format string, args ...interface{}) {
				writer.hook.Fire(level, fmt.Sprintf(format, args...))
			}(level, format, args...)
		}
	}()

	if writer.closed {
		return
	}

	writer.blog.writef(level, format, args...)
}

// Level return logging level threshold
func (writer *ConsoleWriter) Level() Level {
	return writer.blog.Level()
}

// SetLevel set logger level
func (writer *ConsoleWriter) SetLevel(level Level) {
	writer.blog.SetLevel(level)
}

// Colored return whether writer log with color
func (writer *ConsoleWriter) Colored() bool {
	return writer.colored
}

// SetColored set logging color
func (writer *ConsoleWriter) SetColored(colored bool) {
	if colored == writer.colored {
		return
	}

	writer.colored = colored

	initPrefix(colored)
}

// SetHook set hook for logging action
func (writer *ConsoleWriter) SetHook(hook Hook) {
	writer.hook = hook
}

// SetHookLevel set when hook will be called
func (writer *ConsoleWriter) SetHookLevel(level Level) {
	writer.hookLevel = level
}

// Close close console writer
func (writer *ConsoleWriter) Close() {
	if writer.closed {
		return
	}

	writer.blog.flush()
	writer.blog = nil
	writer.closed = true
}

// SetTimeRotated do nothing
func (writer *ConsoleWriter) SetTimeRotated(timeRotated bool) {
	return
}

// SetRetentions do nothing
func (writer *ConsoleWriter) SetRetentions(retentions int64) {
	return
}

// SetRotateSize do nothing
func (writer *ConsoleWriter) SetRotateSize(rotateSize ByteSize) {
	return
}

// SetRotateLines do nothing
func (writer *ConsoleWriter) SetRotateLines(rotateLines int) {
	return
}

// flush buffer to disk
func (writer *ConsoleWriter) flush() {
	writer.blog.flush()
}

// Debug debug
func (writer *ConsoleWriter) Debug(format string) {
	if DEBUG < writer.blog.Level() {
		return
	}

	writer.write(DEBUG, format)
}

// Debugf debugf
func (writer *ConsoleWriter) Debugf(format string, args ...interface{}) {
	if DEBUG < writer.blog.Level() {
		return
	}

	writer.writef(DEBUG, format, args...)
}

// Trace trace
func (writer *ConsoleWriter) Trace(format string) {
	if TRACE < writer.blog.Level() {
		return
	}

	writer.write(TRACE, format)
}

// Tracef tracef
func (writer *ConsoleWriter) Tracef(format string, args ...interface{}) {
	if TRACE < writer.blog.Level() {
		return
	}

	writer.writef(TRACE, format, args...)
}

// Info info
func (writer *ConsoleWriter) Info(format string) {
	if INFO < writer.blog.Level() {
		return
	}

	writer.write(INFO, format)
}

// Infof infof
func (writer *ConsoleWriter) Infof(format string, args ...interface{}) {
	if INFO < writer.blog.Level() {
		return
	}

	writer.writef(INFO, format, args...)
}

// Warn warn
func (writer *ConsoleWriter) Warn(format string) {
	if WARNING < writer.blog.Level() {
		return
	}

	writer.write(WARNING, format)
}

// Warnf warnf
func (writer *ConsoleWriter) Warnf(format string, args ...interface{}) {
	if WARNING < writer.blog.Level() {
		return
	}

	writer.writef(WARNING, format, args...)
}

// Error error
func (writer *ConsoleWriter) Error(format string) {
	if ERROR < writer.blog.Level() {
		return
	}

	writer.write(ERROR, format)
}

// Errorf errorf
func (writer *ConsoleWriter) Errorf(format string, args ...interface{}) {
	if ERROR < writer.blog.Level() {
		return
	}

	writer.writef(ERROR, format, args...)
}

// Critical critical
func (writer *ConsoleWriter) Critical(format string) {
	if CRITICAL < writer.blog.Level() {
		return
	}

	writer.write(CRITICAL, format)
}

// Criticalf criticalf
func (writer *ConsoleWriter) Criticalf(format string, args ...interface{}) {
	if CRITICAL < writer.blog.Level() {
		return
	}

	writer.writef(CRITICAL, format, args...)
}
