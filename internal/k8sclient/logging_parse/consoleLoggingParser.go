package logging_parse

import (
	"io"

	"github.com/go-logfmt/logfmt"
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/sirupsen/logrus"
)

type ConsoleLoggingParse struct {
}

func (ConsoleLoggingParse) Parse(stream io.ReadCloser) ([]entities.LogLine, error) {
	logrus.Debug("parsin console logs")
	d := logfmt.NewDecoder(stream)
	logLineNumber := 1
	loglines := []entities.LogLine{}

	for d.ScanRecord() {
		logLine := entities.LogLine{
			LineNumber: logLineNumber,
			KeyValue:   make(map[string]interface{}),
		}

		for d.ScanKeyval() {
			logLine.KeyValue[string(d.Key())] = string(d.Value())
		}

		if len(logLine.KeyValue) == 0 {
			continue
		}

		loglines = append(loglines, logLine)
		logLineNumber++
	}
	return loglines, nil
}
