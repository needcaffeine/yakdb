package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite
}

func (suite *RepoSuite) SetupAllSuite() {
	// Make sure we start off with an empty list of items
	items = make(Items)
}

func (suite *RepoSuite) TestListOfItemsStartsOutEmpty() {
	suite.Len(items.GetItems(), 0, "The initial list of items should be empty")
}

func (suite *RepoSuite) TestCanFindOneItemById() {
	items.CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})
	i := items.FindOneItemById("1")

	suite.Equal("1", i.Id, "Item ID should be the same")
	suite.Equal("Item 1 Value", i.Value, "Item name should be the same")
}

func (suite *RepoSuite) TestFindNonexistentItemShouldReturnEmptyItem() {
	i := items.FindOneItemById("id-does-not-exist")

	suite.Equal(i, Item{}, "Finding non-existent id should return empty item")
}

func (suite *RepoSuite) TestCreateOneItem() {
	r := items.CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})

	suite.Equal("OK", r.Status, "Status from creating an item should be 'OK'")
	suite.Len(items.GetItems(), 1, "List should have 1 entry after creation")
}

func (suite *RepoSuite) TestDeleteOneItem() {
	items.CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})
	suite.Len(items.GetItems(), 1, "Test should start with 1 item")

	items.DeleteOneItemById("1")
	suite.Len(items.GetItems(), 0, "Items should be empty after deletion")
}

func (suite *RepoSuite) TestDeletionOfNonexistentItemIsFine() {
	items.DeleteOneItemById("1")
}

func (suite *RepoSuite) TestFlushItemsPurgesList() {
	items.CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})
	items.CreateOneItem(Item{"2", "Item 2 Value", time.Now().Unix(), 10})

	suite.Len(items.GetItems(), 2, "Items should contain 2 entries")

	items.FlushItems()

	suite.Len(items.GetItems(), 0, "Items should be empty after flush")
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
