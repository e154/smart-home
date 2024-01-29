// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package bus

import (
	"fmt"
	"go.uber.org/atomic"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTopic(t *testing.T) {

	const topic = "test/topic"

	b := NewTopic(topic)

	var counter = 0
	var wg = sync.WaitGroup{}

	// Test Subscribe
	fn := func(topic string, arg1 string, arg2 string) {
		counter++
		wg.Done()
	}
	wg.Add(1)
	err := b.Subscribe(fn)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Publish
	b.Publish("hello", "world")

	wg.Wait()

	require.Equal(t, counter, 1)

	// ------------------------------------------------------------
	// Test Stat
	stat := b.Stat()
	require.Equal(t, stat.Topic, topic)
	require.Equal(t, stat.Subscribers, 1)

	// ------------------------------------------------------------

	// Test Unsubscribe
	empty, err := b.Unsubscribe(fn)
	if err != nil {
		t.Errorf("Unsubscribe returned an error: %v", err)
	}
	require.Equal(t, true, empty)

	// Test Publish
	b.Publish("hello", "world")

	time.Sleep(time.Second)

	require.Equal(t, 1, counter)

	// ------------------------------------------------------------

	stat = b.Stat()
	require.Equal(t, stat.Topic, topic)
	require.Equal(t, stat.Subscribers, 0)
	// ------------------------------------------------------------

	// Test Subscribe
	fn = func(topic string, arg1 string, arg2 string) {
		counter++
	}
	err = b.Subscribe(fn, false)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	stat = b.Stat()
	require.Equal(t, stat.Subscribers, 1)

	// Test Close
	b.Close()

	stat = b.Stat()
	require.Equal(t, stat.Subscribers, 0)

	// Test Publish
	b.Publish("foo", "bar")

	time.Sleep(time.Second)

	require.Equal(t, 1, counter)

	// ------------------------------------------------------------
	// Test Stat
	stat = b.Stat()
	require.Equal(t, stat.Subscribers, 0)
	// ------------------------------------------------------------

	fn = func(topic string, arg1 string, arg2 string) {
		counter++
	}
	err = b.Subscribe(fn, false)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Close
	b.Close()

	// Test Publish
	b.Publish("hello", "world")

	time.Sleep(time.Second)

	require.Equal(t, 1, counter)

	/// Test Stat
	stat = b.Stat()
	require.Equal(t, stat.Subscribers, 0)

	// ------------------------------------------------------------

	// Test buildHandlerArgs
	args := buildHandlerArgs([]interface{}{topic, "hello", "world"})
	if len(args) != 3 {
		t.Errorf("buildHandlerArgs returned the wrong number of arguments: %v", args)
	}
	if args[0].String() != topic {
		t.Errorf("buildHandlerArgs returned the wrong topic: %v", args[0])
	}
	if args[1].String() != "hello" {
		t.Errorf("buildHandlerArgs returned the wrong arg1: %v", args[1])
	}
	if args[2].String() != "world" {
		t.Errorf("buildHandlerArgs returned the wrong arg2: %v", args[2])
	}

	// Test reflection of buildHandlerArgs
	if reflect.TypeOf(buildHandlerArgs).Kind() != reflect.Func {
		t.Errorf("buildHandlerArgs is not a function")
	}
}

func TestTopic2(t *testing.T) {

	const topic = "test/topic"

	b := NewTopic(topic)

	var counter atomic.Int32
	var wg = sync.WaitGroup{}

	// Test Subscribe
	fn := func(topic string, arg1 string, arg2 string) {
		fmt.Println("fn1")
		counter.Inc()
		wg.Done()
	}

	fn2 := func(topic string, arg1 string, arg2 string) {
		fmt.Println("fn2")
		counter.Inc()
		wg.Done()
	}

	fn3 := func(topic string, arg1 string, arg2 string) {
		fmt.Println("fn3")
		counter.Inc()
		wg.Done()
	}

	wg.Add(3)

	err := b.Subscribe(fn)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}
	err = b.Subscribe(fn2)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}
	err = b.Subscribe(fn3)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Stat
	stat := b.Stat()
	require.Equal(t, 3, stat.Subscribers)

	// Test Publish
	b.Publish("hello", "world")

	wg.Wait()

	require.Equal(t, int32(3), counter.Load())
}

func BenchmarkTopic(b *testing.B) {

	const topic = "test/topic"

	bus := NewTopic(topic)

	var counter atomic.Int32

	// Test Subscribe
	fn := func(topic string, arg1 string, arg2 string) {
		counter.Inc()
	}
	err := bus.Subscribe(fn)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish("hello", "world")
	}

	time.Sleep(time.Second)

	require.Equal(b, int32(b.N), counter.Load())
}
