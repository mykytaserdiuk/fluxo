package pool

import (
	"runtime"
)

type CorePool struct {
	coreCount uint8
	tasks     chan func()
}

func NewCorePool() *CorePool {
	pool := &CorePool{
		coreCount: uint8(runtime.GOMAXPROCS(0)),
		tasks:     make(chan func(), 1024),
	}

	for i := uint8(0); i < pool.coreCount; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *CorePool) New(fn func()) {
	// The logic for the number of pools can be updated
	p.tasks <- fn
}

func (p *CorePool) worker(_ uint8) {
	for task := range p.tasks {
		task()
	}
}
