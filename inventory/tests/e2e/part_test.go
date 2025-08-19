//go:build integration

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
		gomega.Expect(err).ToNot(gomega.HaveOccurred(), "ожидали успешное подключение к gRPC приложению")
		client = inventory_v1.NewInventoryServiceClient(conn)
	})

	ginkgo.AfterEach(func() {
		cancel()
	})

	ginkgo.Describe("Test GetPart", func() {
		ginkgo.BeforeEach(func() {
			ctx, cancel = context.WithCancel(suiteCtx)
			err := env.ClearPartsCollection(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.AfterEach(func() {
			cancel()
		})

		ginkgo.It("должен возвращать деталь по UUID", func() {
			// Setup
			partUUID, err := env.InsertTestPart(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(partUUID).ToNot(gomega.BeEmpty())

			resp, err := client.GetPart(ctx, &inventory_v1.GetPartRequest{Uuid: partUUID})

			// Verify
			gomega.Expect(err).ToNot(gomega.HaveOccurred(), "")
			gomega.Expect(resp.GetPart()).ToNot(gomega.BeNil())
			gomega.Expect(resp.GetPart().GetUuid()).To(gomega.Equal(partUUID))
			gomega.Expect(resp.GetPart().GetName()).ToNot(gomega.BeEmpty())
			gomega.Expect(resp.GetPart().GetDescription()).ToNot(gomega.BeEmpty())
			gomega.Expect(resp.GetPart().GetCreatedAt()).ToNot(gomega.BeNil())
		})
	})

	ginkgo.Describe("Test ListParts", func() {
		ginkgo.BeforeEach(func() {
			ctx, cancel = context.WithCancel(suiteCtx)
			err := env.ClearPartsCollection(ctx)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		ginkgo.AfterEach(func() {
			cancel()
		})

		ginkgo.It("должен возвращать список деталей", func() {
			// Setup
			const quanity = 1
			partUUIDs := make([]string, 0)
			for i := 0; i < quanity; i++ {
				partUUID, err := env.InsertTestPart(ctx)
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(partUUID).ToNot(gomega.BeEmpty())
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
			gomega.Expect(len(resp.GetParts())).To(gomega.Equal(len(partUUIDs)))
			for i := 0; i < quanity; i++ {
				gomega.Expect(resp.GetParts()[i].GetUuid()).To(gomega.Equal(partUUIDs[i]))
			}
		})
	})
})
