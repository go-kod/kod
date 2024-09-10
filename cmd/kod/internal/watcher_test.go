package internal

import (
	"context"
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
		Name: "root.go",
		Op:   fsnotify.Write,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherNotExistFile(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "noexist.go",
		Op:   fsnotify.Write,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherAddDir(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "../internal",
		Op:   fsnotify.Create,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherRemoveDir(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)
	w.EXPECT().Remove("../internal").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "../internal",
		Op:   fsnotify.Remove,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherNonGofile(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan fsnotify.Event, 2)
	events <- fsnotify.Event{
		Name: "test.txt",
		Op:   fsnotify.Write,
	}

	w.EXPECT().Events().Return(events).AnyTimes()
	w.EXPECT().Errors().Return(nil).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
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
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherErrors(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan error, 2)
	events <- fsnotify.ErrEventOverflow
	w.EXPECT().Events().Return(nil).AnyTimes()
	w.EXPECT().Errors().Return(events).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherAddFail(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(errors.New("error"))

	events := make(chan error, 2)
	events <- fsnotify.ErrEventOverflow
	w.EXPECT().Events().Return(nil).AnyTimes()
	w.EXPECT().Errors().Return(events).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	Watch(w, ".", func(fsnotify.Event) {}, true)
}

func TestWatcherErrorsClose(t *testing.T) {
	t.Parallel()

	w := NewMockWatcher(gomock.NewController(t))
	w.EXPECT().Add(".").Return(nil)

	events := make(chan error, 2)
	w.EXPECT().Events().Return(nil).AnyTimes()
	w.EXPECT().Errors().Return(events).AnyTimes()
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
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
	w.EXPECT().Context().Return(context.Background()).AnyTimes()

	time.AfterFunc(time.Second, func() {
		close(events)
	})

	Watch(w, ".", func(fsnotify.Event) {}, true)
}
