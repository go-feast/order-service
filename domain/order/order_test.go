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

func TestMeals(t *testing.T) {
	var (
		validID1   = "8339957e-dd83-4754-bff4-3ec08de40ed9"
		validID2   = "3150bb27-d728-42bd-9676-b567bf6053d9"
		validID3   = "f04c9e40-18ed-4d39-93f1-c5149354759c"
		invalidID1 = "f04c9e40-18ed-4d39-93f1-c5149354759"
		invalidID2 = "1"
		invalidID3 = "asd"
	)

	tests := []struct {
		name     string
		ids      []string
		expected []MealID
		errCount int
	}{
		{
			name: "all valid IDs",
			ids: []string{
				validID1,
				validID2,
				validID3,
			},
			expected: []MealID{
				MealID(validID1),
				MealID(validID2),
				MealID(validID3),
			},
			errCount: 0,
		},
		{
			name: "some invalid IDs",
			ids: []string{
				validID1,
				invalidID1,
				validID1,
			},
			expected: nil,
			errCount: 1,
		},
		{
			name: "all invalid IDs",
			ids: []string{
				invalidID1,
				invalidID2,
				invalidID3,
			},
			expected: nil,
			errCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mealIDs, errs := Meals(tt.ids)

			if tt.errCount == 0 {
				assert.Nil(t, errs)
				assert.ElementsMatch(t, tt.expected, mealIDs)
			} else {
				assert.Nil(t, mealIDs)
				assert.Len(t, errs, tt.errCount)
			}
		})
	}
}
