package metadata

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	metadata "github.com/brunoscheufler/aws-ecs-metadata-go"
)

const EnvKeyFargateMetadataURI = "ECS_CONTAINER_METADATA_URI_V4"

var (
	ErrMetadataUnavailable = errors.New("metadata retrieval is not available")
	ErrUnsupportedMetadata = errors.New("unsupported metadata")
)

// Client is the HTTP client used to interact with the Fargate metadata endpoint.
var Client = http.DefaultClient

// Disabled returns true if metadata retrieval is unsupported.
func Disabled() bool {
	return os.Getenv(EnvKeyFargateMetadataURI) == ""
}

// Enabled returns true if metadata retrieval is supported.
func Enabled() bool {
	return !Disabled()
}

type containerMetadataCache struct {
	sync.Mutex

	metadata *metadata.ContainerMetadataV4
}

func (c *containerMetadataCache) get(ctx context.Context) (*metadata.ContainerMetadataV4, error) {
	c.Lock()
	defer c.Unlock()

	if c.metadata != nil {
		return c.metadata, nil
	}

	data, err := metadata.GetContainer(ctx, Client)
	if err != nil {
		return nil, err
	}

	ok := false
	if c.metadata, ok = data.(*metadata.ContainerMetadataV4); ok {
		return c.metadata, nil
	}

	return nil, fmt.Errorf("%w", ErrUnsupportedMetadata)
}

var container containerMetadataCache

// Container retrieves the service metadata from the underlying provider and returns the results.
func Container(ctx context.Context) (*metadata.ContainerMetadataV4, error) {
	if Disabled() {
		return nil, fmt.Errorf("%w", ErrMetadataUnavailable)
	}

	return container.get(ctx)
}

type taskMetadataCache struct {
	sync.Mutex

	metadata *metadata.TaskMetadataV4
}

func (t *taskMetadataCache) get(ctx context.Context) (*metadata.TaskMetadataV4, error) {
	t.Lock()
	defer t.Unlock()

	if t.metadata != nil {
		return t.metadata, nil
	}

	data, err := metadata.Get(ctx, Client)
	if err != nil {
		return nil, err
	}

	ok := false
	if t.metadata, ok = data.(*metadata.TaskMetadataV4); ok {
		return t.metadata, nil
	}

	return nil, fmt.Errorf("%w", ErrUnsupportedMetadata)
}

var task taskMetadataCache

// Task retrieves the service metadata from the underlying provider and returns the results.
func Task(ctx context.Context) (*metadata.TaskMetadataV4, error) {
	if Disabled() {
		return nil, fmt.Errorf("%w", ErrMetadataUnavailable)
	}

	return task.get(ctx)
}
