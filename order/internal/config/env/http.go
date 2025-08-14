package env

import (
	"net"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
)

type httpEnvConfig struct {
	Host              string `env:"HTTP_HOST,required"`
	Port              string `env:"HTTP_PORT,required"`
	ReadHeaderTimeout string `env:"HTTP_READ_HEADER_TIMEOUT_SEC,required"`
}

type httpConfig struct {
	raw               httpEnvConfig
	readHeaderTimeout int64
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	readHeaderTimeout, err := strconv.ParseInt(raw.ReadHeaderTimeout, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse read http header timeout")
	}

	return &httpConfig{
		raw:               raw,
		readHeaderTimeout: readHeaderTimeout,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(
		cfg.raw.Host,
		cfg.raw.Port,
	)
}

func (cfg *httpConfig) ReadHeaderTimeout() time.Duration {
	return time.Duration(cfg.readHeaderTimeout) * time.Second
}
