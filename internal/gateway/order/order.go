package order

import (
	"context"
	"fmt"
	"github.com/SoulHardd/service-courier/internal/config"
	"github.com/SoulHardd/service-courier/internal/domain"
	"github.com/SoulHardd/service-courier/internal/metrics"
	pb "github.com/SoulHardd/service-courier/proto/order"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gateway struct {
	order order
	cfg   *config.GrpcConfig
}

func NewGateway(o order, cfg *config.GrpcConfig) *Gateway {
	return &Gateway{
		order: o,
		cfg:   cfg,
	}
}

func (g *Gateway) GetOrders(ctx context.Context, cursor time.Time) ([]domain.Order, error) {
	var out *pb.GetOrdersResponse
	err := withRetry(
		ctx,
		func(ctx context.Context) error {
			var err error
			out, err = g.order.GetOrders(
				ctx,
				&pb.GetOrdersRequest{
					From: timestamppb.New(cursor),
				},
			)
			return err
		},
		g.cfg.BaseRetryDelay,
		g.cfg.MaxRetries,
		g.cfg.DelayMultiplier,
	)
	if err != nil {
		return nil, fmt.Errorf("get orders failed: %v", err)
	}
	if out == nil {
		return nil, fmt.Errorf("no orders found")
	}

	orders := make([]domain.Order, 0, len(out.Orders))
	for _, o := range out.Orders {
		var createdAt time.Time
		if o.CreatedAt != nil {
			createdAt = o.CreatedAt.AsTime()
		}
		orders = append(orders, domain.Order{
			ID:        o.Id,
			Status:    domain.OrderStatus(o.Status),
			CreatedAt: createdAt,
		})
	}
	return orders, nil
}

func (g *Gateway) GetOrderById(ctx context.Context, id string) (domain.Order, error) {
	var out *pb.GetOrderByIdResponse
	err := withRetry(
		ctx,
		func(ctx context.Context) error {
			var err error
			out, err = g.order.GetOrderById(
				ctx,
				&pb.GetOrderByIdRequest{Id: id},
			)
			return err
		},
		g.cfg.BaseRetryDelay,
		g.cfg.MaxRetries,
		g.cfg.DelayMultiplier,
	)
	if err != nil {
		return domain.Order{}, fmt.Errorf("get order failed: %v", err)
	}
	if out == nil {
		return domain.Order{}, fmt.Errorf("no order found")
	}

	var createdAt time.Time
	if out.Order.CreatedAt != nil {
		createdAt = out.Order.CreatedAt.AsTime()
	}
	return domain.Order{
		ID:        out.Order.Id,
		Status:    domain.OrderStatus(out.Order.Status),
		CreatedAt: createdAt,
	}, nil
}

func withRetry(ctx context.Context, fn func(context.Context) error, baseDelay time.Duration, maxRetries int, delayMultiplier int) error {
	var err error
	delay := baseDelay
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = fn(ctx)
		if err == nil {
			return nil
		}

		st, ok := status.FromError(err)
		if !ok {
			return err
		}

		switch st.Code() {
		case codes.ResourceExhausted, codes.Unavailable, codes.DeadlineExceeded:
			metrics.GatewayRetriesTotal.Inc()
			if attempt == maxRetries-1 {
				return err
			}
			select {
			case <-time.After(delay):
				delay *= time.Duration(delayMultiplier)
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		default:
			return err
		}
	}
	return err
}
