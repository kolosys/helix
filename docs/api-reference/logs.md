# logs API

Complete API documentation for the logs package.

**Import Path:** `github.com/kolosys/helix/logs`

## Package Documentation

Package logs provides a high-performance, context-aware structured logging library.

Features:
  - Zero-allocation hot paths using sync.Pool
  - Context-aware logging with context.Context
  - Type-safe field builders
  - Multiple output formats (text, JSON, pretty)
  - Sampling for high-volume logs
  - Async logging option
  - Hook system for extensibility
  - Built-in caller information
  - Chained/fluent API

Basic usage:

	log := logs.New()
	log.Info("server started", logs.Int("port", 8080))

With context:

	log.InfoContext(ctx, "request processed", logs.Duration("latency", time.Since(start)))


## Constants

**RequestIDKey**

RequestID is a common field key for request IDs.


```go
const RequestIDKey = "request_id"
```

**TraceIDKey**

TraceID is a common field key for trace IDs.


```go
const TraceIDKey = "trace_id"
```

**UserIDKey**

UserID is a common field key for user IDs.


```go
const UserIDKey = "user_id"
```

## Types

### AlwaysSampler
AlwaysSampler always allows logging.

#### Example Usage

```go
// Create a new AlwaysSampler
alwayssampler := AlwaysSampler{

}
```

#### Type Definition

```go
type AlwaysSampler struct {
}
```

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### Builder
Builder provides a fluent/chainable API for building log entries. It accumulates fields and then emits a log entry when a level method is called.

#### Example Usage

```go
// Create a new Builder
builder := Builder{

}
```

#### Type Definition

```go
type Builder struct {
}
```

## Methods

### Bool

Bool adds a bool field.

```go
func (*Builder) Bool(key string, value bool) *Builder
```

**Parameters:**
- `key` (string)
- `value` (bool)

**Returns:**
- *Builder

### Debug

Debug logs at debug level.

```go
func (*Builder) Debug(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Err

Err adds an error field with key "error".

```go
func (*Builder) Err(err error) *Builder
```

**Parameters:**
- `err` (error)

**Returns:**
- *Builder

### Error

Error logs at error level.

```go
func (*Builder) Error(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Fatal

Fatal logs at fatal level and exits.

```go
func (*ErrorBuilder) Fatal(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Float64

Float64 adds a float64 field.

```go
func (*Builder) Float64(key string, value float64) *Builder
```

**Parameters:**
- `key` (string)
- `value` (float64)

**Returns:**
- *Builder

### Info

Info logs at info level.

```go
func (*ErrorBuilder) Info(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Int

Int adds an int field.

```go
func (*Builder) Int(key string, value int) *Builder
```

**Parameters:**
- `key` (string)
- `value` (int)

**Returns:**
- *Builder

### Int64

Int64 adds an int64 field.

```go
func (*Builder) Int64(key string, value int64) *Builder
```

**Parameters:**
- `key` (string)
- `value` (int64)

**Returns:**
- *Builder

### Log

Log logs at the specified level.

```go
func (*Logger) Log(level Level, msg string, fields ...Field)
```

**Parameters:**
- `level` (Level)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Msg

Msg is an alias for Info (zerolog-style).

```go
func (*Builder) Msg(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Panic

Panic logs at panic level and panics.

```go
func Panic(msg string, fields ...Field)
```

**Parameters:**
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Send

Send logs with an empty message (zerolog-style).

```go
func (*Builder) Send()
```

**Parameters:**
  None

**Returns:**
  None

### Str

Str adds a string field.

```go
func (*Builder) Str(key, value string) *Builder
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
- *Builder

### Trace

Trace logs at trace level.

```go
func (*Builder) Trace(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Uint

Uint adds a uint field.

```go
func (*Builder) Uint(key string, value uint) *Builder
```

**Parameters:**
- `key` (string)
- `value` (uint)

**Returns:**
- *Builder

### Uint64

Uint64 adds a uint64 field.

```go
func (*Builder) Uint64(key string, value uint64) *Builder
```

**Parameters:**
- `key` (string)
- `value` (uint64)

**Returns:**
- *Builder

### Warn

Warn logs at warn level.

```go
func Warn(msg string, fields ...Field)
```

**Parameters:**
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### With

With adds a field to the builder using auto-detection.

```go
func (*Builder) With(key string, value any) *Builder
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- *Builder

### WithContext

WithContext sets the context for the log entry.

```go
func (*Builder) WithContext(ctx context.Context) *Builder
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- *Builder

### WithError

WithError adds an error field.

```go
func (*Builder) WithError(err error) *Builder
```

**Parameters:**
- `err` (error)

**Returns:**
- *Builder

### WithField

WithField adds a typed field to the builder.

```go
func (*Builder) WithField(f Field) *Builder
```

**Parameters:**
- `f` (Field)

**Returns:**
- *Builder

### WithFields

WithFields adds multiple typed fields to the builder.

```go
func WithFields(fields ...Field) Option
```

**Parameters:**
- `fields` (...Field)

**Returns:**
- Option

### CompositeSampler
CompositeSampler combines multiple samplers with AND logic.

#### Example Usage

```go
// Create a new CompositeSampler
compositesampler := CompositeSampler{

}
```

#### Type Definition

```go
type CompositeSampler struct {
}
```

### Constructor Functions

### NewCompositeSampler

NewCompositeSampler creates a sampler that requires all samplers to pass.

```go
func NewCompositeSampler(samplers ...Sampler) *CompositeSampler
```

**Parameters:**
- `samplers` (...Sampler)

**Returns:**
- *CompositeSampler

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### CountSampler
CountSampler logs every Nth occurrence.

#### Example Usage

```go
// Create a new CountSampler
countsampler := CountSampler{

}
```

#### Type Definition

```go
type CountSampler struct {
}
```

### Constructor Functions

### NewCountSampler

NewCountSampler creates a sampler that logs every Nth occurrence.

```go
func NewCountSampler(n int) *CountSampler
```

**Parameters:**
- `n` (int)

**Returns:**
- *CountSampler

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### Entry
Entry represents a log entry.

#### Example Usage

```go
// Create a new Entry
entry := Entry{
    Level: Level{},
    Time: /* value */,
    Message: "example",
    Fields: [],
    Caller: "example",
    Stack: "example",
}
```

#### Type Definition

```go
type Entry struct {
    Level Level
    Time time.Time
    Message string
    Fields []Field
    Caller string
    Stack string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Level | `Level` |  |
| Time | `time.Time` |  |
| Message | `string` |  |
| Fields | `[]Field` |  |
| Caller | `string` |  |
| Stack | `string` |  |

## Methods

### GetField

GetField returns the field with the given key, or an empty field if not found.

```go
func (*Entry) GetField(key string) (Field, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- Field
- bool

### GetString

GetString returns the string value of a field, or empty string if not found.

```go
func (*Entry) GetString(key string) string
```

**Parameters:**
- `key` (string)

**Returns:**
- string

### HasField

HasField returns true if the entry has a field with the given key.

```go
func (*Entry) HasField(key string) bool
```

**Parameters:**
- `key` (string)

**Returns:**
- bool

### ErrorBuilder
ErrorBuilder provides a fluent API for logging errors.

#### Example Usage

```go
// Create a new ErrorBuilder
errorbuilder := ErrorBuilder{

}
```

#### Type Definition

```go
type ErrorBuilder struct {
}
```

## Methods

### Debug

Debug logs at debug level if error is not nil.

```go
func (*ErrorBuilder) Debug(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Error

Error logs at error level if error is not nil.

```go
func (*ErrorBuilder) Error(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Fatal

Fatal logs at fatal level if error is not nil and exits.

```go
func (*Builder) Fatal(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Info

Info logs at info level if error is not nil.

```go
func (*Builder) Info(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Trace

Trace logs at trace level if error is not nil.

```go
func (*Builder) Trace(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Warn

Warn logs at warn level if error is not nil.

```go
func (*Builder) Warn(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### With

With adds a field to the error builder.

```go
func (*Builder) With(key string, value any) *Builder
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- *Builder

### WithField

WithField adds a typed field.

```go
func (*Builder) WithField(f Field) *Builder
```

**Parameters:**
- `f` (Field)

**Returns:**
- *Builder

### WithFields

WithFields adds multiple fields.

```go
func (*Builder) WithFields(fields ...Field) *Builder
```

**Parameters:**
- `fields` (...Field)

**Returns:**
- *Builder

### ErrorHook
ErrorHook collects errors for inspection.

#### Example Usage

```go
// Create a new ErrorHook
errorhook := ErrorHook{

}
```

#### Type Definition

```go
type ErrorHook struct {
}
```

### Constructor Functions

### NewErrorHook

NewErrorHook creates a hook that collects error entries.

```go
func NewErrorHook(maxEntries int) *ErrorHook
```

**Parameters:**
- `maxEntries` (int)

**Returns:**
- *ErrorHook

## Methods

### Clear

Clear clears collected errors.

```go
func (*ErrorHook) Clear()
```

**Parameters:**
  None

**Returns:**
  None

### Errors

Errors returns collected errors.

```go
func (*ErrorHook) Errors() []Entry
```

**Parameters:**
  None

**Returns:**
- []Entry

### Fire

Fire implements Hook.

```go
func (*FuncHook) Fire(entry *Entry)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
  None

### Levels

Levels implements Hook.

```go
func (*FuncHook) Levels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

### Field
Field represents a structured log field.

#### Example Usage

```go
// Create a new Field
field := Field{
    Key: "example",
    Type: FieldType{},
    Int: 42,
    Uint: 42,
    Float: 3.14,
    String: "example",
    Interface: any{},
}
```

#### Type Definition

```go
type Field struct {
    Key string
    Type FieldType
    Int int64
    Uint uint64
    Float float64
    String string
    Interface any
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Key | `string` |  |
| Type | `FieldType` |  |
| Int | `int64` |  |
| Uint | `uint64` |  |
| Float | `float64` |  |
| String | `string` |  |
| Interface | `any` |  |

### Constructor Functions

### Any

Any creates a field with any value.

```go
func Any(key string, value any) Field
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- Field

### Bool

Bool creates a bool field.

```go
func (*Builder) Bool(key string, value bool) *Builder
```

**Parameters:**
- `key` (string)
- `value` (bool)

**Returns:**
- *Builder

### Bytes

Bytes creates a []byte field.

```go
func Bytes(key string, value []byte) Field
```

**Parameters:**
- `key` (string)
- `value` ([]byte)

**Returns:**
- Field

### Duration

Duration creates a time.Duration field.

```go
func Duration(key string, value time.Duration) Field
```

**Parameters:**
- `key` (string)
- `value` (time.Duration)

**Returns:**
- Field

### Err

Err creates an error field with key "error".

```go
func (*Builder) Err(err error) *Builder
```

**Parameters:**
- `err` (error)

**Returns:**
- *Builder

### ErrChain

ErrChain creates a field that unwraps the error chain.

```go
func ErrChain(err error) Field
```

**Parameters:**
- `err` (error)

**Returns:**
- Field

### ErrWithStack

ErrWithStack creates an error field with stack trace.

```go
func ErrWithStack(err error) Field
```

**Parameters:**
- `err` (error)

**Returns:**
- Field

### FieldsFromContext

FieldsFromContext extracts fields from the context.

```go
func FieldsFromContext(ctx context.Context) []Field
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- []Field

### Float32

Float32 creates a float32 field.

```go
func Float32(key string, value float32) Field
```

**Parameters:**
- `key` (string)
- `value` (float32)

**Returns:**
- Field

### Float64

Float64 creates a float64 field.

```go
func (*Builder) Float64(key string, value float64) *Builder
```

**Parameters:**
- `key` (string)
- `value` (float64)

**Returns:**
- *Builder

### Int

Int creates an int field.

```go
func (*Builder) Int(key string, value int) *Builder
```

**Parameters:**
- `key` (string)
- `value` (int)

**Returns:**
- *Builder

### Int16

Int16 creates an int16 field.

```go
func Int16(key string, value int16) Field
```

**Parameters:**
- `key` (string)
- `value` (int16)

**Returns:**
- Field

### Int32

Int32 creates an int32 field.

```go
func Int32(key string, value int32) Field
```

**Parameters:**
- `key` (string)
- `value` (int32)

**Returns:**
- Field

### Int64

Int64 creates an int64 field.

```go
func (*Builder) Int64(key string, value int64) *Builder
```

**Parameters:**
- `key` (string)
- `value` (int64)

**Returns:**
- *Builder

### Int8

Int8 creates an int8 field.

```go
func Int8(key string, value int8) Field
```

**Parameters:**
- `key` (string)
- `value` (int8)

**Returns:**
- Field

### JSON

JSON creates a field that will be JSON-encoded.

```go
func JSON(key string, value any) Field
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- Field

### NamedErr

NamedErr creates an error field with a custom key.

```go
func NamedErr(key string, err error) Field
```

**Parameters:**
- `key` (string)
- `err` (error)

**Returns:**
- Field

### Namespace

Namespace creates a namespace field for grouping.

```go
func Namespace(key string) Field
```

**Parameters:**
- `key` (string)

**Returns:**
- Field

### Object

Object creates a field from any value, attempting to extract structure. For structs, it extracts fields. For other types, it uses Any.

```go
func Object(key string, v any) Field
```

**Parameters:**
- `key` (string)
- `v` (any)

**Returns:**
- Field

### Stack

Stack creates a field containing a stack trace.

```go
func Stack(key string) Field
```

**Parameters:**
- `key` (string)

**Returns:**
- Field

### String

String creates a string field.

```go
func (Level) String() string
```

**Parameters:**
  None

**Returns:**
- string

### Stringer

Stringer creates a field from a fmt.Stringer.

```go
func Stringer(key string, value fmt.Stringer) Field
```

**Parameters:**
- `key` (string)
- `value` (fmt.Stringer)

**Returns:**
- Field

### Strings

Strings creates a string slice field.

```go
func Strings(key string, values []string) Field
```

**Parameters:**
- `key` (string)
- `values` ([]string)

**Returns:**
- Field

### Struct

Struct creates fields from a struct's exported fields. It uses json tags for field names if available, otherwise uses the field name. Nested structs are flattened with dot notation.

```go
func Struct(key string, v any) Field
```

**Parameters:**
- `key` (string)
- `v` (any)

**Returns:**
- Field

### StructFlat

StructFlat creates multiple top-level fields from a struct's exported fields. Unlike Struct, this doesn't nest under a key - fields are added directly.

```go
func StructFlat(v any) []Field
```

**Parameters:**
- `v` (any)

**Returns:**
- []Field

### Time

Time creates a time.Time field.

```go
func Time(key string, value time.Time) Field
```

**Parameters:**
- `key` (string)
- `value` (time.Time)

**Returns:**
- Field

### Uint

Uint creates a uint field.

```go
func (*Builder) Uint(key string, value uint) *Builder
```

**Parameters:**
- `key` (string)
- `value` (uint)

**Returns:**
- *Builder

### Uint16

Uint16 creates a uint16 field.

```go
func Uint16(key string, value uint16) Field
```

**Parameters:**
- `key` (string)
- `value` (uint16)

**Returns:**
- Field

### Uint32

Uint32 creates a uint32 field.

```go
func Uint32(key string, value uint32) Field
```

**Parameters:**
- `key` (string)
- `value` (uint32)

**Returns:**
- Field

### Uint64

Uint64 creates a uint64 field.

```go
func (*Builder) Uint64(key string, value uint64) *Builder
```

**Parameters:**
- `key` (string)
- `value` (uint64)

**Returns:**
- *Builder

### Uint8

Uint8 creates a uint8 field.

```go
func Uint8(key string, value uint8) Field
```

**Parameters:**
- `key` (string)
- `value` (uint8)

**Returns:**
- Field

### V

V creates a field with automatic type detection. This is a convenience function for when you don't want to specify the type.

```go
func V(key string, value any) Field
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- Field

## Methods

### StringValue

StringValue returns the field value as a string.

```go
func (Field) StringValue() string
```

**Parameters:**
  None

**Returns:**
- string

### Value

Value returns the field value as an interface{}.

```go
func (Field) Value() any
```

**Parameters:**
  None

**Returns:**
- any

### FieldType
FieldType represents the type of a field value.

#### Example Usage

```go
// Example usage of FieldType
var value FieldType
// Initialize with appropriate value
```

#### Type Definition

```go
type FieldType uint8
```

### FileHook
FileHook writes entries to a file.

#### Example Usage

```go
// Create a new FileHook
filehook := FileHook{

}
```

#### Type Definition

```go
type FileHook struct {
    *WriterHook
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| **WriterHook | `*WriterHook` |  |

### Constructor Functions

### NewFileHook

NewFileHook creates a hook that writes to a file.

```go
func NewFileHook(path string, formatter Formatter, levels ...Level) (*FileHook, error)
```

**Parameters:**
- `path` (string)
- `formatter` (Formatter)
- `levels` (...Level)

**Returns:**
- *FileHook
- error

## Methods

### Close

Close closes the file.

```go
func (*FileHook) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### FilterHook
FilterHook conditionally fires another hook.

#### Example Usage

```go
// Create a new FilterHook
filterhook := FilterHook{

}
```

#### Type Definition

```go
type FilterHook struct {
}
```

### Constructor Functions

### NewFilterHook

NewFilterHook creates a hook that conditionally fires.

```go
func NewFilterHook(hook Hook, filter func(*Entry) bool) *FilterHook
```

**Parameters:**
- `hook` (Hook)
- `filter` (func(*Entry) bool)

**Returns:**
- *FilterHook

## Methods

### Fire

Fire implements Hook.

```go
func (*FuncHook) Fire(entry *Entry)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
  None

### Levels

Levels implements Hook.

```go
func (*FuncHook) Levels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

### FirstNSampler
FirstNSampler logs only the first N occurrences.

#### Example Usage

```go
// Create a new FirstNSampler
firstnsampler := FirstNSampler{

}
```

#### Type Definition

```go
type FirstNSampler struct {
}
```

### Constructor Functions

### NewFirstNSampler

NewFirstNSampler creates a sampler that logs only the first N occurrences.

```go
func NewFirstNSampler(n int) *FirstNSampler
```

**Parameters:**
- `n` (int)

**Returns:**
- *FirstNSampler

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### Formatter
Formatter formats log entries.

#### Example Usage

```go
// Example implementation of Formatter
type MyFormatter struct {
    // Add your fields here
}

func (m MyFormatter) Format(param1 *Entry) []byte {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Formatter interface {
    Format(entry *Entry) ([]byte, error)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### FuncHook
FuncHook wraps a function as a hook.

#### Example Usage

```go
// Create a new FuncHook
funchook := FuncHook{

}
```

#### Type Definition

```go
type FuncHook struct {
}
```

### Constructor Functions

### NewFuncHook

NewFuncHook creates a hook from a function.

```go
func NewFuncHook(fn func(*Entry), levels ...Level) *FuncHook
```

**Parameters:**
- `fn` (func(*Entry))
- `levels` (...Level)

**Returns:**
- *FuncHook

## Methods

### Fire

Fire implements Hook.

```go
func (*FuncHook) Fire(entry *Entry)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
  None

### Levels

Levels implements Hook.

```go
func (*FuncHook) Levels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

### Hook
Hook is called when a log entry is written.

#### Example Usage

```go
// Example implementation of Hook
type MyHook struct {
    // Add your fields here
}

func (m MyHook) Fire(param1 *Entry)  {
    // Implement your logic here
    return
}

func (m MyHook) Levels() []Level {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Hook interface {
    Fire(entry *Entry)
    Levels() []Level
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### JSONFormatter
JSONFormatter formats logs as JSON.

#### Example Usage

```go
// Create a new JSONFormatter
jsonformatter := JSONFormatter{
    TimestampFormat: "example",
    DisableTimestamp: true,
    TimestampKey: "example",
    LevelKey: "example",
    MessageKey: "example",
    CallerKey: "example",
    StackKey: "example",
    PrettyPrint: true,
    EscapeHTML: true,
}
```

#### Type Definition

```go
type JSONFormatter struct {
    TimestampFormat string
    DisableTimestamp bool
    TimestampKey string
    LevelKey string
    MessageKey string
    CallerKey string
    StackKey string
    PrettyPrint bool
    EscapeHTML bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TimestampFormat | `string` | TimestampFormat is the format for timestamps. Default: time.RFC3339Nano |
| DisableTimestamp | `bool` | DisableTimestamp disables timestamp output. |
| TimestampKey | `string` | TimestampKey is the key for the timestamp field. Default: "time" |
| LevelKey | `string` | LevelKey is the key for the level field. Default: "level" |
| MessageKey | `string` | MessageKey is the key for the message field. Default: "msg" |
| CallerKey | `string` | CallerKey is the key for the caller field. Default: "caller" |
| StackKey | `string` | StackKey is the key for the stack trace field. Default: "stack" |
| PrettyPrint | `bool` | PrettyPrint enables pretty-printed JSON. |
| EscapeHTML | `bool` | EscapeHTML escapes HTML in JSON strings. |

## Methods

### Format

Format formats an entry as JSON.

```go
func (*NamedFormatter) Format(entry *Entry) ([]byte, error)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
- []byte
- error

### Level
Level represents a log level.

#### Example Usage

```go
// Example usage of Level
var value Level
// Initialize with appropriate value
```

#### Type Definition

```go
type Level int
```

### Constructor Functions

### AllLevels

AllLevels returns all log levels.

```go
func AllLevels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

### ParseLevel

ParseLevel parses a string into a Level.

```go
func ParseLevel(s string) Level
```

**Parameters:**
- `s` (string)

**Returns:**
- Level

## Methods

### Color

Color returns the ANSI color code for the level.

```go
func (Level) Color() string
```

**Parameters:**
  None

**Returns:**
- string

### ShortString

ShortString returns a 4-character representation for alignment.

```go
func (Level) ShortString() string
```

**Parameters:**
  None

**Returns:**
- string

### String

String returns the string representation of a level.

```go
func String(key, value string) Field
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
- Field

### LevelHook
LevelHook fires only for specific levels.

#### Example Usage

```go
// Create a new LevelHook
levelhook := LevelHook{

}
```

#### Type Definition

```go
type LevelHook struct {
}
```

### Constructor Functions

### NewLevelHook

NewLevelHook creates a hook that only fires for specific levels.

```go
func NewLevelHook(hook Hook, levels ...Level) *LevelHook
```

**Parameters:**
- `hook` (Hook)
- `levels` (...Level)

**Returns:**
- *LevelHook

## Methods

### Fire

Fire implements Hook.

```go
func (*FuncHook) Fire(entry *Entry)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
  None

### Levels

Levels implements Hook.

```go
func (*FuncHook) Levels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

### LevelSampler
LevelSampler applies different samplers per level.

#### Example Usage

```go
// Create a new LevelSampler
levelsampler := LevelSampler{

}
```

#### Type Definition

```go
type LevelSampler struct {
}
```

### Constructor Functions

### NewLevelSampler

NewLevelSampler creates a sampler with per-level configuration.

```go
func NewLevelSampler(fallback Sampler) *LevelSampler
```

**Parameters:**
- `fallback` (Sampler)

**Returns:**
- *LevelSampler

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### WithLevel

WithLevel sets the sampler for a specific level.

```go
func (*LevelSampler) WithLevel(level Level, sampler Sampler) *LevelSampler
```

**Parameters:**
- `level` (Level)
- `sampler` (Sampler)

**Returns:**
- *LevelSampler

### Logger
Logger is the main logging interface.

#### Example Usage

```go
// Create a new Logger
logger := Logger{

}
```

#### Type Definition

```go
type Logger struct {
}
```

### Constructor Functions

### Component

Component creates a component logger from the default logger.

```go
func Component(name string) *Logger
```

**Parameters:**
- `name` (string)

**Returns:**
- *Logger

### Default

Default returns the default logger.

```go
func Default() *Logger
```

**Parameters:**
  None

**Returns:**
- *Logger

### LoggerFromContext

LoggerFromContext extracts a logger from the context. Returns the default logger if none is set.

```go
func LoggerFromContext(ctx context.Context) *Logger
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- *Logger

### Named

Named creates a named child of the default logger.

```go
func Named(name string) *Logger
```

**Parameters:**
- `name` (string)

**Returns:**
- *Logger

### New

New creates a new Logger with default settings.

```go
func New(opts ...Option) *Logger
```

**Parameters:**
- `opts` (...Option)

**Returns:**
- *Logger

### With

With creates a child of the default logger with additional fields.

```go
func (*Builder) With(key string, value any) *Builder
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- *Builder

## Methods

### AddHook

AddHook adds a hook to the logger.

```go
func (*Logger) AddHook(hook Hook)
```

**Parameters:**
- `hook` (Hook)

**Returns:**
  None

### Build

Build returns a new Builder for constructing a log entry.

```go
func (*Logger) Build() *Builder
```

**Parameters:**
  None

**Returns:**
- *Builder

### CheckErr

CheckErr logs an error and returns true if err is not nil. Useful in if statements: if log.CheckErr(err, "failed") { return }

```go
func (*Logger) CheckErr(err error, msg string, fields ...Field) bool
```

**Parameters:**
- `err` (error)
- `msg` (string)
- `fields` (...Field)

**Returns:**
- bool

### Close

Close closes the logger and flushes any pending async logs.

```go
func (*FileHook) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Component

Component creates a named logger for a specific component. This is an alias for Named with a more semantic name.

```go
func Component(name string) *Logger
```

**Parameters:**
- `name` (string)

**Returns:**
- *Logger

### Ctx

Ctx returns a builder with context set.

```go
func (*Logger) Ctx(ctx context.Context) *Builder
```

**Parameters:**
- `ctx` (context.Context)

**Returns:**
- *Builder

### Debug

Debug logs at debug level.

```go
func (*Builder) Debug(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### DebugContext

DebugContext logs at debug level with context.

```go
func (*Logger) DebugContext(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Debugf

Debugf logs a formatted message at debug level.

```go
func Debugf(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### DebugfContext

DebugfContext logs a formatted message at debug level with context.

```go
func (*Logger) DebugfContext(ctx context.Context, format string, args ...any)
```

**Parameters:**
- `ctx` (context.Context)
- `format` (string)
- `args` (...any)

**Returns:**
  None

### Error

Error logs at error level.

```go
func (*Builder) Error(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### ErrorContext

ErrorContext logs at error level with context.

```go
func (*Logger) ErrorContext(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Errorf

Errorf logs a formatted message at error level.

```go
func Errorf(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### ErrorfContext

ErrorfContext logs a formatted message at error level with context.

```go
func (*Logger) ErrorfContext(ctx context.Context, format string, args ...any)
```

**Parameters:**
- `ctx` (context.Context)
- `format` (string)
- `args` (...any)

**Returns:**
  None

### F

F returns a builder with fields added using key-value pairs. Keys must be strings, values are auto-detected. Example: log.F("user", "john", "age", 30).Info("created")

```go
func (*Logger) F(keyvals ...any) *Builder
```

**Parameters:**
- `keyvals` (...any)

**Returns:**
- *Builder

### Fatal

Fatal logs at fatal level and exits.

```go
func (*Builder) Fatal(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### Fatalf

Fatalf logs a formatted message at fatal level and exits.

```go
func Fatalf(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### FullName

FullName returns the complete hierarchical name of the logger.

```go
func (*Logger) FullName() string
```

**Parameters:**
  None

**Returns:**
- string

### GetLevel

GetLevel returns the current log level.

```go
func (*Logger) GetLevel() Level
```

**Parameters:**
  None

**Returns:**
- Level

### GetName

GetName returns the logger's name, or empty string if not named.

```go
func (*Logger) GetName() string
```

**Parameters:**
  None

**Returns:**
- string

### IfErr

IfErr returns an ErrorBuilder if err is not nil, otherwise returns a no-op builder. This allows for one-liner error logging: log.IfErr(err).With("user", id).Error("failed to create user")

```go
func (*Logger) IfErr(err error) *ErrorBuilder
```

**Parameters:**
- `err` (error)

**Returns:**
- *ErrorBuilder

### Info

Info logs at info level.

```go
func (*Builder) Info(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### InfoContext

InfoContext logs at info level with context.

```go
func (*Logger) InfoContext(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Infof

Infof logs a formatted message at info level.

```go
func Infof(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### InfofContext

InfofContext logs a formatted message at info level with context.

```go
func (*Logger) InfofContext(ctx context.Context, format string, args ...any)
```

**Parameters:**
- `ctx` (context.Context)
- `format` (string)
- `args` (...any)

**Returns:**
  None

### IsEnabled

IsEnabled returns true if the given level is enabled.

```go
func (*Logger) IsEnabled(level Level) bool
```

**Parameters:**
- `level` (Level)

**Returns:**
- bool

### Log

Log logs at a specific level.

```go
func (*Builder) Log(level Level, msg string)
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
  None

### LogContext

LogContext logs at a specific level with context.

```go
func (*Logger) LogContext(ctx context.Context, level Level, msg string, fields ...Field)
```

**Parameters:**
- `ctx` (context.Context)
- `level` (Level)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### LogErr

LogErr logs an error at error level if not nil. This is a simple one-liner for common error logging. log.LogErr(err, "operation failed")

```go
func (*Logger) LogErr(err error, msg string, fields ...Field)
```

**Parameters:**
- `err` (error)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Module

Module creates a named logger for a module. This is an alias for Named with a more semantic name.

```go
func (*Logger) Module(name string) *Logger
```

**Parameters:**
- `name` (string)

**Returns:**
- *Logger

### MustErr

MustErr panics if error is not nil, logging the error first.

```go
func (*Logger) MustErr(err error, msg string, fields ...Field)
```

**Parameters:**
- `err` (error)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### NameParts

NameParts returns the name split into parts.

```go
func (*Logger) NameParts() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Named

Named creates a child logger with a name prefix. Names are joined with dots to create a hierarchy. userLog := log.Named("users")           // [users] authLog := userLog.Named("auth")        // [users.auth] log.Info("action")                      // [users.auth] action

```go
func Named(name string) *Logger
```

**Parameters:**
- `name` (string)

**Returns:**
- *Logger

### Panic

Panic logs at panic level and panics.

```go
func Panic(msg string, fields ...Field)
```

**Parameters:**
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Panicf

Panicf logs a formatted message at panic level and panics.

```go
func Panicf(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### Print

Print logs a message at info level (stdlib log compatibility).

```go
func Print(args ...any)
```

**Parameters:**
- `args` (...any)

**Returns:**
  None

### Printf

Printf logs a formatted message at info level (stdlib log compatibility).

```go
func Printf(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### Println

Println logs a message at info level (stdlib log compatibility).

```go
func Println(args ...any)
```

**Parameters:**
- `args` (...any)

**Returns:**
  None

### Service

Service creates a named logger for a service. This is an alias for Named with a more semantic name.

```go
func (*Logger) Service(name string) *Logger
```

**Parameters:**
- `name` (string)

**Returns:**
- *Logger

### SetFormatter

SetFormatter sets the formatter.

```go
func (*Logger) SetFormatter(f Formatter)
```

**Parameters:**
- `f` (Formatter)

**Returns:**
  None

### SetLevel

SetLevel sets the minimum log level.

```go
func (*Logger) SetLevel(level Level)
```

**Parameters:**
- `level` (Level)

**Returns:**
  None

### SetOutput

SetOutput sets the output writer.

```go
func (*Logger) SetOutput(w io.Writer)
```

**Parameters:**
- `w` (io.Writer)

**Returns:**
  None

### Trace

Trace logs at trace level.

```go
func (*Builder) Trace(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### TraceContext

TraceContext logs at trace level with context.

```go
func (*Logger) TraceContext(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Tracef

Tracef logs a formatted message at trace level.

```go
func Tracef(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### TracefContext

TracefContext logs a formatted message at trace level with context.

```go
func (*Logger) TracefContext(ctx context.Context, format string, args ...any)
```

**Parameters:**
- `ctx` (context.Context)
- `format` (string)
- `args` (...any)

**Returns:**
  None

### Warn

Warn logs at warn level.

```go
func (*Builder) Warn(msg string)
```

**Parameters:**
- `msg` (string)

**Returns:**
  None

### WarnContext

WarnContext logs at warn level with context.

```go
func (*Logger) WarnContext(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
- `ctx` (context.Context)
- `msg` (string)
- `fields` (...Field)

**Returns:**
  None

### Warnf

Warnf logs a formatted message at warn level.

```go
func Warnf(format string, args ...any)
```

**Parameters:**
- `format` (string)
- `args` (...any)

**Returns:**
  None

### WarnfContext

WarnfContext logs a formatted message at warn level with context.

```go
func (*Logger) WarnfContext(ctx context.Context, format string, args ...any)
```

**Parameters:**
- `ctx` (context.Context)
- `format` (string)
- `args` (...any)

**Returns:**
  None

### With

With creates a child logger with additional fields.

```go
func (*Builder) With(key string, value any) *Builder
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
- *Builder

### WrapErr

WrapErr wraps an error with additional context and logs it. Returns the wrapped error for returning from functions. return log.WrapErr(err, "failed to connect", logs.String("host", host))

```go
func (*Logger) WrapErr(err error, msg string, fields ...Field) error
```

**Parameters:**
- `err` (error)
- `msg` (string)
- `fields` (...Field)

**Returns:**
- error

### WrapErrLevel

WrapErrLevel wraps an error and logs at a specific level.

```go
func (*Logger) WrapErrLevel(level Level, err error, msg string, fields ...Field) error
```

**Parameters:**
- `level` (Level)
- `err` (error)
- `msg` (string)
- `fields` (...Field)

**Returns:**
- error

### MetricsHook
MetricsHook tracks log counts by level.

#### Example Usage

```go
// Create a new MetricsHook
metricshook := MetricsHook{

}
```

#### Type Definition

```go
type MetricsHook struct {
}
```

### Constructor Functions

### NewMetricsHook

NewMetricsHook creates a hook that tracks log counts.

```go
func NewMetricsHook() *MetricsHook
```

**Parameters:**
  None

**Returns:**
- *MetricsHook

## Methods

### Count

Count returns the count for a level.

```go
func (*MetricsHook) Count(level Level) uint64
```

**Parameters:**
- `level` (Level)

**Returns:**
- uint64

### Counts

Counts returns all counts.

```go
func (*MetricsHook) Counts() map[Level]uint64
```

**Parameters:**
  None

**Returns:**
- map[Level]uint64

### Fire

Fire implements Hook.

```go
func (*FuncHook) Fire(entry *Entry)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
  None

### Levels

Levels implements Hook.

```go
func (*FuncHook) Levels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

### Reset

Reset resets all counts.

```go
func (*MetricsHook) Reset()
```

**Parameters:**
  None

**Returns:**
  None

### NamedFormatter
NamedFormatter is a formatter that includes the logger name in output.

#### Example Usage

```go
// Create a new NamedFormatter
namedformatter := NamedFormatter{
    Inner: Formatter{},
    Separator: "example",
    Brackets: "example",
}
```

#### Type Definition

```go
type NamedFormatter struct {
    Inner Formatter
    Separator string
    Brackets string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Inner | `Formatter` | Inner is the formatter to wrap. |
| Separator | `string` | Separator is placed between the name and message. Default: " " |
| Brackets | `string` | Brackets wraps the name. Default: "[]" |

## Methods

### Format

Format formats an entry, prefixing with the logger name if present.

```go
func (*NamedFormatter) Format(entry *Entry) ([]byte, error)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
- []byte
- error

### NeverSampler
NeverSampler never allows logging.

#### Example Usage

```go
// Create a new NeverSampler
neversampler := NeverSampler{

}
```

#### Type Definition

```go
type NeverSampler struct {
}
```

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### NoopFormatter
NoopFormatter discards all output.

#### Example Usage

```go
// Create a new NoopFormatter
noopformatter := NoopFormatter{

}
```

#### Type Definition

```go
type NoopFormatter struct {
}
```

## Methods

### Format

Format returns nil.

```go
func (*NamedFormatter) Format(entry *Entry) ([]byte, error)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
- []byte
- error

### OncePerSampler
OncePerSampler logs a message only once per duration.

#### Example Usage

```go
// Create a new OncePerSampler
oncepersampler := OncePerSampler{

}
```

#### Type Definition

```go
type OncePerSampler struct {
}
```

### Constructor Functions

### NewOncePerSampler

NewOncePerSampler creates a sampler that logs each message at most once per period.

```go
func NewOncePerSampler(period time.Duration) *OncePerSampler
```

**Parameters:**
- `period` (time.Duration)

**Returns:**
- *OncePerSampler

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### Option
Option configures a Logger.

#### Example Usage

```go
// Example usage of Option
var value Option
// Initialize with appropriate value
```

#### Type Definition

```go
type Option func(*Logger)
```

### Constructor Functions

### WithAsync

WithAsync enables asynchronous logging.

```go
func WithAsync(bufferSize int) Option
```

**Parameters:**
- `bufferSize` (int)

**Returns:**
- Option

### WithCaller

WithCaller enables caller information in logs.

```go
func WithCaller() Option
```

**Parameters:**
  None

**Returns:**
- Option

### WithCallerDepth

WithCallerDepth sets the caller stack depth.

```go
func WithCallerDepth(depth int) Option
```

**Parameters:**
- `depth` (int)

**Returns:**
- Option

### WithFields

WithFields adds default fields to all log entries.

```go
func (*Builder) WithFields(fields ...Field) *Builder
```

**Parameters:**
- `fields` (...Field)

**Returns:**
- *Builder

### WithFormatter

WithFormatter sets the log formatter.

```go
func WithFormatter(f Formatter) Option
```

**Parameters:**
- `f` (Formatter)

**Returns:**
- Option

### WithHooks

WithHooks adds hooks to the logger.

```go
func WithHooks(hooks ...Hook) Option
```

**Parameters:**
- `hooks` (...Hook)

**Returns:**
- Option

### WithLevel

WithLevel sets the minimum log level.

```go
func (*LevelSampler) WithLevel(level Level, sampler Sampler) *LevelSampler
```

**Parameters:**
- `level` (Level)
- `sampler` (Sampler)

**Returns:**
- *LevelSampler

### WithNamedFormat

WithNamedFormat wraps the current formatter to include logger names.

```go
func WithNamedFormat() Option
```

**Parameters:**
  None

**Returns:**
- Option

### WithOutput

WithOutput sets the output writer.

```go
func WithOutput(w io.Writer) Option
```

**Parameters:**
- `w` (io.Writer)

**Returns:**
- Option

### WithSampler

WithSampler sets a sampler for rate limiting logs.

```go
func WithSampler(s Sampler) Option
```

**Parameters:**
- `s` (Sampler)

**Returns:**
- Option

### WithStackTrace

WithStackTrace enables stack traces for error and above.

```go
func WithStackTrace() Option
```

**Parameters:**
  None

**Returns:**
- Option

### PrettyFormatter
PrettyFormatter formats logs with colors and alignment for development.

#### Example Usage

```go
// Create a new PrettyFormatter
prettyformatter := PrettyFormatter{
    TimestampFormat: "example",
    ShowCaller: true,
    ShowTimestamp: true,
}
```

#### Type Definition

```go
type PrettyFormatter struct {
    TimestampFormat string
    ShowCaller bool
    ShowTimestamp bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TimestampFormat | `string` | TimestampFormat is the format for timestamps. Default: "15:04:05.000" |
| ShowCaller | `bool` | ShowCaller shows caller information. |
| ShowTimestamp | `bool` | ShowTimestamp shows timestamps. |

## Methods

### Format

Format formats an entry in a pretty, colorful format.

```go
func (*NoopFormatter) Format(entry *Entry) ([]byte, error)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
- []byte
- error

### RandomSampler
RandomSampler samples a percentage of logs.

#### Example Usage

```go
// Create a new RandomSampler
randomsampler := RandomSampler{

}
```

#### Type Definition

```go
type RandomSampler struct {
}
```

### Constructor Functions

### NewRandomSampler

NewRandomSampler creates a sampler that logs a percentage of entries. percentage should be between 0 and 100.

```go
func NewRandomSampler(percentage int) *RandomSampler
```

**Parameters:**
- `percentage` (int)

**Returns:**
- *RandomSampler

## Methods

### Sample

Sample implements Sampler using a simple modulo for deterministic "random" sampling.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### RateSampler
RateSampler limits logs to a certain rate per message.

#### Example Usage

```go
// Create a new RateSampler
ratesampler := RateSampler{

}
```

#### Type Definition

```go
type RateSampler struct {
}
```

### Constructor Functions

### NewRateSampler

NewRateSampler creates a sampler that limits log rate per message. rate is the number of logs allowed per window. burst is the initial burst allowance.

```go
func NewRateSampler(rate int, window time.Duration) *RateSampler
```

**Parameters:**
- `rate` (int)
- `window` (time.Duration)

**Returns:**
- *RateSampler

## Methods

### Sample

Sample implements Sampler.

```go
func (*NeverSampler) Sample(level Level, msg string) bool
```

**Parameters:**
- `level` (Level)
- `msg` (string)

**Returns:**
- bool

### WithBurst

WithBurst sets the burst allowance.

```go
func (*RateSampler) WithBurst(burst int) *RateSampler
```

**Parameters:**
- `burst` (int)

**Returns:**
- *RateSampler

### Sampler
Sampler determines if a log entry should be emitted.

#### Example Usage

```go
// Example implementation of Sampler
type MySampler struct {
    // Add your fields here
}

func (m MySampler) Sample(param1 Level, param2 string) bool {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Sampler interface {
    Sample(level Level, msg string) bool
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### TextFormatter
TextFormatter formats logs as text.

#### Example Usage

```go
// Create a new TextFormatter
textformatter := TextFormatter{
    TimestampFormat: "example",
    DisableTimestamp: true,
    DisableColors: true,
    DisableQuoting: true,
    QuoteEmptyFields: true,
    FullTimestamp: true,
    PadLevelText: true,
    FieldSeparator: "example",
    KeyValueSeparator: "example",
}
```

#### Type Definition

```go
type TextFormatter struct {
    TimestampFormat string
    DisableTimestamp bool
    DisableColors bool
    DisableQuoting bool
    QuoteEmptyFields bool
    FullTimestamp bool
    PadLevelText bool
    FieldSeparator string
    KeyValueSeparator string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| TimestampFormat | `string` | TimestampFormat is the format for timestamps. Default: "2006-01-02T15:04:05.000Z07:00" |
| DisableTimestamp | `bool` | DisableTimestamp disables timestamp output. |
| DisableColors | `bool` | DisableColors disables ANSI colors. |
| DisableQuoting | `bool` | DisableQuoting disables quoting of string values. |
| QuoteEmptyFields | `bool` | QuoteEmptyFields quotes empty field values. |
| FullTimestamp | `bool` | FullTimestamp shows full timestamp instead of relative time. |
| PadLevelText | `bool` | PadLevelText pads level text for alignment. |
| FieldSeparator | `string` | FieldSeparator is the separator between fields. Default: " " |
| KeyValueSeparator | `string` | KeyValueSeparator is the separator between key and value. Default: "=" |

## Methods

### Format

Format formats an entry as text.

```go
func (*NamedFormatter) Format(entry *Entry) ([]byte, error)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
- []byte
- error

### WriterHook
WriterHook writes entries to an io.Writer.

#### Example Usage

```go
// Create a new WriterHook
writerhook := WriterHook{

}
```

#### Type Definition

```go
type WriterHook struct {
}
```

### Constructor Functions

### NewWriterHook

NewWriterHook creates a hook that writes to an io.Writer.

```go
func NewWriterHook(w io.Writer, formatter Formatter, levels ...Level) *WriterHook
```

**Parameters:**
- `w` (io.Writer)
- `formatter` (Formatter)
- `levels` (...Level)

**Returns:**
- *WriterHook

## Methods

### Fire

Fire implements Hook.

```go
func (*FuncHook) Fire(entry *Entry)
```

**Parameters:**
- `entry` (*Entry)

**Returns:**
  None

### Levels

Levels implements Hook.

```go
func (*FuncHook) Levels() []Level
```

**Parameters:**
  None

**Returns:**
- []Level

## Functions

### CtxDebug
CtxDebug logs at debug level using the logger from context.

```go
func CtxDebug(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `msg` | `string` | |
| `fields` | `...Field` | |

**Returns:**
None

**Example:**

```go
// Example usage of CtxDebug
result := CtxDebug(/* parameters */)
```

### CtxError
CtxError logs at error level using the logger from context.

```go
func CtxError(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `msg` | `string` | |
| `fields` | `...Field` | |

**Returns:**
None

**Example:**

```go
// Example usage of CtxError
result := CtxError(/* parameters */)
```

### CtxInfo
CtxInfo logs at info level using the logger from context.

```go
func CtxInfo(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `msg` | `string` | |
| `fields` | `...Field` | |

**Returns:**
None

**Example:**

```go
// Example usage of CtxInfo
result := CtxInfo(/* parameters */)
```

### CtxTrace
CtxTrace logs at trace level using the logger from context.

```go
func CtxTrace(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `msg` | `string` | |
| `fields` | `...Field` | |

**Returns:**
None

**Example:**

```go
// Example usage of CtxTrace
result := CtxTrace(/* parameters */)
```

### CtxWarn
CtxWarn logs at warn level using the logger from context.

```go
func CtxWarn(ctx context.Context, msg string, fields ...Field)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `msg` | `string` | |
| `fields` | `...Field` | |

**Returns:**
None

**Example:**

```go
// Example usage of CtxWarn
result := CtxWarn(/* parameters */)
```

### Debug
Debug logs at debug level using the default logger.

```go
func (*Builder) Debug(msg string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Debug
result := Debug(/* parameters */)
```

### Debugf
Debugf logs a formatted message at debug level.

```go
func Debugf(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Debugf
result := Debugf(/* parameters */)
```

### Error
Error logs at error level using the default logger.

```go
func Error(msg string, fields ...Field)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |
| `fields` | `...Field` | |

**Returns:**
None

**Example:**

```go
// Example usage of Error
result := Error(/* parameters */)
```

### Errorf
Errorf logs a formatted message at error level.

```go
func Errorf(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Errorf
result := Errorf(/* parameters */)
```

### Fatal
Fatal logs at fatal level using the default logger and exits.

```go
func (*Builder) Fatal(msg string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Fatal
result := Fatal(/* parameters */)
```

### Fatalf
Fatalf logs a formatted message at fatal level and exits.

```go
func Fatalf(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Fatalf
result := Fatalf(/* parameters */)
```

### Info
Info logs at info level using the default logger.

```go
func (*Builder) Info(msg string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Info
result := Info(/* parameters */)
```

### Infof
Infof logs a formatted message at info level.

```go
func Infof(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Infof
result := Infof(/* parameters */)
```

### Must
Must logs and panics if error is not nil. Useful for initialization code. db := log.Must(sql.Open("postgres", dsn))

```go
func Must(l *Logger, val T, err error) T
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `l` | `*Logger` | |
| `val` | `T` | |
| `err` | `error` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |

**Example:**

```go
// Example usage of Must
result := Must(/* parameters */)
```

### Panic
Panic logs at panic level using the default logger and panics.

```go
func (*Builder) Panic(msg string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Panic
result := Panic(/* parameters */)
```

### Panicf
Panicf logs a formatted message at panic level and panics.

```go
func Panicf(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Panicf
result := Panicf(/* parameters */)
```

### Print
Print logs a message at info level.

```go
func Print(args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Print
result := Print(/* parameters */)
```

### Printf
Printf logs a formatted message at info level.

```go
func Printf(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Printf
result := Printf(/* parameters */)
```

### Println
Println logs a message at info level.

```go
func Println(args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Println
result := Println(/* parameters */)
```

### SetDefault
SetDefault sets the default logger.

```go
func SetDefault(l *Logger)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `l` | `*Logger` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetDefault
result := SetDefault(/* parameters */)
```

### SetDefaultFormatter
SetDefaultFormatter sets the default formatter.

```go
func SetDefaultFormatter(f Formatter)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `f` | `Formatter` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetDefaultFormatter
result := SetDefaultFormatter(/* parameters */)
```

### SetDefaultLevel
SetDefaultLevel sets the default level.

```go
func SetDefaultLevel(level Level)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `level` | `Level` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetDefaultLevel
result := SetDefaultLevel(/* parameters */)
```

### Trace
Trace logs at trace level using the default logger.

```go
func (*Builder) Trace(msg string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Trace
result := Trace(/* parameters */)
```

### Tracef
Tracef logs a formatted message at trace level.

```go
func Tracef(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Tracef
result := Tracef(/* parameters */)
```

### Warn
Warn logs at warn level using the default logger.

```go
func (*Builder) Warn(msg string)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `msg` | `string` | |

**Returns:**
None

**Example:**

```go
// Example usage of Warn
result := Warn(/* parameters */)
```

### Warnf
Warnf logs a formatted message at warn level.

```go
func Warnf(format string, args ...any)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `format` | `string` | |
| `args` | `...any` | |

**Returns:**
None

**Example:**

```go
// Example usage of Warnf
result := Warnf(/* parameters */)
```

### WithContextFields
WithFields adds fields to the context that will be included in all logs.

```go
func WithContextFields(ctx context.Context, fields ...Field) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `fields` | `...Field` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithContextFields
result := WithContextFields(/* parameters */)
```

### WithLogger
WithLogger attaches a logger to the context.

```go
func WithLogger(ctx context.Context, logger *Logger) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `logger` | `*Logger` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithLogger
result := WithLogger(/* parameters */)
```

### WithRequestID
WithRequestID adds a request ID to the context.

```go
func WithRequestID(ctx context.Context, requestID string) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `requestID` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithRequestID
result := WithRequestID(/* parameters */)
```

### WithTraceID
WithTraceID adds a trace ID to the context.

```go
func WithTraceID(ctx context.Context, traceID string) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `traceID` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithTraceID
result := WithTraceID(/* parameters */)
```

### WithUserID
WithUserID adds a user ID to the context.

```go
func WithUserID(ctx context.Context, userID string) context.Context
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | |
| `userID` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `context.Context` | |

**Example:**

```go
// Example usage of WithUserID
result := WithUserID(/* parameters */)
```

## External Links

- [Package Overview](../packages/logs.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/helix/logs)
- [Source Code](https://github.com/kolosys/helix/tree/main/logs)
