package amap

import (
	"sync"
	"time"
)

type AMaper interface {
	Set(key string, val interface{})
	Get(key string, timeout time.Duration) interface{}
}

type Listener struct {
	id  int64
	key string
	out (chan interface{})
}

type AMap struct {
	*sync.Map

	listeners map[string]map[int64]Listener
}

func (m AMap) Set(key string, val interface{}) {
	for _, v := range m.listeners[key] {
		v.out <- val
		m.unsubscribe(v)
	}
	m.Store(key, val)
}

func (m AMap) subscribe(listener Listener) {
	_, ok := m.listeners[listener.key]
	if !ok {
		m.listeners[listener.key] = map[int64]Listener{}
	}
	_, ok = m.listeners[listener.key][listener.id]
	if ok {
		return
	}
	m.listeners[listener.key][listener.id] = listener
}

func (m AMap) unsubscribe(listener Listener) {
	close(listener.out)
	delete(m.listeners[listener.key], listener.id)
}

func (m AMap) Get(key string, timeout time.Duration) interface{} {
	v, ok := m.Load(key)
	if ok {
		return v
	}

	timestamp := time.Now().UnixNano()
	for {
		_, ok := m.listeners[key][timestamp]
		if ok {
			timestamp += 1
		} else {
			break
		}
	}

	listener := Listener{
		id:  timestamp,
		key: key,
		out: make(chan interface{}),
	}
	go func() {
		time.AfterFunc(timeout, func() {
			m.unsubscribe(listener)
		})
	}()

	m.subscribe(listener)

	val, ok := <-listener.out
	if ok {
		return val
	}

	return nil
}

var (
	_ AMaper = AMap{}
)
