package telegram

import (
	"bytes"
	"context"
	"embed"
	"strings"
	"text/template"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/notification/internal/client/http"
	"github.com/andredubov/rocket-factory/notification/internal/config"
	"github.com/andredubov/rocket-factory/notification/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

//go:embed templates/order_paid_notification.tmpl
var orderPaidEventTemplateFS embed.FS

//go:embed templates/order_assembled_notification.tmpl
var orderAssembledEventTemplateFS embed.FS

type (
	orderPaidTemplateData struct {
		UUID            string
		OrderUUID       string
		UserUUID        string
		PaymentMethod   string
		TransactionUUID string
	}

	orderAssembledTemplateData struct {
		UUID         string
		OrderUUID    string
		UserUUID     string
		BuildTimeSec int64
	}
)

var (
	orderPaidEventTemplate      = template.Must(template.ParseFS(orderPaidEventTemplateFS, "templates/order_paid_notification.tmpl"))
	orderAssembledEventTemplate = template.Must(template.ParseFS(orderAssembledEventTemplateFS, "templates/order_assembled_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

// NewService создает новый Telegram сервис
func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

// SendOrderPaidNotification отправляет уведомление об оплате заказа
func (s *service) SendOrderPaidNotification(ctx context.Context, uuid string, event model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(uuid, event)
	if err != nil {
		return err
	}

	chatID := config.AppConfig().Telegram.TelegramChatID()

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", chatID), zap.String("message", message))
	return nil
}

// SendOrderAssembledNotification отправляет уведомление о сборке заказа
func (s *service) SendOrderAssembledNotification(ctx context.Context, uuid string, event model.OrderAssembledEvent) error {
	message, err := s.buildOrderAssembledMessage(uuid, event)
	if err != nil {
		return err
	}

	chatID := config.AppConfig().Telegram.TelegramChatID()

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildOrderPaidMessage создает сообщение об оплате заказа
func (s *service) buildOrderPaidMessage(uuid string, event model.OrderPaidEvent) (string, error) {
	data := orderPaidTemplateData{
		UUID:            uuid,
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   escapeMarkdown(string(event.PaymentMethod)),
		TransactionUUID: event.TransactionUUID,
	}

	var buf bytes.Buffer
	err := orderPaidEventTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildOrderAssembledMessage создает сообщение о сборке заказа
func (s *service) buildOrderAssembledMessage(uuid string, event model.OrderAssembledEvent) (string, error) {
	data := orderAssembledTemplateData{
		UUID:         uuid,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := orderAssembledEventTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// escapeMarkdown экранирует специальные символы Markdown
func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"`", "\\`",
		"[", "\\[",
	)
	return replacer.Replace(text)
}
