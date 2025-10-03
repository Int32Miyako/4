package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFindUsers(t *testing.T) {
	// Тестовый сервер для успешных запросов
	server := httptest.NewServer(http.HandlerFunc(searchHandler))
	defer server.Close()

	tests := []struct {
		name    string
		req     SearchRequest
		IsError bool
	}{
		{
			name:    "Успешный запрос",
			req:     SearchRequest{Limit: 10, Offset: 0, Query: "Boyd", OrderField: "Age", OrderBy: OrderByDesc},
			IsError: false,
		},
		{
			name:    "Лимит больше 25 (должен быть обрезан до 25)",
			req:     SearchRequest{Limit: 50, Offset: 0, Query: "", OrderField: "Name", OrderBy: OrderByAsc},
			IsError: false,
		},
		{
			name:    "Отрицательный лимит",
			req:     SearchRequest{Limit: -1, Offset: 0, Query: "", OrderField: "Name", OrderBy: OrderByAsc},
			IsError: true,
		},
		{
			name:    "Отрицательный оффсет",
			req:     SearchRequest{Limit: 10, Offset: -1, Query: "", OrderField: "Name", OrderBy: OrderByAsc},
			IsError: true,
		},
		{
			name:    "Пусто поле сортировки",
			req:     SearchRequest{Limit: 10, Offset: 0, Query: "", OrderField: "InvalidField", OrderBy: OrderByAsc},
			IsError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchClient := &SearchClient{
				AccessToken: "test-token",
				URL:         server.URL,
			}

			resp, err := searchClient.FindUsers(tt.req)

			if tt.IsError {
				if err == nil {
					t.Errorf("ожидали ошибку, но получили nil")
				}
				return
			}

			if err != nil {
				t.Errorf("неожиданная ошибка: %v", err)
				return
			}

			if resp == nil {
				t.Errorf("ожидали не nil ответ")
			}
		})
	}
}

// Тест ошибки 401 Unauthorized
func TestFindUsers401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	searchClient := &SearchClient{
		AccessToken: "invalid-token",
		URL:         server.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали ошибку 401")
	}
}

// Тест ошибки 500 Internal Server Error
func TestFindUsers500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	searchClient := &SearchClient{
		AccessToken: "test-token",
		URL:         server.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали ошибку 500")
	}
}

// Тест невалидного JSON в ошибке
func TestFindUsersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("invalid json"))
		if err != nil {
			t.Errorf(err.Error())
		}
	}))
	defer server.Close()

	searchClient := &SearchClient{
		AccessToken: "test-token",
		URL:         server.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали ошибку невалидного JSON")
	}
}

// Тест неизвестной ошибки 400
func TestFindUsersUnknownBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errorResp := SearchErrorResponse{Error: "UnknownError"}
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			t.Errorf(err.Error())
		}
	}))
	defer server.Close()

	searchClient := &SearchClient{
		AccessToken: "test-token",
		URL:         server.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали неизвестную ошибку 400")
	}
}

// Тест невалидного JSON в ответе
func TestFindUsersInvalidResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("invalid json"))
		if err != nil {
			t.Errorf(err.Error())
		}
	}))
	defer server.Close()

	searchClient := &SearchClient{
		AccessToken: "test-token",
		URL:         server.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали ошибку невалидного JSON в ответе")
	}
}

// Тест таймаута
func TestFindUsersTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Больше чем таймаут клиента
	}))
	defer server.Close()

	searchClient := &SearchClient{
		AccessToken: "test-token",
		URL:         server.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали ошибку таймаута")
	}
}

// Тест сетевой ошибки
func TestFindUsersNetworkError(t *testing.T) {
	searchClient := &SearchClient{
		AccessToken: "test-token",
		URL:         "http://invalid-url-that-does-not-exist.local",
	}

	_, err := searchClient.FindUsers(SearchRequest{Limit: 10})
	if err == nil {
		t.Errorf("ожидали сетевую ошибку")
	}
}
