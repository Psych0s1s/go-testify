package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")
	assert.NotEmpty(t, responseRecorder.Body.String(), "Expected non-empty response body")
}

func TestMainHandlerWrongCity(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=1&city=berlin", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status code 400")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Expected error message 'wrong city value'")
}

func TestMainHandlerWhenCountExceedsAvailable(t *testing.T) {
	TotalCount := 4

	req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")

	cafesReturned := strings.Split(responseRecorder.Body.String(), ",")

	require.Len(t, cafesReturned, TotalCount, "Expected all available cafes to be returned when the requested count exceeds the available ones")

	expectedCafes := []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"}
	assert.ElementsMatch(t, cafesReturned, expectedCafes, "The returned cafes should match the expected list of all available cafes")
}
