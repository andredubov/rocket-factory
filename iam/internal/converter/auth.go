package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/andredubov/rocket-factory/iam/internal/model"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
	common_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/common/v1"
)

func SessionUUIDToLoginResponse(uuid uuid.UUID) *auth_v1.LoginResponse {
	return &auth_v1.LoginResponse{
		SessionUuid: uuid.String(),
	}
}

func NotificationMethodToProto(notificationMethods []*model.NotificationMethod) []*common_v1.NotificationMethod {
	methods := make([]*common_v1.NotificationMethod, 0, len(notificationMethods))

	for _, notificationMethod := range notificationMethods {
		methods = append(methods, &common_v1.NotificationMethod{
			ProviderName: notificationMethod.ProviderName,
			Target:       notificationMethod.Target,
		})
	}

	return methods
}

func SessionToWhoamiResponse(session *model.Session) (*auth_v1.WhoamiResponse, error) {
	if session == nil {
		return nil, model.ErrSessionNotFound
	}

	if session.UUID == uuid.Nil {
		return nil, model.ErrInvalidSessionUUID
	}

	if session.User.UUID == uuid.Nil {
		return nil, model.ErrInvalidUserUUID
	}

	response := &auth_v1.WhoamiResponse{
		Session: &common_v1.Session{
			Uuid: session.UUID.String(),
		},
		User: &common_v1.User{
			Uuid: session.User.UUID.String(),
		},
	}

	createdAt := timestamppb.New(session.CreatedAt)
	if err := createdAt.CheckValid(); err != nil {
		return nil, err
	}
	response.Session.CreatedAt = createdAt

	updatedAt := timestamppb.New(session.UpdatedAt)
	if err := updatedAt.CheckValid(); err != nil {
		return nil, err
	}
	response.Session.UpdatedAt = updatedAt

	expiresAt := timestamppb.New(session.ExpiresAt)
	if err := expiresAt.CheckValid(); err != nil {
		return nil, err
	}
	response.Session.ExpiresAt = expiresAt

	if session.User.CreatedAt != nil {
		userCreatedAt := timestamppb.New(*session.User.CreatedAt)
		if err := userCreatedAt.CheckValid(); err != nil {
			return nil, err
		}
		response.User.CreatedAt = userCreatedAt
	}

	if session.User.UpdatedAt != nil {
		userUpdatedAt := timestamppb.New(*session.User.UpdatedAt)
		if err := userUpdatedAt.CheckValid(); err != nil {
			return nil, err
		}
		response.User.UpdatedAt = userUpdatedAt
	}

	response.User.Info = &common_v1.UserInfo{
		Login:               session.User.Info.Login,
		Email:               session.User.Info.Email,
		NotificationMethods: NotificationMethodToProto(session.User.Info.NotificationMethods),
	}

	return response, nil
}
