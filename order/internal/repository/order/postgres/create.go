package postgres

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

func (r *ordersRepository) AddOrder(ctx context.Context, order model.Order) error {
	repoOrder := converter.OrderToRepoModel(order)
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	orderBuilderInsert := sq.Insert(ordersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(
			userUUIDTableColumn,
			totalPriceTableColumn,
			transactionUUIDTableColumn,
			paymentMethodTableColumn,
			statusTableColumn,
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
		orderPartsBuilderInsert := sq.Insert(orderPartsTable).
			PlaceholderFormat(sq.Dollar).
			Columns(orderUUIDTableColumn, partUUIDTableColumn).
			Values(repoOrder.OrderUUID, partUUID)

		orderPartsQuery, orderPartsQueryArgs, err := orderPartsBuilderInsert.ToSql()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, orderPartsQuery, orderPartsQueryArgs...); err != nil {
			return err
		}
	}

	committed = true
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
