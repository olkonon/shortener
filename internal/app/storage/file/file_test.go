package file

import (
	"context"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func init() {
	logrus.SetOutput(io.Discard)
}

func TestFileStorage_GetURLByID(t *testing.T) {
	filename := "C8AA7A99-98E3-4D04-AD5D-2ED521F0D027"
	store := NewFileStorage(filename)
	defer func() {
		err := os.Remove(filename)
		require.NoError(t, err)
	}()

	testID1 := common.GenHashedString("https://test.com")
	testID2 := common.GenHashedString("https://test2.com")
	testID3 := common.GenHashedString("https://test3.com")

	err := store.appendToFile(Record{
		URL: "https://test.com",
		ID:  testID1,
	})
	require.NoError(t, err)
	err = store.appendToFile(Record{
		URL: "https://test2.com",
		ID:  testID2,
	})
	require.NoError(t, err)
	err = store.appendToFile(Record{
		URL: "https://test3.com",
		ID:  testID3,
	})
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "Test record exists #1",
			id:      testID1,
			want:    "https://test.com",
			wantErr: false,
		},
		{
			name:    "Test record exists #2",
			id:      testID2,
			want:    "https://test2.com",
			wantErr: false,
		},
		{
			name:    "Test record exists #3",
			id:      testID3,
			want:    "https://test3.com",
			wantErr: false,
		},
		{
			name:    "Test record not exists #1",
			id:      "fd",
			wantErr: true,
		},
		{
			name:    "Test record not exists #2",
			id:      "fd56df",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			fs := NewFileStorage(filename)
			defer func() {
				err := fs.Close()
				require.NoError(t, err)
			}()
			got, err := fs.GetURLByID(context.Background(), test.id)
			if (err != nil) != test.wantErr {
				t.Errorf("GetURLByID() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("GetURLByID() got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestFileStorage_GenIDByURL(t *testing.T) {
	filename := "C8AA7A99-98E3-4D04-AD5D-2ED521F0D027"
	store := NewFileStorage(filename)
	defer func() {
		err := os.Remove(filename)
		require.NoError(t, err)
	}()
	testID1 := common.GenHashedString("https://test.com")
	testID2 := common.GenHashedString("https://test2.com")
	testID3 := common.GenHashedString("https://test3.com")
	err := store.appendToFile(Record{
		URL: "https://test.com",
		ID:  testID1,
	})
	require.NoError(t, err)
	err = store.appendToFile(Record{
		URL: "https://test2.com",
		ID:  testID2,
	})
	require.NoError(t, err)
	err = store.appendToFile(Record{
		URL: "https://test3.com",
		ID:  testID3,
	})
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)

	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "Test generate from existed URL #1",
			url:     "https://test.com",
			want:    testID1,
			wantErr: true,
		},
		{
			name:    "Test generate from existed URL #2",
			url:     "https://test2.com",
			want:    testID2,
			wantErr: true,
		},
		{
			name:    "Test generate from existed URL #3",
			url:     "https://test3.com",
			want:    testID3,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			fs := NewFileStorage(filename)
			defer func() {
				err := fs.Close()
				require.NoError(t, err)
			}()
			got, err := fs.GenIDByURL(context.Background(), test.url)
			if (err != nil) != test.wantErr {
				t.Errorf("GenIDByURL() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("GenIDByURL() got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestFileStorage_BatchSave(t *testing.T) {
	filename := "C8AA7A99-98E3-4D04-AD5D-2ED521F0D027"
	store := NewFileStorage(filename)
	defer func() {
		store.Close()
		err := os.Remove(filename)
		require.NoError(t, err)
	}()
	testURL1 := "https://test.com"
	testID1 := common.GenHashedString(testURL1)
	testURL2 := "https://test2.com"
	testID2 := common.GenHashedString(testURL2)
	testURL3 := "https://test3.com"
	testID3 := common.GenHashedString(testURL3)

	request := []storage.BatchSaveRequest{
		{
			OriginalURL:   testURL1,
			CorrelationID: "1",
		},
		{
			OriginalURL:   testURL2,
			CorrelationID: "2",
		},
		{
			OriginalURL:   testURL3,
			CorrelationID: "3",
		},
	}

	response := []storage.BatchSaveResponse{
		{
			CorrelationID: "1",
			ShortID:       testID1,
		},
		{
			CorrelationID: "2",
			ShortID:       testID2,
		},
		{
			CorrelationID: "3",
			ShortID:       testID3,
		},
	}

	res, err := store.BatchSave(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, res, response)
}
