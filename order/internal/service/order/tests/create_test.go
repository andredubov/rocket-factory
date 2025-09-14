package tests

import (
	"context"
	"testing"

	"github.com/dvln/testify/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/metrics"
	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/service/mocks"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
)

// TestCreateOrder_Success tests the successful order creation scenario.
func TestCreateOrder_Success(t *testing.T) {
	// Инициализируем метрики
	err := metrics.Init(context.Background())
	require.NoError(t, err)

	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		userUUID               = uuid.New()
		partUUIDs              = []uuid.UUID{uuid.New(), uuid.New()}
		expectedTotalPrice     = 300.0
		order                  = model.Order{
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
	inventoryClient.On("ListParts", mock.Anything, partFilter).Return(parts, nil)
	ordersRepository.On("AddOrder", mock.Anything, mock.Anything).Return(nil)

	// Test
	err = ordersService.CreateOrder(ctx, &order)

	// Verify
	require.NoError(t, err)
	inventoryClient.AssertExpectations(t)
	ordersRepository.AssertExpectations(t)

	// Проверка аргументов AddOrder
	if len(ordersRepository.Calls) > 0 {
		args := ordersRepository.Calls[0].Arguments
		addedOrder := args.Get(1).(model.Order)
		require.Equal(t, userUUID, addedOrder.UserUUID)
		require.Equal(t, partUUIDs, addedOrder.PartUUIDs)
		require.Equal(t, model.OrderStatusPending, addedOrder.Status)
		require.Equal(t, expectedTotalPrice, addedOrder.TotalPrice)
		require.NotEqual(t, uuid.Nil, addedOrder.OrderUUID)
	}
}

// TestCreateOrder_EmptyParts tests order creation with empty parts list.
func TestCreateOrder_EmptyParts(t *testing.T) {
	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()

		order = model.Order{
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{},
		}
	)

	// Test
	err := ordersService.CreateOrder(ctx, &order)

	// Verify
	require.Error(t, err)
	require.Equal(t, err, model.ErrOrderHasNoParts)
	ordersRepository.AssertNotCalled(t, "AddOrder")
}

// TestCreateOrder_InventoryClientError tests order creation when inventory service fails.
func TestCreateOrder_InventoryClientError(t *testing.T) {
	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		partUUIDs              = []uuid.UUID{uuid.New(), uuid.New()}

		order = model.Order{
			UserUUID:  uuid.New(),
			PartUUIDs: partUUIDs,
		}
		parts      = []model.Part{}
		partFilter = createPartFilter(partUUIDs)
	)

	// Mock expectations
	inventoryClient.On("ListParts", mock.Anything, partFilter).Return(parts, model.ErrInvalidPartFilter)

	// Test
	err := ordersService.CreateOrder(ctx, &order)

	// Verify
	require.Error(t, err)
	require.ErrorIs(t, err, model.ErrInvalidPartFilter)
	inventoryClient.AssertExpectations(t)
	ordersRepository.AssertNotCalled(t, "AddOrder")
}

// TestCreateOrder_RepositoryError tests order creation when repository fails.
func TestCreateOrder_RepositoryError(t *testing.T) {
	// Инициализируем метрики
	err := metrics.Init(context.Background())
	require.NoError(t, err)

	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		userUUID               = uuid.New()
		partUUIDs              = []uuid.UUID{uuid.New(), uuid.New()}
		expectedTotalPrice     = 300.0

		order = model.Order{
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
	inventoryClient.On("ListParts", mock.Anything, partFilter).Return(parts, nil)
	ordersRepository.On("AddOrder", mock.Anything, mock.Anything).Return(model.ErrOrderAlreadyExists)

	// Test
	err = ordersService.CreateOrder(ctx, &order)

	// Verify
	require.Error(t, err)
	require.ErrorIs(t, err, model.ErrOrderAlreadyExists)
	inventoryClient.AssertExpectations(t)
	ordersRepository.AssertExpectations(t)

	// Проверка аргументов, с которыми был вызван AddOrder
	if len(ordersRepository.Calls) > 0 {
		args := ordersRepository.Calls[0].Arguments
		addedOrder := args.Get(1).(model.Order)
		require.Equal(t, userUUID, addedOrder.UserUUID)
		require.Equal(t, partUUIDs, addedOrder.PartUUIDs)
		require.Equal(t, model.OrderStatusPending, addedOrder.Status)
		require.Equal(t, expectedTotalPrice, addedOrder.TotalPrice)
		require.NotEqual(t, uuid.Nil, addedOrder.OrderUUID)
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
