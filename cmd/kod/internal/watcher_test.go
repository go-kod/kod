package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=watcher.go -destination=mock_watcher_test.go -package=internal

func TestWatcherNormal(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "test",
		Op:   fsnotify.Write,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func() {})
}

func TestWatcherInvalidEvents(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "test",
		Op:   1000,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func() {})
}

func TestWatcherErrors(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan error, 2)
	events <- fsnotify.ErrEventOverflow
	w.EXPECT().Events().Return(nil).AnyTimes()
	w.EXPECT().Errors().Return(events).AnyTimes()

	Watch(w, ".", func() {})
}

func TestWatcherAddFail(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(errors.New("error"))

	events := make(chan error, 2)
	events <- fsnotify.ErrEventOverflow
	w.EXPECT().Events().Return(nil).AnyTimes()
	w.EXPECT().Errors().Return(events).AnyTimes()

	Watch(w, ".", func() {})
}

func TestWatcherErrorsClose(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan error, 2)
	w.EXPECT().Events().Return(nil).AnyTimes()
	w.EXPECT().Errors().Return(events).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func() {})
}

func TestWatcherFilterGenPath(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "kod_gen",
		Op:   fsnotify.Write,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func() {})
}
