package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/xavesen/search-admin/internal/storage"
	"github.com/xavesen/search-admin/internal/utils"
)

func TestPingHandler(t *testing.T) {
	server := NewServer("", &storage.StorageMock{}, nil)

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("Unable to create request, error: %s\n", err)
	}

	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	expectedRespStruct := utils.Response{
		Success: true,
		ErrorMessage: "",
		Data: PingResponse{
			Pong: "pong",
		},
	}
	expectedResp, err := json.Marshal(expectedRespStruct)
	if err != nil {
		t.Fatalf("Unable to marshal expected response, error: %s\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK, "wrong response code")
	assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), string(expectedResp), "wrong body contents")
}