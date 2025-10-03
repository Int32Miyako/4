package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSearchHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(searchHandler))
	defer server.Close()

	tests := []struct {
		name    string
		req     *SearchRequest
		IsError bool
	}{
		{
			name:    "Hilda по Age Desc",
			req:     &SearchRequest{Limit: 10, Offset: 0, Query: "Hilda", OrderField: "Age", OrderBy: OrderByDesc},
			IsError: false,
		},
		{
			name:    "Hilda по Name Asc",
			req:     &SearchRequest{Limit: 5, Offset: 0, Query: "F", OrderField: "Name", OrderBy: OrderByAsc},
			IsError: false,
		},
		{
			name:    "Hilda по Name OrderByAsIs с offset",
			req:     &SearchRequest{Limit: 10, Offset: 1, Query: "F", OrderField: "Name", OrderBy: OrderByAsIs},
			IsError: false,
		},
		{
			name:    "все пользователи по Id Desc",
			req:     &SearchRequest{Limit: 20, Offset: 0, Query: "", OrderField: "Id", OrderBy: OrderByDesc},
			IsError: false,
		},
		{
			name:    "Мария без orderField",
			req:     &SearchRequest{Limit: 15, Offset: 10, Query: "", OrderField: "", OrderBy: OrderByAsc},
			IsError: false,
		},
		{
			name:    "Мария с offset > limit",
			req:     &SearchRequest{Limit: 5, Offset: 10, Query: "Мария", OrderField: "Name", OrderBy: OrderByAsc},
			IsError: false,
		},
	}

	for caseNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqURL := fmt.Sprintf(
				"%s/search?query=%s&orderField=%s&orderBy=%d&limit=%d&offset=%d",
				server.URL,
				url.QueryEscape(tt.req.Query),
				url.QueryEscape(tt.req.OrderField),
				tt.req.OrderBy,
				tt.req.Limit,
				tt.req.Offset,
			)

			resp, err := http.Get(reqURL)

			if err != nil {
				t.Errorf("unexpected error: %d %v", caseNum, err)
			}
			defer func() {
				if err = resp.Body.Close(); err != nil {
					t.Errorf("close body error: %v", err)
				}
			}()

			var users []User
			if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
				t.Errorf("decode error: %v", err)
			}

			if tt.IsError {
				// ожидаем, что вернется ошибка (например 400)
				if resp.StatusCode == http.StatusOK {
					t.Errorf("%s: expected error, got 200 OK", tt.name)
				}
				return // дальше JSON не читаем
			}
		})
	}
}
