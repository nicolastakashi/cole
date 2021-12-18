package logging_parse

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/nicolastakashi/cole/internal/entities"
)

type JsonLoggingParse struct {
}

func (JsonLoggingParse) Parse(stream io.ReadCloser) ([]entities.LogLine, error) {
	logLineNumber := 1
	loglines := []entities.LogLine{}

	for in := bufio.NewScanner(stream); in.Scan(); {
		logLine := entities.LogLine{
			LineNumber: logLineNumber,
			KeyValue:   make(map[string]interface{}),
		}
		json.Unmarshal(in.Bytes(), &logLine.KeyValue)

		if len(logLine.KeyValue) == 0 {
			continue
		}

		loglines = append(loglines, logLine)
		logLineNumber++
	}

	return loglines, nil
}
