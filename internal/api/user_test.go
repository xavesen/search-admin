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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

var getUserByIdTests = []struct {
	testName			string
	storage				*storage.StorageMock
	userId				string
	expectedCode		int
	expectedResponse	utils.Response
}{
	{
		testName: "Returns 200 and user",
		storage: &storage.StorageMock{
			Error: 	nil,
			User:	models.User{
				Id:	"1",
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
			},
		},
		userId: "1",
		expectedCode: http.StatusOK,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.User{
				Id:	"1",
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
			},
		},
	},
	{
		testName: "Returns 404 when no such id in db",
		storage: &storage.StorageMock{
			Error: 	mongo.ErrNoDocuments,
		},
		userId: "2",
		expectedCode: http.StatusNotFound,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "No user with such id",
			Data: nil,
		},
	},
	{
		testName: "Returns 404 when id is invalid",
		storage: &storage.StorageMock{
			Error: 	primitive.ErrInvalidHex,
		},
		userId: "c",
		expectedCode: http.StatusNotFound,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "No user with such id",
			Data: nil,
		},
	},
}

func TestGetUserByIdHandler(t *testing.T) {
	for i, test := range getUserByIdTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.storage, nil)

		path := fmt.Sprintf("/user/%s", test.userId)
		req, err := http.NewRequest(http.MethodGet, path, nil)
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

var createUserTests = []struct {
	testName			string
	storage				*storage.StorageMock
	payload				*models.User
	expectedCode		int
	expectedResponse	utils.Response
}{
	{
		testName: "Returns 201 and user with id when payload is correct",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload: &models.User{
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
		},
		expectedCode: http.StatusCreated,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.User{
				Id:	"1",
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
			},
		},
	},
	{
		testName: "Returns 201 and user with id when payload is correct but no indexes listed",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload: &models.User{
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
				Indexes: []models.Index{
					{
						Name: "aaa",
						Id: "aaa",
					},
					{
						Name: "bbb",
						Id: "bbb",
					},
				},
		},
		expectedCode: http.StatusCreated,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.User{
				Id:	"1",
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
				Indexes: []models.Index{
					{
						Name: "aaa",
						Id: "aaa",
					},
					{
						Name: "bbb",
						Id: "bbb",
					},
				},
			},
		},
	},
	{
		testName: "Returns 400 with empty payload",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: login is required, password is required, index_limit is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without login",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
			Password: "12345",
			IndexLimit: 5,
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: login is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without index_limit",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload: &models.User{
			Login: "mary",
			Password: "12345",
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: index_limit is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without password",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
			Login: "mary",
			IndexLimit: 5,
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: password is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without id in index",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Name: "aaa",
				},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes id is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without name in index",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Id: "aaa",
				},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes name is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 with empty index",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes id is required, name is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 with one of indexes empty",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Id: "aaa",
					Name: "aaa",
				},
				{},
				{
					Id: "bbb",
					Name: "bbb",
				},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes id is required, name is required",
			Data: nil,
		},
	},
}

func TestCreateUserHandler(t *testing.T) {
	for i, test := range createUserTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.storage, nil)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(marshaledPayload))
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

var deleteUserTests = []struct {
	testName			string
	storage				*storage.StorageMock
	userId				string
	expectedCode		int
	expectedResponse	utils.Response
}{
	{
		testName: "Returns 200",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "1",
		expectedCode: http.StatusOK,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: nil,
		},
	},
	{
		testName: "Returns 404 when no user with such id in db",
		storage: &storage.StorageMock{
			Error: 	mongo.ErrNoDocuments,
		},
		userId: "1",
		expectedCode: http.StatusNotFound,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "No user with such id",
			Data: nil,
		},
	},
	{
		testName: "Returns 404 when id is invalid",
		storage: &storage.StorageMock{
			Error: 	primitive.ErrInvalidHex,
		},
		userId: "1",
		expectedCode: http.StatusNotFound,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "No user with such id",
			Data: nil,
		},
	},
	{
		testName: "Returns 500 when db returns an error",
		storage: &storage.StorageMock{
			Error: 	errors.New("random error"),
		},
		userId: "1",
		expectedCode: http.StatusInternalServerError,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
}

func TestDeleteUserHandler(t *testing.T) {
	for i, test := range deleteUserTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.storage, nil)

		path := fmt.Sprintf("/user/%s", test.userId)
		req, err := http.NewRequest(http.MethodDelete, path, nil)
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

var updateUserTests = []struct {
	testName			string
	storage				*storage.StorageMock
	payload				*models.User
	userId				string
	expectedCode		int
	expectedResponse	utils.Response
}{
	{
		testName: "Returns 200 without indexes in payload",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload: &models.User{
			Login: "mary",
			Password: "12345",
			IndexLimit: 5,
		},
		expectedCode: http.StatusOK,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.User{
				Id:	"66d8420df6e5311a791e0a08",
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
			},
		},
	},
	{
		testName: "Returns 200 when indexes are in payload",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload: &models.User{
			Login: "mary",
			Password: "12345",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Name: "aaa",
					Id: "aaa",
				},
				{
					Name: "bbb",
					Id: "bbb",
				},
			},
		},
		expectedCode: http.StatusOK,
		expectedResponse: utils.Response{
			Success: true,
			ErrorMessage: "",
			Data: models.User{
				Id:	"66d8420df6e5311a791e0a08",
				Login: "mary",
				Password: "12345",
				IndexLimit: 5,
				Indexes: []models.Index{
					{
						Name: "aaa",
						Id: "aaa",
					},
					{
						Name: "bbb",
						Id: "bbb",
					},
				},
			},
		},
	},
	{
		testName: "Returns 400 with empty payload",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: login is required, password is required, index_limit is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without login",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Password: "12345",
			IndexLimit: 5,
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: login is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without index_limit",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "12345",
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: index_limit is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without password",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			IndexLimit: 5,
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: password is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without id in index",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Name: "aaa",
				},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes id is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 without name in index",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Id: "aaa",
				},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes name is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 with empty index",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes id is required, name is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 400 with one of indexes empty",
		storage: &storage.StorageMock{
			Error: 	nil,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "aaa",
			IndexLimit: 5,
			Indexes: []models.Index{
				{
					Id: "aaa",
					Name: "aaa",
				},
				{},
				{
					Id: "bbb",
					Name: "bbb",
				},
			},
		},
		expectedCode: http.StatusBadRequest,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Bad request: in indexes id is required, name is required",
			Data: nil,
		},
	},
	{
		testName: "Returns 404 when no such id in db",
		storage: &storage.StorageMock{
			Error: 	mongo.ErrNoDocuments,
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "12345",
			IndexLimit: 5,
		},
		expectedCode: http.StatusNotFound,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "No user with such id",
			Data: nil,
		},
	},
	{
		testName: "Returns 404 when id is invalid",
		storage: &storage.StorageMock{
			Error: 	primitive.ErrInvalidHex,
		},
		userId: "c",
		payload:  &models.User{
			Login: "mary",
			Password: "12345",
			IndexLimit: 5,
		},
		expectedCode: http.StatusNotFound,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "No user with such id",
			Data: nil,
		},
	},
	{
		testName: "Returns 500 when db returns an error",
		storage: &storage.StorageMock{
			Error: 	errors.New("random error"),
		},
		userId: "66d8420df6e5311a791e0a08",
		payload:  &models.User{
			Login: "mary",
			Password: "12345",
			IndexLimit: 5,
		},
		expectedCode: http.StatusInternalServerError,
		expectedResponse: utils.Response{
			Success: false,
			ErrorMessage: "Internal server error",
			Data: nil,
		},
	},
}

func TestUpdateUserHandler(t *testing.T) {
	for i, test := range updateUserTests {
		fmt.Printf("Running test #%d: %s\n", i+1, test.testName)

		server := NewServer("", test.storage, nil)

		marshaledPayload, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Unable to marshal payload, error: %s\n", err)
		}

		path := fmt.Sprintf("/user/%s", test.userId)
		req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(marshaledPayload))
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