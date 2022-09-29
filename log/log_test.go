package log_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/mdonahue-godaddy/go-http-server/log"
)

type LogSuite struct {
	suite.Suite
}

func TestLogSuite(t *testing.T) {
	suite.Run(t, &LogSuite{})
}

type testCtx struct {
	context.Context
}

func (t *testCtx) Set(key string, value interface{}) {
	t.Context = context.WithValue(t.Context, key, value)
}

func (s *LogSuite) TestContextAssociation() {
	ctx := &testCtx{context.Background()}

	logger := log.CtxOrDefault(ctx)
	s.Equal(log.DefaultLogger, *logger)

	newLogger := log.NewLogger()
	log.Associate(ctx, &newLogger)

	logger = log.CtxOrDefault(ctx)
	s.Equal(newLogger, *logger)
}

func (s *LogSuite) TestLevels() {
	s.Equal(log.DefaultLogger.Debug(), log.Debug())
	s.Equal(log.DefaultLogger.Info(), log.Info())
	s.Equal(log.DefaultLogger.Warn(), log.Warn())
	s.Equal(log.DefaultLogger.Error(), log.Error())

	err := fmt.Errorf("boom")
	s.Equal(log.DefaultLogger.Err(err), log.Err(err))
}

func Example_wrappedEvent_Tags() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger()

	logger.Info().Tags(log.ApplicationTag, log.SecurityTag).Msg("testing")
	// Output: {"tags":["application","security"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func setup() (teardown func()) {
	ow := log.DefaultWriter
	ot := log.DefaultTimestampFunc

	log.DefaultWriter = os.Stdout
	log.DefaultTimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}

	log.EnableECS()

	return func() {
		log.DefaultWriter = ow
		log.DefaultTimestampFunc = ot

		log.EnableECS()
	}
}

func ExampleClient() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger(log.Client("1.2.3.4", "arn:asws:iam::111111111111:user/iam-user"))

	logger.Info().Msg("testing")
	// Output: {"client":{"ip":"1.2.3.4","user":{"id":"arn:asws:iam::111111111111:user/iam-user"}},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func ExampleCloud() {
	teardown := setup()
	defer teardown()

	os.Setenv("AWS_REGION", "us-test-1")
	logger := log.NewLogger(log.Cloud("112233445566"))

	logger.Info().Msg("testing")
	// Output: {"cloud":{"provider":"aws","region":"us-test-1","availability_zone":"","service":{"name":"fargate"},"account":{"id":"112233445566"}},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func ExampleService() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger(log.Service("example", "unittest"))

	logger.Info().Msg("testing")
	// Output: {"ecs":{"version":"1.12.0"},"service":{"name":"example","version":"0.0.0-local","environment":"unittest"},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func ExampleLogger_WithRequest() {
	teardown := setup()
	defer teardown()

	var (
		method    = "GET"
		requestId = "some-request-id"
		path      = "/some/path"
	)

	logger := log.NewLogger()
	logger = logger.WithRequest(method, path, requestId)

	logger.Info().Msg("testing")
	// Output: {"url":{"path":"/some/path"},"http":{"request":{"id":"some-request-id","method":"GET"}},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func ExampleLogger_WithResponse() {
	teardown := setup()
	defer teardown()

	var (
		method    = "GET"
		requestId = "some-request-id"
		path      = "/some/path"
	)

	local := log.NewLogger()
	logger := local.WithResponse(method, path, requestId, 200)

	logger.Info().Msg("has req/resp")

	// verify context for local is unmodified
	local.Info().Msg("no req/resp")

	// Output: {"url":{"path":"/some/path"},"http":{"request":{"id":"some-request-id","method":"GET"},"response":{"status_code":200}},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"has req/resp"}
	// {"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"no req/resp"}
}

func Example_wrappedEvent_ECSEvent() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger()

	logger.Info().ECSEvent(
		log.Database,
		log.Success,
		log.CreationType, log.AllowedType).Msg("testing")
	// Output: {"event":{"kind":"event","category":"database","type":["creation","allowed"],"outcome":"success"},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func Example_wrappedEvent_ECSTimedEvent() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger()

	logger.Info().ECSTimedEvent(
		log.Database,
		log.Success,
		time.Millisecond,
		log.CreationType, log.AllowedType).Msg("testing")
	// Output: {"event":{"kind":"event","category":"database","type":["creation","allowed"],"outcome":"success","duration":1000000},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func Example_wrappedEvent_Fields() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger()

	logger.Info().
		Dict("dict", zerolog.Dict()).
		Str("str", "val").
		Strs("strs", []string{"a", "b"}).
		Int("int", 0).
		Interface("interface", struct{}{}).
		Msg("testing")
	// Output: {"dict":{},"str":"val","strs":["a","b"],"int":0,"interface":{},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}
