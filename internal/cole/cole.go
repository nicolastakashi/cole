package cole

import (
	"context"
	"time"

	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/nicolastakashi/cole/internal/k8sclient/logging_parse"
	"github.com/nicolastakashi/cole/internal/loghandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type Cole struct {
	Ctx           context.Context
	Scmd          command.Server
	Client        k8sclient.Client
	LastSinceTime *time.Time
	LogHandler    loghandler.Handler
	Timer         *time.Timer
	Out           chan bool
}

func (c *Cole) UpdateLastSinceTime() {
	lastSinceTime := time.Now()
	c.LastSinceTime = &lastSinceTime
	logrus.Infof("updating last sinc time %v", c.LastSinceTime)
}

var lastSuccessfulSync = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "cole",
	Name:      "last_success_sync_timestamp_seconds",
	Help:      "Unix timestamp of the last successful dashboard sync in seconds",
})

var syncSuccessTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "cole",
		Name:      "sync_total_success",
		Help:      "Total number of successful sync operations",
	},
)

var syncErrorTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "cole",
		Name:      "sync_total_error",
		Help:      "Total number of sync operations with errors",
	},
)

func (cole *Cole) Start() error {
	for {
		select {
		case <-cole.Timer.C:
			if err := cole.run(); err != nil {
				syncErrorTotal.Inc()
				return err
			} else {
				logrus.Info("sixth sense updated")
				syncSuccessTotal.Inc()
				lastSuccessfulSync.SetToCurrentTime()
			}
			cole.Timer.Reset(30 * time.Second)
			if cole.Out != nil {
				cole.Out <- true
			}
		case <-cole.Ctx.Done():
			logrus.Info("cole, please stop...")
			return nil
		}
	}
}

func (c *Cole) run() error {
	pods, err := c.Client.ListPods(c.Scmd.Namespace, c.Scmd.LabelSelector)

	logrus.Debugf("listed pods %v", len(pods))
	if err != nil {
		logrus.Errorf("error to list pods %v", err)
		return err
	}

	logs := []entities.LogLine{}

	for _, pod := range pods {
		logrus.Debugf("getting logs %v", pod.Name)
		lr, err := c.Client.GetPodLogs(c.Scmd.Namespace, c.Scmd.Container, pod, *c.LastSinceTime)
		if err != nil {
			logrus.Errorf("error to get pod %v logs %v", pod.Name, err)
			return err
		}

		stream, err := lr.Stream(c.Ctx)

		if err != nil {
			logrus.Errorf("error to read logs %v", err)
			return err
		}

		defer stream.Close()

		lgs, err := logging_parse.Get(c.Scmd.LogFormat).Parse(stream)

		if err != nil {
			return err
		}

		logs = append(logs, lgs...)
		logrus.Debugf("available logs %v", len(logs))
	}

	c.UpdateLastSinceTime()

	for _, log := range logs {
		c.LogHandler.Handle(log)
	}

	logrus.Debug("i'm ready to tell you my secret now.")

	return nil
}
