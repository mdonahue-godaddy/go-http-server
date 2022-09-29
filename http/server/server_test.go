package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/config"
	"github.com/mdonahue-godaddy/go-http-server/http/server"
	"github.com/mdonahue-godaddy/go-http-server/shared"
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
		actual := server.NewServer(tc.ServiceName, tc.Cfg, nil)

		assert.Equal(tc.Cfg, actual.GetConfig(), tc.Description)
		assert.Equal(tc.ServiceName, actual.GetServiceName(), tc.Description)
		assert.Equal(false, actual.IsInitialized(), tc.Description)
		assert.Equal(false, actual.IsShuttingDown(), tc.Description)
	}
}

func Test_generateHtmlBodyFromTemplate(t *testing.T) {
	assert := assert.New(t)

	cfg := &config.Settings{}
	template := "<!DOCTYPE HTML><html lang='en-us'><head><title>** {{page_title}} **</title></head><body>{{page_body}}</body></html>"
	svc := server.NewServer("TestServiceName", cfg, &template)

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
		actual := svc.GenerateHtmlBodyFromTemplate(tc.PageTitle, tc.PageBody)

		assert.Equal(tc.Expected, actual, tc.Description)
	}
}
