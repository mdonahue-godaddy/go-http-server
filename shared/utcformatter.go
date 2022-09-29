package shared

import (
	log "github.com/sirupsen/logrus"
)

// UTCFormatter - used to format to UTC
type UTCFormatter struct {
	Formatter log.Formatter
}

// Format - formats time to UTC
func (u UTCFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}
