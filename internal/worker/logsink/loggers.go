// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package logsink

import (
	"context"
	"io"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/loggo/v2"
	"github.com/juju/worker/v4"
	"gopkg.in/tomb.v2"

	corelogger "github.com/juju/juju/core/logger"
	internallogger "github.com/juju/juju/internal/logger"
)

var (
	fallbackLogger = internallogger.GetLogger("logsink")
)

type modelLogger struct {
	tomb tomb.Tomb

	bufferedLogWriter *bufferedLogWriterCloser
	bufferedLogCloser func() error

	loggerContext corelogger.LoggerContext
}

// ModelLoggerConfig holds the configuration for a model logger.
type ModelLoggerConfig struct {
	MachineID     string
	NewLogWriter  corelogger.LogWriterForModelFunc
	BufferSize    int
	FlushInterval time.Duration
	Clock         clock.Clock
}

// NewModelLogger returns a new model logger instance.
func NewModelLogger(
	ctx context.Context,
	key corelogger.LoggerKey,
	config ModelLoggerConfig,
) (worker.Worker, error) {
	// Create a newLogWriter for the model.
	logger, err := config.NewLogWriter(ctx, key)
	if err != nil {
		return nil, errors.Annotatef(err, "getting logger for model %q", key.ModelName)
	}

	// Create a buffered log writer for the model, so that it correctly handles
	// the flushing of the logs to disk.
	bufferedLogWriter := &bufferedLogWriterCloser{
		BufferedLogWriter: corelogger.NewBufferedLogWriter(
			logger,
			config.BufferSize,
			config.FlushInterval,
			config.Clock,
		),
		closer:    logger,
		modelUUID: key.ModelUUID,
		machineID: config.MachineID,
	}

	// Create a new logger context for the model. This will use the buffered
	// log writer to write the logs to disk.
	loggerContext := internallogger.LoggerContext(corelogger.INFO)

	w := &modelLogger{
		bufferedLogWriter: bufferedLogWriter,
		bufferedLogCloser: sync.OnceValue(func() error {
			return bufferedLogWriter.Close()
		}),

		loggerContext: loggerContext,
	}

	if err := w.AddWriter("model-sink", bufferedLogWriter); err != nil {
		return nil, errors.Annotatef(err, "adding model-sink writer")
	}

	w.tomb.Go(w.loop)

	return w, nil
}

// Log writes the given log records to the logger's storage.
func (d *modelLogger) Log(records []corelogger.LogRecord) error {
	return d.bufferedLogWriter.Log(records)
}

// GetLogger returns a logger with the given name and tags.
func (d *modelLogger) GetLogger(name string, tags ...string) corelogger.Logger {
	return d.loggerContext.GetLogger(name, tags...)
}

// ConfigureLoggers configures loggers according to the given string
// specification, which specifies a set of modules and their associated
// logging levels. Loggers are colon- or semicolon-separated; each
// module is specified as <modulename>=<level>.  White space outside of
// module names and levels is ignored. The root module is specified
// with the name "<root>".
//
// An example specification:
//
//	<root>=ERROR; foo.bar=WARNING
//
// Label matching can be applied to the loggers by providing a set of labels
// to the function. If a logger has a label that matches the provided labels,
// then the logger will be configured with the provided level. If the logger
// does not have a label that matches the provided labels, then the logger
// will not be configured. No labels will configure all loggers in the
// specification.
func (d *modelLogger) ConfigureLoggers(specification string) error {
	return d.loggerContext.ConfigureLoggers(specification)
}

// ResetLoggerLevels iterates through the known logging modules and sets the
// levels of all to UNSPECIFIED, except for <root> which is set to WARNING.
// If labels are provided, then only loggers that have the provided labels
// will be reset.
func (d *modelLogger) ResetLoggerLevels() {
	d.loggerContext.ResetLoggerLevels()
}

// Config returns the current configuration of the Loggers. Loggers
// with UNSPECIFIED level will not be included.
func (d *modelLogger) Config() corelogger.Config {
	return d.loggerContext.Config()
}

// AddWriter adds a writer to the list to be called for each logging call.
// The name cannot be empty, and the writer cannot be nil. If an existing
// writer exists with the specified name, an error is returned.
//
// Note: we're relying on loggo.Writer here, until we do model level logging.
// Deprecated: This will be removed in the future and is only here whilst
// we cut things across.
func (d *modelLogger) AddWriter(name string, writer loggo.Writer) error {
	return d.loggerContext.AddWriter(name, writer)
}

// Close closes the model logger.
func (d *modelLogger) Close() error {
	return d.bufferedLogCloser()
}

// Kill stops the model logger.
func (d *modelLogger) Kill() {
	d.tomb.Kill(nil)
}

// Wait waits for the model logger to stop.
func (d *modelLogger) Wait() error {
	return d.tomb.Wait()
}

func (d *modelLogger) loop() error {
	// Close the buffered log writer when the model logger is stopped or killed.
	defer func() {
		_ = d.bufferedLogCloser()
	}()

	// Wait for the heat death of the universe.
	<-d.tomb.Dying()
	return tomb.ErrDying
}

type bufferedLogWriterCloser struct {
	*corelogger.BufferedLogWriter
	closer io.Closer

	modelUUID string
	machineID string
}

func (l *bufferedLogWriterCloser) Write(entry loggo.Entry) {
	err := l.Log([]corelogger.LogRecord{{
		Time:      entry.Timestamp,
		Entity:    "controller-" + l.machineID,
		Module:    entry.Module,
		Location:  filepath.Base(entry.Filename) + strconv.Itoa(entry.Line),
		Level:     corelogger.Level(entry.Level),
		Message:   entry.Message,
		Labels:    entry.Labels,
		ModelUUID: l.modelUUID,
	}})

	if err != nil {
		fallbackLogger.Warningf(context.Background(), "writing model logs failed for model %q, %v", l.modelUUID, err)
	}
}

func (b *bufferedLogWriterCloser) Close() error {
	err := errors.Trace(b.BufferedLogWriter.Flush())
	_ = b.closer.Close()
	return err
}
