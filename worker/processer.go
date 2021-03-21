package worker

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

type jobs []*Job

//go:generate goption -p . -c processer
type Processer struct {
	jobs jobs

	workers []*Worker

	wait sync.WaitGroup

	// controls
	QueueMutex sync.RWMutex
	started    atomic.Bool

	Result Result
}

func Default() *Processer {
	return &Processer{
		jobs:    make(jobs, 0, 4),
		workers: make([]*Worker, 0, 4),
	}
}

type Result struct {
	errs chan error
}

func (this *Result) NextError() (err error) {
	return <-this.errs
}

// go func(){}
func (this *Processer) Go(jobName string, fn func()) *Processer {
	return this.Call(jobName, fn)
}

// go func(args){}
func (this *Processer) Call(jobName string, fn interface{}, args ...interface{}) *Processer {

	this.QueueMutex.Lock()
	this.jobs = append(this.jobs, &Job{name: jobName, fn: fn})
	this.QueueMutex.Unlock()

	return this
}

func (this *Processer) Run() *Processer {
	this.Start()
	return this.Wait()
}

func (this *Processer) Wait() *Processer {
	this.wait.Wait()

	return this
}

func (this *Processer) Start() *Processer {

	if startOk := this.started.CAS(false, true); !startOk {
		panic(`concurrent start`)
		return this
	}

	this.Result.errs = make(chan error, len(this.jobs))

	for _, job := range this.jobs {

		this.wait.Add(1)

		// todo goroutine pool
		go func(job *Job) {
			var err error

			defer func() {
				this.Result.errs <- err

				this.wait.Done()
			}()

			defer func() {
				if ex := recover(); ex != nil {
					err = fmt.Errorf("%+v", ex)
					err = errors.Wrapf(err, "JobPanic[%+v]", job.name)
				}
			}()

			err = job.Call()
			if err != nil {
				err = errors.Wrapf(err, "JobFailed[%+v]", job.name)
			}

		}(job)
	}

	return this
}

func callFunc(f interface{}, args ...interface{}) (out []interface{}) {
	vf := reflect.ValueOf(f)
	if vf.Kind() != reflect.Func {
		panic(123)
	}

	in := make([]reflect.Value, len(args))
	for idx, arg := range args {
		in[idx] = reflect.ValueOf(arg)
	}

	outv := vf.Call(in)
	for _, v := range outv {
		out = append(out, v.Interface())
	}

	return
}
