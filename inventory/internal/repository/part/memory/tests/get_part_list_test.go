package tests

import (
	"context"
	"sync"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// TestGetPartList_NoFilter verifies that GetPartList returns all parts when no filter is applied.
// It checks that the correct number of parts is returned and no error occurs.
func (s *InventoryRepositorySuite) TestGetPartList_NoFilter() {
	// Test

	ctx := context.Background()

	result, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{})

	// Verify
	s.Require().NoError(err)
	s.Require().Len(result, 20)
}

// TestGetPartList_FilterByUUIDs verifies filtering parts by UUIDs works correctly.
// Tests that only parts with matching UUIDs are returned, and non-existent UUIDs are ignored.
func (s *InventoryRepositorySuite) TestGetPartList_FilterByUUIDs() {
	// Setup
	var (
		ctx   = context.Background()
		uuids = []string{s.parts[0].Uuid, s.parts[1].Uuid, "non-existent-uuid"}
	)

	// Test
	result, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{UUIDs: uuids})

	// Verify
	s.Require().NoError(err)
	s.Require().Len(result, 2)
	s.Require().Contains([]string{result[0].Uuid, result[1].Uuid}, s.parts[0].Uuid)
	s.Require().Contains([]string{result[0].Uuid, result[1].Uuid}, s.parts[1].Uuid)
}

// TestGetPartList_FilterByNames verifies filtering parts by names functions properly.
// Ensures the repository returns all parts matching the provided names and ignores non-existent names.
func (s *InventoryRepositorySuite) TestGetPartList_FilterByNames() {
	// Setup
	var (
		ctx   = context.Background()
		names = []string{s.parts[0].Name, s.parts[1].Name, "non-existent-name"}
	)

	// Test
	result, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{Names: names})

	// Verify
	s.Require().NoError(err)
	s.Require().True(len(result) >= 2) // At least the two we filtered for
	for _, part := range result {
		s.Require().Contains(names[:2], part.Name)
	}
}

// TestGetPartList_FilterByCategories verifies filtering by part categories works as expected.
// Checks that only parts belonging to the specified categories are returned.
func (s *InventoryRepositorySuite) TestGetPartList_FilterByCategories() {
	// Setup
	var (
		ctx        = context.Background()
		categories = []model.PartCategory{model.PartCategoryEngine, model.PartCategoryFuel}
	)

	// Test
	result, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{Categories: categories})

	// Verify
	s.Require().NoError(err)
	s.Require().True(len(result) > 0)
	for _, part := range result {
		s.Require().Contains(categories, part.Category)
	}
}

// TestGetPartList_FilterByCountries verifies filtering by manufacturer countries functions correctly.
// Tests that only parts manufactured in the specified countries are returned.
func (s *InventoryRepositorySuite) TestGetPartList_FilterByCountries() {
	// Setup
	var (
		ctx       = context.Background()
		countries = []string{"USA", "Germany"}
	)

	// Test
	result, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{ManufacturerCountries: countries})

	// Verify
	s.Require().NoError(err)
	s.Require().True(len(result) > 0)
	for _, part := range result {
		s.Require().Contains(countries, part.Manufacturer.Country)
	}
}

// TestGetPartList_FilterByTags verifies filtering by tags works properly.
// Ensures parts are correctly filtered when they contain any of the specified tags.
func (s *InventoryRepositorySuite) TestGetPartList_FilterByTags() {
	// Setup
	var (
		ctx  = context.Background()
		tags = []string{s.parts[0].Tags[0], s.parts[1].Tags[1]}
	)

	// Test
	result, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{Tags: tags})

	// Verify
	s.Require().NoError(err)
	s.Require().True(len(result) > 0)
	for _, part := range result {
		hasMatchingTag := false
		for _, tag := range tags {
			for _, partTag := range part.Tags {
				if tag == partTag {
					hasMatchingTag = true
					break
				}
			}
			if hasMatchingTag {
				break
			}
		}
		s.Require().True(hasMatchingTag)
	}
}

// TestGetPartList_CombinedFilters verifies that multiple filters can be combined.
// Tests that the repository correctly applies all specified filters simultaneously.
func (s *InventoryRepositorySuite) TestGetPartList_CombinedFilters() {
	// Setup
	var (
		ctx    = context.Background()
		filter = model.PartFilter{
			Categories:            []model.PartCategory{model.PartCategoryEngine},
			ManufacturerCountries: []string{"USA"},
		}
	)

	// Test
	result, err := s.inventoryRepository.GetPartList(ctx, filter)

	// Verify
	s.Require().NoError(err)
	for _, part := range result {
		s.Require().Equal(model.PartCategoryEngine, part.Category)
		s.Require().Equal("USA", part.Manufacturer.Country)
	}
}

// TestGetPartList_ConcurrentAccess verifies thread-safe behavior during concurrent access.
// Ensures the repository can handle multiple simultaneous read operations without errors.
func (s *InventoryRepositorySuite) TestGetPartList_ConcurrentAccess() {
	// Test
	var (
		wg      sync.WaitGroup
		ctx     = context.Background()
		results = make(chan []model.Part, 5)
		errs    = make(chan error, 5)
	)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := s.inventoryRepository.GetPartList(ctx, model.PartFilter{})
			if err != nil {
				errs <- err
				return
			}
			results <- res
		}()
	}

	wg.Wait()
	close(results)
	close(errs)

	// Verify
	for res := range results {
		s.Require().Len(res, 20)
	}
	for err := range errs {
		s.Require().NoError(err)
	}
}

// TestGetPartList_EmptyRepository verifies behavior when querying an empty repository.
// Checks that an empty result is returned without errors when no parts exist.
func (s *InventoryRepositorySuite) TestGetPartList_EmptyRepository() {
	// Setup
	var (
		ctx       = context.Background()
		emptyRepo = memory.NewInventoryRepository()
	)

	// Test
	result, err := emptyRepo.GetPartList(ctx, model.PartFilter{})

	// Verify
	s.Require().NoError(err)
	s.Require().Empty(result)
}
