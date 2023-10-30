package router

import (
	"bytes"
	"encoding/json"
	"github.com/olkonon/shortener/internal/app/api"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/handler"
	"github.com/olkonon/shortener/internal/app/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	logrus.SetOutput(io.Discard)
}

func TestRouter_POST(t *testing.T) {
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
			name: "Test fail URL #1",
			body: "12324",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Test fail URL #2",
			body: "http:h32ogewfrnophgeprge",
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
				body:       "http://example.com/" + memory.MockID2,
			},
		},
		{
			name:    "Test right URL #2",
			body:    "http://test.com/test?v=3",
			baseURL: "http://example.com",
			want: want{
				statusCode: http.StatusCreated,
				body:       "http://example.com/" + memory.MockID1,
			},
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()
			store := memory.NewMockStorage()
			defer func() {
				err := store.Close()
				require.NoError(t, err)
			}()
			r := New(handler.New(store, test.baseURL))
			r.ServeHTTP(w, request)
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

func TestRouter_POST_JSON(t *testing.T) {
	type want struct {
		statusCode int
		body       api.AddURLResponse
		json       bool
	}
	tests := []struct {
		name        string
		body        api.AddURLRequest
		contentType string
		gzip        bool
		baseURL     string
		want        want
	}{
		{
			name:        "Test fail URL #1",
			body:        api.AddURLRequest{URL: "12324"},
			contentType: handler.ContentTypeApplicationJSON,
			want: want{
				statusCode: http.StatusBadRequest,
				json:       false,
			},
		},
		{
			name:        "Test fail URL #2",
			body:        api.AddURLRequest{URL: "http:h32ogewfrnophgeprge"},
			contentType: handler.ContentTypeApplicationJSON,
			want: want{
				statusCode: http.StatusBadRequest,
				json:       false,
			},
		},
		{
			name:        "Test right URL #1",
			body:        api.AddURLRequest{URL: "http://test.com/test"},
			contentType: handler.ContentTypeApplicationJSON,
			baseURL:     "http://example.com",
			want: want{
				statusCode: http.StatusCreated,
				json:       true,
				body:       api.AddURLResponse{Result: "http://example.com/" + memory.MockID2},
			},
		},
		{
			name:        "Test right URL #2",
			body:        api.AddURLRequest{URL: "http://test.com/test?v=3"},
			contentType: handler.ContentTypeApplicationJSON,
			baseURL:     "http://example.com",
			want: want{
				statusCode: http.StatusCreated,
				json:       true,
				body:       api.AddURLResponse{Result: "http://example.com/" + memory.MockID1},
			},
		},
		{
			name:        "Test bad content-type 12",
			body:        api.AddURLRequest{URL: "http://test.com/test?v=3"},
			contentType: "text/plain",
			baseURL:     "http://example.com",
			want: want{
				statusCode: http.StatusBadRequest,
				json:       false,
			},
		},
		{
			name:        "Test GZIP compression #1",
			body:        api.AddURLRequest{URL: "http://test.com/test?v=3"},
			contentType: handler.ContentTypeApplicationJSON,
			gzip:        true,
			baseURL:     "http://example.com",
			want: want{
				statusCode: http.StatusCreated,
				json:       false,
			},
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			reqBody, err := json.Marshal(test.body)
			require.NoError(t, err)
			if test.gzip {
				reqBody, err = common.CompressGzip(reqBody)
				require.NoError(t, err)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(reqBody))
			request.Header.Set(handler.ContentTypeHeader, test.contentType)
			if test.gzip {
				request.Header.Set(handler.ContentEncodingHeader, "gzip")
			}
			w := httptest.NewRecorder()
			store := memory.NewMockStorage()
			defer func() {
				// Кажется бессмысленным усложнением, но кто знает как поменяется реализация в будущем?
				err := store.Close()
				require.NoError(t, err)
			}()
			r := New(handler.New(store, test.baseURL))
			r.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			userResult, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if test.want.json {
				responseData := api.AddURLResponse{}
				err = json.Unmarshal(userResult, &responseData)
				require.NoError(t, err)
				assert.Equal(t, test.want.body, responseData)
			}
		}
		t.Run(test.name, f)
	}
}

func TestRouter_GET(t *testing.T) {
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
		{
			name: "Test right URL #1",
			url:  "/" + memory.MockID2,
			want: want{
				statusCode: 307,
				location:   "http://test.com/test",
			},
		},
		{
			name: "Test right URL #2",
			url:  "/" + memory.MockID1,
			want: want{
				statusCode: 307,
				location:   "http://test.com/test?v=3",
			},
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()
			store := memory.NewMockStorage()
			defer func() {
				err := store.Close()
				require.NoError(t, err)
			}()
			r := New(handler.New(store, common.DefaultBaseURL))
			r.ServeHTTP(w, request)
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

func TestRouter_ServeHTTP(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		reqURL string
		method string
		body   string
		want   want
	}{
		{
			name:   "Test Bad method PUT",
			method: http.MethodPut,
			body:   "12324",
			reqURL: "/",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "Test Bad method OPTION",
			method: http.MethodOptions,
			body:   "12324",
			reqURL: "/",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:   "Test fail URL #1",
			body:   "http:h32ogewfrnophgeprge",
			method: http.MethodPost,
			reqURL: "/",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "Test right URL #1",
			body:   "http://test.com/test",
			method: http.MethodPost,
			reqURL: "/",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:   "Test right URL #2",
			body:   "http://test.com/test?v=3",
			method: http.MethodPost,
			reqURL: "/",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:   "Test right GET URL #1",
			body:   "",
			method: http.MethodGet,
			reqURL: "/" + memory.MockID2,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
			},
		},
		{
			name:   "Test wrong GET URL #1",
			body:   "",
			method: http.MethodGet,
			reqURL: "/rfdsgd34rt43",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.reqURL, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			store := memory.NewMockStorage()
			defer func() {
				err := store.Close()
				require.NoError(t, err)
			}()
			h := handler.New(store, common.DefaultBaseURL)
			r := New(h)
			r.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		}
		t.Run(test.name, f)
	}
}
