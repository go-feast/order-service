package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewID(t *testing.T) {
	var (
		correctUUID         = uuid.NewString()
		incorrectLengthUUID = "cf9c3dc5-281b-478b-93f4-a46c99f0c22"
	)

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		want    ID
		wantErr bool
	}{
		{
			name:    "correct uuid",
			args:    args{id: correctUUID},
			want:    ID(correctUUID),
			wantErr: false,
		},
		{
			name:    "incorrect length of uuid",
			args:    args{id: incorrectLengthUUID},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewID(tt.args.id)

			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
