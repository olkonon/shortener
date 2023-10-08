package memory

import "testing"

func TestInMemory_GetURLByID(t *testing.T) {
	type fields struct {
		storeByID  map[string]string
		storeByURL map[string]string
	}
	type args struct {
		ID string
	}
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
				storeByID:  map[string]string{"fwrefw": "https://test.com"},
				storeByURL: map[string]string{"https://test.com": "fwrefw"},
			},
			args: struct{ ID string }{ID: "fwrefw"},
			want: "https://test.com",
		},
		{
			name: "record not exist",
			fields: fields{
				storeByID:  map[string]string{"fwrefw": "https://test.com"},
				storeByURL: map[string]string{"https://test.com": "fwrefw"},
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
				storeByID:  test.fields.storeByID,
				storeByURL: test.fields.storeByURL,
			}
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "generate from existed URL",
			fields: struct {
				storeByID  map[string]string
				storeByURL map[string]string
			}{
				storeByID:  map[string]string{"fwrefw": "https://test.com"},
				storeByURL: map[string]string{"https://test.com": "fwrefw"},
			},
			args:    struct{ url string }{url: "https://test.com"},
			want:    "fwrefw",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ims := &InMemory{
				storeByID:  test.fields.storeByID,
				storeByURL: test.fields.storeByURL,
			}
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
