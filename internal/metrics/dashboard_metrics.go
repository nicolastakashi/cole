package metrics

import (
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var dashboardViewTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "cole",
		Name:      "dashboard_view_total",
		Help:      "Total number of views of a dashboard",
	},
	[]string{"dashboard_uid", "org_id", "user_id", "user_name"},
)

var dashboardLastView = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "cole",
		Name:      "dashboard_last_view_seconds",
		Help:      "dashboard last view in seconds",
	},
	[]string{"dashboard_uid", "org_id"},
)

type MetricCollector interface {
	Collect(logLine entities.LogLine)
}

type DashboardMetrics struct{}

func (DashboardMetrics) Collect(logLine entities.LogLine) {
	dl := entities.NewDashboardLog(logLine)

	dashboardViewTotal.WithLabelValues(dl.DashboardUid, dl.OrgId, dl.UserId, dl.UserName).Inc()
	dashboardLastView.WithLabelValues(dl.DashboardUid, dl.OrgId).SetToCurrentTime()

	logrus.Info("collecting dashboard metrics")
}
