// Package logs provides a high-performance, context-aware structured logging library.
//
// Features:
//   - Zero-allocation hot paths using sync.Pool
//   - Context-aware logging with context.Context
//   - Type-safe field builders
//   - Multiple output formats (text, JSON, pretty)
//   - Sampling for high-volume logs
//   - Async logging option
//   - Hook system for extensibility
//   - Built-in caller information
//   - Chained/fluent API
//
// Basic usage:
//
//	log := logs.New()
//	log.Info("server started", logs.Int("port", 8080))
//
// With context:
//
//	log.InfoContext(ctx, "request processed", logs.Duration("latency", time.Since(start)))
package logs

import (
	"context"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Logger is the main logging interface.
type Logger struct {
	output      io.Writer
	level       atomic.Int32
	formatter   Formatter
	hooks       []Hook
	fields      []Field
	callerDepth int
	addCaller   bool
	addStack    bool
	async       bool
	asyncCh     chan *Entry
	asyncWg     sync.WaitGroup
	mu          sync.RWMutex
	entryPool   *sync.Pool
	closed      atomic.Bool
	sampler     Sampler
}

// Options configures a Logger.
type Options struct {
	// Output is the writer where logs are written.
	// Default is os.Stdout.
	Output io.Writer

	// Level is the minimum log level.
	// Default is InfoLevel.
	Level Level

	// Formatter sets the log output format.
	// Default is TextFormatter.
	Formatter Formatter

	// AddCaller enables caller information in logs.
	// Default is false.
	AddCaller bool

	// CallerDepth sets the caller stack depth.
	// Default is 2.
	CallerDepth int

	// AddStack enables stack traces for error and above.
	// Default is false.
	AddStack bool

	// AsyncBufferSize enables asynchronous logging with the specified buffer size.
	// If 0, synchronous logging is used.
	// If > 0, async logging is enabled with the specified buffer size.
	// Default is 0 (synchronous).
	AsyncBufferSize int

	// Hooks are additional hooks to add to the logger.
	Hooks []Hook

	// Fields are default fields to include in all log entries.
	Fields []Field

	// Sampler is used for rate limiting logs.
	Sampler Sampler
}

// applyDefaults applies default values to nil or zero-valued options.
func (o *Options) applyDefaults() {
	if o.Output == nil {
		o.Output = os.Stdout
	}
	if o.Formatter == nil {
		o.Formatter = &TextFormatter{}
	}
	if o.CallerDepth == 0 {
		o.CallerDepth = 2
	}
	// Level defaults to InfoLevel (0), but 0 is also a valid level (DebugLevel)
	// so we can't distinguish between "not set" and "explicitly set to DebugLevel"
	// We'll handle this in the New function
}

// New creates a new Logger with the provided options.
// If opts is nil, default options will be used.
func New(opts *Options) *Logger {
	if opts == nil {
		opts = &Options{}
	}

	// Apply defaults
	opts.applyDefaults()

	l := &Logger{
		output:      opts.Output,
		formatter:   opts.Formatter,
		callerDepth: opts.CallerDepth,
		addCaller:   opts.AddCaller,
		addStack:    opts.AddStack,
		hooks:       opts.Hooks,
		fields:      opts.Fields,
		sampler:     opts.Sampler,
		entryPool: &sync.Pool{
			New: func() any {
				return &Entry{
					Fields: make([]Field, 0, 8),
				}
			},
		},
	}

	// Set level (default to InfoLevel if not specified)
	if opts.Level == 0 {
		l.level.Store(int32(InfoLevel))
	} else {
		l.level.Store(int32(opts.Level))
	}

	// Enable async if buffer size is set
	if opts.AsyncBufferSize > 0 {
		l.async = true
		l.asyncCh = make(chan *Entry, opts.AsyncBufferSize)
		l.asyncWg.Add(1)
		go l.asyncWorker()
	}

	return l
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level Level) {
	l.level.Store(int32(level))
}

// GetLevel returns the current log level.
func (l *Logger) GetLevel() Level {
	return Level(l.level.Load())
}

// SetOutput sets the output writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	l.output = w
	l.mu.Unlock()
}

// SetFormatter sets the formatter.
func (l *Logger) SetFormatter(f Formatter) {
	l.mu.Lock()
	l.formatter = f
	l.mu.Unlock()
}

// AddHook adds a hook to the logger.
func (l *Logger) AddHook(hook Hook) {
	l.mu.Lock()
	l.hooks = append(l.hooks, hook)
	l.mu.Unlock()
}

// With creates a child logger with additional fields.
func (l *Logger) With(fields ...Field) *Logger {
	child := &Logger{
		output:      l.output,
		formatter:   l.formatter,
		hooks:       l.hooks,
		callerDepth: l.callerDepth,
		addCaller:   l.addCaller,
		addStack:    l.addStack,
		async:       l.async,
		asyncCh:     l.asyncCh,
		entryPool:   l.entryPool,
		sampler:     l.sampler,
		fields:      make([]Field, 0, len(l.fields)+len(fields)),
	}
	child.level.Store(l.level.Load())
	child.fields = append(child.fields, l.fields...)
	child.fields = append(child.fields, fields...)
	return child
}

// Close closes the logger and flushes any pending async logs.
func (l *Logger) Close() error {
	if l.closed.CompareAndSwap(false, true) {
		if l.async && l.asyncCh != nil {
			close(l.asyncCh)
			l.asyncWg.Wait()
		}
	}
	return nil
}

// asyncWorker processes async log entries.
func (l *Logger) asyncWorker() {
	defer l.asyncWg.Done()
	for entry := range l.asyncCh {
		l.writeEntry(entry)
		l.releaseEntry(entry)
	}
}

// getEntry gets an entry from the pool.
func (l *Logger) getEntry() *Entry {
	e := l.entryPool.Get().(*Entry)
	e.Time = time.Now()
	e.Fields = e.Fields[:0]
	e.Caller = ""
	e.Stack = ""
	return e
}

// releaseEntry returns an entry to the pool.
func (l *Logger) releaseEntry(e *Entry) {
	e.Message = ""
	e.Caller = ""
	e.Stack = ""
	e.Fields = e.Fields[:0]
	l.entryPool.Put(e)
}

// log logs a message at the given level.
func (l *Logger) log(level Level, msg string, fields []Field) {
	if Level(l.level.Load()) < level {
		return
	}

	// Check sampler
	if l.sampler != nil && !l.sampler.Sample(level, msg) {
		return
	}

	e := l.getEntry()
	e.Level = level
	e.Message = msg

	// Add default fields
	e.Fields = append(e.Fields, l.fields...)
	// Add call-site fields
	e.Fields = append(e.Fields, fields...)

	// Add caller info
	if l.addCaller {
		e.Caller = getCaller(l.callerDepth + 1)
	}

	// Add stack trace for errors
	if l.addStack && level <= ErrorLevel {
		e.Stack = getStack()
	}

	// Run hooks
	l.mu.RLock()
	for _, hook := range l.hooks {
		levels := hook.Levels()
		if len(levels) == 0 {
			// Fire for all levels
			hook.Fire(e)
		} else {
			// Check if level matches
			for _, lvl := range levels {
				if lvl == level {
					hook.Fire(e)
					break
				}
			}
		}
	}
	l.mu.RUnlock()

	if l.async && l.asyncCh != nil && !l.closed.Load() {
		// Clone entry for async processing
		clone := l.getEntry()
		*clone = *e
		clone.Fields = make([]Field, len(e.Fields))
		copy(clone.Fields, e.Fields)

		select {
		case l.asyncCh <- clone:
		default:
			// Channel full, write synchronously
			l.writeEntry(e)
		}
		l.releaseEntry(e)
	} else {
		l.writeEntry(e)
		l.releaseEntry(e)
	}
}

// logContext logs with context.
func (l *Logger) logContext(ctx context.Context, level Level, msg string, fields []Field) {
	// Check if context has logger fields
	if ctxFields := FieldsFromContext(ctx); len(ctxFields) > 0 {
		allFields := make([]Field, 0, len(ctxFields)+len(fields))
		allFields = append(allFields, ctxFields...)
		allFields = append(allFields, fields...)
		l.log(level, msg, allFields)
		return
	}
	l.log(level, msg, fields)
}

// writeEntry formats and writes the entry.
func (l *Logger) writeEntry(e *Entry) {
	l.mu.RLock()
	output := l.output
	formatter := l.formatter
	l.mu.RUnlock()

	data, err := formatter.Format(e)
	if err != nil {
		return
	}
	output.Write(data)
}

// Trace logs at trace level.
func (l *Logger) Trace(msg string, fields ...Field) {
	l.log(TraceLevel, msg, fields)
}

// Debug logs at debug level.
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields)
}

// Info logs at info level.
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields)
}

// Warn logs at warn level.
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields)
}

// Error logs at error level.
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields)
}

// Fatal logs at fatal level and exits.
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields)
	if l.async {
		l.Close()
	}
	os.Exit(1)
}

// Panic logs at panic level and panics.
func (l *Logger) Panic(msg string, fields ...Field) {
	l.log(PanicLevel, msg, fields)
	panic(msg)
}

// TraceContext logs at trace level with context.
func (l *Logger) TraceContext(ctx context.Context, msg string, fields ...Field) {
	l.logContext(ctx, TraceLevel, msg, fields)
}

// DebugContext logs at debug level with context.
func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	l.logContext(ctx, DebugLevel, msg, fields)
}

// InfoContext logs at info level with context.
func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	l.logContext(ctx, InfoLevel, msg, fields)
}

// WarnContext logs at warn level with context.
func (l *Logger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	l.logContext(ctx, WarnLevel, msg, fields)
}

// ErrorContext logs at error level with context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	l.logContext(ctx, ErrorLevel, msg, fields)
}

// Log logs at a specific level.
func (l *Logger) Log(level Level, msg string, fields ...Field) {
	l.log(level, msg, fields)
}

// LogContext logs at a specific level with context.
func (l *Logger) LogContext(ctx context.Context, level Level, msg string, fields ...Field) {
	l.logContext(ctx, level, msg, fields)
}

// IsEnabled returns true if the given level is enabled.
func (l *Logger) IsEnabled(level Level) bool {
	return Level(l.level.Load()) >= level
}

// getCaller returns the caller's file and line.
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}

	// Get just the filename
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	buf := make([]byte, 0, len(short)+12)
	buf = append(buf, short...)
	buf = append(buf, ':')
	buf = appendInt(buf, line)
	return string(buf)
}

// getStack returns a stack trace.
func getStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// appendInt appends an int to a byte slice.
func appendInt(buf []byte, n int) []byte {
	if n < 0 {
		buf = append(buf, '-')
		n = -n
	}
	if n < 10 {
		return append(buf, byte('0'+n))
	}
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}
	return append(buf, digits[i:]...)
}

// Default logger
var defaultLogger = New(nil)

// SetDefault sets the default logger.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Default returns the default logger.
func Default() *Logger {
	return defaultLogger
}

// SetDefaultFormatter sets the default formatter.
func SetDefaultFormatter(f Formatter) {
	defaultLogger.SetFormatter(f)
}

// SetDefaultLevel sets the default level.
func SetDefaultLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// Package-level functions that use the default logger

// Trace logs at trace level using the default logger.
func Trace(msg string, fields ...Field) { defaultLogger.Trace(msg, fields...) }

// Debug logs at debug level using the default logger.
func Debug(msg string, fields ...Field) { defaultLogger.Debug(msg, fields...) }

// Info logs at info level using the default logger.
func Info(msg string, fields ...Field) { defaultLogger.Info(msg, fields...) }

// Warn logs at warn level using the default logger.
func Warn(msg string, fields ...Field) { defaultLogger.Warn(msg, fields...) }

// Error logs at error level using the default logger.
func Error(msg string, fields ...Field) { defaultLogger.Error(msg, fields...) }

// Fatal logs at fatal level using the default logger and exits.
func Fatal(msg string, fields ...Field) { defaultLogger.Fatal(msg, fields...) }

// Panic logs at panic level using the default logger and panics.
func Panic(msg string, fields ...Field) { defaultLogger.Panic(msg, fields...) }

// With creates a child of the default logger with additional fields.
func With(fields ...Field) *Logger { return defaultLogger.With(fields...) }
