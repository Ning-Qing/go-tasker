package gotasker

import (
	"fmt"
)

type Worker interface {
	Register(name string, builder WorkBuilder) error
}

type WorkBuilder func() Work

type worker struct {
	builders map[string]WorkBuilder
}

func NewWorker() Worker {
	w := &worker{
		builders: make(map[string]WorkBuilder, 0),
	}
	return w
}

func (w *worker) Register(name string, builder WorkBuilder) error {
	if _, ok := w.builders[name]; ok {
		return fmt.Errorf("work [%s] already exist", name)
	}
	w.builders[name] = builder
	return nil
}
