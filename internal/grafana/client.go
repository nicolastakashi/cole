package grafana

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DashboardInfo struct {
	UID           string
	IsStared      bool
	Version       float64
	SchemaVersion float64
	Timezone      string
}

type GrafanaConfig struct {
	GrafanaApiPoolTime *time.Timer
	Address            string `yaml:"address"`
	ApiKey             string `yaml:"apiKey"`
}

type GrafanaClient struct {
	Api *gapi.Client
}

var search_dashboard_or_folder_latency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: "cole",
		Name:      "search_dashboard_or_folder_duration_seconds",
		Help:      "Duration of search dashboard or folder request in seconds.",
		Buckets:   prometheus.LinearBuckets(0.01, 0.05, 10),
	},
)

var get_dashboard_latency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: "cole",
		Name:      "get_dashboard_duration_seconds",
		Help:      "Get dashboard request duration in seconds.",
		Buckets:   prometheus.LinearBuckets(0.01, 0.05, 10),
	},
)

var search_dashboard_or_folder_error_total = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "cole",
	Name:      "search_dashboard_or_folder_error_total",
})

var get_dashboard_error_total = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "cole",
	Name:      "get_dashboard_error_total",
})

func (gc *GrafanaConfig) ReadConfigFile(grafanaApiConfigFile string) error {
	if grafanaApiConfigFile != "" {
		file, err := os.ReadFile(grafanaApiConfigFile)
		if err != nil {
			logrus.Error("error to read grafana api config file")
			return err
		}

		err = yaml.Unmarshal(file, &gc)
		if err != nil {
			logrus.Error("error to parse grafana api config file")
			return err
		}
	}
	return nil
}

func (config GrafanaConfig) NewClient() (GrafanaClient, error) {
	client, err := gapi.New(config.Address, gapi.Config{
		APIKey: config.ApiKey,
	})

	if err != nil {
		return GrafanaClient{}, err
	}

	return GrafanaClient{
		Api: client,
	}, nil
}

func (gc GrafanaClient) GetDashboardsInfo() ([]DashboardInfo, error) {
	start := time.Now()
	dashboards, err := gc.Api.FolderDashboardSearch(url.Values{
		"type": []string{"dash-db"},
	})

	elapsedSeconds := time.Since(start).Seconds()
	search_dashboard_or_folder_latency.Observe(elapsedSeconds)

	if err != nil {
		search_dashboard_or_folder_error_total.Inc()
		logrus.Error(err)
		return nil, err
	}

	dashboardsInfos := []DashboardInfo{}

	for _, dashboardSearchResponse := range dashboards {
		start := time.Now()

		dashboard, err := gc.Api.DashboardByUID(dashboardSearchResponse.UID)

		if err != nil {

			if strings.Contains(err.Error(), "status: 404") {
				logrus.Warn(err)
				continue
			}

			get_dashboard_error_total.Inc()
			logrus.Error(fmt.Printf("%s %s", err, dashboardSearchResponse.UID))
			continue
		}

		elapsedSeconds := time.Since(start).Seconds()
		get_dashboard_latency.Observe(elapsedSeconds)

		di := DashboardInfo{
			UID:      dashboardSearchResponse.UID,
			IsStared: dashboard.Meta.IsStarred,
		}

		if version, ok := dashboard.Model["version"].(float64); ok {
			di.Version = version
		}

		if schemaVersion, ok := dashboard.Model["schemaVersion"].(float64); ok {
			di.SchemaVersion = schemaVersion
		}

		if timezone, ok := dashboard.Model["timezone"].(string); ok {
			di.Timezone = timezone
		}

		dashboardsInfos = append(dashboardsInfos, di)
	}

	return dashboardsInfos, nil
}
