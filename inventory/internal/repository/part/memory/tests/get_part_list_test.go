package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// TestGetPartList_NoFilter verifies that GetPartList returns all parts when no filter is applied.
func TestGetPartList_NoFilter(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		parts               = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Part 1",
				Description:   "description 1",
				Price:         123.45,
				StockQuantity: 12,
				Category:      model.PartCategoryEngine,
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Part 2",
				Description:   "description 2",
				Price:         123.45,
				StockQuantity: 12,
				Category:      model.PartCategoryWing,
			},
		}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expectedParts, err := inventoryRepository.GetPartList(ctx, model.PartFilter{})

	// Verify
	require.NoError(t, err)
	require.Equal(t, len(expectedParts), len(parts))
	for _, part := range parts {
		require.Contains(t, expectedParts, part)
	}
}

// TestGetPartList_FilterByUUIDs verifies filtering parts by UUIDs works correctly.
func TestGetPartList_FilterByUUIDs(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		uuidOne             = uuid.New().String()
		uuidTwo             = uuid.New().String()

		parts = []model.Part{
			{
				Uuid:        uuidOne,
				Name:        "Part 1",
				Description: "description 1",
			},
			{
				Uuid:        uuidTwo,
				Name:        "Part 2",
				Description: "description 2",
			},
		}

		partFilter = model.PartFilter{
			UUIDs: []string{
				uuidOne,
				uuidTwo,
				"non-existent-uuid",
			},
		}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expectedParts, err := inventoryRepository.GetPartList(ctx, partFilter)

	// Verify
	require.NoError(t, err)
	require.Equal(t, len(expectedParts), len(parts))
	for _, part := range parts {
		require.Contains(t, expectedParts, part)
	}
}

// TestGetPartList_FilterByNames verifies filtering parts by names functions properly.
func TestGetPartList_FilterByNames(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		partNameOne         = "Part 1"
		partNameTwo         = "Part 2"

		parts = []model.Part{
			{
				Uuid:        gofakeit.UUID(),
				Name:        partNameOne,
				Description: "description 1",
			},
			{
				Uuid:        gofakeit.UUID(),
				Name:        partNameTwo,
				Description: "description 2",
			},
		}

		partFilter = model.PartFilter{
			Names: []string{
				partNameOne,
				partNameTwo,
				"non-existent-name",
			},
		}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expectedParts, err := inventoryRepository.GetPartList(ctx, partFilter)

	// Verify
	require.NoError(t, err)
	require.Equal(t, len(expectedParts), len(parts))
	for _, part := range parts {
		require.Contains(t, expectedParts, part)
	}
}

// TestGetPartList_FilterByCategories verifies filtering by part categories works as expected.
func TestGetPartList_FilterByCategories(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		partCategoryOne     = model.PartCategoryEngine
		partCategoryTwo     = model.PartCategoryWing

		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Part 1",
				Description:   "description 1",
				Price:         123.45,
				StockQuantity: 12,
				Category:      partCategoryOne,
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Part 2",
				Description:   "description 2",
				Price:         123.45,
				StockQuantity: 12,
				Category:      partCategoryTwo,
			},
		}

		partFilter = model.PartFilter{
			Categories: []model.PartCategory{
				partCategoryOne,
				partCategoryTwo,
				model.PartCategoryUnknown,
			},
		}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expectedParts, err := inventoryRepository.GetPartList(ctx, partFilter)

	// Verify
	require.NoError(t, err)
	require.Equal(t, len(expectedParts), len(parts))
	for _, part := range parts {
		require.Contains(t, expectedParts, part)
	}
}

// TestGetPartList_FilterByCountries verifies filtering by manufacturer countries functions correctly.
func TestGetPartList_FilterByCountries(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		countryOne          = "Russia"
		countryTwo          = "USA"

		parts = []model.Part{
			{
				Uuid:         gofakeit.UUID(),
				Name:         "Part 1",
				Description:  "description 1",
				Manufacturer: model.Manufacturer{Country: countryOne},
			},
			{
				Uuid:         gofakeit.UUID(),
				Name:         "Part 2",
				Description:  "description 2",
				Manufacturer: model.Manufacturer{Country: countryTwo},
			},
		}

		partFilter = model.PartFilter{ManufacturerCountries: []string{
			countryOne,
			countryTwo,
			"unknowm country",
		}}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expectedParts, err := inventoryRepository.GetPartList(ctx, partFilter)

	// Verify
	require.NoError(t, err)
	require.Equal(t, len(expectedParts), len(parts))
	for _, part := range parts {
		require.Contains(t, expectedParts, part)
	}
}

// TestGetPartList_FilterByTags verifies filtering by tags works properly.
func TestGetPartList_FilterByTags(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		tagOne              = "tagOne"
		tagTwo              = "tagTwo"
		tagTen              = "tagThree"

		parts = []model.Part{
			{
				Uuid:        gofakeit.UUID(),
				Name:        "Part 1",
				Description: "description 1",
				Tags:        []string{tagOne, tagTwo},
			},
			{
				Uuid:        gofakeit.UUID(),
				Name:        "Part 2",
				Description: "description 2",
				Tags:        []string{tagOne, tagTen},
			},
		}

		partFilter = model.PartFilter{Tags: []string{
			tagTwo,
			tagTen,
			tagOne,
			"unknowm tag",
		}}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expectedParts, err := inventoryRepository.GetPartList(ctx, partFilter)

	// Verify
	require.NoError(t, err)
	require.Equal(t, len(expectedParts), len(parts))
	for _, part := range parts {
		require.Contains(t, expectedParts, part)
	}
}

// TestGetPartList_CombinedFilters verifies that multiple filters can be combined.
func TestGetPartList_CombinedFilters(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		countryOne          = "Russia"
		countryTwo          = "USA"
		partCategoryOne     = model.PartCategoryEngine
		partCategoryTwo     = model.PartCategoryWing

		parts = []model.Part{
			{
				Uuid:         gofakeit.UUID(),
				Name:         "Part 1",
				Description:  "description 1",
				Manufacturer: model.Manufacturer{Country: countryOne},
				Category:     partCategoryOne,
			},
			{
				Uuid:         gofakeit.UUID(),
				Name:         "Part 2",
				Description:  "description 2",
				Manufacturer: model.Manufacturer{Country: countryTwo},
				Category:     partCategoryTwo,
			},
		}

		partFilter = model.PartFilter{
			ManufacturerCountries: []string{countryOne, countryTwo, "unknowm country"},
			Categories:            []model.PartCategory{partCategoryOne, partCategoryTwo, model.PartCategoryUnknown},
		}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	expcetedParts, err := inventoryRepository.GetPartList(ctx, partFilter)

	// Verify
	require.NoError(t, err)
	for _, part := range expcetedParts {
		require.Contains(t, expcetedParts, part)
	}
}

// TestGetPartList_ConcurrentAccess verifies thread-safe behavior during concurrent access.
func TestGetPartList_ConcurrentAccess(t *testing.T) {
	// Test
	var (
		inventoryRepository = memory.NewInventoryRepository()

		wg      sync.WaitGroup
		ctx     = context.Background()
		results = make(chan []model.Part, 5)
		errs    = make(chan error, 5)

		countryOne      = "Russia"
		countryTwo      = "USA"
		partCategoryOne = model.PartCategoryEngine
		partCategoryTwo = model.PartCategoryWing

		parts = []model.Part{
			{
				Uuid:         gofakeit.UUID(),
				Name:         "Part 1",
				Description:  "description 1",
				Manufacturer: model.Manufacturer{Country: countryOne},
				Category:     partCategoryOne,
			},
			{
				Uuid:         gofakeit.UUID(),
				Name:         "Part 2",
				Description:  "description 2",
				Manufacturer: model.Manufacturer{Country: countryTwo},
				Category:     partCategoryTwo,
			},
		}

		partFilter = model.PartFilter{
			ManufacturerCountries: []string{countryOne, countryTwo, "unknowm country"},
			Categories:            []model.PartCategory{partCategoryOne, partCategoryTwo, model.PartCategoryUnknown},
		}
	)

	// Insert data into Inventory repository
	for _, part := range parts {
		err := inventoryRepository.AddPart(ctx, part)
		require.NoError(t, err)
	}

	// Test
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			retrivedParts, err := inventoryRepository.GetPartList(ctx, partFilter)
			if err != nil {
				errs <- err
				return
			}
			results <- retrivedParts
		}()
	}

	wg.Wait()
	close(results)
	close(errs)

	// Verify
	for res := range results {
		require.Equal(t, len(res), len(parts))
	}
	for err := range errs {
		require.NoError(t, err)
	}
}

// TestGetPartList_EmptyRepository verifies behavior when querying an empty repository.
func TestGetPartList_EmptyRepository(t *testing.T) {
	// Setup
	var (
		ctx       = context.Background()
		emptyRepo = memory.NewInventoryRepository()
	)

	// Test
	result, err := emptyRepo.GetPartList(ctx, model.PartFilter{})

	// Verify
	require.NoError(t, err)
	require.Empty(t, result)
}
