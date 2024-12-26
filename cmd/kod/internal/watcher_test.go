package internal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/require"
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

func TestWatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("watch directory", func(t *testing.T) {
		dir := t.TempDir()
		mockWatcher := NewMockWatcher(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mockWatcher.EXPECT().Context().Return(ctx).AnyTimes()
		mockWatcher.EXPECT().Add(dir).Return(nil)
		mockWatcher.EXPECT().Events().Return(make(chan fsnotify.Event)).AnyTimes()
		mockWatcher.EXPECT().Errors().Return(make(chan error)).AnyTimes()

		var called bool
		Watch(mockWatcher, dir, func(event fsnotify.Event) {
			called = true
		}, true)

		require.False(t, called) // No events happened
	})

	t.Run("watch file changes", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "test.go")
		require.NoError(t, os.WriteFile(file, []byte("package test"), 0o644))

		mockWatcher := NewMockWatcher(ctrl)
		events := make(chan fsnotify.Event, 1)
		errors := make(chan error)

		events <- fsnotify.Event{
			Name: file,
			Op:   fsnotify.Write,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mockWatcher.EXPECT().Context().Return(ctx).AnyTimes()
		mockWatcher.EXPECT().Add(dir).Return(nil)
		mockWatcher.EXPECT().Events().Return(events).AnyTimes()
		mockWatcher.EXPECT().Errors().Return(errors).AnyTimes()

		var eventReceived fsnotify.Event
		Watch(mockWatcher, dir, func(event fsnotify.Event) {
			eventReceived = event
		}, true)

		time.Sleep(100 * time.Millisecond)
		require.Equal(t, file, eventReceived.Name)
		require.Equal(t, fsnotify.Write, eventReceived.Op)
	})

	t.Run("watch errors", func(t *testing.T) {
		mockWatcher := NewMockWatcher(ctrl)
		events := make(chan fsnotify.Event)
		errors := make(chan error, 1)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mockWatcher.EXPECT().Context().Return(ctx).AnyTimes()
		mockWatcher.EXPECT().Events().Return(events).AnyTimes()
		mockWatcher.EXPECT().Errors().Return(errors).AnyTimes()

		Watch(mockWatcher, "invalid", func(event fsnotify.Event) {}, true)
	})
}
