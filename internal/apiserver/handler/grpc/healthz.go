// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package grpc

import (
	"context"
	"miniblog/internal/pkg/log"
	"time"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	apiv1 "miniblog/pkg/api/apiserver/v1"
)

// 健康服务检查.
func (h *Handler) Healthz(ctx context.Context, rp *emptypb.Empty) (*apiv1.HealthzResponse, error) {
	log.W(ctx).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	return &apiv1.HealthzResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
