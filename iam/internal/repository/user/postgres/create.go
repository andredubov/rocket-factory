package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	"github.com/andredubov/rocket-factory/iam/internal/repository/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (u *usersRepository) Create(ctx context.Context, user model.User) error {
	return WithTx(ctx, u.pool, func(tx pgx.Tx) error {
		repoUser := converter.UserToRepoModel(user)

		insertQueryBuilder := sq.Insert(tableUsers).
			PlaceholderFormat(sq.Dollar).
			Columns(
				columnUUID,
				columnLogin,
				columnEmail,
				columnPassword,
				columnCreatedAt,
				columnUpdatedAt,
			).
			Values(
				repoUser.UUID,
				repoUser.Info.Login,
				repoUser.Info.Email,
				repoUser.Password,
				repoUser.CreatedAt,
				repoUser.UpdatedAt,
			)

		query, args, err := insertQueryBuilder.ToSql()
		if err != nil {
			logger.Error(ctx, "Failed to build user insert query", zap.Error(err))
			return err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			logger.Error(ctx, "Failed to execute user insert query", zap.Error(err))
			return err
		}

		if repoUser.Info.NotificationMethods == nil {
			logger.Info(ctx, "No notification methods to insert")
			return nil
		}

		for _, notificationMethod := range repoUser.Info.NotificationMethods {
			insertQueryBuilder := sq.Insert(tableNotificationMethods).
				PlaceholderFormat(sq.Dollar).
				Columns(
					columnUserUUID,
					columnProviderName,
					columnTarget,
				).
				Values(
					repoUser.UUID,
					notificationMethod.ProviderName,
					notificationMethod.Target,
				)

			query, args, err := insertQueryBuilder.ToSql()
			if err != nil {
				logger.Error(ctx, "Failed to build notification method insert query", zap.Error(err))
				return err
			}

			_, err = tx.Exec(ctx, query, args...)
			if err != nil {
				logger.Error(ctx, "Failed to execute notification method insert query", zap.Error(err))
				return err
			}
		}

		return nil
	})
}
