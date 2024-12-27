package interceptor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIf(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name      string
		condition Condition
		expected  bool
	}{
		{
			name: "Condition is true",
			condition: func(ctx context.Context, info CallInfo) bool {
				// TODO: Implement condition logic
				return true
			},
			expected: false,
		},
		{
			name: "Condition is false",
			condition: func(ctx context.Context, info CallInfo) bool {
				// TODO: Implement condition logic
				return false
			},
			expected: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock invoker function
			invoker := func(ctx context.Context, info CallInfo, req, reply []interface{}) error {
				// TODO: Implement invoker logic
				return nil
			}

			// Create a mock interceptor function
			interceptor := func(ctx context.Context, info CallInfo, req, reply []interface{}, invoker HandleFunc) error {
				// TODO: Implement interceptor logic
				return errors.New("not implemented")
			}

			// Call the If function with the test case inputs
			err := If(interceptor, tc.condition)(context.Background(), CallInfo{}, []interface{}{}, []interface{}{}, invoker)

			// Check if the result matches the expected value
			if (err == nil) != tc.expected {
				t.Errorf("Expected error: %v, but got: %v", tc.expected, err)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name       string
		conditions []Condition
		expected   bool
	}{
		{
			name: "All conditions are true",
			conditions: []Condition{
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
			},
			expected: false,
		},
		{
			name: "At least one condition is false",
			conditions: []Condition{
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
			},
			expected: true,
		},
		{
			name: "At least one condition is false",
			conditions: []Condition{
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
			},
			expected: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock invoker function
			invoker := func(ctx context.Context, info CallInfo, req, reply []interface{}) error {
				// TODO: Implement invoker logic
				return nil
			}

			// Create a mock interceptor function
			interceptor := func(ctx context.Context, info CallInfo, req, reply []interface{}, invoker HandleFunc) error {
				// TODO: Implement interceptor logic
				return errors.New("not implemented")
			}

			// Create the combined condition using And function
			combinedCondition := And(tc.conditions[0], tc.conditions[1], tc.conditions[2:]...)

			// Call the If function with the combined condition
			err := If(interceptor, combinedCondition)(context.Background(), CallInfo{}, []interface{}{}, []interface{}{}, invoker)

			// Check if the result matches the expected value
			if (err == nil) != tc.expected {
				t.Errorf("Expected error: %v, but got: %v", tc.expected, err)
			}
		})
	}
}

func TestOr(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name       string
		conditions []Condition
		expected   bool
	}{
		{
			name: "At least one condition is true",
			conditions: []Condition{
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
			},
			expected: false,
		},
		{
			name: "All conditions are false",
			conditions: []Condition{
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
			},
			expected: true,
		},
		{
			name: "The last condition is false",
			conditions: []Condition{
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
				func(ctx context.Context, info CallInfo) bool {
					return false
				},
				func(ctx context.Context, info CallInfo) bool {
					return true
				},
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock invoker function
			invoker := func(ctx context.Context, info CallInfo, req, reply []interface{}) error {
				// TODO: Implement invoker logic
				return nil
			}

			// Create a mock interceptor function
			interceptor := func(ctx context.Context, info CallInfo, req, reply []interface{}, invoker HandleFunc) error {
				// TODO: Implement interceptor logic
				return errors.New("not implemented")
			}

			// Create the combined condition using Or function
			combinedCondition := Or(tc.conditions[0], tc.conditions[1], tc.conditions[2:]...)

			// Call the If function with the combined condition
			err := If(interceptor, combinedCondition)(context.Background(), CallInfo{}, []interface{}{}, []interface{}{}, invoker)

			// Check if the result matches the expected value
			if (err == nil) != tc.expected {
				t.Errorf("Expected error: %v, but got: %v", tc.expected, err)
			}
		})
	}
}

func TestNot(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name      string
		condition Condition
		expected  bool
	}{
		{
			name: "Condition is true",
			condition: func(ctx context.Context, info CallInfo) bool {
				return true
			},
			expected: false,
		},
		{
			name: "Condition is false",
			condition: func(ctx context.Context, info CallInfo) bool {
				return false
			},
			expected: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the Not function with the test case input
			result := Not(tc.condition)(context.Background(), CallInfo{})

			// Check if the result matches the expected value
			if result != tc.expected {
				t.Errorf("Expected result: %v, but got: %v", tc.expected, result)
			}
		})
	}
}

func TestIsMethod(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		method   string
		info     CallInfo
		expected bool
	}{
		{
			name:   "Method matches",
			method: "test",
			info: CallInfo{
				FullMethod: "test",
			},
			expected: true,
		},
		{
			name:   "Method does not match",
			method: "test",
			info: CallInfo{
				FullMethod: "other",
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the IsMethod function with the test case input
			result := IsMethod(tc.method)(context.Background(), tc.info)

			// Check if the result matches the expected value
			if result != tc.expected {
				t.Errorf("Expected result: %v, but got: %v", tc.expected, result)
			}
		})
	}
}

func TestSingleton(t *testing.T) {
	name := "FullMethod"

	initCount := 0
	callCount := 0

	// Create a mock invoker function
	initFn := func() Interceptor {
		initCount++
		return func(ctx context.Context, info CallInfo, req, reply []interface{}, invoker HandleFunc) error {
			callCount++
			return nil
		}
	}

	// Create the singleton interceptor using the test case input
	singletonInterceptor := SingletonByFullMethod(initFn)

	err := singletonInterceptor(context.Background(), CallInfo{
		FullMethod: name,
	}, nil, nil, nil)

	require.Nil(t, err)
	require.Equal(t, 1, initCount)
	require.Equal(t, 1, callCount)

	err = singletonInterceptor(context.Background(), CallInfo{
		FullMethod: name,
	}, nil, nil, nil)

	require.Nil(t, err)
	require.Equal(t, 1, initCount)
	require.Equal(t, 2, callCount)

	err = singletonInterceptor(context.Background(), CallInfo{
		FullMethod: "no cache",
	}, nil, nil, nil)

	require.Nil(t, err)
	require.Equal(t, 2, initCount)
	require.Equal(t, 3, callCount)
}
