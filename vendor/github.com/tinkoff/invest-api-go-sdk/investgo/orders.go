package investgo

import (
	"context"

	pb "github.com/tinkoff/invest-api-go-sdk/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type OrdersServiceClient struct {
	conn     *grpc.ClientConn
	config   Config
	logger   Logger
	ctx      context.Context
	pbClient pb.OrdersServiceClient
}

// PostOrder - Метод выставления биржевой заявки
func (os *OrdersServiceClient) PostOrder(req *PostOrderRequest) (*PostOrderResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.PostOrder(os.ctx, &pb.PostOrderRequest{
		Quantity:     req.Quantity,
		Price:        req.Price,
		Direction:    req.Direction,
		AccountId:    req.AccountId,
		OrderType:    req.OrderType,
		OrderId:      req.OrderId,
		InstrumentId: req.InstrumentId,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &PostOrderResponse{
		PostOrderResponse: resp,
		Header:            header,
	}, err
}

// Buy - Метод выставления поручения на покупку инструмента
func (os *OrdersServiceClient) Buy(req *PostOrderRequestShort) (*PostOrderResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.PostOrder(os.ctx, &pb.PostOrderRequest{
		Quantity:     req.Quantity,
		Price:        req.Price,
		Direction:    pb.OrderDirection_ORDER_DIRECTION_BUY,
		AccountId:    req.AccountId,
		OrderType:    req.OrderType,
		OrderId:      req.OrderId,
		InstrumentId: req.InstrumentId,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &PostOrderResponse{
		PostOrderResponse: resp,
		Header:            header,
	}, err
}

// Sell - Метод выставления поручения на продажу инструмента
func (os *OrdersServiceClient) Sell(req *PostOrderRequestShort) (*PostOrderResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.PostOrder(os.ctx, &pb.PostOrderRequest{
		Quantity:     req.Quantity,
		Price:        req.Price,
		Direction:    pb.OrderDirection_ORDER_DIRECTION_SELL,
		AccountId:    req.AccountId,
		OrderType:    req.OrderType,
		OrderId:      req.OrderId,
		InstrumentId: req.InstrumentId,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &PostOrderResponse{
		PostOrderResponse: resp,
		Header:            header,
	}, err
}

// CancelOrder - Метод отмены биржевой заявки
func (os *OrdersServiceClient) CancelOrder(accountId, orderId string) (*CancelOrderResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.CancelOrder(os.ctx, &pb.CancelOrderRequest{
		AccountId: accountId,
		OrderId:   orderId,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &CancelOrderResponse{
		CancelOrderResponse: resp,
		Header:              header,
	}, err
}

// GetOrderState - Метод получения статуса торгового поручения
func (os *OrdersServiceClient) GetOrderState(accountId, orderId string) (*GetOrderStateResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.GetOrderState(os.ctx, &pb.GetOrderStateRequest{
		AccountId: accountId,
		OrderId:   orderId,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &GetOrderStateResponse{
		OrderState: resp,
		Header:     header,
	}, err
}

// GetOrders - Метод получения списка активных заявок по счёту
func (os *OrdersServiceClient) GetOrders(accountId string) (*GetOrdersResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.GetOrders(os.ctx, &pb.GetOrdersRequest{
		AccountId: accountId,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &GetOrdersResponse{
		GetOrdersResponse: resp,
		Header:            header,
	}, err
}

// ReplaceOrder - Метод изменения выставленной заявки
func (os *OrdersServiceClient) ReplaceOrder(req *ReplaceOrderRequest) (*PostOrderResponse, error) {
	var header, trailer metadata.MD
	resp, err := os.pbClient.ReplaceOrder(os.ctx, &pb.ReplaceOrderRequest{
		AccountId:      req.AccountId,
		OrderId:        req.OrderId,
		IdempotencyKey: req.NewOrderId,
		Quantity:       req.Quantity,
		Price:          req.Price,
		PriceType:      req.PriceType,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		header = trailer
	}
	return &PostOrderResponse{
		PostOrderResponse: resp,
		Header:            header,
	}, err
}
