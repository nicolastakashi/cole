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

func TestJsonParseLogging(t *testing.T) {
	in := `
	✔ Downloaded marcusolsson-json-datasource v1.3.0 zip successfully

	Please restart Grafana after installing plugins. Refer to Grafana documentation for instructions if necessary.

	{"logger":"settings","lvl":"warn","msg":"falling back to legacy setting of 'min_interval_seconds'; please use the configuration option in the unified_alerting section if Grafana 8 alerts are enabled.","t":"2021-12-18T09:20:12.278092279Z"}
	{"logger":"settings","lvl":"warn","msg":"falling back to legacy setting of 'min_interval_seconds'; please use the configuration option in the unified_alerting section if Grafana 8 alerts are enabled.","t":"2021-12-18T09:20:12.278262959Z"}
	{"file":"/usr/share/grafana/conf/defaults.ini","logger":"settings","lvl":"info","msg":"Config loaded from","t":"2021-12-18T09:20:12.278352264Z"}
	`
	reader := io.NopCloser(strings.NewReader(in))

	logLines, _ := logging_parse.Get("json").Parse(reader)

	assert.Equal(t, len(logLines), 3)
}
