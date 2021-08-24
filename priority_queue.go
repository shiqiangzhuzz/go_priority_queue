package priorityq

import (
	"container/heap"
	"errors"
	"sync"
)

type PriorityQueueOption func(*PriorityQueue)
type LessFunc func(v1 interface{}, v2 interface{}) bool

var (
	_               = heap.Interface(&store{})
	defaultCapacity = 1024
	defaultLessFunc = func(v1 interface{}, v2 interface{}) bool { return false }
	defaultPriority = 0
)

type PriorityQueue struct {
	s        *store
	lock     sync.RWMutex
	capacity int
}

// 构造函数
func NewPriorityQueue(opts ...PriorityQueueOption) *PriorityQueue {
	pq := &PriorityQueue{
		s: &store{
			less:   defaultLessFunc,
			lookup: map[interface{}]*item{},
		},
		capacity: defaultCapacity,
	}

	for _, opt := range opts {
		opt(pq)
	}

	return pq
}

// 配置
func WithCapacity(capacity int) PriorityQueueOption {
	return func(pq *PriorityQueue) {
		pq.capacity = capacity
	}
}

func WithLessFunc(less LessFunc) PriorityQueueOption {
	return func(pq *PriorityQueue) {
		pq.s.less = less
	}
}

// 业务函数
// 添加节点，节点存在，则更新，不存在则添加
func (pq *PriorityQueue) Add(x interface{}, priorities ...int) error {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	priority := defaultPriority
	if len(priorities) != 0 {
		priority = priorities[0]
	}

	if it, ok := pq.s.lookup[x]; ok {
		it.priority = priority
		heap.Fix(pq.s, it.index)
		return nil
	}

	if len(pq.s.items) >= pq.capacity {
		return errors.New("overflow")
	}

	it := &item{value: x, index: len(pq.s.items), priority: priority}
	heap.Push(pq.s, it)
	return nil
}

// 删除节点，存在就删除，不存在报错
func (pq *PriorityQueue) Delete(x interface{}) error {
	if it, ok := pq.s.lookup[x]; ok {
		_ = heap.Remove(pq.s, it.index)
		return nil
	}

	return errors.New("value not found")
}

// 返回头部元素并从队列中删除
func (pq *PriorityQueue) Pop() (interface{}, error) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	if len(pq.s.items) == 0 {
		return nil, errors.New("underflow")
	}

	it := heap.Pop(pq.s)
	if it == nil {
		return nil, errors.New("value was removed")
	}

	return it.(*item).value, nil
}

// 返回头部元素以及优先级，不删除元素
func (pq *PriorityQueue) Peek() (interface{}, int, error) {
	pq.lock.RLock()
	defer pq.lock.RUnlock()

	if len(pq.s.items) == 0 {
		return nil, -1, errors.New("underflow")
	}

	it := pq.s.items[0]
	return it.value, it.priority, nil
}

// 更新优先级，存在就更新，不存在报错
func (pq *PriorityQueue) UpdatePriority(x interface{}, priority int) error {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	it, ok := pq.s.lookup[x]
	if !ok {
		return errors.New("value not found")
	}

	it.priority = priority
	heap.Fix(pq.s, it.index)
	return nil
}

// heap.Interface 接口实现
type item struct {
	value    interface{}
	priority int
	index    int
}

type store struct {
	items  []*item
	lookup map[interface{}]*item
	less   LessFunc
}

func (s store) Len() int { return len(s.items) }

func (s store) Less(i, j int) bool {
	if s.items[i].priority == s.items[j].priority {
		return s.less(s.items[i].value, s.items[j].value)
	}

	return s.items[i].priority > s.items[j].priority
}

func (s store) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
	s.items[i].index = i
	s.items[j].index = j
}

func (s *store) Push(x interface{}) {
	data := x.(*item)
	s.lookup[data.value] = data
	s.items = append(s.items, data)
}

func (s *store) Pop() interface{} {
	size := len(s.items)
	it := s.items[size-1]

	delete(s.lookup, it.value)

	s.items[size-1].index = -1
	s.items[size-1] = nil
	s.items = s.items[0 : size-1]
	return it
}
