package m

import (
	"sync"
	"time"
)

type Maper interface {
	Set(key string, val interface{})
	Get(key string, timeout time.Duration) interface{}
}

type Map struct {
	listeners map[string]map[int64](chan interface{})
	sync.Map
}

func (m *Map) Set(key string, val interface{}) {
	for idx, listener := range m.listeners[key] {
		listener <- val
		m.unsubscribe(key, idx)
	}
	m.Store(key, val)
}

func (m *Map) Get(key string, timeout time.Duration) interface{} {
	v, ok := m.Load(key)
	if ok {
		return v
	}

	idx, out := m.subscribe(key)
	go func() {
		time.AfterFunc(timeout, func() {
			m.unsubscribe(key, idx)
		})
	}()

	val, ok := <-out
	if ok {
		return val
	}

	return nil
}

func (m *Map) subscribe(key string) (int64, chan interface{}) {
	idx := time.Now().UnixNano()
	for {
		_, ok := m.listeners[key][idx]
		if ok {
			idx += 1
		} else {
			break
		}
	}

	if m.listeners == nil {
		m.listeners = map[string]map[int64]chan interface{}{}
	}
	_, ok := m.listeners[key]
	if !ok {
		m.listeners[key] = map[int64]chan interface{}{}
	}
	listener, ok := m.listeners[key][idx]
	if ok {
		return idx, listener
	}
	out := make(chan interface{})
	m.listeners[key][idx] = out
	return idx, out
}

func (m *Map) unsubscribe(key string, id int64) {
	listener, ok := m.listeners[key][id]
	if ok {
		close(listener)
		delete(m.listeners[key], id)
	}
}

var (
	_ Maper = &Map{}
)
