package mapx

var _ MapInterface = &HashMap{}

type HashMap map[interface{}]interface{}

func MakeHashMap(size int) HashMap {
	return make(HashMap, size)
}
func (m HashMap) Load(key interface{}) (value interface{}, exist bool) {
	value, exist = m[key]
	return
}

func (m HashMap) Store(key, value interface{}) {
	m[key] = value
}

func (m HashMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, exist := m.Load(key)
	if exist {
		return actual, true
	}

	m.Store(key, value)
	return value, false
}

func (m HashMap) Delete(key interface{}) {
	delete(m, key)
}

func (m HashMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	for k, v := range m {
		shouldContinue := f(k, v)
		if !shouldContinue {
			break
		}
	}
}

func (m HashMap) Len() (length int) {
	return len(m)
}
