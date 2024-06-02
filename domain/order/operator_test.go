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

func TestStateOperator_trySetState(t *testing.T) {
	testCases := []struct { //nolint:govet
		name string

		operatorState  State
		replacingState State

		wantErr bool

		setted      bool
		expectedErr error
	}{
		{
			name:           "OK",
			operatorState:  Delivering,
			replacingState: Delivered,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "set state to the closed order",
			operatorState:  Closed,
			replacingState: Delivering,

			wantErr:     true,
			setted:      false,
			expectedErr: ErrOrderClosed,
		},
		{
			name:           "cancel order",
			operatorState:  Created,
			replacingState: Canceled,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "canceling setted order",
			operatorState:  Canceled,
			replacingState: Canceled,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "canceling closed order",
			operatorState:  Closed,
			replacingState: Canceled,

			wantErr:     true,
			setted:      false,
			expectedErr: ErrOrderClosed,
		},
		{
			name:           "set past state",
			operatorState:  Finished,
			replacingState: Cooking,

			wantErr:     true,
			setted:      false,
			expectedErr: ErrInvalidState,
		},
		{
			name:           "setting state to the setted order",
			operatorState:  Canceled,
			replacingState: Cooking,

			wantErr:     true,
			setted:      false,
			expectedErr: ErrOrderCanceled,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			operator := createOperator(t)

			operator.o.state = tc.operatorState

			canceled, err := operator.trySetState(tc.replacingState)
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.setted, canceled)
		})
	}
}

func TestStateOperator_OrderAny(t *testing.T) {
	operator := createOperator(t)

	testCases := []struct { //nolint:govet
		name string

		operatorState  State
		replacingState State

		changingStateFunc func() (bool, error)

		wantErr bool

		setted      bool
		expectedErr error
	}{
		{
			name:           "from created to paid",
			operatorState:  Created,
			replacingState: Paid,

			changingStateFunc: operator.PayOrder,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from paid to cooking",
			operatorState:  Paid,
			replacingState: Cooking,

			changingStateFunc: operator.CookOrder,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from cooking to finished",
			operatorState:  Cooking,
			replacingState: Finished,

			changingStateFunc: operator.OrderFinished,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from finished to waiting for courier",
			operatorState:  Finished,
			replacingState: WaitingForCourier,

			changingStateFunc: operator.WaitForCourier,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from waiting for courier to courier took",
			operatorState:  WaitingForCourier,
			replacingState: CourierTook,

			changingStateFunc: operator.CourierTookOrder,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from courier took to delivering",
			operatorState:  CourierTook,
			replacingState: Delivering,

			changingStateFunc: operator.DeliveringOrder,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from delivering to delivered",
			operatorState:  Delivering,
			replacingState: Delivered,

			changingStateFunc: operator.OrderDelivered,

			wantErr: false,
			setted:  true,
		},
		{
			name:           "from delivered to closed",
			operatorState:  Delivered,
			replacingState: Closed,

			changingStateFunc: operator.CloseOrder,

			wantErr: false,
			setted:  true,
		},
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			operator.o.state = tc.operatorState

			setted, err := tc.changingStateFunc()
			if tc.wantErr {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.setted, setted)
		})
	}
}

func TestStateOperator_FromAnyStateToCanceled(t *testing.T) {
	operator := createOperator(t)
	states := []State{
		Created,
		Paid,
		Cooking,
		Finished,
		WaitingForCourier,
		CourierTook,
		Delivering,
		Delivered,
	}

	for _, state := range states {
		operator.o.state = state

		canceled, err := operator.CancelOrder()

		assert.True(t, canceled)
		assert.NoError(t, err)
	}
}

func TestStateOperator_FromAnyStateToClosed(t *testing.T) {
	operator := createOperator(t)
	states := []State{
		Created,
		Paid,
		Cooking,
		Finished,
		WaitingForCourier,
		CourierTook,
		Delivering,
		Delivered,
	}

	for _, state := range states {
		operator.o.state = state

		canceled, err := operator.CloseOrder()

		assert.True(t, canceled)
		assert.NoError(t, err)
	}
}
