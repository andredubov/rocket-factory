package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

func (r *ordersRepository) AddOrder(ctx context.Context, order model.Order) error {
	return WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		repoOrder := converter.OrderToRepoModel(order)

		orderBuilderInsert := sq.Insert(OrdersTable).
			PlaceholderFormat(sq.Dollar).
			Columns(
				UserUUIDTableColumn,
				TotalPriceTableColumn,
				TransactionUUIDTableColumn,
				PaymentMethodTableColumn,
				StatusTableColumn,
			).
			Values(
				repoOrder.UserUUID,
				repoOrder.TotalPrice,
				repoOrder.PaymentInfo.TransactionUUID,
				repoOrder.PaymentInfo.PaymentMethod,
				repoOrder.Status,
			)

		orderQuery, orderQueryArgs, err := orderBuilderInsert.ToSql()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, orderQuery, orderQueryArgs...); err != nil {
			return err
		}

		for _, partUUID := range order.PartUUIDs {
			orderPartsBuilderInsert := sq.Insert(OrderPartsTable).
				PlaceholderFormat(sq.Dollar).
				Columns(OrderUUIDTableColumn, PartUUIDTableColumn).
				Values(repoOrder.OrderUUID, partUUID)

			orderPartsQuery, orderPartsQueryArgs, err := orderPartsBuilderInsert.ToSql()
			if err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, orderPartsQuery, orderPartsQueryArgs...); err != nil {
				return err
			}
		}
		return nil
	})
}
