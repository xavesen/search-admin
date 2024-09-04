package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/xavesen/search-admin/internal/models"
	"github.com/xavesen/search-admin/internal/storage"
	"github.com/xavesen/search-admin/internal/utils"
)

var createFilterTests = []struct {
	testName			string
	storage				*storage.StorageMock
	payload				*models.Filter
	expectedCode		int
	expectedResponse	utils.Response
}{
	{
		testName: "Returns 201 and filter with id when payload is correct",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload: &models.Filter{
				Regex: "^[a-zA-Z]+$",
		},
		expectedCode: http.StatusCreated,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.Filter{
				Id:	"1",
				Regex: "^[a-zA-Z]+$",
			},
		},
	},
	{
		testName: "Returns 400 with empty payload",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.Filter{
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: regex is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 with wrong regex",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload: &models.Filter{
			Regex: "+++",
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: regex must be a regular expression accepted by RE2",
			Data: nil,
		},
	},
	{
		testName: "Returns 500 when db returns an error",
		storage: &storage.StorageMock{
			Error: 	errors.New("random error"),
		},
		payload: &models.Filter{
				Regex: "^[a-zA-Z]+$",
		},
		expectedCode: http.StatusInternalServerError,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
}

func TestCreateFilterHandler(t *testing.T) {
	for i, test := range createFilterTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.storage, nil)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/filter", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatalf("Unable to create request, error: %s\n", err)
		}

		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)

		expectedResp, err := json.Marshal(test.expectedResponse)
		if err != nil {
			t.Fatalf("Unable to marshal expected response, error: %s\n", err)
		}

		assert.Equal(t, rr.Code, test.expectedCode, "wrong response code")
		assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), string(expectedResp), "wrong body contents")
	}
}