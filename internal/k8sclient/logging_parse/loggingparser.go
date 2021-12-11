package logging_parse

import (
	"io"

	"github.com/nicolastakashi/cole/internal/entities"
)

type LoggingParser interface {
	Parse(stream io.ReadCloser) ([]entities.LogLine, error)
}

func Get(lg string) LoggingParser {
	if lg == "json" {
		return JsonLoggingParse{}
	}
	return ConsoleLoggingParse{}
}
