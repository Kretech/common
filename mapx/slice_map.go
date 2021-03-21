package mapx

var _ MapInterface = &SliceMap{}

type sliceElem struct {
	key   interface{}
	value interface{}

	deleted bool
}

func (e *sliceElem) isAlive() bool {
	return !e.deleted
}

type SliceMap []*sliceElem

func MakeSliceMap(size int) *SliceMap {
	slice := make(SliceMap, 0, size)
	return &slice
}

func (m SliceMap) Load(k interface{}) (v interface{}, exist bool) {
	//for i := len(m) - 1; i >= 0; i-- {
	//	e := m[i]
	for _, e := range m {
		if e.isAlive() && keyEqual(e.key, k) {
			return e.value, true
		}
	}
	return nil, false
}

func (m *SliceMap) Store(key, value interface{}) {
	//for i := len(*m) - 1; i >= 0; i-- {
	//	e := (*m)[i]
	for _, e := range *m {
		if keyEqual(e.key, key) {
			e.value = value
			e.deleted = false
			return
		}
	}

	*m = append(*m, &sliceElem{key, value, false})
}

func (m *SliceMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	for _, e := range *m {
		if keyEqual(e.key, key) {
			if e.isAlive() {
				return e.value, true
			}

			e.value = value
			e.deleted = false
			return e.value, false
		}
	}

	*m = append(*m, &sliceElem{key, value, false})

	return value, false
}

func (m SliceMap) Delete(key interface{}) {
	for _, e := range m {
		if keyEqual(e.key, key) {
			e.deleted = true
			return
		}
	}
}

func (m SliceMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	for _, e := range m {
		if !e.isAlive() {
			continue
		}

		shouldContinue := f(e.key, e.value)
		if !shouldContinue {
			break
		}
	}
}

func (m SliceMap) Len() int {
	return len(m)
}

func keyEqual(a, b interface{}) bool {
	return a == b
}
