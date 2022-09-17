package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/pkg/config"
	"github.com/mdonahue-godaddy/go-http-server/pkg/shared"
)

// CreateTestContext -
func CreateTestContext(actionName, actionType string) context.Context {
	ctx := shared.CreateContext(context.Background(), actionName, actionType)
	return ctx
}

func Test_NewServer(t *testing.T) {
	assert := assert.New(t)

	cfg := config.Settings{}

	testCases := []struct {
		ServiceName string
		Cfg         *config.Settings
		Expected    bool
		Description string
	}{
		{
			ServiceName: "MyServiceName",
			Cfg:         nil,
			Expected:    true,
			Description: "nil config",
		},
		{
			ServiceName: "MyServiceName",
			Cfg:         &cfg,
			Expected:    true,
			Description: "config",
		},
	}

	for _, tc := range testCases {
		actual := NewServer(tc.ServiceName, tc.Cfg)

		assert.Equal(tc.Cfg, actual.config, tc.Description)
		assert.Equal(tc.ServiceName, actual.serviceName, tc.Description)
		assert.Equal(false, actual.isInitialized, tc.Description)
		assert.Equal(false, actual.isShuttingDown, tc.Description)
	}
}

func Test_generateHtmlBodyFromTemplate(t *testing.T) {
	assert := assert.New(t)

	svc := &Server{
		responseTemplateFile: defaultResponseTemplateFile,
	}

	testCases := []struct {
		PageTitle   string
		PageBody    string
		Expected    string
		Description string
	}{
		{
			PageTitle:   "",
			PageBody:    "",
			Expected:    "<!DOCTYPE HTML><html lang='en-us'><head><title>**  **</title></head><body></body></html>",
			Description: "empty title and body",
		},
		{
			PageTitle:   "",
			PageBody:    "body",
			Expected:    "<!DOCTYPE HTML><html lang='en-us'><head><title>**  **</title></head><body>body</body></html>",
			Description: "empty title with body",
		},
		{
			PageTitle:   "title",
			PageBody:    "",
			Expected:    "<!DOCTYPE HTML><html lang='en-us'><head><title>** title **</title></head><body></body></html>",
			Description: "title and empty body",
		},
		{
			PageTitle:   "title",
			PageBody:    "body",
			Expected:    "<!DOCTYPE HTML><html lang='en-us'><head><title>** title **</title></head><body>body</body></html>",
			Description: "title and body",
		},
	}

	for _, tc := range testCases {
		actual := svc.generateHtmlBodyFromTemplate(tc.PageTitle, tc.PageBody)

		assert.Equal(tc.Expected, actual, tc.Description)
	}
}
