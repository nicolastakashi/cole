package loghandler

import (
	"github.com/nicolastakashi/cole/internal/entities"
)

func New() *LogHandler {
	dlh := DashboardLogHandler{}

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
		return
	}

	lh.next.Handle(ll)
}
