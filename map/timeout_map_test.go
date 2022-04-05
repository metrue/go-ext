package m

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	t.Run("set/get", func(t *testing.T) {
		var m Map
		m.Set("k1", "v1")
		v := m.Get("k1", 0)
		if v != "v1" {
			t.Fatalf("should get %s but got %s", "v1", v)
		}
	})

	t.Run("set/get with timeout", func(t *testing.T) {
		var m Map
		start := time.Now()
		v := m.Get("k1", time.Second*1)
		end := time.Now()
		duration := end.Sub(start)
		if duration < time.Second*1 {
			t.Fatalf("should get %s but got %s", time.Second*1, duration)
		}
		if v != nil {
			t.Fatalf("should get %v but got %s", nil, v)
		}
	})

	t.Run("set/get with value and timeout", func(t *testing.T) {
		var m Map
		var v interface{}

		go func() {
			start := time.Now()
			v = m.Get("k1", time.Second*1)
			end := time.Now()
			duration := end.Sub(start)
			if duration < time.Millisecond*500 {
				log.Fatalf("should get %s but got %s", time.Second*1, duration)
			}
		}()

		time.AfterFunc(time.Millisecond*500, func() {
			m.Set("k1", "v1")
		})

		time.Sleep(500*time.Millisecond + 50*time.Millisecond)

		if v != "v1" {
			t.Fatalf("should get %v but got %s", "v1", v)
		}
	})
	t.Run("map", func(t *testing.T) {
		var m sync.Map
		m.Load("a")
	})

	t.Run("thread safe", func(t *testing.T) {
		var m Map

		go func() {
			for {
				m.Set("key1", "val1")
			}
		}()

		go func() {
			for {
				m.Set("key1", "val2")
			}
		}()

		go func() {
			for {
				m.Set("notfound", "val2")
			}
		}()

		go func() {
			for {
				_ = m.Get("notfound", time.Second*1)
			}
		}()
		time.Sleep(2 * time.Second)
	})
}
