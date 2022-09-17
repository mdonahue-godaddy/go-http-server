package shared

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_UTCFormatter(t *testing.T) {
	assert := assert.New(t)

	jsonFormatter := log.JSONFormatter{
		TimestampFormat: TimestampFormat,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
			log.FieldKeyFunc:  "caller",
		},
	}
	fmtr := UTCFormatter{
		&jsonFormatter,
	}

	logger := log.Logger{
		Out:          nil,
		Hooks:        map[log.Level][]log.Hook{},
		Formatter:    nil,
		ReportCaller: false,
		Level:        0,
		ExitFunc: func(int) {
		},
	}

	entry := log.Entry{Logger: &logger, Data: log.Fields{}, Time: time.Now(), Level: log.DebugLevel}

	actual, err := fmtr.Format(&entry)
	assert.Nil(err, "error is nil")
	assert.NotNil(actual, "results is not nil")
}
