package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
)

const (
	SessionCookieName   = "session-uuid"
	HeaderUserUUID      = "X-User-Uuid"
	HeaderUserLogin     = "X-User-Login"
	HeaderContentType   = "content-type"
	HeaderAuthStatus    = "X-Auth-Status"
	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"
	ContentTypeJSON     = "application/json"
	AuthStatusDenied    = "denied"
)

func (a *AuthImplementation) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	logger.Debug(ctx, "🔐 External Authorization Check called")

	sessionUUID, err := a.extractSessionUUID(req)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("❌ Session extraction failed: %v", err))
		return a.denyRequest("Missing or invalid session", 403), nil
	}

	logger.Debug(ctx, fmt.Sprintf("✅ Extracted session_uuid: %s", sessionUUID))

	// Проверим валидность сессии через Whoami
	logger.Debug(ctx, fmt.Sprintf("🔍 Calling Whoami for session: %s", sessionUUID))
	whoamiResp, err := a.Whoami(ctx, &auth_v1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("❌ Whoami failed: %v", err))
		return a.denyRequest("Invalid session", 403), nil
	}

	logger.Debug(ctx, fmt.Sprintf("✅ Whoami successful. User: %s", whoamiResp.User.Uuid))

	// Проверим создание ответа
	response := a.allowRequest(whoamiResp)
	logger.Debug(ctx, "✅ Auth check passed, allowing request")

	return response, nil
}

func (a *AuthImplementation) extractSessionUUID(req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	if cookieHeader, ok := headers[HeaderCookie]; ok && cookieHeader != "" {
		sessionUUID := a.extractSessionFromCookies(cookieHeader)
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	return "", fmt.Errorf("session uuid not found in cookies")
}

func (a *AuthImplementation) extractSessionFromCookies(cookieHeader string) string {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Add(HeaderCookie, cookieHeader)

	if cookie, err := req.Cookie(SessionCookieName); err == nil {
		var sessionUUID string
		sessionUUID, err = url.QueryUnescape(cookie.Value)
		if err != nil {
			return cookie.Value
		}

		return sessionUUID
	}

	return ""
}

func (a *AuthImplementation) allowRequest(whoamiResp *auth_v1.WhoamiResponse) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserUUID,
				Value: whoamiResp.User.Uuid,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserLogin,
				Value: whoamiResp.User.Info.Login,
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: 0},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}

func (a *AuthImplementation) denyRequest(message string, statusCode int32) *authv3.CheckResponse {
	code := codes.Unauthenticated
	if statusCode == 403 {
		code = codes.PermissionDenied
	}

	// Коды gRPC стандартизированы и всегда находятся в безопасном диапазоне для int32
	statusCodeInt32 := int32(code) // nolint:gosec

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: statusCodeInt32},
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
