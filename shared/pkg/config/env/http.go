package env

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/andredubov/rocket-factory/shared/pkg/config"
	"github.com/pkg/errors"
)

const (
	httpHostEnvName              = "HTTP_HOST"
	httpPortEnvName              = "HTTP_PORT"
	httpReadHeaderTimeoutEnvName = "HTTP_READ_HEADER_TIMEOUT_SEC"
)

type httpConfig struct {
	host              string
	port              string
	readHeaderTimeout time.Duration
}

// NewHTTPConfig returns an instance of httpConfig struct
func NewHTTPConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("%s", "http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "http port not found")
	}

	readHeaderTimeoutStr := os.Getenv(httpReadHeaderTimeoutEnvName)
	if len(readHeaderTimeoutStr) == 0 {
		return nil, errors.New("read http header timeout not found")
	}

	readHeaderTimeout, err := strconv.ParseInt(readHeaderTimeoutStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse read http header timeout")
	}

	return &httpConfig{
		host:              host,
		port:              port,
		readHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
	}, nil
}

// Address returns grpc server address
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// ReadHeaderTimeout returns http read header timeout
func (cfg *httpConfig) ReadHeaderTimeout() time.Duration {
	return cfg.readHeaderTimeout
}
