package grpc

import apiv1 "miniblog/pkg/api/apiserver/v1"

type Handler struct {
	apiv1.UnimplementedMiniBlogServer
}

func NewHandler() *Handler {
	return &Handler{}
}
