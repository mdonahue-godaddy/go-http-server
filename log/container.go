package log

import (
	"context"
	"strings"

	"github.com/rs/zerolog"

	"github.com/mdonahue-godaddy/go-http-server/metadata"
)

func Container() Option {
	if metadata.Disabled() {
		return func(l zerolog.Logger) zerolog.Logger {
			return l // nop
		}
	}

	metadata, err := metadata.Container(context.Background())
	if err != nil {
		panic(err)
	}
	parts := strings.Split(metadata.Image, ":")
	name, tags := parts[0], []string{parts[1]}

	return func(l zerolog.Logger) zerolog.Logger {
		return l.With().
			Dict("container", zerolog.Dict().
				Str("id", metadata.DockerID).
				Str("name", metadata.Name).
				Interface("labels", metadata.Labels).
				Str("runtime", "fargate").
				Dict("image", zerolog.Dict().
					Str("name", name).
					Strs("tag", tags))).Logger() // key is singular but its value is plural? Yep, see: https://www.elastic.co/guide/en/ecs/current/ecs-container.html#field-container-image-tag
	}
}
