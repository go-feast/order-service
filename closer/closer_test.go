package closer

import "testing"

type testError string

func (t testError) Error() string { return string(t) }

func TestCloseFunc_Close(t *testing.T) {
	var err testError = "test error"

	noError := func() error { return nil }
	testErr := func() error { return err }

	tests := []struct {
		f           CloseFunc
		expectedErr error
		name        string
		wantErr     bool
	}{
		{
			name:    "nil",
			f:       noError,
			wantErr: false,
		},
		{
			name:        "err",
			f:           testErr,
			expectedErr: err,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
