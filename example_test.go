package safeonce_test

import (
	"fmt"
	"github.com/koofr/safeonce"
)

func ExampleOnce() {
	var safeOnce safeonce.SafeOnce

	err := safeOnce.Do(func() error {
		return fmt.Errorf("Initialization error")
	})

	fmt.Println(err)

	err = safeOnce.Do(func() error {
		fmt.Println("Initialization successful")
		return nil
	})

	fmt.Println(err)

	err = safeOnce.Do(func() error {
		fmt.Println("This will not be executed")
		return nil
	})

	fmt.Println(err)

	// Output:
	// Initialization error
	// Initialization successful
	// <nil>
	// <nil>
}
