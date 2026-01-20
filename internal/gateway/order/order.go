package order

import (
	"avito/internal/domain"
	pb "avito/proto/order"
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gateway struct {
	order order
}

func NewGateway(o order) *Gateway {
	return &Gateway{order: o}
}

func (g *Gateway) GetOrders(ctx context.Context, cursor time.Time) ([]domain.Order, error) {
	out, err := g.order.GetOrders(ctx, &pb.GetOrdersRequest{From: timestamppb.New(cursor)})
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
	out, err := g.order.GetOrderById(ctx, &pb.GetOrderByIdRequest{Id: id})
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
