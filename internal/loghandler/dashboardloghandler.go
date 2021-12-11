package loghandler

import (
	"regexp"

	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var regexDashboardPath = regexp.MustCompile(`\/api\/dashboards\/uid\/.+`)

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

type DashboardLogHandler struct {
}

func (dlh *DashboardLogHandler) Handle(ll k8sclient.LogLine) {
	path := ll.KeyValue["path"]

	if !regexDashboardPath.MatchString(path) {
		return
	}

	dl := entities.NewDashboardLog(ll)

	if ll.KeyValue["status"] == "200" {
		dashboardViewTotal.WithLabelValues(dl.DashboardUid, dl.OrgId, dl.UserId, dl.UserName).Inc()
		dashboardLastView.WithLabelValues(dl.DashboardUid, dl.OrgId).SetToCurrentTime()
	}
}
