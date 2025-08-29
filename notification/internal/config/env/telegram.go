package env

import (
	"github.com/caarlos0/env/v11"
)

type telegramEnvConfig struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	TelegramChatID   int64  `env:"TELEGRAM_CHAT_ID,required"`
}

type telegramConfig struct {
	raw telegramEnvConfig
}

func NewTelegramConfig() (*telegramConfig, error) {
	var raw telegramEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &telegramConfig{raw: raw}, nil
}

func (cfg *telegramConfig) TelegramBotToken() string {
	return cfg.raw.TelegramBotToken
}

func (cfg *telegramConfig) TelegramChatID() int64 {
	return cfg.raw.TelegramChatID
}
