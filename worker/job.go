package worker

import (
	"reflect"
)

type Job struct {
	name string
	fn   interface{}
	args []interface{}
}

func (this *Job) Call() (err error) {
	switch f := this.fn.(type) {

	case func() error:
		return f()

	case func():
		f()
		return nil

	default: // call by reflect

		vf := reflect.ValueOf(this.fn)
		if vf.Kind() == reflect.Func {
			out := callFunc(this.fn, this.args)

			if len(out) < 1 {
				return
			}

			if err, ok := out[0].(error); ok {
				return err
			}

			// ignore other outputs
			return
		}
	}

	return
}
