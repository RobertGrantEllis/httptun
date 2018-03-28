package portreg

import (
	"sync"
)

type PortRegistry interface {
	Allocate() int
	Release(int)
}

type portRegistry struct {
	allocated   map[int]bool
	unallocated chan int
	mutex       *sync.Mutex
}

func New(min, max int) PortRegistry {

	if min > max {
		min, max = max, min
	}

	num := max - min + 1

	allocated := make(map[int]bool, num)
	unallocated := make(chan int, num)

	for i := min; i <= max; i++ {
		unallocated <- i
	}

	return &portRegistry{
		allocated:   allocated,
		unallocated: unallocated,
		mutex:       &sync.Mutex{},
	}
}

func (pr *portRegistry) Allocate() int {

	return <-pr.unallocated
}

func (pr *portRegistry) Release(port int) {

	pr.mutex.Lock()
	if pr.allocated[port] {
		delete(pr.allocated, port)
		pr.unallocated <- port
	}
	pr.mutex.Unlock()
}
