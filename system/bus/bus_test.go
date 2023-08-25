package bus

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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
	stats, total, err := b.Stat()
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
	stats, total, err = b.Stat()
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
	stats, total, err = b.Stat()
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
	stats, total, err = b.Stat()
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
	stats, total, err := b.Stat()
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
