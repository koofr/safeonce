package safeonce_test

import (
	"fmt"
	"testing"

	. "github.com/koofr/safeonce"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, safeOnce *SafeOnce, o *one, c chan bool) {
	safeOnce.Do(func() error { o.Increment(); return nil })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	safeOnce := new(SafeOnce)
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, safeOnce, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("safeOnce failed outside run: %d is not 1", *o)
	}
}

func TestSafeOncePanic(t *testing.T) {
	var safeOnce SafeOnce
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("SafeOnce.Do did not panic")
			}
		}()
		err := safeOnce.Do(func() error {
			panic("failed")
			return nil
		})
		if err != nil {
			t.Fatalf("SafeOnce.Do error should be nil")
		}
	}()

	err := safeOnce.Do(func() error {
		t.Fatalf("SafeOnce.Do called twice")
		return nil
	})

	if err != nil {
		t.Fatalf("SafeOnce.Do error should be nil")
	}
}

func TestSafeOnceError(t *testing.T) {
	var safeOnce SafeOnce

	err := safeOnce.Do(func() error {
		return fmt.Errorf("SafeOnce.Do error")
	})

	if err == nil {
		t.Fatalf("SafeOnce.Do error should not be nil")
	}

	calledTwice := false

	err = safeOnce.Do(func() error {
		calledTwice = true
		return nil
	})

	if !calledTwice {
		t.Fatalf("SafeOnce.Do should be called twice")
	}

	if err != nil {
		t.Fatalf("SafeOnce.Do error should be nil")
	}
}
