package cole

import (
	"context"
	"time"

	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/nicolastakashi/cole/internal/k8sclient/logging_parse"
	"github.com/nicolastakashi/cole/internal/loghandler"
	"github.com/sirupsen/logrus"
)

type Cole struct {
	Ctx           context.Context
	Scmd          command.Server
	Client        k8sclient.Client
	LastSinceTime time.Time
	LogHandler    loghandler.Handler
	Timer         *time.Timer
	Out           chan bool
}

func (cole Cole) Start() error {
	for {
		select {
		case <-cole.Timer.C:
			if err := cole.run(); err != nil {
				// syncErrorTotal.Inc()
				return err
			} else {
				logrus.Info("done")
				// syncSuccessTotal.Inc()
				// lastSuccessfulSync.SetToCurrentTime()
			}
			cole.Timer.Reset(30 * time.Second)
			if cole.Out != nil {
				cole.Out <- true
			}
		case <-cole.Ctx.Done():
			logrus.Info("shut down cole")
			return nil
		}
	}
}

func (c *Cole) run() error {
	pods, err := c.Client.ListPods(c.Scmd.Namespace, c.Scmd.LabelSelector)

	if err != nil {
		logrus.Errorf("error to lost pods %v", err)
		return err
	}

	logs := []entities.LogLine{}

	for _, pod := range pods {
		lr, err := c.Client.GetPodLogs(c.Scmd.Namespace, pod, c.LastSinceTime)
		if err != nil {
			logrus.Errorf("error to get pod %v logs %v", pod.Name, err)
			return err
		}

		stream, err := lr.Stream(c.Ctx)

		if err != nil {
			return err
		}

		defer stream.Close()

		lgs, err := logging_parse.Get(c.Scmd.LogType).Parse(stream)

		if err != nil {
			return err
		}

		logs = append(logs, lgs...)
	}

	c.LastSinceTime = time.Now()

	for _, log := range logs {
		c.LogHandler.Handle(log)
	}

	return nil
}
