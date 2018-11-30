package main

import (
	"context"
	"encoding/json"

	"github.com/govinda-attal/hello-kafka/pkg/core/status"

	"github.com/govinda-attal/hello-kafka/pkg/example"
)

type GreeterHandler struct {
	srv *Greeter
}

func NewGreeterHandler(g *Greeter) *GreeterHandler {
	return &GreeterHandler{srv: g}
}

func (gh *GreeterHandler) Hello(ctx context.Context, data []byte) ([]byte, error) {
	var rq example.HelloRq
	err := json.Unmarshal(data, &rq)
	if err != nil {
		return nil, status.ErrBadRequest.WithError(err)
	}
	rs, err := gh.srv.Hello(rq)
	if err == nil {
		return json.Marshal(&rs)
	}
	return nil, err
}
