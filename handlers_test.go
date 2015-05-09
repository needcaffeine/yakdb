package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func (suite *HandlerSuite) SetupTest() {
	// Let's seed our list with an initial item to make the tests a little
	// cleaner
	items["1"] = Item{"1", "Item 1 Value", time.Now().Unix(), 0}
}

func (suite *RepoSuite) TestJsonResponseSetsAppropriateHeaders() {
	w := httptest.NewRecorder()

	handle := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	}

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

func (suite *HandlerSuite) TestListOfItemsStartsOutEmpty() {
	items = make(map[string]Item)
	w := httptest.NewRecorder()

	ItemsList(w, nil, nil)

	i := make(map[string]Item)
	json.Unmarshal(w.Body.Bytes(), &i)

	suite.Len(i, 0, "List of items should be empty")
}

func (suite *HandlerSuite) TestCanGetListOfItems() {
	w := httptest.NewRecorder()
	ItemsList(w, nil, nil)

	i := make(map[string]Item)
	json.Unmarshal(w.Body.Bytes(), &i)

	suite.Len(i, 1, "List of items should have 1 entry")
	suite.Equal("Item 1 Value", i["1"].Value, "First item should have value of 'Item 1 Value'")
}

func (suite *HandlerSuite) TestCanGetOneItem() {
	params := httprouter.Params{}
	params = append(params, httprouter.Param{"itemId", "1"})

	w := httptest.NewRecorder()
	ItemsGet(w, nil, params)

	i := Item{}
	json.Unmarshal(w.Body.Bytes(), &i)

	suite.Equal("Item 1 Value", i.Value, "First item should have value of 'Item 1 Value'")
}

func (suite *HandlerSuite) TestCanCreateOneItem() {
	itemJson := `{"id": "2", "value": "Item 2 Value"}`

	reader := strings.NewReader(itemJson)
	request, err := http.NewRequest("PUT", "/items", reader)
	if err != nil {
		suite.Error(err)
	}

	w := httptest.NewRecorder()
	ItemsPut(w, request, nil)

	suite.Len(items, 2, "Items should have 2 entries now")
}

func (suite *HandlerSuite) TestMissingValueGeneratesFailure() {
	itemJson := `{"id": "2"}`

	reader := strings.NewReader(itemJson)
	request, err := http.NewRequest("PUT", "/items", reader)
	if err != nil {
		suite.Error(err)
	}

	w := httptest.NewRecorder()
	ItemsPut(w, request, nil)

	suite.Len(items, 1, "Items should have 1 entry still")
	suite.Equal(http.StatusBadRequest, w.Code, "Missing value should return 400 bad request")
}

func (suite *HandlerSuite) TestCanDeleteOneItem() {
	params := httprouter.Params{}
	params = append(params, httprouter.Param{"itemId", "1"})

	w := httptest.NewRecorder()
	ItemsDelete(w, nil, params)

	r := Result{}
	json.Unmarshal(w.Body.Bytes(), &r)

	suite.Equal("OK", r.Status, "Flush should return 'OK'")
	suite.Len(items, 0, "List should be empty")
}

func (suite *HandlerSuite) TestCanDeleteAllItems() {
	w := httptest.NewRecorder()
	ItemsFlush(w, nil, nil)

	r := Result{}
	json.Unmarshal(w.Body.Bytes(), &r)

	suite.Equal("OK", r.Status, "Flush should return 'OK'")
	suite.Len(items, 0, "List should be empty")
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
