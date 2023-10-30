package memory

import (
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/sirupsen/logrus"
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
			got, err := ims.GetURLByID(test.args.ID)
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
			got, err := ims.GenIDByURL(test.args.url)
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
