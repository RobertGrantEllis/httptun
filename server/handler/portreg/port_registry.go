package portreg

import (
	"sync"

	"github.com/pkg/errors"
)

type PortRegistry interface {
	Allocate() (int, error)
	Release(int)
}

type portRegistry struct {
	allocated   map[int]bool
	unallocated []int
	mutex       *sync.Mutex
}

func New(min, max int) PortRegistry {

	if min > max {
		min, max = max, min
	}

	num := max - min + 1

	allocated := make(map[int]bool, num)
	unallocated := make([]int, num)

	for i := min; i <= max; i++ {
		unallocated = append(unallocated, i)
	}

	return &portRegistry{
		allocated:   allocated,
		unallocated: unallocated,
		mutex:       &sync.Mutex{},
	}
}

func (pr *portRegistry) Allocate() (int, error) {

	var port int

	pr.mutex.Lock()

	if len(pr.unallocated) > 0 {
		port, pr.unallocated = pr.unallocated[0], pr.unallocated[1:]
		pr.allocated[port] = true
	}

	pr.mutex.Unlock()

	if port == 0 {
		return 0, errors.New(`cannot allocate port`)
	}

	return port, nil
}

func (pr *portRegistry) Release(port int) {

	pr.mutex.Lock()
	if pr.allocated[port] {
		delete(pr.allocated, port)
		pr.unallocated = append(pr.unallocated, port)
	}
	pr.mutex.Unlock()
}
