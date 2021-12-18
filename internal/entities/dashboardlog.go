package entities

import (
	"fmt"
	"strings"
)

type DashBoardLog struct {
	DashboardUid string
	OrgId        string
	UserId       string
	UserName     string
}

func NewDashboardLog(ll LogLine) DashBoardLog {
	path := fmt.Sprintf("%v", ll.KeyValue["path"])
	orgId := fmt.Sprintf("%v", ll.KeyValue["orgId"])
	duid := "unknow"
	uid := "unknow"
	uname := "unknow"

	if splitedPath := strings.Split(path, "/"); len(splitedPath) > 3 {
		duid = splitedPath[4]
	}

	if value, ok := ll.KeyValue["uname"]; ok {
		uname = fmt.Sprintf("%v", value)
	}

	if value, ok := ll.KeyValue["userId"]; ok {
		uid = fmt.Sprintf("%v", value)
	}

	return DashBoardLog{
		DashboardUid: duid,
		OrgId:        orgId,
		UserId:       uid,
		UserName:     uname,
	}
}
