package grpc

import (
	"context"

	"time"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	apiv1 "miniblog/pkg/api/apiserver/v1"
)

// 健康服务检查
func (h *Handler) Healthz(ctx context.Context, rp *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
