package cole

import (
	"context"
	"time"

	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, scmd command.Server) error {
	client, err := k8sclient.New(scmd.KubeConfig)

	if err != nil {
		return err
	}
	t := time.NewTimer(1 * time.Millisecond)
	lastSinceTime := time.Now()

	for {
		select {
		case <-t.C:
			if err := start(ctx, scmd, client, &lastSinceTime); err != nil {
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

func start(ctx context.Context, scmd command.Server, client *k8sclient.K8sClient, lastSinceTime *time.Time) error {
	pods, err := client.ListPods(ctx, scmd.Namespace, scmd.LabelSelector)

	if err != nil {
		return err
	}

	for _, pod := range pods {
		client.GetPodLogs(ctx, scmd.Namespace, pod, *lastSinceTime)
		*lastSinceTime = time.Now()
	}
	return nil
}
