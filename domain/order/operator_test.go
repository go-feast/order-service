package order

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createOperator(t *testing.T) *StateOperator {
	order, err := NewOrder(
		uuid.NewString(),
		uuid.NewString(),
		uuid.NewString(),
		[]string{
			uuid.NewString(), uuid.NewString(),
		},
		0.0,
		0.0,
	)

	assert.NoError(t, err)

	return NewStateOperator(order)
}

func TestNewStateOperator(t *testing.T) {
	t.Run("assert NewStateOperator sets order", func(t *testing.T) {
		order, err := NewOrder(
			uuid.NewString(),
			uuid.NewString(),
			uuid.NewString(),
			[]string{
				uuid.NewString(), uuid.NewString(),
			},
			0.0,
			0.0,
		)

		assert.NoError(t, err)

		operator := NewStateOperator(order)

		assert.EqualExportedValues(t, operator.o, order)
	})
}

func TestStateOperator_setState(t *testing.T) {
	t.Run("assert setState sets state", func(t *testing.T) {
		operator := createOperator(t)

		actual := &operator.o.state

		assert.Equal(t, Created, *actual)

		operator.setState(Canceled)

		assert.Equal(t, Canceled, *actual)
	})
}

func TestStateOperator_getStateIndexes(t *testing.T) {
	invalidState := State{"invalid.state", nil}

	testCases := []struct { //nolint:govet
		name           string
		operatorState  State
		replacingState State

		expectedOperatorStateIndex  int
		expectedReplacingStateIndex int

		wantErr     bool
		expectedErr error
	}{
		{
			name:                        "OK",
			operatorState:               Created,
			replacingState:              Closed,
			expectedOperatorStateIndex:  0,
			expectedReplacingStateIndex: 9,

			wantErr: false,
		},
		{
			name:                        "invalid replacing state",
			operatorState:               Created,
			replacingState:              invalidState,
			expectedOperatorStateIndex:  -1,
			expectedReplacingStateIndex: -1,

			wantErr:     true,
			expectedErr: ErrNoSuchState,
		},
		{
			name:                        "invalid operator state",
			operatorState:               invalidState,
			replacingState:              Closed,
			expectedOperatorStateIndex:  -1,
			expectedReplacingStateIndex: -1,

			wantErr:     true,
			expectedErr: ErrNoSuchState,
		},
		{
			name:                        "invalid operator and replacing state",
			operatorState:               invalidState,
			replacingState:              invalidState,
			expectedOperatorStateIndex:  -1,
			expectedReplacingStateIndex: -1,

			wantErr:     true,
			expectedErr: ErrNoSuchState,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			operator := createOperator(t)

			operator.o.state = tc.operatorState

			index, state, err := operator.getStateIndexes(tc.replacingState)

			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedOperatorStateIndex, index)
			assert.Equal(t, tc.expectedReplacingStateIndex, state)
		})
	}
}

func TestStateOperator_trySetState(t *testing.T) {
	testCases := []struct { //nolint:govet
		name string

		operatorState  State
		replacingState State

		wantErr bool

		canceled    bool
		expectedErr error
	}{
		{
			name:           "OK",
			operatorState:  Delivering,
			replacingState: Delivered,

			wantErr:  false,
			canceled: true,
		},
		{
			name:           "set state to the closed order",
			operatorState:  Closed,
			replacingState: Delivering,

			wantErr:     true,
			canceled:    false,
			expectedErr: ErrOrderClosed,
		},
		{
			name:           "cancel order",
			operatorState:  Created,
			replacingState: Canceled,

			wantErr:  false,
			canceled: true,
		},
		{
			name:           "canceling canceled order",
			operatorState:  Canceled,
			replacingState: Canceled,

			wantErr:  false,
			canceled: true,
		},
		{
			name:           "canceling closed order",
			operatorState:  Closed,
			replacingState: Canceled,

			wantErr:     true,
			canceled:    false,
			expectedErr: ErrOrderClosed,
		},
		{
			name:           "set past state",
			operatorState:  Finished,
			replacingState: Cooking,

			wantErr:     true,
			canceled:    false,
			expectedErr: ErrCannotSetState,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			operator := createOperator(t)

			operator.o.state = tc.operatorState

			_, canceled, err := operator.trySetState(tc.replacingState)
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.canceled, canceled)
		})
	}
}

func TestStateOperator_CancelOrder(t *testing.T) {
	testCases := []struct { //nolint:govet
		name string

		operatorState  State
		replacingState State

		wantErr bool

		canceled    bool
		expectedErr error
	}{
		{
			name:           "OK",
			operatorState:  Delivering,
			replacingState: Canceled,

			wantErr:  false,
			canceled: true,
		},
		{
			name:           "cancel canceled",
			operatorState:  Canceled,
			replacingState: Canceled,

			wantErr:  false,
			canceled: true,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			operator := createOperator(t)

			operator.o.state = tc.operatorState

			_, canceled, err := operator.CancelOrder()
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.canceled, canceled)
		})
	}
}
