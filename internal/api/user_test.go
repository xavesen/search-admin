package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/magiconair/properties/assert"
	"github.com/xavesen/search-admin/internal/models"
	"github.com/xavesen/search-admin/internal/storage"
	"github.com/xavesen/search-admin/internal/utils"
)

var getAllUsersTests = []struct {
	testName			string
	storage				*storage.StorageMock
	expectedCode		int
	expectedResponse	utils.Response
}{
	{
		testName: "Return 200 and all users",
		storage: &storage.StorageMock{
			Error: 	nil,
			Users:	[]models.User{
				{
					Id:	"1",
					Login: "mary",
					Password: "12345",
					IndexLimit: 5,
				},
				{
					Id:	"2",
					Login: "dane",
					Password: "qwerty",
					IndexLimit: 4,
				},
				{
					Id:	"3",
					Login: "linda",
					Password: "password",
					IndexLimit: 1,
				},
			},
		},
		expectedCode: http.StatusOK,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: []models.User{
				{
					Id:	"1",
					Login: "mary",
					Password: "12345",
					IndexLimit: 5,
				},
				{
					Id:	"2",
					Login: "dane",
					Password: "qwerty",
					IndexLimit: 4,
				},
				{
					Id:	"3",
					Login: "linda",
					Password: "password",
					IndexLimit: 1,
				},
			},
		},
	},
	{
		testName: "Return 200 but there are no users in db",
		storage: &storage.StorageMock{
			Error: 	nil,
			Users:	[]models.User{},
		},
		expectedCode: http.StatusOK,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: []models.User{},
		},
	},
	{
		testName: "Return 500 when DB returns an error",
		storage: &storage.StorageMock{
			Error: 	errors.New("random error"),
			Users:	[]models.User{
				{
					Id:	"1",
					Login: "mary",
					Password: "12345",
					IndexLimit: 5,
				},
				{
					Id:	"2",
					Login: "dane",
					Password: "qwerty",
					IndexLimit: 4,
				},
				{
					Id:	"3",
					Login: "linda",
					Password: "password",
					IndexLimit: 1,
				},
			},
		},
		expectedCode: http.StatusInternalServerError,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},

}

func TestGetAllUsersHandler(t *testing.T) {
	for i, test := range getAllUsersTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.storage, nil)

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
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