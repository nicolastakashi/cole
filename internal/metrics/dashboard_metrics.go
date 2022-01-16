package metrics

import (
	"fmt"

	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/grafana"
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
	[]string{"dashboard_uid", "org_id", "user_id"},
)

var dashboardInfo = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "cole",
		Name:      "dashboard_info",
		Help:      "Metric with constant 1 value that holds information about the dashboard",
	},
	[]string{"dashboard_uid", "is_stared", "version", "schema_version", "timezone"},
)

type MetricCollector interface {
	Collect(logLine entities.LogLine)
}

type DashboardMetrics struct {
	Scmd command.Server
}

func (dm DashboardMetrics) Collect(logLine entities.LogLine) {

	dl := entities.NewDashboardLog(logLine, dm.Scmd.IncludeUname)

	dashboardViewTotal.WithLabelValues(dl.DashboardUid, dl.OrgId, dl.UserId, dl.UserName).Inc()
	dashboardLastView.WithLabelValues(dl.DashboardUid, dl.OrgId, dl.UserId).SetToCurrentTime()

	logrus.Info("collecting dashboard metrics")
}

func (dm DashboardMetrics) CollectFromGrafanaApi(dashbordsInfo []grafana.DashboardInfo) {
	for _, di := range dashbordsInfo {
		dashboardInfo.WithLabelValues(di.UID,
			fmt.Sprintf("%v", di.IsStared),
			fmt.Sprintf("%v", di.Version),
			fmt.Sprintf("%v", di.SchemaVersion),
			di.Timezone).Set(1)
	}
	logrus.Debug("dashboar_info collected")
}
