package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	//"strings"
	//"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) SetupTest() {
	// Let's seed our list with an initial item to make the tests a little
	// cleaner
	var items = make(Items)
	items["1"] = Item{
		Data: map[string]interface{}{
			"value": "Item 1 Value",
		},
		Created: time.Now().Unix(),
	}

	var collections = make(Collections)
	collections["1"] = Collection{
		Items: items,
	}
}

func (suite *RepoSuite) TestJsonResponseSetsAppropriateHeaders() {
	w := httptest.NewRecorder()

	handle := func(w http.ResponseWriter, r *http.Request, p map[string]string) {}

	handle = JsonResponse(handle)
	handle(w, nil, nil)
	suite.Equal("application/json; charset=UTF-8", w.Header().Get("Content-Type"), "Content-Type should be application/json; UTF-8")
}

func (suite *HandlerSuite) TestIndexReturnsHelp() {
	w := httptest.NewRecorder()

	Index(w, nil, nil)

	usageLines := []string{
		"This is yakdb, a highly performant key-value store written in Go.",
		"Usage:",
		"------",
		"List all items: GET /items",
		"Get an item: GET /items/{itemid}",
		"Put an item: PUT /items",
		"Delete an item: DELETE /items/{itemid}",
		"Delete all items: DELETE /items",
		"More documentation: https://github.com/needcaffeine/yakdb"}

	body := w.Body.String()

	for _, s := range usageLines {
		suite.Contains(body, s)
	}
}

func (suite *HandlerSuite) TestListOfCollectionsStartsOutEmpty() {
	w := httptest.NewRecorder()

	CollectionsList(w, nil, nil)

	c := make(Collections)
	json.Unmarshal(w.Body.Bytes(), &c)

	suite.Len(c, 0, "List of collections should be empty.")
}
