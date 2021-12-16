package loghandler

import (
	"regexp"

	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/metrics"
)

var regexDashboardPath = regexp.MustCompile(`\/api\/dashboards\/uid\/.+`)

type DashboardLogHandler struct {
	DashboardMetrics metrics.MetricCollector
}

func (dlh *DashboardLogHandler) Handle(logLine entities.LogLine) {
	path := logLine.KeyValue["path"]

	if !regexDashboardPath.MatchString(path) {
		return
	}

	if logLine.KeyValue["status"] == "200" {
		dlh.DashboardMetrics.Collect(logLine)
	}
}
