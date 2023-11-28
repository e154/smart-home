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
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBus(t *testing.T) {

	const topic = "test/topic"

	b := NewBus()

	var counter = 0
	var wg = sync.WaitGroup{}

	// Test Subscribe
	fn := func(topic string, arg1 string, arg2 string) {
		counter++
		wg.Done()
	}
	wg.Add(1)
	err := b.Subscribe(topic, fn)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Publish
	b.Publish(topic, "hello", "world")

	wg.Wait()

	require.Equal(t, counter, 1)

	// ------------------------------------------------------------

	// Test Stat
	stats, total, err := b.Stat(context.Background(), 999, 0, "", "")
	if err != nil {
		t.Errorf("Stat returned an error: %v", err)
	}
	if total != 1 {
		t.Errorf("Stat returned a non-zero total: %d", total)
	}
	if len(stats) != 1 {
		t.Errorf("Stat returned a non-empty stats slice: %v", stats)
	}

	require.Equal(t, stats[0].Topic, topic)
	require.Equal(t, stats[0].Subscribers, 1)

	// ------------------------------------------------------------

	// Test Unsubscribe
	err = b.Unsubscribe(topic, fn)
	if err != nil {
		t.Errorf("Unsubscribe returned an error: %v", err)
	}

	// Test Publish
	b.Publish(topic, "hello", "world")

	time.Sleep(time.Second)

	require.Equal(t, counter, 1)

	// ------------------------------------------------------------

	// Test Stat
	stats, total, err = b.Stat(context.Background(), 999, 0, "", "")
	if err != nil {
		t.Errorf("Stat returned an error: %v", err)
	}
	if total != 0 {
		t.Errorf("Stat returned a non-zero total: %d", total)
	}
	if len(stats) != 0 {
		t.Errorf("Stat returned a non-empty stats slice: %v", stats)
	}

	// ------------------------------------------------------------

	wg.Add(1)
	err = b.Subscribe(topic, fn)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Close
	b.CloseTopic(topic)

	// Test Publish
	b.Publish(topic, "hello", "world")

	time.Sleep(time.Second)

	require.Equal(t, counter, 1)

	// Test Stat
	stats, total, err = b.Stat(context.Background(), 999, 0, "", "")
	if err != nil {
		t.Errorf("Stat returned an error: %v", err)
	}
	if total != 0 {
		t.Errorf("Stat returned a non-zero total: %d", total)
	}
	if len(stats) != 0 {
		t.Errorf("Stat returned a non-empty stats slice: %v", stats)
	}

	// ------------------------------------------------------------

	wg.Add(1)
	err = b.Subscribe(topic, fn)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Close
	b.Purge()

	// Test Publish
	b.Publish(topic, "hello", "world")

	time.Sleep(time.Second)

	require.Equal(t, counter, 1)

	// Test Stat
	stats, total, err = b.Stat(context.Background(), 999, 0, "", "")
	if err != nil {
		t.Errorf("Stat returned an error: %v", err)
	}
	if total != 0 {
		t.Errorf("Stat returned a non-zero total: %d", total)
	}
	if len(stats) != 0 {
		t.Errorf("Stat returned a non-empty stats slice: %v", stats)
	}

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

func TestBus2(t *testing.T) {

	const topic = "test/topic"

	b := NewBus()

	var counter int32
	var wg = sync.WaitGroup{}

	// Test Subscribe
	fn := func(topic string, arg1 string, arg2 string) {
		fmt.Println("fn1")
		atomic.AddInt32(&counter, 1)
		wg.Done()
	}

	fn2 := func(topic string, arg1 string, arg2 string) {
		fmt.Println("fn2")
		atomic.AddInt32(&counter, 1)
		wg.Done()
	}

	fn3 := func(topic string, arg1 string, arg2 string) {
		fmt.Println("fn3")
		atomic.AddInt32(&counter, 1)
		wg.Done()
	}

	wg.Add(3)

	err := b.Subscribe(topic, fn)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}
	err = b.Subscribe(topic, fn2)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}
	err = b.Subscribe(topic, fn3)
	if err != nil {
		t.Errorf("Subscribe returned an error: %v", err)
	}

	// Test Stat
	stats, total, err := b.Stat(context.Background(), 999, 0, "", "")
	if err != nil {
		t.Errorf("Stat returned an error: %v", err)
	}
	if total != 1 {
		t.Errorf("Stat returned a non-zero total: %d", total)
	}
	if len(stats) != 1 {
		t.Errorf("Stat returned a non-empty stats slice: %v", stats)
	}

	// Test Publish
	b.Publish(topic, "hello", "world")

	wg.Wait()

	require.Equal(t, int32(3), counter)
}

func TestBus3(t *testing.T) {

	bus := NewBus()

	var counter int32
	var wg = sync.WaitGroup{}

	const n = 10000

	wg.Add(n)
	fn := func(_ string, msg interface{}) {
		//fmt.Println("msg", msg)
		atomic.AddInt32(&counter, 1)
		wg.Done()
	}

	for i := 0; i < n; i++ {
		_ = bus.Subscribe(fmt.Sprintf("foo/bar/%d", i), fn)
	}

	time.Sleep(time.Second)

	stat, total, err := bus.Stat(context.Background(), 999, 0, "", "")
	require.NoError(t, err)
	require.Equal(t, len(stat), n)
	require.Equal(t, total, int64(n))
	require.Equal(t, counter, int32(0))

	for i := 0; i < n; i++ {
		bus.Publish(fmt.Sprintf("foo/bar/%d", i), i)
	}

	time.Sleep(time.Second)

	for i := 0; i < n; i++ {
		_ = bus.Unsubscribe(fmt.Sprintf("foo/bar/%d", i), fn)
	}
	time.Sleep(time.Second)

	stat, total, err = bus.Stat(context.Background(), 999, 0, "", "")
	require.NoError(t, err)
	require.Equal(t, len(stat), 0)
	require.Equal(t, total, int64(0))
	require.Equal(t, counter, int32(n))

	for i := 0; i < n; i++ {
		bus.Publish(fmt.Sprintf("foo/bar/%d", i), i)
	}

	wg.Wait()

	require.Equal(t, counter, int32(n))
}

func TestBus4(t *testing.T) {

	bus := NewBus()

	var counter1 int32
	var counter2 int32
	var wg1 = sync.WaitGroup{}
	var wg2 = sync.WaitGroup{}

	const n = 1

	wg1.Add(n)
	fn1 := func(_ string, msg interface{}) {
		//fmt.Println("msg", msg)
		atomic.AddInt32(&counter1, 1)
		wg1.Done()
	}
	wg2.Add(n)
	fn2 := func(_ string, msg interface{}) {
		//fmt.Println("msg", msg)
		atomic.AddInt32(&counter2, 1)
		wg2.Done()
	}

	for i := 0; i < n; i++ {
		_ = bus.Subscribe(fmt.Sprintf("foo/bar/%d", i), fn1)
		_ = bus.Subscribe(fmt.Sprintf("foo/bar/%d", i), fn2)
	}

	time.Sleep(time.Second)

	stat, total, err := bus.Stat(context.Background(), 999, 0, "", "")
	require.NoError(t, err)
	require.Equal(t, len(stat), n)
	require.Equal(t, total, int64(n))
	require.Equal(t, counter1, int32(0))
	require.Equal(t, counter2, int32(0))

	for i := 0; i < n; i++ {
		bus.Publish(fmt.Sprintf("foo/bar/%d", i), i)
	}

	wg1.Wait()
	wg2.Wait()

	require.Equal(t, counter1, int32(n))
	require.Equal(t, counter2, int32(n))

	for i := 0; i < n; i++ {
		_ = bus.Unsubscribe(fmt.Sprintf("foo/bar/%d", i), fn1)
	}
	time.Sleep(time.Second)

	stat, total, err = bus.Stat(context.Background(), 999, 0, "", "")
	require.NoError(t, err)
	require.Equal(t, len(stat), n)
	require.Equal(t, total, int64(n))

	wg2.Add(n)
	for i := 0; i < n; i++ {
		bus.Publish(fmt.Sprintf("foo/bar/%d", i), i)
	}

	wg2.Wait()

	require.Equal(t, counter1, int32(n))
	require.Equal(t, counter2, int32(n*2))

	stat, total, err = bus.Stat(context.Background(), 999, 0, "", "")
	require.NoError(t, err)
	require.Equal(t, len(stat), n)
	require.Equal(t, total, int64(n))
}

func BenchmarkBus(b *testing.B) {

	const topic = "test/topic"

	bus := NewBus()

	var counter int32

	// Test Subscribe
	fn := func(topic string, arg1 string, arg2 string) {
		atomic.AddInt32(&counter, 1)
	}
	err := bus.Subscribe(topic, fn)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(topic, "hello", "world")
	}

	time.Sleep(time.Second)

	require.Equal(b, int32(b.N), counter)
}
