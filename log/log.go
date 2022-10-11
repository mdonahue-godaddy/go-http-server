package log

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmzerolog"

	"github.com/mdonahue-godaddy/go-http-server/metadata"
	"github.com/mdonahue-godaddy/go-http-server/version"
)

const (
	DefaultEnvironment     = "local"
	DefaultSource          = "pstore-api"
	DefaultTimeFieldFormat = "2006-01-02T15:04:05.999Z" // RFC3339 at millisecond resolution in zulu timezone
)

var (
	DefaultLogger                  = NewLogger(Service(DefaultSource, DefaultEnvironment))
	DefaultTimestampFunc           = func() time.Time { return time.Now().UTC() }
	DefaultWriter        io.Writer = zerolog.MultiLevelWriter(os.Stdout, &apmzerolog.Writer{})
)

// EnableECS configures the global settings
func EnableECS() {
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = ""
	zerolog.TimeFieldFormat = DefaultTimeFieldFormat
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.TimestampFunc = DefaultTimestampFunc
	zerolog.LevelFieldName = ""
	zerolog.CallerSkipFrameCount = 4
	zerolog.ErrorStackMarshaler = apmzerolog.MarshalErrorStack
}

type Option func(l zerolog.Logger) zerolog.Logger

func Client(ip string, user string) Option {
	return func(l zerolog.Logger) zerolog.Logger {
		return SetClient(ip, user)(l.With()).Logger()
	}
}

func SetClient(ip string, user string) func(zerolog.Context) zerolog.Context {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.
			Dict("client", zerolog.Dict().
				Str("ip", ip).
				Dict("user", zerolog.Dict().
					Str("id", user)))
	}
}

// Cloud returns an Option that can be used to add an ECS "cloud" field to a new logger's context.
// Note that if metadata retrieval fails for any reason some values will be missing or empty.
func Cloud(accountId string) Option {
	var (
		region           string = os.Getenv("AWS_REGION")
		availabilityZone string
	)

	if task, err := metadata.Task(context.Background()); err == nil {
		availabilityZone = task.AvailabilityZone
	}

	return func(l zerolog.Logger) zerolog.Logger {
		return l.With().
			Dict("cloud", zerolog.Dict().
				Str("provider", "aws").
				Str("region", region).
				Str("availability_zone", availabilityZone).
				Dict("service", zerolog.Dict().
					Str("name", "fargate")).
				Dict("account", zerolog.Dict().
					Str("id", accountId))).Logger()
	}
}

func Service(name, environment string) Option {
	ver := version.Version
	if len(version.Version) == 0 {
		ver = "0.0.0-local"
	}

	return func(l zerolog.Logger) zerolog.Logger {
		return l.With().
			Dict("ecs", zerolog.Dict().
				Str("version", "1.12.0")).
			Dict("service", zerolog.Dict().
				Str("name", name).
				Str("version", ver).
				Str("environment", environment)).
			Logger()
	}
}

type nestedLevelHook struct{}

func (h nestedLevelHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Dict("log", zerolog.Dict().
			Str("level", level.String()).
			Str("logger", "github.com/rs/zerolog"))
	}
}

type Logger struct {
	zerolog.Logger

	tags []tag
}

func NewLogger(opts ...Option) Logger {
	l := zerolog.New(DefaultWriter).With().Timestamp().Logger().Hook(nestedLevelHook{})
	for _, option := range opts {
		l = option(l)
	}

	return Logger{Logger: l}
}

func (l *Logger) SetDefaultTags(tags ...tag) {
	l.tags = tags
}

func (l *Logger) SetTraceContext(ctx context.Context) {
	l.Logger = l.Logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		tx := apm.TransactionFromContext(ctx)
		if tx == nil {
			return
		}

		traceContext := tx.TraceContext()

		e.Dict("trace", zerolog.Dict().
			Hex("id", traceContext.Trace[:]))

		e.Dict("transaction", zerolog.Dict().
			Hex("id", traceContext.Span[:]))

		if span := apm.SpanFromContext(ctx); span != nil {
			spanTraceContext := span.TraceContext()

			e.Dict("span", zerolog.Dict().
				Hex("id", spanTraceContext.Span[:]))
		}
	}))
}

func (l *Logger) WithContext(ctx context.Context) context.Context {
	if lp, ok := ctx.Value(logContextKey).(*Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	} else if l.GetLevel() == zerolog.Disabled {
		// Do not store disabled logger.
		return ctx
	}
	return context.WithValue(ctx, logContextKey, l)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l Logger) Debug() *wrappedEvent {
	return l.newWrappedEvent(l.Logger.Debug())
}

func (l *Logger) newWrappedEvent(e *zerolog.Event) *wrappedEvent {
	return &wrappedEvent{
		Event: e,
		tags:  l.tags,
	}
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l Logger) Info() *wrappedEvent {
	return l.newWrappedEvent(l.Logger.Info())
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l Logger) Warn() *wrappedEvent {
	return l.newWrappedEvent(l.Logger.Warn())
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l Logger) Error() *wrappedEvent {
	return l.newWrappedEvent(l.Logger.Error())
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func (l Logger) Err(err error) *wrappedEvent {
	if err != nil {
		return l.Error().Err(err)
	}

	return l.Info()
}

func (l Logger) WithRequest(method string, path string, requestId string) Logger {
	return Logger{
		Logger: l.With().
			Dict("url", zerolog.Dict().
				Str("path", path)).
			Dict("http", zerolog.Dict().
				Dict("request", zerolog.Dict().
					Str("id", requestId).
					Str("method", method))).Logger(),

		tags: l.tags,
	}
}

func (l Logger) WithResponse(method string, path string, requestId string, status int) Logger {
	return Logger{
		Logger: l.With().
			Dict("url", zerolog.Dict().
				Str("path", path)).
			Dict("http", zerolog.Dict().
				Dict("request", zerolog.Dict().
					Str("id", requestId).
					Str("method", method)).
				Dict("response", zerolog.Dict().
					Int("status_code", status))).Logger(),

		tags: l.tags,
	}
}

type tag string

const (
	ApplicationTag tag = "application"
	AsherahTag     tag = "asherah"
	MetricsTag     tag = "metrics"
	RequestTag     tag = "request"
	SecurityTag    tag = "security"
)

// Supported event categories are below, for full list see:
// https://www.elastic.co/guide/en/ecs/current/ecs-allowed-values-event-category.html
type eventCategory string

var (
	// Authentication event category.
	// Expected event types include: start, end, info
	Authentication eventCategory = "authentication"

	// Configuration event category.
	// Expected event types include: access, change, creation, deletion, info
	Configuration eventCategory = "configuration"

	// Database event category.
	// Expected event types include: access, change, info, error
	Database eventCategory = "database"

	// Process event category.
	// Expected event types include: access, change, end, info, start
	Process eventCategory = "process"

	// Web event category.
	// Expected event types include: access, error, info
	Web eventCategory = "web"
)

// Supported event types are below, for full list see:
// https://www.elastic.co/guide/en/ecs/current/ecs-event.html#field-event-type
type eventType string

var (
	AccessType   eventType = "access"
	CreationType eventType = "creation"
	DeletionType eventType = "deletion"

	InformationType eventType = "info"

	AllowedType eventType = "allowed"
	DeniedType  eventType = "denied"
	ErrorType   eventType = "error"

	StartType eventType = "start"
	EndType   eventType = "end"
)

type EventTypes []eventType

// Supported event outcomes are below, for full list see:
// https://www.elastic.co/guide/en/ecs/current/ecs-event.html#field-event-outcome
type eventOutcome string

var (
	Failure       eventOutcome = "failure"
	Success       eventOutcome = "success"
	Unknown       eventOutcome = "unknown"
	NotApplicable eventOutcome = "n/a"
)

type wrappedEvent struct {
	*zerolog.Event

	tags []tag
}

// Interface adds the field key with i marshaled using reflection.
func (w *wrappedEvent) Interface(key string, i interface{}) *wrappedEvent {
	w.Event.Interface(key, i)
	return w
}

// Err adds the field "error" with serialized err to the wrapped *zerolog.Event context.
// If err is nil, no field is added.
func (w *wrappedEvent) Err(err error) *wrappedEvent {
	if err != nil {
		w.Event.Dict("error", zerolog.Dict().Str("message", err.Error()))
	}

	return w
}

// Tags sets the tags to be added to the log context when w.Msg is called.
func (w *wrappedEvent) Tags(tags ...tag) *wrappedEvent {
	w.tags = tags
	return w
}

// ECSEvent adds ECS "event" field to the wrapped *zerolog.Event context.
func (l *wrappedEvent) ECSEvent(category eventCategory, outcome eventOutcome, types ...eventType) *wrappedEvent {
	return l.ECSTimedEvent(category, outcome, -1, types...)
}

// ECSTimedEvent adds ECS "event" field with a duration to the wrapped *zerolog.Event context.
func (w *wrappedEvent) ECSTimedEvent(category eventCategory, outcome eventOutcome, duration time.Duration, types ...eventType) *wrappedEvent {
	var strs []string

	switch len(types) {
	case 0:
		strs = []string{string(InformationType)}
	default:
		strs = make([]string, len(types))
		for i := range types {
			strs[i] = string(types[i])
		}
	}

	detail := zerolog.Dict().
		Str("kind", "event").
		Str("category", string(category)).
		Strs("type", strs)

	if outcome != NotApplicable {
		detail = detail.Str("outcome", string(outcome))
	}

	if duration > 0 {
		detail = detail.Int64("duration", duration.Nanoseconds())
	}

	return w.Dict("event", detail)
}

// ECSMetric adds ECS "event" field with a duration to the wrapped *zerolog.Event context.
func (w *wrappedEvent) ECSMetric(types ...eventType) *wrappedEvent {
	var strs []string

	switch len(types) {
	case 0:
		strs = []string{string(InformationType)}
	default:
		strs = make([]string, len(types))
		for i := range types {
			strs[i] = string(types[i])
		}
	}

	detail := zerolog.Dict().
		Str("kind", "metric").
		Strs("category", []string{string(Database), string(Process), string(Web)}).
		Strs("type", strs)

	return w.Dict("event", detail)
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (w *wrappedEvent) Dict(key string, dict *zerolog.Event) *wrappedEvent {
	w.Event.Dict(key, dict)
	return w
}

func (w *wrappedEvent) Msg(message string) {
	w.withTags().Event.Msg(message)
}

// Str adds the field key with val as a string to the wrapped event context.
func (w *wrappedEvent) Str(key, val string) *wrappedEvent {
	w.Event.Str(key, val)
	return w
}

// Strs adds the field key with vals as a []string to the wrapped event context.
func (w *wrappedEvent) Strs(key string, vals []string) *wrappedEvent {
	w.Event.Strs(key, vals)
	return w
}

// Int adds the field key with i as a int to the wrapped event context.
func (w *wrappedEvent) Int(key string, i int) *wrappedEvent {
	w.Event.Int(key, i)
	return w
}

func (w *wrappedEvent) withTags() *wrappedEvent {
	if len(w.tags) == 0 {
		w.tags = []tag{ApplicationTag}
	}

	strs := make([]string, len(w.tags))
	for i := range w.tags {
		strs[i] = string(w.tags[i])
	}

	w.Event.Strs("tags", strs)

	return w
}

type contextKey string

const logContextKey contextKey = "github.com/gdcorp-appservices/pstore/log;contextkey=log"

type ContextSetter interface {
	context.Context
	Set(key string, value interface{})
}

func Associate(ctx ContextSetter, log *Logger) {
	ctx.Set(string(logContextKey), log)
}

func CtxOrDefault(ctx context.Context) *Logger {
	logger := &DefaultLogger

	var key interface{}

	switch ctx.(type) {
	case ContextSetter:
		key = string(logContextKey)
	default:
		key = logContextKey
	}

	if l, ok := ctx.Value(key).(*Logger); ok {
		logger = l
	}

	return logger
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *wrappedEvent {
	return DefaultLogger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *wrappedEvent {
	return DefaultLogger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *wrappedEvent {
	return DefaultLogger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *wrappedEvent {
	return DefaultLogger.Error()
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *wrappedEvent {
	return DefaultLogger.Err(err)
}
