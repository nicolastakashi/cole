package loghandler

import "github.com/nicolastakashi/cole/internal/k8sclient"

type Handler interface {
	Handle(ll k8sclient.LogLine)
}
