package lck_test

import (
	"fmt"

	"github.com/reiver/go-lck"
)

func ExampleMap_string_int() {

	var store lck.Map[string,int]

	fmt.Printf("INITIAL-STORE-LEN: %d\n", store.Len())

	store.Set("once", 1)
	store.Set("twice", 2)
	store.Set("thrice", 3)
	store.Set("fource", 4)

	fmt.Printf("STORE-LEN-AFTER-SETTING: %d\n", store.Len())

	value, found := store.Get("twice")
	fmt.Printf("STORE[twice]: %d, %t\n", value, found)

	fmt.Println()
	fmt.Println("KEY-VALUES:")
	store.For(func(key string, value int){
		fmt.Println()
		fmt.Printf("- KEY: %q\n", key)
		fmt.Printf("- VALUE: %d\n", value)
	})

	store.Unset("twice")

	fmt.Println()
	fmt.Println("KEY-VALUES:")
	store.For(func(key string, value int){
		fmt.Println()
		fmt.Printf("- KEY: %q\n", key)
		fmt.Printf("- VALUE: %d\n", value)
	})

	// Output:
	// INITIAL-STORE-LEN: 0
	// STORE-LEN-AFTER-SETTING: 4
	// STORE[twice]: 2, true
	// 
	// KEY-VALUES:
	// 
	// - KEY: "fource"
	// - VALUE: 4
	// 
	// - KEY: "once"
	// - VALUE: 1
	//
	// - KEY: "thrice"
	// - VALUE: 3
	// 
	// - KEY: "twice"
	// - VALUE: 2
	// 
	// KEY-VALUES:
	// 
	// - KEY: "fource"
	// - VALUE: 4
	// 
	// - KEY: "once"
	// - VALUE: 1
	//
	// - KEY: "thrice"
	// - VALUE: 3

}
