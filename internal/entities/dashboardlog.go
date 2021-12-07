package entities

import (
	"strings"

	"github.com/nicolastakashi/cole/internal/k8sclient"
)

type DashBoardLog struct {
	DashboardUid string
	OrgId        string
	UserId       string
	UserName     string
}

func NewDashboardLog(ll k8sclient.LogLine) DashBoardLog {
	path := ll.KeyValue["path"]
	orgId := ll.KeyValue["orgId"]
	duid := "unknow"
	uid := "unknow"
	uname := "unknow"

	if splitedPath := strings.Split(path, "/"); len(splitedPath) > 3 {
		duid = splitedPath[4]
	}

	if value, ok := ll.KeyValue["uname"]; ok {
		uname = value
	}

	if value, ok := ll.KeyValue["userId"]; ok {
		uid = value
	}

	return DashBoardLog{
		DashboardUid: duid,
		OrgId:        orgId,
		UserId:       uid,
		UserName:     uname,
	}
}
