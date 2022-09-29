package log_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/mdonahue-godaddy/go-http-server/log"
	"github.com/mdonahue-godaddy/go-http-server/metadata"
)

func ExampleContainer() {
	teardown := setupContainerMetadata()
	defer teardown()

	logger := log.NewLogger(log.Container())

	logger.Info().Msg("testing")
	// Output: {"container":{"id":"ea32192c8553fbff06c9340478a2ff089b2bb5646fb718b4ee206641c9086d66","name":"curl","labels":{"com.amazonaws.ecs.cluster":"default","com.amazonaws.ecs.container-name":"curl","com.amazonaws.ecs.task-arn":"arn:aws:ecs:us-west-2:111122223333:task/default/8f03e41243824aea923aca126495f665","com.amazonaws.ecs.task-definition-family":"curltest","com.amazonaws.ecs.task-definition-version":"24"},"runtime":"fargate","image":{"name":"111122223333.dkr.ecr.us-west-2.amazonaws.com/curltest","tag":["latest"]}},"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}

func setupContainerMetadata() func() {
	key := metadata.EnvKeyFargateMetadataURI
	teardown := setup()
	orig := os.Getenv(key)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, metadata.ExampleECSContainerMetadata)
	}))

	metadata.Client = ts.Client()

	os.Setenv(key, ts.URL)

	return func() {
		os.Setenv(key, orig)

		ts.Close()

		teardown()
	}
}

// ExampleContainer_nop demonstrates the use of the Container option when container metadata is unavailable.
func ExampleContainer_nop() {
	teardown := setup()
	defer teardown()

	logger := log.NewLogger(log.Container())

	logger.Info().Msg("testing")
	// Output: {"tags":["application"],"@timestamp":"2008-01-08T17:05:05Z","log":{"level":"info","logger":"github.com/rs/zerolog"},"message":"testing"}
}
