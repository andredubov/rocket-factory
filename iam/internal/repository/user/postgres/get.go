package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/andredubov/rocket-factory/iam/internal/model"
)

func (u *usersRepository) Get(ctx context.Context, filter model.Filter) (*model.User, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	// Получаем основную информацию о пользователе
	user, err := u.getUser(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	// Получаем информацию о методах уведомления пользователя
	user.Info.NotificationMethods, err = u.getUserProviderMethods(ctx, tx, user.UUID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true

	return user, nil
}

func (u *usersRepository) getUser(ctx context.Context, tx pgx.Tx, filter model.Filter) (*model.User, error) {
	if filter.UserUUID == nil && filter.UserLogin == nil {
		return nil, model.ErrInvalidUserFilter
	}

	userQueryBuilder := sq.Select(
		columnUUID,
		columnLogin,
		columnEmail,
		columnPassword,
		columnCreatedAt,
		columnUpdatedAt,
	).
		From(tableUsers).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	if filter.UserUUID != nil {
		userUUID, err := uuid.Parse(*filter.UserUUID)
		if err != nil {
			return nil, model.ErrInvalidUserUUID
		}
		userQueryBuilder = userQueryBuilder.Where(sq.Eq{columnUUID: userUUID})
	}

	if filter.UserLogin != nil {
		userQueryBuilder = userQueryBuilder.Where(sq.Eq{columnLogin: *filter.UserLogin})
	}

	query, args, err := userQueryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var (
		user      model.User
		createdAt pgtype.Timestamptz
		updatedAt pgtype.Timestamptz
	)

	err = tx.QueryRow(ctx, query, args...).Scan(
		&user.UUID,
		&user.Info.Login,
		&user.Info.Email,
		&user.Password,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to execute user query: %w", err)
	}

	if createdAt.Valid {
		user.CreatedAt = &createdAt.Time
	}

	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}

	return &user, nil
}

func (u *usersRepository) getUserProviderMethods(ctx context.Context, tx pgx.Tx, userUUID uuid.UUID) ([]*model.NotificationMethod, error) {
	methodsQueryBuilder := sq.Select(
		columnProviderName,
		columnTarget,
	).
		From(tableNotificationMethods).
		Where(sq.Eq{columnUserUUID: userUUID}).
		PlaceholderFormat(sq.Dollar)

	methods := make([]*model.NotificationMethod, 0)

	query, args, err := methodsQueryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute notification methods query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var method model.NotificationMethod
		err = rows.Scan(
			&method.ProviderName,
			&method.Target,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification methods: %w", err)
		}
		methods = append(methods, &method)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing results: %w", err)
	}

	return methods, nil
}
