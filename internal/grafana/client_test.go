package grafana_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/nicolastakashi/cole/internal/grafana"
	"github.com/stretchr/testify/assert"
)

type mockServerCall struct {
	code int
	body string
}

type mockServer struct {
	upcomingCalls []mockServerCall
	executedCalls []mockServerCall
	server        *httptest.Server
}

func (m *mockServer) Close() {
	m.server.Close()
}

func gapiTestToolsFromCalls(t *testing.T, calls []mockServerCall) *gapi.Client {
	t.Helper()

	mock := &mockServer{
		upcomingCalls: calls,
	}

	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := mock.upcomingCalls[0]
		if len(calls) > 1 {
			mock.upcomingCalls = mock.upcomingCalls[1:]
		} else {
			mock.upcomingCalls = nil
		}
		w.WriteHeader(call.code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, call.body)
		mock.executedCalls = append(mock.executedCalls, call)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(mock.server.URL)
		},
	}

	httpClient := &http.Client{Transport: tr}

	client, err := gapi.New("http://my-grafana.com", gapi.Config{APIKey: "my-key", Client: httpClient})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		mock.Close()
	})

	return client
}

func TestGetDashboardInfoWithModelProperty(t *testing.T) {
	const getFolderDashboardSearchResponse = `[
		{
			"id":1,
			"uid": "cIBgcSjkk",
			"title":"Production Overview",
			"url": "/d/cIBgcSjkk/production-overview",
			"type":"dash-db",
			"tags":["prod"],
			"isStarred":true,
			"folderId": 2,
			"folderUid": "000000163",
			"folderTitle": "Folder",
			"folderUrl": "/dashboards/f/000000163/folder",
			"uri":"db/production-overview"
		}
	]`

	const dashboard = `
	{
		"id":1,
		"uid": "cIBgcSjkk",
		"title":"Production Overview",
		"url": "/d/cIBgcSjkk/production-overview",
		"type":"dash-db",
		"tags":["prod"],
		"isStarred":true,
		"folderId": 2,
		"folderUid": "000000163",
		"folderTitle": "Folder",
		"folderUrl": "/dashboards/f/000000163/folder",
		"uri":"db/production-overview",
		"dashboard": {
			"version": 1,
			"schemaVersion": 36,
			"timezone": "utc"
		}
	}`

	gc := grafana.GrafanaClient{
		Api: gapiTestToolsFromCalls(t, []mockServerCall{{200, getFolderDashboardSearchResponse}, {200, dashboard}}),
	}

	dashboardInfos, err := gc.GetDashboardsInfo()

	assert.Nil(t, err)
	assert.Len(t, dashboardInfos, 1)
	assert.NotNil(t, dashboardInfos[0].Version)
	assert.NotNil(t, dashboardInfos[0].SchemaVersion)
	assert.NotNil(t, dashboardInfos[0].Timezone)
}

func TestGetDashboardInfoWithoutModelProperty(t *testing.T) {
	const getFolderDashboardSearchResponse = `[
		{
			"id":1,
			"uid": "cIBgcSjkk",
			"title":"Production Overview",
			"url": "/d/cIBgcSjkk/production-overview",
			"type":"dash-db",
			"tags":["prod"],
			"isStarred":true,
			"folderId": 2,
			"folderUid": "000000163",
			"folderTitle": "Folder",
			"folderUrl": "/dashboards/f/000000163/folder",
			"uri":"db/production-overview"
		}
	]`

	const dashboard = `
	{
		"id":1,
		"uid": "cIBgcSjkk",
		"title":"Production Overview",
		"url": "/d/cIBgcSjkk/production-overview",
		"type":"dash-db",
		"tags":["prod"],
		"isStarred":true,
		"folderId": 2,
		"folderUid": "000000163",
		"folderTitle": "Folder",
		"folderUrl": "/dashboards/f/000000163/folder",
		"uri":"db/production-overview"
	}`

	gc := grafana.GrafanaClient{
		Api: gapiTestToolsFromCalls(t, []mockServerCall{{200, getFolderDashboardSearchResponse}, {200, dashboard}}),
	}

	dashboardInfos, err := gc.GetDashboardsInfo()

	assert.Nil(t, err)
	assert.Len(t, dashboardInfos, 1)
	assert.Equal(t, dashboardInfos[0].Version, float64(0))
	assert.Equal(t, dashboardInfos[0].SchemaVersion, float64(0))
	assert.Equal(t, dashboardInfos[0].Timezone, "")
}
