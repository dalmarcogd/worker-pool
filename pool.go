package worker_pool

import (
	"fmt"
	"reflect"
)

type Pool interface {
	Delegate(args ...interface{})
	Start(maxWorkers int, fn interface{}) error
}

type pool struct {
	queue chan []reflect.Value
}

func (p *pool) Delegate(args ...interface{}) {

	fArgs := make([]reflect.Value, 0)
	for _, arg := range args {
		fArgs = append(fArgs, reflect.ValueOf(arg))
	}

	go func() {
		p.queue <- fArgs
	}()
}

func (p *pool) Start(maxWorkers int, fn interface{}) error {
	if maxWorkers < 1 {
		return fmt.Errorf("Invalid number of workers: %d", maxWorkers)
	}

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	for i := 1; i <= maxWorkers; i++ {
		h := reflect.ValueOf(fn)

		go func() {
			for args := range p.queue {
				h.Call(args)
			}
		}()
	}

	return nil
}

func New(queueLength int) Pool {
	return &pool{
		queue: make(chan []reflect.Value, queueLength),
	}
}
