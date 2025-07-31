package tests

import (
	"context"

	"github.com/dvln/testify/mock"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestCreateOrder_Success tests the successful order creation scenario.
func (s *OrdersServiceSuite) TestCreateOrder_Success() {
	var (
		ctx                = context.Background()
		userUUID           = uuid.New()
		partUUIDs          = []uuid.UUID{uuid.New(), uuid.New()}
		expectedTotalPrice = 300.0
		order              = model.Order{
			UserUUID:  userUUID,
			PartUUIDs: partUUIDs,
		}
		parts = []model.Part{
			{Uuid: partUUIDs[0].String(), Price: 100},
			{Uuid: partUUIDs[1].String(), Price: 200},
		}
		partFilter = createPartFilter(partUUIDs)
	)

	// Mock expectations
	s.inventoryClient.On("ListParts", ctx, partFilter).Return(parts, nil)
	s.ordersRepository.On("AddOrder", ctx, mock.Anything).Return(nil)

	// Test
	err := s.ordersService.CreateOrder(ctx, order)

	// Verify
	s.NoError(err)
	s.inventoryClient.AssertExpectations(s.T())
	s.ordersRepository.AssertExpectations(s.T())

	// Проверка аргументов AddOrder
	if len(s.ordersRepository.Calls) > 0 {
		args := s.ordersRepository.Calls[0].Arguments
		addedOrder := args.Get(1).(model.Order)
		s.Equal(userUUID, addedOrder.UserUUID)
		s.Equal(partUUIDs, addedOrder.PartUUIDs)
		s.Equal(model.OrderStatusPending, addedOrder.Status)
		s.Equal(expectedTotalPrice, addedOrder.TotalPrice)
		s.NotEqual(uuid.Nil, addedOrder.OrderUUID)
	}
}

// TestCreateOrder_EmptyParts tests order creation with empty parts list.
func (s *OrdersServiceSuite) TestCreateOrder_EmptyParts() {
	var (
		ctx   = context.Background()
		order = model.Order{
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{},
		}
	)

	// Test
	err := s.ordersService.CreateOrder(ctx, order)

	// Verify
	s.Error(err)
	s.Equal(err, model.ErrOrderHasNoParts)
	s.ordersRepository.AssertNotCalled(s.T(), "AddOrder")
}

// TestCreateOrder_InventoryClientError tests order creation when inventory service fails.
func (s *OrdersServiceSuite) TestCreateOrder_InventoryClientError() {
	var (
		ctx       = context.Background()
		partUUIDs = []uuid.UUID{uuid.New(), uuid.New()}
		order     = model.Order{
			UserUUID:  uuid.New(),
			PartUUIDs: partUUIDs,
		}
		parts      = []model.Part{}
		partFilter = createPartFilter(partUUIDs)
	)

	// Mock expectations
	s.inventoryClient.On("ListParts", ctx, partFilter).Return(parts, model.ErrInvalidPartFilter)

	// Test
	err := s.ordersService.CreateOrder(ctx, order)

	// Verify
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidPartFilter)
	s.inventoryClient.AssertExpectations(s.T())
	s.ordersRepository.AssertNotCalled(s.T(), "AddOrder")
}

// TestCreateOrder_RepositoryError tests order creation when repository fails.
func (s *OrdersServiceSuite) TestCreateOrder_RepositoryError() {
	var (
		ctx                = context.Background()
		userUUID           = uuid.New()
		partUUIDs          = []uuid.UUID{uuid.New(), uuid.New()}
		expectedTotalPrice = 300.0
		order              = model.Order{
			UserUUID:  userUUID,
			PartUUIDs: partUUIDs,
		}
		parts = []model.Part{
			{Uuid: partUUIDs[0].String(), Price: 100},
			{Uuid: partUUIDs[1].String(), Price: 200},
		}
		partFilter = createPartFilter(partUUIDs)
	)

	// Mock expectations
	s.inventoryClient.On("ListParts", ctx, partFilter).Return(parts, nil)
	s.ordersRepository.On("AddOrder", ctx, mock.Anything).Return(model.ErrOrderAlreadyExists)

	// Test
	err := s.ordersService.CreateOrder(ctx, order)

	// Verify
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderAlreadyExists)
	s.inventoryClient.AssertExpectations(s.T())
	s.ordersRepository.AssertExpectations(s.T())

	// Проверка аргументов, с которыми был вызван AddOrder
	if len(s.ordersRepository.Calls) > 0 {
		args := s.ordersRepository.Calls[0].Arguments
		addedOrder := args.Get(1).(model.Order)
		s.Equal(userUUID, addedOrder.UserUUID)
		s.Equal(partUUIDs, addedOrder.PartUUIDs)
		s.Equal(model.OrderStatusPending, addedOrder.Status)
		s.Equal(expectedTotalPrice, addedOrder.TotalPrice)
		s.NotEqual(uuid.Nil, addedOrder.OrderUUID)
	}
}

// createPartFilter creates a PartFilter from given UUIDs.
func createPartFilter(partUUIDs []uuid.UUID) model.PartFilter {
	filter := model.PartFilter{}
	for _, uuid := range partUUIDs {
		filter.UUIDs = append(filter.UUIDs, uuid.String())
	}
	return filter
}
