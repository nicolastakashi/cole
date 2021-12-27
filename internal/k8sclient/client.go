package k8sclient

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	ListPods(namespace string, labelSelector string) ([]v1.Pod, error)
	GetPodLogs(namespace string, container string, pod v1.Pod, sinceTime time.Time) (*rest.Request, error)
}

type KClient struct {
	ClientSet kubernetes.Interface
	Ctx       context.Context
}

func New(ctx context.Context, kubeConfig string) (*KClient, error) {
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

	return &KClient{
		Ctx:       ctx,
		ClientSet: clientset,
	}, nil
}

func (c KClient) ListPods(namespace string, labelSelector string) ([]v1.Pod, error) {
	pods, err := c.ClientSet.CoreV1().Pods(namespace).List(c.Ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		logrus.Errorf("error to get grafana pods %v", err)
		return nil, err
	}

	return pods.Items, nil
}

func (c KClient) GetPodLogs(namespace string, container string, pod v1.Pod, sinceTime time.Time) (*rest.Request, error) {
	rc := c.ClientSet.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		Container: container,
		SinceTime: &metav1.Time{
			Time: sinceTime,
		},
	})

	return rc, nil
}
