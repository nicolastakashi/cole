package cole

import (
	"context"
	"time"

	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/nicolastakashi/cole/internal/loghandler"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, scmd command.Server) error {
	client, err := k8sclient.New(scmd.KubeConfig)

	if err != nil {
		return err
	}
	t := time.NewTimer(1 * time.Millisecond)
	lastSinceTime := time.Now().Add(time.Duration(-24) * time.Hour)
	logHandler := loghandler.New()

	for {
		select {
		case <-t.C:
			if err := start(ctx, scmd, client, &lastSinceTime, logHandler); err != nil {
				// syncErrorTotal.Inc()
				return err
			} else {
				logrus.Info("done")
				// syncSuccessTotal.Inc()
				// lastSuccessfulSync.SetToCurrentTime()
			}
			t.Reset(30 * time.Second)
		case <-ctx.Done():
			logrus.Info("shut down cole")
			return nil
		}
	}
}

func start(ctx context.Context, scmd command.Server, client *k8sclient.K8sClient, lastSinceTime *time.Time, lh loghandler.LogHandler) error {
	pods, err := client.ListPods(ctx, scmd.Namespace, scmd.LabelSelector)

	if err != nil {
		logrus.Errorf("error to lost pods ", err)
		return err
	}

	logs := []k8sclient.LogLine{}

	for _, pod := range pods {
		lgs, err := client.GetPodLogs(ctx, scmd.Namespace, pod, *lastSinceTime)
		if err != nil {
			logrus.Errorf("error to get pod %v logs ", pod.Name, err)
			return err
		}
		logs = append(logs, lgs...)
	}

	*lastSinceTime = time.Now()
	for _, log := range logs {
		lh.Handle(log)
	}

	return nil
}
