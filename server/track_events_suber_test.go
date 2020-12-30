package server

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestTrackEventsSuber(t *testing.T) {
	defer goleak.VerifyNone(t)

	in := make(chan TrackEvent)

	var closeOnce sync.Once

	closeInput := func() {
		closeOnce.Do(func() {
			close(in)
		})
	}

	defer closeInput()

	s := newTrackEventsSuber(in)

	sub1, err := s.Subscribe("a")
	assert.NoError(t, err)
	assert.NotNil(t, sub1)

	sub2, err := s.Subscribe("a")
	assert.NoError(t, err)
	assert.NotNil(t, sub1)

	select {
	case _, ok := <-sub1:
		assert.False(t, ok, "sub1 should have been closed and replaced")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for sub1 to close")
	}

	sub3, err := s.Subscribe("b")
	assert.NoError(t, err)
	assert.NotNil(t, sub1)

	ev := TrackEvent{
		TrackInfo: TrackInfo{},
		Type:      TrackEventTypeAdd,
	}

	select {
	case in <- ev:
	case <-time.After(time.Second):
		assert.Fail(t, "timed out while sending event")
	}

	select {
	case recv, ok := <-sub2:
		assert.Equal(t, ev, recv)
		assert.True(t, ok, "sub2 should not have been closed yet")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for sub2 event")
	}

	select {
	case recv, ok := <-sub3:
		assert.Equal(t, ev, recv)
		assert.True(t, ok, "sub3 should not have been closed yet")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for sub3 event")
	}

	err = s.Unsubscribe("b")
	assert.NoError(t, err)

	select {
	case _, ok := <-sub3:
		assert.False(t, ok, "sub3 should have been closed")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for sub3 to close")
	}

	closeInput()

	select {
	case _, ok := <-sub2:
		assert.False(t, ok, "sub2 should have been closed")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for sub2 to close")
	}

	sub4, err := s.Subscribe("c")
	assert.Nil(t, sub4)
	assert.Equal(t, io.ErrClosedPipe, errors.Cause(err))
}