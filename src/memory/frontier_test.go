package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/thiagolcmelo/webcrawler/src/memory"
)

func TestMemoryFrontier_Push(t *testing.T) {
	mf := memory.NewMemoryFrontier()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case val := <-mf.Pop():
				if val != "value" {
					t.Errorf("unexpected value %s", val)
				}
				return
			case <-ctx.Done():
				t.Error("push didn't go through")
			}
		}
	}()

	mf.Push("value")
}

func TestMemoryFrontier_Pop(t *testing.T) {
	mf := memory.NewMemoryFrontier()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		expectedValues := []string{"value1", "value2", "value3"}
		receivedValues := map[string]bool{}

		for {
			select {
			case val := <-mf.Pop():
				receivedValues[val] = true
				everythingFound := true
				for _, expected := range expectedValues {
					if _, ok := receivedValues[expected]; !ok {
						everythingFound = false
						break
					}
				}
				if everythingFound {
					return
				}
			case <-ctx.Done():
				t.Error("not all values were popped")
			}
		}
	}()

	mf.Push("value1")
	mf.Push("value2")
	mf.Push("value3")
}
