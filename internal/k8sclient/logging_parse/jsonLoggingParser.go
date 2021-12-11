package logging_parse

import (
	"io"

	"github.com/nicolastakashi/cole/internal/entities"
)

type JsonLoggingParse struct {
}

func (JsonLoggingParse) Parse(stream io.ReadCloser) ([]entities.LogLine, error) {
	return nil, nil
}
