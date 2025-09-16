//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

const (
	testsTimeout = 5 * time.Minute
)

var (
	env         *TestEnvironment
	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "inventory Service Integration Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	err := logger.Init(context.Background(), logger.Config{
		Level:      "info",
		AsJSON:     true,
		EnableOTLP: false,
	})
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}
	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)
	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = ginkgo.AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
