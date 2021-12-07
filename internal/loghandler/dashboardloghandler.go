package loghandler

import (
	"regexp"
	"strings"

	"github.com/nicolastakashi/cole/internal/k8sclient"
)

var regexDashboardPath = regexp.MustCompile(`\/api\/dashboards\/uid\/.+`)

type DashboardLogHandler struct {
	next handler
}

func (dlh *DashboardLogHandler) Handle(ll k8sclient.LogLine) {
	path := ll.KeyValue["path"]

	if !regexDashboardPath.MatchString(path) {
		return
	}

	splitedPath := strings.Split(path, "/")

	if len(splitedPath) < 4 {
		return
	}

	dashboardUid := splitedPath[4]
	orgId := ll.KeyValue["orgId"]
	uid := ll.KeyValue["userId"]
	uname := ll.KeyValue["uname"]

	print(dashboardUid, uid, uname, orgId)

	dlh.next.Handle(ll)
}
