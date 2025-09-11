package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (r *ordersRepository) AddOrder(ctx context.Context, order model.Order) error {
	return WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		repoOrder := converter.OrderToRepoModel(order)

		orderBuilderInsert := sq.Insert(OrdersTable).
			PlaceholderFormat(sq.Dollar).
			Columns(
				UUIDTableColumn,
				UserUUIDTableColumn,
				TotalPriceTableColumn,
				StatusTableColumn,
			).
			Values(
				repoOrder.OrderUUID,
				repoOrder.UserUUID,
				repoOrder.TotalPrice,
				repoOrder.Status,
			)

		if repoOrder.PaymentInfo != nil {
			orderBuilderInsert = orderBuilderInsert.
				Columns(
					TransactionUUIDTableColumn,
					PaymentMethodTableColumn,
				).
				Values(
					repoOrder.PaymentInfo.TransactionUUID,
					string(repoOrder.PaymentInfo.PaymentMethod),
				)
		}

		orderQuery, orderQueryArgs, err := orderBuilderInsert.ToSql()
		if err != nil {
			logger.Error(ctx, "Failed to build order query", zap.Error(err))
			return err
		}

		if _, err := tx.Exec(ctx, orderQuery, orderQueryArgs...); err != nil {
			logger.Error(ctx, "Failed to exect order query", zap.Error(err))
			return err
		}

		for _, partUUID := range order.PartUUIDs {
			orderPartsBuilderInsert := sq.Insert(OrderPartsTable).
				PlaceholderFormat(sq.Dollar).
				Columns(OrderUUIDTableColumn, PartUUIDTableColumn).
				Values(repoOrder.OrderUUID, partUUID)

			orderPartsQuery, orderPartsQueryArgs, err := orderPartsBuilderInsert.ToSql()
			if err != nil {
				logger.Error(ctx, "Failed to build order query", zap.Error(err))
				return err
			}

			if _, err := tx.Exec(ctx, orderPartsQuery, orderPartsQueryArgs...); err != nil {
				logger.Error(ctx, "Failed to exect order query", zap.Error(err))
				return err
			}
		}

		return nil
	})
}
