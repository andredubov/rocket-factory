//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	middlewaregrpc "github.com/andredubov/rocket-factory/platform/pkg/middleware/grpc"
	auth_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/auth/v1"
	common_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/common/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	user_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/user/v1"
)

var _ = ginkgo.Describe("InventoryService", func() {
	var (
		ctx         context.Context
		cancel      context.CancelFunc
		clients     *TestClients
		sessionUUID string
	)

	ginkgo.BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		var err error
		clients, err = NewTestClients(ctx, env)
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешное подключение к сервисам: IAM, Inventory.")

		username := customUsername(6)
		password := generateSimplePassword(6)

		userResponse, err := clients.UserClient.Register(ctx, &user_v1.RegisterRequest{
			Info: &common_v1.UserInfo{
				Login: username,
				Email: gofakeit.Email(),
				NotificationMethods: []*common_v1.NotificationMethod{
					{
						ProviderName: "telegram",
						Target:       "@username",
					},
				},
			},
			Password: password,
		})
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешнyю регистрацию пользователя в IAM сервисе")
		gomega.Expect(userResponse.GetUserUuid()).ToNot(gomega.BeEmpty())

		authResponse, err := clients.AuthClient.Login(ctx, &auth_v1.LoginRequest{
			Login:    username,
			Password: password,
		})
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешную аутентификацию пользователя в IAM сервисе")
		gomega.Expect(authResponse.GetSessionUuid()).ToNot(gomega.BeEmpty())

		sessionUUID = authResponse.GetSessionUuid()
	})

	ginkgo.AfterEach(func() {
		err := clients.CloseWithContext(ctx)
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешное закрытие клиентов сервисов IAM, Inventory")
		cancel()
	})

	ginkgo.Describe("Test GetPart", func() {
		ginkgo.BeforeEach(func() {
			ctx, cancel = context.WithCancel(suiteCtx)
			err := env.ClearPartsCollection(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			ctx = middlewaregrpc.AddSessionUUIDToContext(ctx, sessionUUID)
			ctx = middlewaregrpc.ForwardSessionUUIDToGRPC(ctx)
		})

		ginkgo.AfterEach(func() {
			cancel()
		})

		ginkgo.It("Должен возвращать деталь по UUID", func() {
			// Setup
			partUUID, err := env.InsertTestPart(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешное создание тестовой детали")
			gomega.Expect(partUUID).ToNot(gomega.BeEmpty())

			// Test
			resp, err := clients.InventoryClient.GetPart(ctx, &inventory_v1.GetPartRequest{Uuid: partUUID})

			// Verify
			gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешное получение детали по UUID")
			gomega.Expect(resp.GetPart()).ToNot(gomega.BeNil())
			gomega.Expect(resp.GetPart().GetUuid()).To(gomega.Equal(partUUID))
			gomega.Expect(resp.GetPart().GetName()).ToNot(gomega.BeEmpty())
			gomega.Expect(resp.GetPart().GetDescription()).ToNot(gomega.BeEmpty())
			gomega.Expect(resp.GetPart().GetPrice()).To(gomega.BeNumerically(">=", 0))
			gomega.Expect(resp.GetPart().GetStockQuantity()).To(gomega.BeNumerically(">=", 0))
			gomega.Expect(resp.GetPart().GetCategory()).Should(gomega.BeElementOf([]inventory_v1.Category{
				inventory_v1.Category_CATEGORY_ENGINE,
				inventory_v1.Category_CATEGORY_FUEL,
				inventory_v1.Category_CATEGORY_PORTHOLE,
				inventory_v1.Category_CATEGORY_WING,
				inventory_v1.Category_CATEGORY_UNSPECIFIED,
			}))
			gomega.Expect(resp.GetPart().GetCreatedAt()).ToNot(gomega.BeNil())
			gomega.Expect(resp.GetPart().GetUpdatedAt()).ToNot(gomega.BeNil())
		})

		ginkgo.It("должен возвращать ошибку для несуществующего UUID", func() {
			nonExistentUUID := gofakeit.UUID()
			resp, err := clients.InventoryClient.GetPart(ctx, &inventory_v1.GetPartRequest{Uuid: nonExistentUUID})

			gomega.Expect(err).To(gomega.HaveOccurred(), "ожидали ошибку для несуществующей детали")
			gomega.Expect(status.Code(err)).To(gomega.Equal(codes.NotFound))
			gomega.Expect(resp).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("Test ListParts", func() {
		ginkgo.BeforeEach(func() {
			ctx, cancel = context.WithCancel(suiteCtx)
			err := env.ClearPartsCollection(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			ctx = middlewaregrpc.AddSessionUUIDToContext(ctx, sessionUUID)
			ctx = middlewaregrpc.ForwardSessionUUIDToGRPC(ctx)
		})

		ginkgo.AfterEach(func() {
			cancel()
		})

		ginkgo.It("Должен возвращать список деталей", func() {
			// Setup
			const quantity = 3
			partUUIDs := make([]string, 0)
			for i := 0; i < quantity; i++ {
				partUUID, err := env.InsertTestPart(ctx)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(partUUID).ToNot(gomega.BeEmpty())
				partUUIDs = append(partUUIDs, partUUID)
			}

			// Test
			resp, err := clients.InventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
				Filter: &inventory_v1.PartsFilter{
					Uuids: partUUIDs,
				},
			})

			// Verify
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.GetParts()).ToNot(gomega.BeNil())
			gomega.Expect(len(resp.GetParts())).To(gomega.Equal(len(partUUIDs)))

			// Создаем map для проверки наличия всех ожидаемых UUID
			expectedUUIDs := make(map[string]bool)
			for _, uuid := range partUUIDs {
				expectedUUIDs[uuid] = true
			}

			// Проверяем, что все возвращенные детали соответствуют ожидаемым UUID
			returnedUUIDs := make(map[string]bool)
			for _, part := range resp.GetParts() {
				returnedUUIDs[part.GetUuid()] = true
				gomega.Expect(expectedUUIDs[part.GetUuid()]).To(gomega.BeTrue(), "неожиданный UUID детали: %s", part.GetUuid())
			}

			// Проверяем, что все ожидаемые UUID присутствуют в ответе
			for _, uuid := range partUUIDs {
				gomega.Expect(returnedUUIDs[uuid]).To(gomega.BeTrue(), "ожидали UUID %s в ответе", uuid)
			}
		})
	})
})
