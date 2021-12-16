package loghandler_test

import (
	"testing"

	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/loghandler"
	"github.com/stretchr/testify/assert"
)

type FakeDashboardMetrics struct {
	Collected bool
}

func (fm *FakeDashboardMetrics) Collect(logLine entities.LogLine) {
	fm.Collected = true
}

func TestDoNotHandleLogWithoutPath(t *testing.T) {
	dashboardMetrics := FakeDashboardMetrics{
		Collected: false,
	}

	dlh := loghandler.DashboardLogHandler{
		DashboardMetrics: &dashboardMetrics,
	}

	dlh.Handle(entities.LogLine{
		LineNumber: 1,
		KeyValue: map[string]string{
			"cenas": "xablau",
		},
	})

	assert.False(t, dashboardMetrics.Collected)
}

func TestDoNotHandleLogThatPathDoesNotMatch(t *testing.T) {
	dashboardMetrics := FakeDashboardMetrics{
		Collected: false,
	}

	dlh := loghandler.DashboardLogHandler{
		DashboardMetrics: &dashboardMetrics,
	}

	dlh.Handle(entities.LogLine{
		LineNumber: 1,
		KeyValue: map[string]string{
			"path": "/bananas",
		},
	})

	assert.False(t, dashboardMetrics.Collected)
}

func TestDoNotCollectMetricFromLogThatIsNot200(t *testing.T) {
	dashboardMetrics := FakeDashboardMetrics{
		Collected: false,
	}

	dlh := loghandler.DashboardLogHandler{
		DashboardMetrics: &dashboardMetrics,
	}

	dlh.Handle(entities.LogLine{
		LineNumber: 1,
		KeyValue: map[string]string{
			"path":   "/api/dashboards/uid/dashboard_uid",
			"status": "500",
		},
	})

	assert.False(t, dashboardMetrics.Collected)
}

func TestCollectMetrics(t *testing.T) {
	dashboardMetrics := FakeDashboardMetrics{
		Collected: false,
	}

	dlh := loghandler.DashboardLogHandler{
		DashboardMetrics: &dashboardMetrics,
	}

	dlh.Handle(entities.LogLine{
		LineNumber: 1,
		KeyValue: map[string]string{
			"path":   "/api/dashboards/uid/dashboard_uid",
			"status": "200",
		},
	})

	assert.True(t, dashboardMetrics.Collected)
}
