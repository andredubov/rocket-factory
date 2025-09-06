package redis

import "time"

type Option func(*Config)

func WithNetworkName(network string) Option {
	return func(c *Config) {
		c.NetworkName = network
	}
}

func WithContainerName(containerName string) Option {
	return func(c *Config) {
		c.ContainerName = containerName
	}
}

func WithImageName(image string) Option {
	return func(c *Config) {
		c.ImageName = image
	}
}

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.ConnectionTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.IdleTimeout = timeout
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(c *Config) {
		c.MaxIdle = maxIdle
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
