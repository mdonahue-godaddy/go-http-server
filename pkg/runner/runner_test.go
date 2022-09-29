package runner_test

import (
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/pkg/runner"
)

func RestoreEnvironment(pairs []string) {
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		os.Setenv(parts[0], parts[1])
	}
}

func Test_GetAllEnvironmentVariables(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		FargateCluster string
		FargateService string
		Expected       string
		Description    string
	}{
		{
			FargateCluster: "FargateCluster",
			FargateService: "FargateService",
			Expected:       "'FARGATE_CLUSTER'='FargateCluster', 'FARGATE_SERVICE'='FargateService'",
			Description:    "test Good",
		},
		{
			FargateCluster: "",
			FargateService: "",
			Expected:       "'FARGATE_CLUSTER'='', 'FARGATE_SERVICE'=''",
			Description:    "test empty",
		},
	}

	// Save and Clear environment variables
	savedEnviron := os.Environ()
	os.Clearenv()

	for _, tc := range testCases {
		os.Setenv("FARGATE_CLUSTER", tc.FargateCluster)
		os.Setenv("FARGATE_SERVICE", tc.FargateService)

		actual := runner.GetAllEnvironmentVariables()

		assert.NotEmpty(actual, tc.Description)
		assert.Len(actual, 2, tc.Description)
		key := "FARGATE_CLUSTER"
		assert.Equal(tc.FargateCluster, actual[key], tc.Description)
		key = "FARGATE_SERVICE"
		assert.Equal(tc.FargateService, actual[key], tc.Description)
	}

	// Restore environment variables
	RestoreEnvironment(savedEnviron)
}

func Test_SetupDefaultLogrusConfig(t *testing.T) {
	assert := assert.New(t)

	runner.SetupDefaultLogrusConfig()

	level := log.GetLevel()

	assert.Equal(log.InfoLevel, level, "Log Level")
	assert.Equal(log.StandardLogger().Out, os.Stdout, "Log Our is nor os.Stdout")
}
