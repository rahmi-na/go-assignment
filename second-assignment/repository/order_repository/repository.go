package order_repository

import (
	"second-assignment/entity"
	"second-assignment/pkg/errs"
)

type Repository interface {
	ReadOrderById(orderId int) (*entity.Order, errs.Error)
	CreateOrderWithItems(orderPayload entity.Order, itemPayload []entity.Item) errs.Error
	ReadOrders() ([]OrderItemMapped, errs.Error)
	UpdateOrder(orderPayload entity.Order, itemPayload []entity.Item) errs.Error
	DeleteOrderById(orderId int) errs.Error
}
