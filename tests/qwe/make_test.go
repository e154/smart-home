package main

import "testing"

func BenchmarkProcessUsingChannels(b *testing.B) {

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ch := MakeTasks(countTasks)
		b.StartTimer()
		ProcessUsingChannels(ch, countWorkers)
	}
}

func BenchmarkProcessUsingMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ch := MakeTasks(countTasks)
		b.StartTimer()
		ProcessUsingMutex(ch, countWorkers)
	}
}
