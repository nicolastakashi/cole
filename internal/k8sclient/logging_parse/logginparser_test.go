package logging_parse_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/nicolastakashi/cole/internal/k8sclient/logging_parse"
	"github.com/stretchr/testify/assert"
)

func TestLoggingParseGetConsole(t *testing.T) {
	loggingParse := logging_parse.Get("console")
	assert.True(t, fmt.Sprintf("%T", loggingParse) == fmt.Sprintf("%T", logging_parse.ConsoleLoggingParse{}))
}

func TestLoggingParseGetJson(t *testing.T) {
	loggingParse := logging_parse.Get("json")
	assert.True(t, fmt.Sprintf("%T", loggingParse) == fmt.Sprintf("%T", logging_parse.JsonLoggingParse{}))
}

func TestConsoleParseLogging(t *testing.T) {
	in := `
		id=1 dur=1.001s
		id=1 path=/path/to/file err="file not found"
		`
	reader := io.NopCloser(strings.NewReader(in))

	logLines, _ := logging_parse.Get("console").Parse(reader)

	assert.Equal(t, len(logLines), 2)

}
