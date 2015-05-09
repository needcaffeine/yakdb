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
	items = make(map[string]Item)
}

func (suite *RepoSuite) TestListOfItemsStartsOutEmpty() {
	suite.Len(GetItems(), 0, "The initial list of items should be empty")
}

func (suite *RepoSuite) TestCanFindOneItemById() {
	CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})
	i := FindOneItemById("1")

	suite.Equal("1", i.Id, "Item ID should be the same")
	suite.Equal("Item 1 Value", i.Value, "Item name should be the same")
}

func (suite *RepoSuite) TestFindNonexistentItemShouldReturnEmptyItem() {
	i := FindOneItemById("id-does-not-exist")

	suite.Equal(i, Item{}, "Finding non-existent id should return empty item")
}

func (suite *RepoSuite) TestCreateOneItem() {
	r := CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})

	suite.Equal("OK", r.Status, "Status from creating an item should be 'OK'")
	suite.Len(GetItems(), 1, "List should have 1 entry after creation")
}

func (suite *RepoSuite) TestDeleteOneItem() {
	CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})
	suite.Len(GetItems(), 1, "Test should start with 1 item")

	DeleteOneItemById("1")
	suite.Len(GetItems(), 0, "Items should be empty after deletion")
}

func (suite *RepoSuite) TestDeletionOfNonexistentItemIsFine() {
	DeleteOneItemById("1")
}

func (suite *RepoSuite) TestFlushItemsPurgesList() {
	CreateOneItem(Item{"1", "Item 1 Value", time.Now().Unix(), 0})
	CreateOneItem(Item{"2", "Item 2 Value", time.Now().Unix(), 10})

	suite.Len(GetItems(), 2, "Items should contain 2 entries")

	FlushItems()

	suite.Len(GetItems(), 0, "Items should be empty after flush")
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
