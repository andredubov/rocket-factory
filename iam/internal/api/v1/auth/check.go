package api

import (
	"context"
	"fmt"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

const (
	HeaderContentType = "content-type"
	HeaderAuthStatus  = "x-auth-status"

	HeaderCookie             = "cookie"
	HeaderAuthorization      = "authorization"
	HeaderSessionUUIDPrimary = "session-uuid"
	HeaderSessionUUIDAlt     = "x-session-uuid"

	ContentTypeJSON = "application/json"

	AuthStatusDenied = "denied"
)

// Check реализует External Authorization для Envoy
func (a *AuthImplementation) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	logger.Info(ctx, "🔒 External Authorization Check called")

	sessionUUIDStr, err := a.extractSessionUUID(ctx, req)
	if err != nil {
		logger.Error(ctx, "❌ Session extraction failed", zap.Error(err))
		return a.denyRequest("Missing or invalid session", 401), nil
	}

	logger.Info(ctx, "🔑 Extracted session-uuid", zap.String("value", sessionUUIDStr))

	sessionUUID, err := uuid.Parse(sessionUUIDStr)
	if err != nil {
		logger.Error(ctx, "❌ Invalid session UUID format", zap.Error(err))
		return a.denyRequest("Invalid session format", 401), nil
	}

	sessionUser, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		logger.Error(ctx, "❌ Whoami failed", zap.Error(err))
		return a.denyRequest("Invalid session", 403), nil
	}

	logger.Info(ctx, "✅ Authorization successful for user", zap.String("user's login", sessionUser.User.Info.Login))

	return a.allowRequest(sessionUUID), nil
}

// extractSessionUUID извлекает session UUID из различных источников в запросе
func (a *AuthImplementation) extractSessionUUID(ctx context.Context, req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	if sessionHeader, ok := headers[HeaderSessionUUIDAlt]; ok && sessionHeader != "" {
		logger.Info(ctx, "✅ Session header", zap.String(HeaderSessionUUIDAlt, sessionHeader))
		return sessionHeader, nil
	}

	logger.Error(ctx, "❌ Session uuid not found in any headers or cookies")

	return "", fmt.Errorf("session uuid not found in any headers or cookies")
}

// allowRequest создает разрешающий ответ с пользовательскими заголовками
func (a *AuthImplementation) allowRequest(sessionUUID uuid.UUID) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderSessionUUIDPrimary,
				Value: sessionUUID.String(),
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}

func (a *AuthImplementation) denyRequest(message string, statusCode int32) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: typev3.StatusCode(statusCode),
				},
				Body: fmt.Sprintf(`{"error": "%s", "timestamp": "%s"}`,
					message, time.Now().Format(time.RFC3339)),
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderContentType,
							Value: ContentTypeJSON,
						},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderAuthStatus,
							Value: AuthStatusDenied,
						},
					},
				},
			},
		},
	}
}
