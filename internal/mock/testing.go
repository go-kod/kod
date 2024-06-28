package mock

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockTestingT struct {
	FailNowCalled bool
	testing.TB
}

var _ testing.TB = (*mockTestingT)(nil)

func (m *mockTestingT) FailNow() {
	// register the method is called
	m.FailNowCalled = true
	// exit, as normal behavior
	runtime.Goexit()
}

func ExpectFailure(tb testing.TB, fn func(tt testing.TB)) {
	var wg sync.WaitGroup

	// create a mock structure for TestingT
	mockT := &mockTestingT{TB: tb}
	// setup the barrier
	wg.Add(1)
	// start a co-routine to execute the test function f
	// and release the barrier at its end
	go func() {
		defer wg.Done()
		fn(mockT)
	}()

	// wait for the barrier.
	wg.Wait()
	// verify fail now is invoked
	require.True(tb, mockT.FailNowCalled)
}
