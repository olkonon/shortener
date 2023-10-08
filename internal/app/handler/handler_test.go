package handler

import (
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_POST(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		body    string
		baseURL string
		want    want
	}{
		{
			name:    "Test fail URL #1",
			body:    "12324",
			baseURL: common.DefaultBaseURL,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Test fail URL #2",
			body:    "http:h32ogewfrnophgeprge",
			baseURL: common.DefaultBaseURL,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:    "Test right URL #1",
			body:    "http://test.com/test",
			baseURL: "http://example.com",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/rfdsgd",
			},
		},
		{
			name:    "Test right URL #2",
			body:    "http://test.com/test?v=3",
			baseURL: "http://example.com",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/srewfrEW",
			},
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()
			h := New(memory.NewMockStorage(), test.baseURL)
			h.POST(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			userResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.want.body, string(userResult))
		}
		t.Run(test.name, f)
	}
}

func TestHandler_GET(t *testing.T) {
	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Test not exists URL #1",
			url:  "/12324",
			want: want{
				statusCode: 404,
			},
		},
		{
			name: "Test ot exists URL #2",
			url:  "/httph32ogewfrnophgeprge",
			want: want{
				statusCode: 404,
			},
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()
			h := New(memory.NewMockStorage(), common.DefaultBaseURL)
			h.GET(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)

			userResult := result.Header.Get("Location")
			err := result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.want.location, userResult)
		}
		t.Run(test.name, f)
	}
}
