package krouter

import (
	"fmt"

	at "github.com/govinda-attal/hello-kafka/pkg/asynctrans"
	"github.com/govinda-attal/hello-kafka/pkg/core/status"
)

type RouteGroup struct {
	group    string
	handlers map[string]at.MsgHandler
}

// GetMsgHandler returns the handler for the route.
func (rg *RouteGroup) GetMsgHandler(msgName string) (at.MsgHandler, error) {
	h, ok := rg.handlers[msgName]
	if !ok {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("handler for given message %s not found", msgName))
	}
	return h, nil
}

func (rg *RouteGroup) Invoke(msgName string, handler at.MsgHandler) *RouteGroup {
	rg.handlers[msgName] = handler
	return rg
}
