package priorityq

import "testing"

func BenchmarkPQ(b *testing.B) {
	type Node struct {
		prioriy int
		value   int
	}

	pq := NewPriorityQueue(
		WithLessFunc(func(i1, i2 interface{}) bool {
			n1 := i1.(Node)
			n2 := i2.(Node)

			return n1.prioriy > n2.prioriy
		}),
		WithCapacity(10),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pq.Add(Node{value: i, prioriy: i})
		_, _, _ = pq.Peek()
		_, _ = pq.Pop()
	}
}
