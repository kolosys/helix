package logs

import (
	"io"
	"os"
	"sync"
)

// Hook is called when a log entry is written.
type Hook interface {
	// Fire is called for each log entry.
	Fire(entry *Entry)

	// Levels returns the levels this hook should fire on.
	// If nil or empty, fires on all levels.
	Levels() []Level
}

// LevelHook fires only for specific levels.
type LevelHook struct {
	hook   Hook
	levels map[Level]bool
}

// NewLevelHook creates a hook that only fires for specific levels.
func NewLevelHook(hook Hook, levels ...Level) *LevelHook {
	h := &LevelHook{
		hook:   hook,
		levels: make(map[Level]bool),
	}
	for _, l := range levels {
		h.levels[l] = true
	}
	return h
}

// Fire implements Hook.
func (h *LevelHook) Fire(entry *Entry) {
	if h.levels[entry.Level] {
		h.hook.Fire(entry)
	}
}

// Levels implements Hook.
func (h *LevelHook) Levels() []Level {
	levels := make([]Level, 0, len(h.levels))
	for l := range h.levels {
		levels = append(levels, l)
	}
	return levels
}

// WriterHook writes entries to an io.Writer.
type WriterHook struct {
	writer    io.Writer
	formatter Formatter
	levels    []Level
	mu        sync.Mutex
}

// NewWriterHook creates a hook that writes to an io.Writer.
func NewWriterHook(w io.Writer, formatter Formatter, levels ...Level) *WriterHook {
	return &WriterHook{
		writer:    w,
		formatter: formatter,
		levels:    levels,
	}
}

// Fire implements Hook.
func (h *WriterHook) Fire(entry *Entry) {
	data, err := h.formatter.Format(entry)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.writer.Write(data)
	h.mu.Unlock()
}

// Levels implements Hook.
func (h *WriterHook) Levels() []Level {
	return h.levels
}

// FileHook writes entries to a file.
type FileHook struct {
	*WriterHook
	file *os.File
}

// NewFileHook creates a hook that writes to a file.
func NewFileHook(path string, formatter Formatter, levels ...Level) (*FileHook, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &FileHook{
		WriterHook: NewWriterHook(file, formatter, levels...),
		file:       file,
	}, nil
}

// Close closes the file.
func (h *FileHook) Close() error {
	return h.file.Close()
}

// ErrorHook collects errors for inspection.
type ErrorHook struct {
	errors []Entry
	mu     sync.Mutex
	max    int
}

// NewErrorHook creates a hook that collects error entries.
func NewErrorHook(maxEntries int) *ErrorHook {
	return &ErrorHook{
		errors: make([]Entry, 0, maxEntries),
		max:    maxEntries,
	}
}

// Fire implements Hook.
func (h *ErrorHook) Fire(entry *Entry) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.errors) >= h.max {
		// Remove oldest entry
		copy(h.errors, h.errors[1:])
		h.errors = h.errors[:len(h.errors)-1]
	}

	h.errors = append(h.errors, *entry)
}

// Levels implements Hook.
func (h *ErrorHook) Levels() []Level {
	return []Level{ErrorLevel, FatalLevel, PanicLevel}
}

// Errors returns collected errors.
func (h *ErrorHook) Errors() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()

	result := make([]Entry, len(h.errors))
	copy(result, h.errors)
	return result
}

// Clear clears collected errors.
func (h *ErrorHook) Clear() {
	h.mu.Lock()
	h.errors = h.errors[:0]
	h.mu.Unlock()
}

// MetricsHook tracks log counts by level.
type MetricsHook struct {
	counts map[Level]uint64
	mu     sync.RWMutex
}

// NewMetricsHook creates a hook that tracks log counts.
func NewMetricsHook() *MetricsHook {
	return &MetricsHook{
		counts: make(map[Level]uint64),
	}
}

// Fire implements Hook.
func (h *MetricsHook) Fire(entry *Entry) {
	h.mu.Lock()
	h.counts[entry.Level]++
	h.mu.Unlock()
}

// Levels implements Hook.
func (h *MetricsHook) Levels() []Level {
	return nil // All levels
}

// Count returns the count for a level.
func (h *MetricsHook) Count(level Level) uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.counts[level]
}

// Counts returns all counts.
func (h *MetricsHook) Counts() map[Level]uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[Level]uint64, len(h.counts))
	for k, v := range h.counts {
		result[k] = v
	}
	return result
}

// Reset resets all counts.
func (h *MetricsHook) Reset() {
	h.mu.Lock()
	for k := range h.counts {
		h.counts[k] = 0
	}
	h.mu.Unlock()
}

// FilterHook conditionally fires another hook.
type FilterHook struct {
	hook   Hook
	filter func(*Entry) bool
}

// NewFilterHook creates a hook that conditionally fires.
func NewFilterHook(hook Hook, filter func(*Entry) bool) *FilterHook {
	return &FilterHook{
		hook:   hook,
		filter: filter,
	}
}

// Fire implements Hook.
func (h *FilterHook) Fire(entry *Entry) {
	if h.filter(entry) {
		h.hook.Fire(entry)
	}
}

// Levels implements Hook.
func (h *FilterHook) Levels() []Level {
	return h.hook.Levels()
}

// FuncHook wraps a function as a hook.
type FuncHook struct {
	fn     func(*Entry)
	levels []Level
}

// NewFuncHook creates a hook from a function.
func NewFuncHook(fn func(*Entry), levels ...Level) *FuncHook {
	return &FuncHook{
		fn:     fn,
		levels: levels,
	}
}

// Fire implements Hook.
func (h *FuncHook) Fire(entry *Entry) {
	h.fn(entry)
}

// Levels implements Hook.
func (h *FuncHook) Levels() []Level {
	return h.levels
}

