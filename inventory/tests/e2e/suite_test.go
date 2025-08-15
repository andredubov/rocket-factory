package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/testcontainers/path"
)

const (
	testsTimeout         = 5 * time.Minute
	inventoryEnvFilePath = "/deploy/compose/inventory/.env"
)

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "inventory Service Integration Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Загружаем .env файл и устанавливаем переменные в окружение
	envVars, err := godotenv.Read(path.GetProjectRoot() + inventoryEnvFilePath)
	if err != nil {
		logger.Fatal(suiteCtx, "Не удалось загрузить .env файл", zap.Error(err))
	}

	// Устанавливаем переменные в окружение процесса
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

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
