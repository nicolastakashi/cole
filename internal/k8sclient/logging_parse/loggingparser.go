package logging_parse

import (
	"io"

	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/sirupsen/logrus"
)

type LoggingParser interface {
	Parse(stream io.ReadCloser) ([]entities.LogLine, error)
}

func Get(lg string) LoggingParser {
	logrus.Debug("getting log parser %v", lg)
	if lg == "json" {
		return JsonLoggingParse{}
	}
	return ConsoleLoggingParse{}
}
