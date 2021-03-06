package loghandler

import (
	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/metrics"
	"github.com/sirupsen/logrus"
)

func New(scmd command.Server) *LogHandler {
	dlh := DashboardLogHandler{
		DashboardMetrics: metrics.DashboardMetrics{
			Scmd: scmd,
		},
	}

	return &LogHandler{
		next: &dlh,
	}
}

type LogHandler struct {
	next Handler
}

func (lh *LogHandler) Handle(ll entities.LogLine) {
	// filter only http logs
	if _, ok := ll.KeyValue["status"]; !ok {
		logrus.Debug("discarting logs without status")
		return
	}

	lh.next.Handle(ll)
}
