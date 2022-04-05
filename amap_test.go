package amap

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestAmap(t *testing.T) {
	t.Run("set/get", func(t *testing.T) {
		m := AMap{Map: &sync.Map{}, listeners: map[string]map[int64]Listener{}}
		m.Set("k1", "v1")
		v := m.Get("k1", 0)
		if v != "v1" {
			t.Fatalf("should get %s but got %s", "v1", v)
		}
	})

	t.Run("set/get with timeout", func(t *testing.T) {
		m := AMap{Map: &sync.Map{}, listeners: map[string]map[int64]Listener{}}
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
		m := AMap{Map: &sync.Map{}, listeners: map[string]map[int64]Listener{}}
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
}
