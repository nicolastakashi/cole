package cole_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicolastakashi/cole/internal/cole"
	"github.com/nicolastakashi/cole/internal/command"
	"github.com/nicolastakashi/cole/internal/entities"
	"github.com/nicolastakashi/cole/internal/grafana"
	"github.com/nicolastakashi/cole/internal/k8sclient"
	"github.com/nicolastakashi/cole/internal/loghandler"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
	clienttesting "k8s.io/client-go/testing"
)

func buildCole(clientSet kubernetes.Interface) *cole.Cole {
	ctx := context.TODO()
	lastSinceTime := time.Now()
	scmd := command.Server{
		Namespace:     "grafana",
		LabelSelector: "name=grafana",
	}
	return &cole.Cole{
		Ctx:           ctx,
		Scmd:          scmd,
		LastSinceTime: &lastSinceTime,
		LogHandler:    loghandler.New(scmd),
		Timer:         time.NewTimer(1 * time.Millisecond),
		GrafanaConfig: grafana.GrafanaConfig{
			GrafanaApiPoolTime: time.NewTimer(1 * time.Millisecond),
		},
		Client: k8sclient.KClient{
			ClientSet: clientSet,
			Ctx:       ctx,
		},
		Out: make(chan bool, 1),
	}
}

func TestStartHandleListPodError(t *testing.T) {
	clientSet := testclient.NewSimpleClientset(&v1.PodList{})

	clientSet.PrependReactor("list", "pods", func(clienttesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.PodList{}, errors.New("error to list pods")
	})

	cole := buildCole(clientSet)

	err := cole.Start()

	assert.NotNil(t, err)
}

func TestEnsureUpdateLastSinceTime(t *testing.T) {
	podList := &v1.PodList{
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
		},
	}
	clientSet := testclient.NewSimpleClientset(podList)

	clientSet.PrependReactor("list", "pods", func(clienttesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, podList, nil
	})

	cole := buildCole(clientSet)
	lastSinceTime := cole.LastSinceTime

	go cole.Start()

	<-cole.Out
	cole.Ctx.Done()

	assert.True(t, cole.LastSinceTime.After(*lastSinceTime))
}

func TestEnsureLogHandlerIsCalled(t *testing.T) {
	podList := &v1.PodList{
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
		},
	}
	clientSet := testclient.NewSimpleClientset(podList)

	clientSet.PrependReactor("list", "pods", func(clienttesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, podList, nil
	})

	cole := buildCole(clientSet)
	flh := &FakeLogHandler{
		Called: false,
	}

	cole.LogHandler = flh

	go cole.Start()

	<-cole.Out
	cole.Ctx.Done()

	assert.True(t, flh.Called)
}

type FakeLogHandler struct {
	Called bool
}

func (flh *FakeLogHandler) Handle(ll entities.LogLine) {
	flh.Called = true
}
