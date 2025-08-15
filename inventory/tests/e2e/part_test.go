package integration

import (
	"context"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ = ginkgo.Describe("InventoryService", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		client inventory_v1.InventoryServiceClient
	)

	ginkgo.BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		// Verify
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		client = inventory_v1.NewInventoryServiceClient(conn)
	})

	ginkgo.AfterEach(func() {
		cancel()
	})

	ginkgo.Describe("Test GetPart", func() {
		ginkgo.It("должен возвращать деталь по UUID", func() {
			// Setup
			partUUID, err := env.InsertTestPart(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			resp, err := client.GetPart(ctx, &inventory_v1.GetPartRequest{Uuid: partUUID})

			// Verify
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.GetPart()).ToNot(gomega.BeNil())
			gomega.Expect(resp.GetPart().GetUuid()).To(gomega.Equal(partUUID))
			gomega.Expect(resp.GetPart().GetName()).ToNot(gomega.BeEmpty())
			gomega.Expect(resp.GetPart().GetDescription()).ToNot(gomega.BeEmpty())
			gomega.Expect(resp.GetPart().GetCreatedAt()).ToNot(gomega.BeNil())
		})
	})

	ginkgo.Describe("Test ListParts", func() {
		ginkgo.It("должен возвращать список деталей", func() {
			// Setup
			partUUIDs := make([]string, 0)
			for i := 0; i < 2; i++ {
				partUUID, err := env.InsertTestPart(ctx)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				partUUIDs = append(partUUIDs, partUUID)
			}

			// Test
			resp, err := client.ListParts(ctx, &inventory_v1.ListPartsRequest{
				Filter: &inventory_v1.PartsFilter{
					Uuids: partUUIDs,
				},
			})

			// Verify
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.GetParts()).ToNot(gomega.BeNil())
			gomega.Expect(len(resp.GetParts())).To(gomega.BeNumerically("==", len(partUUIDs)))
			for i := 0; i < 2; i++ {
				gomega.Expect(resp.GetParts()[i].GetUuid()).ToNot(gomega.Equal(partUUIDs[i]))
			}
		})
	})
})
