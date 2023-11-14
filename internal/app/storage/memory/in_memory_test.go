package memory

import (
	"context"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func init() {
	logrus.SetOutput(io.Discard)
}

func TestInMemory_GetURLByID(t *testing.T) {
	type fields struct {
		storeByID map[string]string
	}
	type args struct {
		ID string
	}

	testURL := "https://test.com"
	testID := common.GenHashedString(testURL)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "record exist",
			fields: fields{
				storeByID: map[string]string{testID: testURL},
			},
			args: struct{ ID string }{ID: testID},
			want: testURL,
		},
		{
			name: "record not exist",
			fields: fields{
				storeByID: map[string]string{testID: testURL},
			},
			args:    struct{ ID string }{ID: "fwrefw3"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ims := &InMemory{
				storeByID: test.fields.storeByID,
			}
			defer func() {
				err := ims.Close()
				require.NoError(t, err)
			}()
			got, err := ims.GetURLByID(context.Background(), test.args.ID)
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

func TestInMemory_GenIDByURL(t *testing.T) {
	type fields struct {
		storeByID  map[string]string
		storeByURL map[string]string
	}
	type args struct {
		url string
	}

	testURL := "https://test.com"
	testID := common.GenHashedString(testURL)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "generate from existed URL",
			fields: fields{
				storeByID: map[string]string{testID: testURL},
			},
			args:    struct{ url string }{url: testURL},
			want:    testID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ims := &InMemory{
				storeByID: test.fields.storeByID,
			}
			defer func() {
				err := ims.Close()
				require.NoError(t, err)
			}()
			got, err := ims.GenIDByURL(context.Background(), test.args.url)
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

func TestInMemory_BatchSave(t *testing.T) {
	ims := NewInMemory()
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

	res, err := ims.BatchSave(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, res, response)
}
