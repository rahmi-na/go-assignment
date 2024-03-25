package order_service

import (
	"fmt"
	"net/http"
	"second-assignment/dto"
	"second-assignment/entity"
	"second-assignment/pkg/errs"
	"second-assignment/repository/item_repository"
	"second-assignment/repository/order_repository"
)

type orderService struct {
	OrderRepo order_repository.Repository
	ItemRepo  item_repository.Repository
}

type Service interface {
	CreateOrderWithItems(newOrderRequest dto.NewOrderRequestDto) (*dto.NewOrderResponseDto, errs.Error)
	GetOrders() (*dto.GetOrdersResponseDto, errs.Error)
	UpdateOrder(orderId int, newOrderRequest dto.NewOrderRequestDto) (*dto.NewOrderResponseDto, errs.Error)
	DeleteOrder(orderId int) (*dto.BaseResponse, errs.Error)
}

func NewService(orderRepo order_repository.Repository, itemRepo item_repository.Repository) Service {
	return &orderService{
		OrderRepo: orderRepo,
		ItemRepo:  itemRepo,
	}
}

func (os *orderService) UpdateOrder(orderId int, newOrderRequest dto.NewOrderRequestDto) (*dto.NewOrderResponseDto, errs.Error) {

	_, err := os.OrderRepo.ReadOrderById(orderId)

	if err != nil {
		return nil, err
	}

	itemCodes := []string{}

	for _, eachItem := range newOrderRequest.Items {
		itemCodes = append(itemCodes, eachItem.ItemCode)
	}

	items, err := os.ItemRepo.GetItemsByCodes(itemCodes)

	if err != nil {
		return nil, err
	}

	for _, eachItemFromRequest := range newOrderRequest.Items {
		isFound := false

		for _, eachItem := range items {

			if eachItem.OrderId != uint(orderId) {

				return nil, errs.NewBadRequest(fmt.Sprintf("item with item code %s doesn't belong to the order with id %d", eachItem.ItemCode, orderId))

			}

			if eachItemFromRequest.ItemCode == eachItem.ItemCode {
				isFound = true
				break
			}
		}

		if !isFound {
			return nil, errs.NewNotFoundError(fmt.Sprintf("item with item code %s was not found", eachItemFromRequest.ItemCode))
		}
	}

	itemPayload := []entity.Item{}

	for _, eachItemFromRequest := range newOrderRequest.Items {
		item := entity.Item{
			ItemCode:    eachItemFromRequest.ItemCode,
			Description: eachItemFromRequest.Description,
			Quantity:    eachItemFromRequest.Quantity,
		}

		itemPayload = append(itemPayload, item)

	}

	orderPayload := entity.Order{
		OrderId:      uint(orderId),
		OrderedAt:    newOrderRequest.OrderedAt,
		CustomerName: newOrderRequest.CustomerName,
	}

	err = os.OrderRepo.UpdateOrder(orderPayload, itemPayload)

	if err != nil {
		return nil, err
	}

	response := dto.NewOrderResponseDto{
		BaseResponse: dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "order successfully updated",
		},
	}

	return &response, nil
}

func (os *orderService) GetOrders() (*dto.GetOrdersResponseDto, errs.Error) {
	orders, err := os.OrderRepo.ReadOrders()

	if err != nil {
		return nil, err
	}

	orderResult := []dto.OrderWithItems{}

	for _, eachOrder := range orders {
		order := dto.OrderWithItems{
			OrderId:      int(eachOrder.Order.OrderId),
			CustomerName: eachOrder.Order.CustomerName,
			OrderedAt:    eachOrder.Order.OrderedAt,
			CreatedAt:    eachOrder.Order.CreatedAt,
			UpdatedAt:    eachOrder.Order.UpdatedAt,
			Items:        []dto.GetItemResponseDto{},
		}

		for _, eachItem := range eachOrder.Items {
			item := dto.GetItemResponseDto{
				ItemId:      int(eachItem.ItemId),
				ItemCode:    eachItem.ItemCode,
				Quantity:    eachItem.Quantity,
				Description: eachItem.Description,
				OrderId:     int(eachItem.OrderId),
				CreatedAt:   eachItem.CreatedAt,
				UpdatedAt:   eachItem.UpdatedAt,
			}

			order.Items = append(order.Items, item)
		}

		orderResult = append(orderResult, order)

	}

	response := dto.GetOrdersResponseDto{
		BaseResponse: dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "orders successfully fetched",
		},
		Data: orderResult,
	}

	return &response, nil
}

func (os *orderService) CreateOrderWithItems(newOrderRequest dto.NewOrderRequestDto) (*dto.NewOrderResponseDto, errs.Error) {

	orderPayload := entity.Order{
		OrderedAt:    newOrderRequest.OrderedAt,
		CustomerName: newOrderRequest.CustomerName,
	}

	itemPayload := []entity.Item{}

	for _, eachItem := range newOrderRequest.Items {
		item := entity.Item{
			ItemCode:    eachItem.ItemCode,
			Description: eachItem.Description,
			Quantity:    eachItem.Quantity,
		}

		itemPayload = append(itemPayload, item)

	}

	err := os.OrderRepo.CreateOrderWithItems(orderPayload, itemPayload)

	if err != nil {
		return nil, err
	}

	response := dto.NewOrderResponseDto{
		BaseResponse: dto.BaseResponse{
			StatusCode: http.StatusCreated,
			Message:    "new order successfully created",
		},
	}

	return &response, nil

}

func (os *orderService) DeleteOrder(orderId int) (*dto.BaseResponse, errs.Error) {
	_, err := os.OrderRepo.ReadOrderById(orderId)
	if err != nil {
		return nil, err
	}

	err = os.OrderRepo.DeleteOrderById(orderId)
	if err != nil {
		return nil, err
	}

	response := &dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("Order with ID %d deleted successfully", orderId),
	}

	return response, nil
}
