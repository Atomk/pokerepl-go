package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache(time.Duration(5 * time.Second))
	var val []byte
	var ok bool

	// Get nonexistent key
	val, ok = cache.Get("example")
	if val != nil || ok != false {
		t.Errorf("actual %v/%v, expected <nil>/false", val, ok)
	}

	key := "example"
	expected := []byte("https://example.com")
	cache.Add(key, expected)
	val, ok = cache.Get(key)
	if ok != true {
		t.Errorf("entry exists but `ok` is false")
	}
	if len(val) != len(expected) {
		t.Errorf("expected %v, got %v", expected, val)
	}
	for i := range len(val) {
		if val[i] != expected[i] {
			t.Errorf("expected %v, got %v (difference at index %d)", expected, val, i)
		}
	}

	// Add a nil value
	key = "nulltest"
	cache.Add(key, nil)
	val, ok = cache.Get(key)
	if val != nil || ok != true {
		t.Errorf("actual %v/%v, expected <nil>/true", val, ok)
	}
}

func TestCacheReapLoop(t *testing.T) {
	reapInterval := time.Duration(300 * time.Millisecond)
	cache := NewCache(reapInterval)

	cache.Add("a", []byte("Apple"))
	time.Sleep(80 * time.Millisecond)
	cache.Add("b", []byte("Banana"))
	time.Sleep(80 * time.Millisecond)
	cache.Add("c", []byte("Coconut"))
	time.Sleep(80 * time.Millisecond)
	// Time passed: ~240 ms

	_, ok := cache.Get("a")
	if !ok {
		t.Errorf("first entry deleted before scheduled time")
		return
	}

	time.Sleep(80 * time.Millisecond)
	// The cache reap should trigger after 300 ms, and we waited in total ~320 ms,
	// the first entry added to the cache is now expired and should have been removed.
	_, okA := cache.Get("a")
	_, okB := cache.Get("b")
	_, okC := cache.Get("c")
	if okA {
		t.Errorf("first entry not deleted within scheduled time")
		return
	}
	if !okB || !okC {
		t.Errorf("entries `b` or `c` deleted before scheduled time")
		return
	}
}
