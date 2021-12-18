package loghandler

import (
	"fmt"
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

	if !regexDashboardPath.MatchString(fmt.Sprintf("%v", path)) {
		return
	}

	if status := logLine.KeyValue["status"]; fmt.Sprintf("%v", status) == "200" {
		dlh.DashboardMetrics.Collect(logLine)
	}
}
