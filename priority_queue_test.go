package priorityq

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPriorityQueue(t *testing.T) {
	const (
		PriorityQueueCapacity = 10
	)

	type Node struct {
		ID  int
		Val string
	}

	Convey("测试1：上溢和下溢", t, func() {
		pq := NewPriorityQueue(
			WithCapacity(PriorityQueueCapacity),
			WithLessFunc(func(v1, v2 interface{}) bool {
				n1 := v1.(*Node)
				n2 := v2.(*Node)

				return n1.ID > n2.ID
			}),
		)

		Convey("上溢", func() {
			for i := 0; i < PriorityQueueCapacity; i++ {
				priority := 1
				if i > PriorityQueueCapacity/4 {
					priority = 2
				} else if i > PriorityQueueCapacity/2 {
					priority = 3
				}

				pq.Add(&Node{ID: i, Val: fmt.Sprintf("node%d", i)}, priority)
			}

			err := pq.Add(&Node{
				ID:  PriorityQueueCapacity,
				Val: fmt.Sprintf("node%d", PriorityQueueCapacity),
			})

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "overflow")
		})

		Convey("下溢", func() {
			_, err := pq.Pop()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "underflow")

			_, _, err = pq.Peek()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "underflow")
		})
	})

	Convey("测试2：添加、更新、删除数据", t, func() {
		pq := NewPriorityQueue(
			WithCapacity(1),
		)

		Convey("添加数据", func() {
			pq.Add(Node{ID: 0, Val: "node0"})

			_, pre, _ := pq.Peek()
			So(pre, ShouldEqual, 0)
		})

		Convey("更新数据优先级", func() {
			pq.Add(Node{ID: 0, Val: "node0"})
			_, pre, _ := pq.Peek()
			So(pre, ShouldEqual, 0)

			// 通过ADD更新
			pq.Add(Node{ID: 0, Val: "node0"}, 1)
			_, pre, _ = pq.Peek()
			So(pre, ShouldEqual, 1)

			// 更新接口更新
			pq.UpdatePriority(Node{ID: 0, Val: "node0"}, 2)
			_, pre, _ = pq.Peek()
			So(pre, ShouldEqual, 2)

			// 更新接口更新一个不存在的数据
			err := pq.UpdatePriority(Node{ID: 1, Val: "node1"}, 1)
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "value not found")
		})

		Convey("删除数据", func() {
			pq.Add(Node{ID: 0, Val: "node0"})

			err := pq.Delete(Node{ID: 0, Val: "node0"})
			So(err, ShouldBeEmpty)

			err = pq.Delete(Node{ID: 0, Val: "node0"})
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "value not found")

			// POP删除
			pq.Add(Node{ID: 1, Val: "node1"})

			val, err := pq.Pop()
			So(err, ShouldBeEmpty)
			So(val.(Node).Val, ShouldEqual, "node1")

			_, err = pq.Pop()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "underflow")
		})

	})

}

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
