package timecache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEntriesFound(t *testing.T) {
	tc := NewTimeCache(time.Minute)

	tc.Add("test")

	if !tc.Has("test") {
		t.Fatal("should have this key")
	}
}

func TestEntriesExpire(t *testing.T) {
	tc := NewTimeCache(time.Second)
	for i := 0; i < 11; i++ {
		tc.Add(fmt.Sprint(i))
		time.Sleep(time.Millisecond * 100)
	}

	if tc.Has(fmt.Sprint(0)) {
		t.Fatal("should have dropped this from the cache already")
	}
}

func TestMapEntryStateRace(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	tc := NewTimeCache(time.Second)

	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			tc.Add(fmt.Sprint(i))
			time.Sleep(time.Millisecond * 10)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 10; i < 20; i++ {
			tc.Add(fmt.Sprint(i))
			time.Sleep(time.Millisecond * 10)

		}
		wg.Done()
	}()

	wg.Wait()

	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	default:
	}

	for i := 0; i < 20; i++ {
		if !(tc.Has(fmt.Sprint(i))) {
			t.Fatalf("time cache missing expected element %s", fmt.Sprint(i))
		}
	}

	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	default:
	}
}
