package file

import (
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/sirupsen/logrus"
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
			got, err := fs.GetURLByID(test.id)
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
			wantErr: false,
		},
		{
			name:    "Test generate from existed URL #2",
			url:     "https://test2.com",
			want:    testID2,
			wantErr: false,
		},
		{
			name:    "Test generate from existed URL #3",
			url:     "https://test3.com",
			want:    testID3,
			wantErr: false,
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
			got, err := fs.GenIDByURL(test.url)
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
