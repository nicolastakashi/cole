package k8sclient

import (
	"context"
	"time"

	"github.com/go-logfmt/logfmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	client kubernetes.Clientset
}

type LogLine struct {
	LineNumber int
	KeyValue   map[string]string
}

func New(kubeConfig string) (*K8sClient, error) {
	var config *rest.Config = nil
	var err error = nil

	if kubeConfig != "" {
		logrus.Info("using kube config file")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		logrus.Info("using in cluster config")
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		logrus.Errorf("error while rest client config: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		logrus.Errorf("error while creating k8s client: %v", err)
		return nil, err
	}

	return &K8sClient{
		client: *clientset,
	}, nil
}

func (kc *K8sClient) ListPods(ctx context.Context, namespace string, labelSelector string) (map[string]v1.Pod, error) {
	pods, err := kc.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		logrus.Errorf("error to get grafana pods %v", err)
		return nil, err
	}

	mPods := map[string]v1.Pod{}

	for _, item := range pods.Items {
		mPods[item.Name] = item
	}

	return mPods, nil
}

func (kc *K8sClient) GetPodLogs(ctx context.Context, namespace string, pod v1.Pod, sinceTime time.Time) ([]LogLine, error) {
	rc := kc.client.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: sinceTime,
		},
	})
	stream, err := rc.Stream(ctx)

	if err != nil {
		return nil, err
	}

	defer stream.Close()

	d := logfmt.NewDecoder(stream)
	loglines := []LogLine{}
	logLineNumber := 1

	for d.ScanRecord() {
		logLine := LogLine{
			LineNumber: logLineNumber,
			KeyValue:   make(map[string]string),
		}

		for d.ScanKeyval() {
			logLine.KeyValue[string(d.Key())] = string(d.Value())
		}

		if len(logLine.KeyValue) == 0 {
			continue
		}

		loglines = append(loglines, logLine)
		logLineNumber++
	}
	return loglines, nil
}
