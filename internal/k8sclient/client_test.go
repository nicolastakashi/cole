package k8sclient_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	testclient "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
)

func TestListPodsHandleError(t *testing.T) {
	clientSet := testclient.NewSimpleClientset(&v1.PodList{})

	clientSet.PrependReactor("list", "pods", func(clienttesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.PodList{}, errors.New("error to list pods")
	})

	client := k8sclient.KClient{
		Ctx:       context.TODO(),
		ClientSet: clientSet,
	}

	_, err := client.ListPods("", "")

	assert.NotNil(t, err)
	assert.Equal(t, "error to list pods", err.Error())
}

func TestListPodsFromNameSpace(t *testing.T) {
	clientSet := testclient.NewSimpleClientset(&v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "grafana",
					Namespace: "grafana",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "grafana",
					Namespace: "default",
				},
			},
		},
	})

	client := k8sclient.KClient{
		Ctx:       context.TODO(),
		ClientSet: clientSet,
	}

	pods, _ := client.ListPods("grafana", "")

	assert.Equal(t, 1, len(pods))
}

func TestListPodsWithSelector(t *testing.T) {
	clientSet := testclient.NewSimpleClientset(&v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "grafana",
					Namespace: "grafana",
					Labels: map[string]string{
						"name": "grafana",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gitana",
					Namespace: "grafana",
				},
			},
		},
	})

	client := k8sclient.KClient{
		Ctx:       context.TODO(),
		ClientSet: clientSet,
	}

	pods, _ := client.ListPods("grafana", "name=grafana")

	assert.Equal(t, 1, len(pods))
}

func TestGetPodLogsSuccess(t *testing.T) {
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana",
			Namespace: "grafana",
		},
	}

	clientSet := testclient.NewSimpleClientset(&pod)

	client := k8sclient.KClient{
		Ctx:       context.TODO(),
		ClientSet: clientSet,
	}

	logs, _ := client.GetPodLogs("grafana", pod, time.Now())

	assert.Equal(t, 1, len(logs))
}
