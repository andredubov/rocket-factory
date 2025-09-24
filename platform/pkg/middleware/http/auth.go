package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/andredubov/rocket-factory/platform/pkg/middleware/grpc"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
	common_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-UUID"

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = auth_v1.AuthServiceClient

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	iamClient IAMClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// client (X-Session-UUID) -> auth middleware (add session_uuid in ctx (incomming)) -> order api (outgoing)-> auth interceptor ->inventory
// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем session UUID из разных заголовков
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.Whoami(r.Context(), &auth_v1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		// Также добавляем пользователя в контекст
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.User)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*common_v1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
