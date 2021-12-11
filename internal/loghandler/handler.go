package loghandler

import (
	"github.com/nicolastakashi/cole/internal/entities"
)

type Handler interface {
	Handle(ll entities.LogLine)
}
