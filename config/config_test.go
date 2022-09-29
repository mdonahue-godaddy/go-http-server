package config_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mdonahue-godaddy/go-http-server/config"
)

const (
	testsDir string = "tests"
)

func Test_LoadSettings_Good(t *testing.T) {
	settingsFileName := filepath.Join(testsDir, "good.json")

	actual, err := config.LoadSettings(settingsFileName)

	assert.Nil(t, err, "Error should be nil.")
	assert.NotNil(t, actual, "Settings should NOT be nil.")
	assert.Equal(t, "12s", actual.Service.GracefulShutdownDelaySeconds, "Settings.Service.GracefulShutdownDelaySeconds")
	assert.Equal(t, "0.0.0.0", actual.Service.HTTP.Server.IPv4Address, "Settings.Service.HTTP.Server.IPv4Address")
	assert.Equal(t, uint16(8433), actual.Service.HTTP.Server.Port, "Settings.Service.HTTP.Server.Port")
	assert.Equal(t, "0.0.0.0", actual.Metrics.HTTP.Server.IPv4Address, "Settings.Metrics.HTTP.Server.Address")
	assert.Equal(t, uint16(8081), actual.Metrics.HTTP.Server.Port, "Settings.Metrics.HTTP.Server.Port")
	assert.Equal(t, true, actual.Metrics.HealthcheckEnabled, "Settings.Metrics.HealthcheckEnabled")
	assert.Equal(t, true, actual.Metrics.PPRofEnabled, "Settings.Metrics.PPRofEnabled")
	assert.Equal(t, "debug", actual.Logging.Level, "Settings.Logging.Level")
}

func Test_LoadSettings_Empty(t *testing.T) {
	settingsFileName := filepath.Join(testsDir, "empty.json")

	actual, err := config.LoadSettings(settingsFileName)

	assert.Nil(t, err, "Error should be nil.")
	assert.NotNil(t, actual, "Settings should NOT be nil.")
	assert.Equal(t, "", actual.Service.GracefulShutdownDelaySeconds, "Settings.Service.GracefulShutdownDelaySeconds")
	assert.Equal(t, "", actual.Service.HTTP.Server.IPv4Address, "Settings.Service.HTTP.Server.IPv4Address")
	assert.Equal(t, uint16(0), actual.Service.HTTP.Server.Port, "Settings.Service.HTTP.Server.Port")
	assert.Equal(t, "", actual.Metrics.HTTP.Server.IPv4Address, "Settings.Metrics.HTTP.Server.Address")
	assert.Equal(t, uint16(0), actual.Metrics.HTTP.Server.Port, "Settings.Metrics.HTTP.Server.Port")
	assert.Equal(t, false, actual.Metrics.HealthcheckEnabled, "Settings.Metrics.HealthcheckEnabled")
	assert.Equal(t, false, actual.Metrics.PPRofEnabled, "Settings.Metrics.PPRofEnabled")
	assert.Equal(t, "", actual.Logging.Level, "Settings.Logging.Level")
}

func Test_LoadSettings_BadFile(t *testing.T) {
	settingsFileName := filepath.Join(testsDir, "bad.json")

	actual, err := config.LoadSettings(settingsFileName)

	assert.NotNil(t, err, "Error should NOT be nil.")
	assert.Nil(t, actual, "Settings should be nil.")
	assert.Equal(t, "invalid character 'B' looking for beginning of object key string", err.Error(), "Error should be 'bad json'")
}

func Test_LoadSettings_MissingFile(t *testing.T) {
	settingsFileName := filepath.Join(testsDir, "i.do.not.exist.jsaon")

	actual, err := config.LoadSettings(settingsFileName)

	assert.NotNil(t, err, "Error should NOT be nil.")
	assert.Nil(t, actual, "Settings should be nil.")
	assert.Equal(t, "open tests/i.do.not.exist.jsaon: no such file or directory", err.Error(), "Error should be 'bad json'")
}
