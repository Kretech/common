package worker

import (
	"fmt"
	"testing"
)

func TestProcesser_DemoStart(t *testing.T) {
	var err error
	wg := new(Processer)

	result := wg.
		Go(`a func`, func() {}).
		Call(`a func with error`, func() error { return nil }).
		Call(`a func with args`, func(int) error { return nil }, 2).
		Start().
		Wait()

	fmt.Println(result, err)
}

func TestProcesser_Run(t *testing.T) {
	var err error
	wg := Default()

	done := make([]bool, 3)

	result := wg.
		Go(`a func`, func() { done[1] = true }).
		Call(`a func with error`, func() error { done[2] = true; return nil }).
		Call(`a func with args`, func(int) error { done[0] = true; return nil }, 2).
		Run()

	fmt.Println(result, done, err)
}
