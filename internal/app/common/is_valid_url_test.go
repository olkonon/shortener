package common

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func init() {
	logrus.SetOutput(io.Discard)
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		isValidURL bool
	}{
		{
			name:       "Test fail URL #1",
			data:       "example.com",
			isValidURL: false,
		},
		{
			name:       "Test fail URL #2",
			data:       "http::example.com/test",
			isValidURL: false,
		},
		{
			name:       "Test fail URL #3",
			data:       "example.com/test",
			isValidURL: false,
		},
		{
			name:       "Test fail URL #4",
			data:       "http:///test",
			isValidURL: false,
		},
		{
			name:       "Test right URL #1",
			data:       "http://test.com/test",
			isValidURL: true,
		},
		{
			name:       "Test right URL #2",
			data:       "http://test.com/test?v=3",
			isValidURL: true,
		},
	}
	for _, tt := range tests {
		test := tt
		f := func(t *testing.T) {
			assert.Equal(t, test.isValidURL, IsValidURL(test.data))
		}
		t.Run(test.name, f)
	}
}
