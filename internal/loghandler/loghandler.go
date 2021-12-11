package loghandler

import "github.com/nicolastakashi/cole/internal/k8sclient"

func New() *LogHandler {
	dlh := DashboardLogHandler{}

	return &LogHandler{
		next: &dlh,
	}
}

type LogHandler struct {
	next Handler
}

func (lh *LogHandler) Handle(ll k8sclient.LogLine) {
	// filter only http logs
	if _, ok := ll.KeyValue["status"]; !ok {
		return
	}

	lh.next.Handle(ll)
}
