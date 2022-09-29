package metadata_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/mdonahue-godaddy/go-http-server/metadata"
)

type MetadataSuite struct {
	suite.Suite

	origURI       string
	restoreNeeded bool
	server        *httptest.Server
}

func TestMetadataSuite(t *testing.T) {
	suite.Run(t, &MetadataSuite{})
}

func (s *MetadataSuite) SetupTest() {
	s.origURI, s.restoreNeeded = os.LookupEnv(metadata.EnvKeyFargateMetadataURI)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, metadata.ExampleECSContainerMetadata)
	}))
	mux.Handle("/task", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, metadata.ExampleECSTaskMetadata)
	}))

	s.server = httptest.NewServer(mux)

	metadata.Client = s.server.Client()

	os.Setenv(metadata.EnvKeyFargateMetadataURI, s.server.URL)
}

func (s *MetadataSuite) TearDownTest() {
	if s.restoreNeeded {
		os.Setenv(metadata.EnvKeyFargateMetadataURI, s.origURI)
	}

	s.server.Close()
}

func (s *MetadataSuite) TestDisabled() {
	s.False(metadata.Disabled(), "metadata.Disabled() expected to be false")
	s.True(metadata.Enabled(), "metadata.Enabled() expected to be true")

	os.Setenv(metadata.EnvKeyFargateMetadataURI, "")
	s.True(metadata.Disabled(), "metadata.Disabled() expected to be true")
	s.False(metadata.Enabled(), "metadata.Enabled() expected to be false")
}

func (s *MetadataSuite) TestContainer() {
	data, err := metadata.Container(context.Background())
	s.NoError(err)

	s.Equal("curl", data.Name)
	s.NotEmpty(data.Labels)
}

func (s *MetadataSuite) TestTask() {
	data, err := metadata.Task(context.Background())
	s.NoError(err)

	s.Equal("curltest", data.Family)
	s.NotEmpty(data.Containers)
}

func (s *MetadataSuite) TestMetadataUnavailableError() {
	os.Setenv(metadata.EnvKeyFargateMetadataURI, "")

	_, err := metadata.Container(context.Background())
	s.ErrorIs(err, metadata.ErrMetadataUnavailable)

	_, err = metadata.Task(context.Background())
	s.ErrorIs(err, metadata.ErrMetadataUnavailable)
}
